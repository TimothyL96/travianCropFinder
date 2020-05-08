package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	fmt.Println("******************************")
	fmt.Println("This is travian crop finder")
	fmt.Println("Use it at your own risk!")
	fmt.Println("******************************")

	// Start
	err := CreateError()

	// New reader
	r := bufio.NewReader(os.Stdin)

	// Create and initialize Crawler
	crawler := CreateCrawler()
	crawler.InitializeCrawler(err, r)

	var continueSearch string

	for {
		crawler.InitializeScrapMap(err, r)

		// Start scraping
		crawler.ScrapMap(err)

		// Wait for result
		crawler.collector.Wait()

		// Print the location results
		PrintResults(crawler.results)

		fmt.Print("Do you want to continue search? (yes/no): ")
		continueSearch, _ = r.ReadString('\n')
		if strings.Contains(strings.ToLower(continueSearch), "n") {
			fmt.Println("The end\nPress enter to quit")
			break
		}
	}

	_, _ = r.ReadString('\n')
}

func CreateCrawler() Crawler {
	return Crawler{
		collector:      nil,
		url:            "",
		xAxis:          0,
		yAxis:          0,
		searchDistance: 0,
		name:           "",
		password:       "",
	}
}

func CreateLocation(xAxis, yAxis, nrOfLumber, nrOfClay, nrOfIron, nrOfCrop int, isOasis bool) Location {
	return Location{
		xAxis:      xAxis,
		yAxis:      yAxis,
		nrOfLumber: nrOfLumber,
		nrOfClay:   nrOfClay,
		nrOfIron:   nrOfIron,
		nrOfCrop:   nrOfCrop,
		isOasis:    isOasis,
	}
}

func CreateError() Error {
	return Error{}
}

func PrintResults(locations []Location) {
	fmt.Println("Printing results:")
	fmt.Println("--------------------------")
	for _, v := range locations {
		PrintResult(v)
	}
}

func PrintResult(location Location) {
	fmt.Println("--------------------------")
	fmt.Println("X-Axis:", location.xAxis)
	fmt.Println("Y-Axis:", location.yAxis)
	fmt.Println("Oasis:", location.isOasis)
	if location.nrOfLumber != 0 {
		fmt.Println("Lumber:", location.Lumber())
	}
	if location.nrOfClay != 0 {
		fmt.Println("Clay:", location.Clay())
	}
	if location.nrOfIron != 0 {
		fmt.Println("Iron:", location.Iron())
	}
	if location.nrOfCrop != 0 {
		fmt.Println("Cropper:", location.Crop())
	}
	fmt.Println("--------------------------")
}
