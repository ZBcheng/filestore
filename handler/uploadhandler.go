package handler

import (
	"bufio"
	"database/sql"
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
	pg "filestore/db/postgres"
	"filestore/meta"
	"filestore/util"

	"github.com/garyburd/redigo/redis"
	"github.com/gin-gonic/gin"
)

// UploadInfo : initial info struct
type UploadInfo struct {
	fileMeta   meta.FileMeta
	UploadID   string
	ChunkSize  int64
	ChunkCount int
}

var pgConn *sql.DB

func init() {
	pgConn = pg.DBConn()
}

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

	uploadInfo := UploadInfo{
		fileMeta: meta.FileMeta{
			FileName: fHead.Filename,
			FileHash: util.MD5([]byte(fHead.Filename)),
			Location: fpath,
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

	exists, err := meta.FileExists(uploadInfo.fileMeta)
	if err != nil {
		fmt.Println("Failed to check file")
	}

	if exists {
		c.JSON(http.StatusOK, gin.H{
			"message": "File already exists, upload done!",
		})
		return
	}

	initialUpload(&uploadInfo)

	bfReader := bufio.NewReader(file)

	buf := make([]byte, uploadInfo.ChunkSize)

	for i := 0; i < uploadInfo.ChunkCount; i++ {
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
			return
		}
		wg.Add(1)

		bufCopied := make([]byte, uploadInfo.ChunkSize)
		copy(bufCopied, buf)

		go func(b []byte, curIdx int) {
			defer wg.Done()
			uploadPart(b, curIdx, &uploadInfo)
		}(bufCopied[:n], i)
	}

	wg.Wait()

	if err = completeUpload(&uploadInfo); err != nil {
		fmt.Println("Failed to complete upload, err: ", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to compete upload",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Upload file complete",
	})
}

// initialUpload : 初始化分块上传
func initialUpload(up *UploadInfo) {
	rConn := rPool.RedisPool().Get()
	defer rConn.Close()

	uploadID := "bee" + fmt.Sprintf("%x", time.Now().UnixNano())

	rConn.Do("HSET", "MP_"+uploadID, "chunkcount", up.ChunkCount)
	rConn.Do("HSET", "MP_"+uploadID, "filehash", up.fileMeta.FileHash)
	rConn.Do("HSET", "MP_"+uploadID, "filesize", up.fileMeta.FileSize)

	up.UploadID = uploadID

}

// uploadPart : 上传块文件
func uploadPart(buf []byte, chunkIndex int, up *UploadInfo) (err error) {

	rConn := rPool.RedisPool().Get()
	index := strconv.Itoa(chunkIndex)
	defer rConn.Close()

	fpath := up.fileMeta.Location + up.UploadID + "/" + index
	os.MkdirAll(path.Dir(fpath), 0744)
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

	fmt.Println("chkidx_" + index + " upload success")
	return nil
}

// completeUpload : 合并分块
func completeUpload(up *UploadInfo) (err error) {
	fmt.Println("Upload complete!")
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

		fd.Write(b)

		if err = os.Remove(fpath); err != nil {
			return err
		}
	}

	fmt.Println("Write file complete!")
	rConn.Do("DEL", "MP_"+up.UploadID)

	stmt, err := pgConn.Prepare("INSERT INTO tbl_file values($1, $2, $3, $4, $5)")
	if err != nil {
		return err
	}

	defer stmt.Close()

	if _, err = stmt.Exec(up.fileMeta.FileHash, up.fileMeta.FileName,
		up.fileMeta.FileSize, up.fileMeta.Location, up.fileMeta.UploadAt); err != nil {
		return err
	}

	return nil
}
