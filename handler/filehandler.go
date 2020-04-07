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

	rPool "filestore/cache/redis"
	"filestore/meta"
	"filestore/util"

	"github.com/garyburd/redigo/redis"
	"github.com/gin-gonic/gin"
)

// MultipartUploadInfo : initial info struct
type MultipartUploadInfo struct {
	Filehash   string
	Filesize   int
	UploadID   string
	ChunkSize  int64
	ChunkCount int
}

var lcoalPath string

// UploadHandler : 文件上传
func UploadHandler(c *gin.Context) {
	var wg sync.WaitGroup

	fpath := "/Users/zhangbicheng/Desktop/"
	file, fHead, err := c.Request.FormFile("file")

	if err != nil {
		fmt.Println("Failed to form file, err: ", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to form file",
		})
		return
	}

	defer file.Close()

	fileMeta := meta.FileMeta{
		FileName:  fHead.Filename,
		FileHash:  util.MD5([]byte(fHead.Filename)),
		Location:  fpath,
		FileSize:  fHead.Size,
		ChunkSize: 5 * 1024 * 1024,
		UploadAt:  time.Now().Format("2006-01-02 15:04:05"),
	}

	if fileMeta.FileSize < fileMeta.ChunkSize {
		fileMeta.ChunkCount = 1
	} else {
		fileMeta.ChunkCount = int(math.Ceil(float64(fileMeta.FileSize) / float64(fileMeta.ChunkSize)))
	}

	uploadID := initialMultipartUpload(fileMeta)

	bfReader := bufio.NewReader(file)

	buf := make([]byte, fileMeta.ChunkSize)

	for i := 0; i < fileMeta.ChunkCount; i++ {
		n, err := bfReader.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": "Failed to divide file",
				})
				return
			}
		}
		if n <= 0 {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Failed to divide file",
			})
		}
		wg.Add(1)

		bufCopied := make([]byte, fileMeta.ChunkSize)
		copy(bufCopied, buf)

		go func(b []byte, curIdx int) {
			defer wg.Done()
			uploadPart(b, uploadID, curIdx, fileMeta.Location)
		}(bufCopied[:n], i)
	}

	wg.Wait()

	if err = completeUpload(fileMeta, uploadID); err != nil {
		fmt.Println("Failed to complete upload, err: ", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to compete upload",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Multi-part upload compolete",
	})
	fmt.Println("Multi-part upload complete")

}

// initialMultipartUpload : 初始化分块上传
func initialMultipartUpload(fmeta meta.FileMeta) (uploadID string) {
	rConn := rPool.RedisPool().Get()
	defer rConn.Close()

	uploadID = "bee" + fmt.Sprintf("%x", time.Now().UnixNano())

	rConn.Do("HSET", "MP_"+uploadID, "chunkcount", fmeta.ChunkCount)
	rConn.Do("HSET", "MP_"+uploadID, "filehash", fmeta.FileHash)
	rConn.Do("HSET", "MP_"+uploadID, "filesize", fmeta.FileSize)

	return uploadID
}

func uploadPart(buf []byte, uploadID string, chunkIndex int, location string) (err error) {

	rConn := rPool.RedisPool().Get()
	index := strconv.Itoa(chunkIndex)
	defer rConn.Close()

	fpath := location + uploadID + "/" + index
	os.MkdirAll(path.Dir(fpath), 0744)
	fd, err := os.Create(fpath)

	if err != nil {
		return err
	}

	defer fd.Close()
	if _, err = fd.Write(buf); err != nil {
		return err
	}
	if _, err = rConn.Do("HSET", "MP_"+uploadID, "chkidx_"+index, 1); err != nil {
		return err
	}

	fmt.Println("chkidx_" + index + " upload success")
	return nil
}

// completeUpload : 通知上传合并
func completeUpload(fMeta meta.FileMeta, uploadID string) (err error) {
	fmt.Println("complete")
	rConn := rPool.RedisPool().Get()
	defer rConn.Close()

	data, err := redis.Values(rConn.Do("HGETALL", "MP_"+uploadID))
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

	fmt.Println("totalCount: ", totalCount)
	fmt.Println("chunkCount: ", chunkCount)

	if totalCount != chunkCount {
		return errors.New("invalid request")
	}

	fd, err := os.Create(fMeta.Location + uploadID + "/" + fMeta.FileName)

	if err != nil {
		return err
	}

	defer fd.Close()

	for i := 0; i < chunkCount; i++ {
		fpath := "/Users/zhangbicheng/Desktop/" + uploadID + "/" + strconv.Itoa(i)
		b, err := ioutil.ReadFile(fpath)
		if err != nil {
			return err
		}
		fd.Write(b)
	}

	return nil
}
