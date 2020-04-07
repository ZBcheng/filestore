package main

import (
	"fmt"
	"strings"
)

func main() {
	fileName := "testtxt"
	fileType := strings.Split(fileName, ".")
	fmt.Println(fileType[len(fileType)-1])
}
