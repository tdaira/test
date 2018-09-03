package main

import (
	"fmt"
	"github.com/adjust/redismq"
	"github.com/tdaira/test/app/crawler"
)

func main() {
	siteList, err := crawler.GetSiteList()
	if err != nil {
		panic(err)
	}
	for _, site := range siteList {
		jsonBytes, err := site.ToByteArray()
		if err != nil {
			panic(err)
		}
		fmt.Println(string(jsonBytes))
	}
	workQueue := redismq.CreateQueue(
		"localhost", "6379", "", 0, "crawl_site")
	for _, site := range siteList {
		byteJson, err := site.ToByteArray()
		if err != nil {
			panic(err)
		}
		err = workQueue.Put(string(byteJson))
		if err != nil {
			panic(err)
		}
	}
}
