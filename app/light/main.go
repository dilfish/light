package main

import (
	"github.com/dilfish/light"
    "net/http"
)

func main() {
    mux, err := light.Handler("../../index.html", "../../ico.ico")
    if err != nil {
        panic(err)
    }
    panic(http.ListenAndServe(":80", mux))
}
