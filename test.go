package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

func multipartUpload(filename string, tgtURL string, chunkSize int) error {
	f, err := os.Open(filename)
	if err != nil {
		fmt.Println("Failed to open file, err: ", err.Error())
		return err
	}

	defer f.Close()

	bfRd := bufio.NewReader(f)
	index := 0

	ch := make(chan int)
	buf := make([]byte, chunkSize) // 指定每次读取chunkSize大小

	for {
		n, err := bfRd.Read(buf)
		if n <= 0 {
			break
		}

		index++

		bufCopied := make([]byte, 5*1048576)
		copy(bufCopied, buf)

		go func(b []byte, curIdx int) {
			fmt.Printf("upload_size: %d\n", len(b))

			resp, err := http.Post(
				tgtURL+"&index="+strconv.Itoa(curIdx),
				"multipart/form-data",
				bytes.NewReader(b))

			if err != nil {
				fmt.Println("Failed to upload part, err: ", err.Error())
			}

			body, err := ioutil.ReadAll(resp.Body)
			fmt.Println("%+v %+v\n", string(body), err)
			resp.Body.Close()

			ch <- curIdx
		}(bufCopied[:n], index)

		if err != nil {
			if err == io.EOF {
				break
			} else {
				fmt.Println("Failed to upload multipart, err: ", err.Error())
			}
		}

		for idx := 0; idx < index; idx++ {
			select {
			case res := <-ch:
				fmt.Println(res)
			}
		}
	}

	return nil
}

func main() {

	rConn.Do("SET", "name", "blue")
	rConn.Do("SET", "mykey", "superWang", "EX", "5")
	rConn.Do("PING")
	username := "admin"
	token := "54eefa7dbd5bcf852c52fecd816f2a315c61832c"
	filehash := "dfa39cac093a7a9c94d25130671ec474d51a2995"

	resp, err := http.PostForm(
		"http://localhost:7000/file/mpupload",
		url.Values{
			"username": {username},
			"token":    {token},
			"filehash": {filehash},
			"filesize": {"132489256"},
		})

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}

	uploadID := jsonit.Get(body, "data").Get("UploadID").ToString()
	chunkSize := jsonit.Get(body, "data").Get("ChunkSize").ToInt()
	fmt.Printf("uploadid: %s chunksize: %d\n", uploadID, chunkSize)

	filename := ""
	tgtURL = "http://localhost:7000/file/mpupload/uppart?" +
		"username=admin&token=" + token + "&uploadid=" + uploadID

	multipartUpload(filename, tURL, chunkSize)

	resp, err = http.PostForm(
		"http://localhost:7000/file/mpupload/complete",
		url.Values{
			"username": {username},
			"token":    {token},
			"filehash": {filehash},
			"filesize": {"132489256"},
			"filename": {"go1.10.3.linux-amd64.tar.gz"},
			"uploadid": {uploadID},
		})

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}

	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}
	fmt.Printf("complete result: %s\n", string(body))
}
