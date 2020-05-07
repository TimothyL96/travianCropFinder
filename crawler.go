package main

import (
	"bufio"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
)

type Crawler struct {
	collector      *colly.Collector
	url            string
	xAxis          int
	yAxis          int
	searchDistance int
	name           string
	password       string
	results        []Location
}

func (c *Crawler) InitializeCrawler(err Error, r *bufio.Reader) {
	err.Set("Error initializing crawler")

	c.collector = colly.NewCollector(
		colly.Async(false), // No async at the start as login is required
		colly.AllowURLRevisit(),
	)

	err.error = c.collector.Limit(
		&colly.LimitRule{
			DomainGlob:  "*",
			Delay:       100 * time.Millisecond,
			Parallelism: 10,
		},
	)
	if err.error != nil {
		panic(err.Error())
	}

	// Set error handler
	c.collector.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	// Set request handler
	c.collector.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	// On response
	c.collector.OnResponse(func(r *colly.Response) {
		if r.Request.URL.String() == c.getLoginURL() {
			err.Set("Login failed!\nPlease check your server URL, login ID and password.")
			panic(err.Error())
		}
		fmt.Println("Visited", r.Request.URL)
	})

	// Set scraping
	c.collector.OnHTML("#tileDetails", c.RetrieveMapDetails)

	// Request for user data
	c.getInputs(err, r)
}

func (c *Crawler) getInputs(err Error, r *bufio.Reader) {
	err.Set("Failed to get input for crawler")

	c.getURL(err, r)
	c.getEmailOrUsername(err, r)
	c.getPassword(err, r)

	// Authenticate after getting ID and password
	c.AuthenticateUser(err)
}

func (c *Crawler) getURL(err Error, r *bufio.Reader) {
	err.Set("Failed to read URL")

	var url string
	fmt.Print("Please input your Travian server URL (Ex: ts2.travian.com): ")

	url, err.error = r.ReadString('\n')
	if err.error != nil {
		panic(err.Error())
	}

	c.setURL(url)
}

func (c *Crawler) getXAxis(err Error, r *bufio.Reader) {
	err.Set("Failed to read X-axis")

	var xAxisStr string
	fmt.Print("Please input your search location starting X-axis (Ex: -23): ")

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
	fmt.Print("Please input your search location starting Y-axis (Ex: 15): ")

	yAxisStr, err.error = r.ReadString('\n')
	if err.error != nil {
		panic(err.Error())
	}

	yAxis := StringToInt(err, strings.TrimSpace(yAxisStr))

	c.setYAxis(yAxis)
}

func (c *Crawler) getEmailOrUsername(err Error, r *bufio.Reader) {
	err.Set("Failed to read name or username")

	var name string
	fmt.Print("Please input your account email or username: ")

	name, err.error = r.ReadString('\n')
	if err.error != nil {
		panic(err.Error())
	}

	c.setName(name)
}

func (c *Crawler) getPassword(err Error, r *bufio.Reader) {
	err.Set("Failed to read password")

	var password string
	fmt.Print("Please input your account password: ")

	password, err.error = r.ReadString('\n')
	if err.error != nil {
		panic(err.Error())
	}

	c.setPassword(password)
}

func (c *Crawler) getSearchDistance(err Error, r *bufio.Reader) {
	err.Set("Failed to get search distance")

	var distanceStr string
	fmt.Print("Please input your intended search distance from the center: ")

	distanceStr, err.error = r.ReadString('\n')
	if err.error != nil {
		panic(err.Error())
	}

	distance := StringToInt(err, strings.TrimSpace(distanceStr))

	c.setSearchDistance(err, distance)
}

func (c *Crawler) setURL(url string) {
	url = strings.TrimSpace(url)

	// Remove http if exist
	if len(url) > 7 && url[:7] == "http://" {
		url = url[7:]
	}

	if len(url) <= 4 || url[:5] != "https" {
		// Check if https exist in front of URL
		// If no, prepend it
		url = "https://" + url
	}

	// Remove extra slash at the end
	if len(url) > 0 && url[len(url)-1:] == "/" {
		url = url[:len(url)-1]
	}

	c.url = url
}

func (c *Crawler) setXAxis(xAxis int) {
	c.xAxis = xAxis
}

func (c *Crawler) setYAxis(yAxis int) {
	c.yAxis = yAxis
}

func (c *Crawler) setName(name string) {
	c.name = strings.TrimSpace(name)
}

func (c *Crawler) setPassword(password string) {
	c.password = strings.TrimSpace(password)
}

func (c *Crawler) setSearchDistance(err Error, distance int) {
	if distance < 0 {
		err.Set("Search distance cannot be negative!")
		panic(err.Error())
	}

	c.searchDistance = distance
}

func (c *Crawler) AuthenticateUser(err Error) {
	err.Set("Failed to login user", c.name, "for server", c.url)

	fmt.Println("Logging in...")
	loginUrl := c.getLoginURL()

	err.error = c.collector.Post(loginUrl,
		map[string]string{
			"name":     c.name,
			"password": c.password,
		})

	if err.error != nil {
		panic(err.Error())
	}

	// Turn on async capability
	c.collector.Async = true
}

func (c *Crawler) getLoginURL() string {
	return c.url + "/login.php"
}

func (c *Crawler) RetrieveMapDetails(e *colly.HTMLElement) {
	var xAxis, yAxis, nrOfLumber, nrOfClay, nrOfIron, nrOfCrop int
	var isOasis bool

	// Check if this is an oasis
	if e.DOM.HasClass("oasis") {
		isOasis = true
	}

	// fmt.Println("URL", e.Request.URL)
	xAxis = StringToInt(Error{}, e.Request.URL.Query().Get("x"))
	yAxis = StringToInt(Error{}, e.Request.URL.Query().Get("y"))

	// Get land resources distributions
	queryAll := "#map_details #distribution tr"
	querySingleResValue := ".val"
	occupiedLand := false

	if !e.DOM.Add(queryAll).ChildrenFiltered("td").HasClass("val") {
		occupiedLand = true
		queryAll += " td"
	}

	var resType int

	e.ForEach(queryAll, func(i int, e1 *colly.HTMLElement) {
		var resNrStr, resTypeClass string
		var resNr int

		if occupiedLand {
			resNrStr = e1.Text
			resTypeClass, _ = e1.DOM.ChildrenFiltered("i").Attr("class")
		} else {
			resNrStr = e1.ChildText(querySingleResValue)
			resTypeClass, _ = e1.DOM.ChildrenFiltered(".ico").ChildrenFiltered("i").Attr("class")
		}

		switch resTypeClass {
		case "r1":
			resType = 1
		case "r2":
			resType = 2
		case "r3":
			resType = 3
		case "r4":
			resType = 4
		}

		// Normal land
		if !isOasis {
			resNr = StringToInt(Error{}, strings.TrimSpace(resNrStr))
		} else {
			// Oasis
			reg, _ := regexp.Compile("[^a-zA-Z0-9]+")
			resNrStr = reg.ReplaceAllString(resNrStr, "")

			// resNrStr = strings.ReplaceAll(resNrStr, "%", "")
			resNr = StringToInt(Error{}, strings.TrimSpace(resNrStr))
		}

		switch resType {
		case 1:
			nrOfLumber = resNr
		case 2:
			nrOfClay = resNr
		case 3:
			nrOfIron = resNr
		case 4:
			nrOfCrop = resNr
		}
	})

	if resType >= 1 && resType <= 4 {
		l := CreateLocation(xAxis,
			yAxis,
			nrOfLumber,
			nrOfClay,
			nrOfIron,
			nrOfCrop,
			isOasis,
		)

		c.results = append(c.results, l)
	}
}

func (c *Crawler) InitializeScrapMap(err Error, r *bufio.Reader) {
	c.getXAxis(err, r)
	c.getYAxis(err, r)
	c.getSearchDistance(err, r)
}

func (c *Crawler) ScrapMap(err Error) {
	err.Set("Function scrap map error")

	// Traverse X-Axis from left to right
	for x := c.xAxis - c.searchDistance; x <= c.xAxis+c.searchDistance; x++ {
		// Traverse Y-Axis from down to up
		for y := c.yAxis - c.searchDistance; y <= c.yAxis+c.searchDistance; y++ {
			// Visit the center map
			err.error = c.collector.Visit(fmt.Sprintf("%s/position_details.php?x=%d&y=%d", c.url, x, y))
			if err.error != nil {
				panic(err.Error())
			}
		}
	}
}
