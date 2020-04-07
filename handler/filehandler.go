package handler

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"sync"
	"time"

	rPool "moviesite-filestore/cache/redis"
	"moviesite-filestore/meta"
	"moviesite-filestore/util"

	"github.com/garyburd/redigo/redis"
	"github.com/gin-gonic/gin"
)

// MultipartUploadInfo : initial info struct
type MultipartUploadInfo struct {
	Filehash   string
	Filesize   int
	UploadID   string
	ChunkSize  int
	ChunkCount int
}

var lcoalPath string

// MultipartUploadHandler : 分块上传
func MultipartUploadHandler(c *gin.Context) {
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
		FileName: fHead.Filename,
		Location: fpath + fHead.Filename,
		UploadAt: time.Now().Format("2006-01-02 15:04:05"),
	}

	fileInfo, err := os.Stat(fileMeta.Location)
	if err != nil {
		fmt.Println("Failed to get fileInfo, err: ", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to get file info",
		})
		return
	}

	fileSize := fileInfo.Size()
	fileMeta.FileSize = fileSize
	fileMeta.FileHash = util.MD5([]byte(fileMeta.FileName))

	var chunkSize int64
	chunkSize = 5 * 1024 * 1024
	var chunkCount int

	if fileSize < chunkSize {
		chunkCount = 1
	} else {
		chunkCount = int(math.Ceil(float64(fileSize) / float64(chunkSize)))
	}

	uploadID := initialMultipartUpload(fileMeta, chunkCount)

	f, err := os.Open(fileMeta.Location)

	if err != nil {
		fmt.Println("Failed to open file, err: ", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to open file",
		})
		return
	}

	defer f.Close()

	bfReader := bufio.NewReader(f)

	buf := make([]byte, chunkSize)

	for i := 0; i < chunkCount; i++ {
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

		bufCopied := make([]byte, chunkSize)
		copy(bufCopied, buf)

		go func(b []byte, curIdx int) {
			defer wg.Done()
			uploadPart(b, uploadID, i)
		}(bufCopied[:n], i)
	}

	wg.Wait()

	if err = completeUpload(uploadID); err != nil {
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
func initialMultipartUpload(fmeta meta.FileMeta, chunkCount int) (uploadID string) {
	rConn := rPool.RedisPool().Get()
	defer rConn.Close()

	uploadID = "bee" + fmt.Sprintf("%x", time.Now().UnixNano())

	rConn.Do("HSET", "MP_"+uploadID, "chunkcount", chunkCount)
	rConn.Do("HSET", "MP_"+uploadID, "filehash", fmeta.FileHash)
	rConn.Do("HSET", "MP_"+uploadID, "filesize", fmeta.FileSize)

	return uploadID
}

func uploadPart(buf []byte, uploadID string, chunkIndex int) (err error) {

	rConn := rPool.RedisPool().Get()
	index := strconv.Itoa(chunkIndex)
	defer rConn.Close()

	fpath := "/Users/zhangbicheng/Desktop/" + uploadID + "/" + index
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
func completeUpload(uploadID string) (err error) {
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

	return nil
}
