package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/gocolly/colly/v2"
)

type Crawler struct {
	collector colly.Collector
	url       string
	xAxis     int
	yAxis     int
	email     string
	password  string
}

func (c *Crawler) InitializeCrawler(err Error) {
	c.getInputs(err)
}

func (c *Crawler) getInputs(err Error) {
	err.Set("Failed to get input for crawler")

	// New reader
	r := bufio.NewReader(os.Stdin)

	c.getURL(err, r)
	c.getXAxis(err, r)
	c.getYAxis(err, r)
	c.getEmail(err, r)
	c.getPassword(err, r)
}

func (c *Crawler) getURL(err Error, r *bufio.Reader) {
	err.Set("Failed to read URL")

	var url string
	fmt.Println("Please input your Travian server URL (Ex: ts2.travian.com) :")

	url, err.error = r.ReadString('\n')
	if err.error != nil {
		panic(err.Error())
	}

	c.setURL(url)
}

func (c *Crawler) getXAxis(err Error, r *bufio.Reader) {
	err.Set("Failed to read X-axis")

	var xAxisStr string
	fmt.Println("Please input your search location X-axis (Ex: -23) :")

	xAxisStr, err.error = r.ReadString('\n')
	if err.error != nil {
		panic(err.Error())
	}

	xAxis := StringToInt(err, strings.TrimSpace(xAxisStr))

	c.setXAxis(xAxis)
}

func (c *Crawler) getYAxis(err Error, r *bufio.Reader) {
	err.Set("Failed to read Y-axis")

	var yAxisStr string
	fmt.Println("Please input your search location Y-axis (Ex: 15) :")

	yAxisStr, err.error = r.ReadString('\n')
	if err.error != nil {
		panic(err.Error())
	}

	yAxis := StringToInt(err, strings.TrimSpace(yAxisStr))

	c.setYAxis(yAxis)
}

func (c *Crawler) getEmail(err Error, r *bufio.Reader) {
	err.Set("Failed to read email")

	var email string
	fmt.Println("Please input your account email :")

	email, err.error = r.ReadString('\n')
	if err.error != nil {
		panic(err.Error())
	}

	c.setEmail(err, email)
}

func (c *Crawler) getPassword(err Error, r *bufio.Reader) {
	err.Set("Failed to read password")

	var password string
	fmt.Println("Please input your account password :")

	password, err.error = r.ReadString('\n')
	if err.error != nil {
		panic(err.Error())
	}

	c.setPassword(password)
}

func (c *Crawler) setURL(url string) {
	url = strings.TrimSpace(url)

	// Check if https exist in front of URL
	// If no, prepend it
	if len(url) <= 4 || (url[:4] != "http" && url[:5] != "https") {
		url = "https://" + url
	}

	c.url = url
}

func (c *Crawler) setXAxis(xAxis int) {
	c.xAxis = xAxis
}

func (c *Crawler) setYAxis(yAxis int) {
	c.yAxis = yAxis
}

func (c *Crawler) setEmail(err Error, email string) {
	err.Set("Invalid email address", email)

	validEmail := regexp.MustCompile("(?:[a-z0-9!#$%&'*+/=?^_`{|}~-]+(?:\\.[a-z0-9!#$%&'*+/=?^_`{|}~-]+)*|\"(?:[\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x21\\x23-\\x5b\\x5d-\\x7f]|\\\\[\\x01-\\x09\\x0b\\x0c\\x0e-\\x7f])*\")@(?:(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\\.)+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?|\\[(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?|[a-z0-9-]*[a-z0-9]:(?:[\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x21-\\x5a\\x53-\\x7f]|\\\\[\\x01-\\x09\\x0b\\x0c\\x0e-\\x7f])+)])")
	if !validEmail.MatchString(email) {
		panic(err.Error())
	}

	c.email = strings.TrimSpace(email)
}

func (c *Crawler) setPassword(password string) {
	c.password = strings.TrimSpace(password)
}
