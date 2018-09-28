package main

import (
	"github.com/dilfish/light"
    "net/http"
)

func main() {
	go light.ReportIP()
    mux, err := light.Handler()
    if err != nil {
        panic(err)
    }
    panic(http.ListenAndServe(":80", mux))
}
