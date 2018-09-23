package main

import (
	"github.com/dilfish/light"
)

func main() {
	go light.ReportIP()
	panic(light.Handler())
}
