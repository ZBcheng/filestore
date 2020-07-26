package main

import "github.com/zbcheng/filestore/app/api"

func main() {
	router := api.RegisterRouter()
	router.Run(":7000")

}
