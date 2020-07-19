package handler

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/arstd/log"
	"github.com/zbcheng/filestore/conf"
	rPool "github.com/zbcheng/filestore/drivers/redis"
	"github.com/zbcheng/filestore/models"
	repo "github.com/zbcheng/filestore/repository"
	"github.com/zbcheng/filestore/util"

	"github.com/garyburd/redigo/redis"
	"github.com/gin-gonic/gin"
)

// UploadInfo : initial info struct
type UploadInfo struct {
	fileMeta   models.FileMeta
	UploadID   string
	ChunkSize  int64
	ChunkCount int
}

// UploadHandler : 文件上传
func UploadHandler(c *gin.Context) {
	var wg sync.WaitGroup

	dstPath := conf.Load().DstPath.Path
	file, fHead, err := c.Request.FormFile("file")

	if err != nil {
		fmt.Println("Failed to form file, err: ", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  "Failed to form file",
			"err":  1,
			"data": "",
		})
		return
	}

	defer file.Close()

	uploadInfo := UploadInfo{
		fileMeta: models.FileMeta{
			FileName: fHead.Filename,
			FileHash: util.Sha1([]byte(fHead.Filename)),
			Location: dstPath,
			FileSize: fHead.Size,
			UploadAt: time.Now().Format("2006-01-02 15:04:05"),
		},
		ChunkSize: 5 * 1024 * 1024,
	}

	if uploadInfo.fileMeta.FileSize < uploadInfo.ChunkSize {
		uploadInfo.ChunkCount = 1
	} else {
		uploadInfo.ChunkCount = int(math.Ceil(float64(uploadInfo.fileMeta.FileSize) / float64(uploadInfo.ChunkSize)))
	}

	exists, err := util.FileExists(uploadInfo.fileMeta.Location + uploadInfo.UploadID + "/" + uploadInfo.fileMeta.FileName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  "Failed to check file path!",
			"err":  1,
			"data": "",
		})
		return
	}

	if exists {
		c.JSON(http.StatusOK, gin.H{
			"msg":  "File already exists, upload done!",
			"err":  0,
			"data": uploadInfo.fileMeta.FileHash,
		})
		return
	}

	if uploadSuc := initialUpload(&uploadInfo); !uploadSuc {
		log.Error("Failed to init upload")
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  "Failed to init upload",
			"err":  1,
			"data": "",
		})
		return
	}

	bfReader := bufio.NewReader(file)

	buf := make([]byte, uploadInfo.ChunkSize)

	var uploadPartSuc = true

	for i := 0; i < uploadInfo.ChunkCount; i++ {
		n, err := bfReader.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg":  "Failed to divide file",
					"err":  1,
					"data": "",
				})
				return
			}
		}
		if n <= 0 {
			c.JSON(http.StatusInternalServerError, gin.H{
				"msg":  "Failed to divide file",
				"err":  1,
				"data": "",
			})
			return
		}
		wg.Add(1)

		bufCopied := make([]byte, uploadInfo.ChunkSize)
		copy(bufCopied, buf)

		go func(b []byte, curIdx int) {
			defer wg.Done()
			if err = uploadPart(b, curIdx, &uploadInfo); err != nil {
				log.Error("Failed to upload part: ", curIdx, err)
				uploadPartSuc = false
				return
			}
		}(bufCopied[:n], i)

	}

	wg.Wait()

	if !uploadPartSuc {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  "Failed to upload part",
			"err":  1,
			"data": "",
		})
		return
	}

	if err = completeUpload(&uploadInfo); err != nil {
		log.Error("Failed to complete upload, err: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  "Failed to compete upload",
			"err":  1,
			"data": "",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg":  "Upload file complete",
		"err":  1,
		"data": uploadInfo.fileMeta.FileHash,
	})
}

// initialUpload : 初始化分块上传
func initialUpload(up *UploadInfo) bool {

	var err error
	rConn := rPool.RedisPool().Get()
	defer rConn.Close()

	uploadID := "bee" + fmt.Sprintf("%x", time.Now().UnixNano())

	if _, err = rConn.Do("HSET", "MP_"+uploadID, "chunkcount", up.ChunkCount); err != nil {
		log.Error("Failed to hset chunkcount:", up.ChunkCount)
		return false
	}
	if _, err = rConn.Do("HSET", "MP_"+uploadID, "filehash", up.fileMeta.FileHash); err != nil {
		log.Error("Failed to hset filehash:", up.fileMeta.FileHash)
		return false
	}
	if _, err = rConn.Do("HSET", "MP_"+uploadID, "filesize", up.fileMeta.FileSize); err != nil {
		log.Error("Failed to hset filesize:", up.fileMeta.FileSize)
		return false
	}

	up.UploadID = uploadID

	return true
}

// uploadPart : 上传块文件
func uploadPart(buf []byte, chunkIndex int, up *UploadInfo) (err error) {

	rConn := rPool.RedisPool().Get()
	index := strconv.Itoa(chunkIndex)
	defer rConn.Close()

	fpath := up.fileMeta.Location + up.UploadID + "/" + index

	if err = os.MkdirAll(path.Dir(fpath), 0744); err != nil {
		return err
	}
	fd, err := os.Create(fpath)

	if err != nil {
		return err
	}

	defer fd.Close()
	if _, err = fd.Write(buf); err != nil {
		return err
	}
	if _, err = rConn.Do("HSET", "MP_"+up.UploadID, "chkidx_"+index, 1); err != nil {
		return err
	}

	log.Debug("chkidx_" + index + " upload success")
	return nil
}

// completeUpload : 合并分块
func completeUpload(up *UploadInfo) (err error) {
	rConn := rPool.RedisPool().Get()
	defer rConn.Close()

	data, err := redis.Values(rConn.Do("HGETALL", "MP_"+up.UploadID))
	if err != nil {
		return err
	}

	totalCount := 0
	chunkCount := 0

	for i := 0; i < len(data); i += 2 {
		k := string(data[i].([]byte))
		v := string(data[i+1].([]byte))

		if k == "chunkcount" {
			totalCount, _ = strconv.Atoi(v)
		} else if strings.HasPrefix(k, "chkidx_") && v == "1" {
			chunkCount++
		}
	}

	if totalCount != chunkCount {
		return errors.New("invalid request")
	}

	fd, err := os.Create(up.fileMeta.Location + up.UploadID + "/" + up.fileMeta.FileName)

	if err != nil {
		return err
	}

	defer fd.Close()

	for i := 0; i < chunkCount; i++ {
		fpath := up.fileMeta.Location + up.UploadID + "/" + strconv.Itoa(i)
		b, err := ioutil.ReadFile(fpath)
		if err != nil {
			return err
		}

		if _, err = fd.Write(b); err != nil {
			return err
		}

		if err = os.Remove(fpath); err != nil {
			return err
		}
	}

	if _, err = rConn.Do("DEL", "MP_"+up.UploadID); err != nil {
		fmt.Println("Failed to DEL:", err)
		return err
	}

	repo.StoreFileMeta(up.fileMeta)

	return nil
}
