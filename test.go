package main

import (
	"fmt"
	"filestore/meta"
)

func main() {
	fMeta := meta.FileMeta{
		FileHash: "abbedc9d588b1a533f2da86d54dabb52",
		FileName: "lo.png",
		FileSize: 36430,
		Location: "/Users/zhangbicheng/Desktop/",
		UploadAt: "2020-04-10 11:15:05",
	}
	result, err := meta.FileExists(fMeta)
	if err != nil {
		fmt.Println("Error: ", err)
	}

	fmt.Println(result)
	// meta.RemoveFileMetaDB("5c1c585301f75038c9339d043af5cc6b")
}