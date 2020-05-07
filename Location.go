package main

import (
	"strconv"
)

type Location struct {
	xAxis      int
	yAxis      int
	nrOfLumber int
	nrOfClay   int
	nrOfIron   int
	nrOfCrop   int
	isOasis    bool
}

func (l Location) Lumber() string {
	str := strconv.Itoa(l.nrOfLumber)

	if l.isOasis {
		return str + "%"
	}

	return str
}

func (l Location) Clay() string {
	str := strconv.Itoa(l.nrOfClay)

	if l.isOasis {
		return str + "%"
	}

	return str
}

func (l Location) Iron() string {
	str := strconv.Itoa(l.nrOfIron)

	if l.isOasis {
		return str + "%"
	}

	return str
}

func (l Location) Crop() string {
	str := strconv.Itoa(l.nrOfCrop)

	if l.isOasis {
		return str + "%"
	}

	return str
}
