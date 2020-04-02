package handler

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
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
	filehash := util.MD5([]byte(fileMeta.FileName))

	chunkSize := 5 * 1024 * 1024

	uploadID := initialMultipartUpload(fileSize, chunkSize, filehash)

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
	index := 0

	ch := make(chan int)
	buf := make([]byte, chunkSize)

	for {
		n, err := bfReader.Read(buf)
		if n <= 0 {
			break
		}

		index++

		bufCopied := make([]byte, chunkSize)
		copy(bufCopied, buf)

		go func(b []byte, curIdx int) {
			fmt.Printf("upload size: %d\n", len(b))

			_, err := http.Post(
				"http://127.0.0.1:7000/mpupload/uppart?uploadid="+uploadID+"&index="+strconv.Itoa(curIdx),
				"multipart/form-data",
				bytes.NewReader(b))

			if err != nil {
				fmt.Println(err)
			}

			ch <- curIdx
		}(bufCopied[:n], index)

		if err != nil {
			if err == io.EOF {
				break
			} else {
				fmt.Println("Failed to molti-part upload file, err: ", err.Error())
			}
		}

		for idx := 0; idx < index; idx++ {
			select {
			case res := <-ch:
				fmt.Println(res)
			}
		}

		if err = completeUpload(uploadID, filehash, fileMeta.FileName); err != nil {
			fmt.Println("Failed to complete upload, err: ", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Failed to complete upload",
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Multi-part upload compolete",
	})

}

// initialMultipartUpload : 初始化分块上传
func initialMultipartUpload(fileSize int64, chunkCount int, fileHash string) (uploadID string) {
	rConn := rPool.RedisPool().Get()
	defer rConn.Close()

	uploadID = "bee" + fmt.Sprintf("%x", time.Now().UnixNano())

	rConn.Do("HSET", "MP_"+uploadID, "chunkcount", chunkCount)
	rConn.Do("HSET", "MP_"+uploadID, "filehash", fileHash)
	rConn.Do("HSET", "MP_"+uploadID, "filesize", fileSize)

	return uploadID
}

func UploadPartHandler(c *gin.Context) {

	uploadID := c.Query("uploadid")
	chunkIndex := c.Query("index")
	rConn := rPool.RedisPool().Get()
	defer rConn.Close()

	fpath := "/Users/zhangbicheng/Desktop/" + uploadID + "/" + chunkIndex
	os.MkdirAll(path.Dir(fpath), 0744)
	fd, err := os.Create(fpath)

	if err != nil {
		fmt.Println(err)
	}

	defer fd.Close()

	buf := make([]byte, 1024*1024)

	for {
		n, err := c.Request.Body.Read(buf)
		fd.Write(buf[:n])
		if err != nil {
			break
		}
	}

	rConn.Do("HSET", "MP_"+uploadID, "chkidx_"+chunkIndex, 1)
}

// completeUpload : 通知上传合并
func completeUpload(uploadID string, fileHash string, fileName string) (err error) {
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

	if totalCount != chunkCount {
		return errors.New("invalid request")
	}

	return nil
}
