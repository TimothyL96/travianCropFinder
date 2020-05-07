package main

import (
	"strconv"
)

func StringToInt(err Error, str string) (i int) {
	err.Set("Failed to convert string", str, "to integer")

	i, err.error = strconv.Atoi(str)
	if err.error != nil {
		panic(err.Error())
	}

	return i
}
