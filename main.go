package main

import (
	"fmt"

	"github.com/gocolly/colly/v2"
)

func main() {
	err := CreateError()

	// Create and initialize Crawler
	crawler := CreateCrawler()
	crawler.InitializeCrawler(err)

	// debug
	debug(crawler)
}

func CreateCrawler() Crawler {
	return Crawler{
		collector: colly.Collector{},
		url:       "",
		xAxis:     0,
		yAxis:     0,
	}
}

func CreateError() Error {
	return Error{}
}

func debug(crawler Crawler) {
	fmt.Println("URL:", crawler.url)
	fmt.Println("X-Axis:", crawler.xAxis)
	fmt.Println("Y-Axis:", crawler.yAxis)
	fmt.Println("Email:", crawler.email)
	fmt.Println("Password:", crawler.password)
}
