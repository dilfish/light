package main

import (
	"net/http"

	"github.com/dilfish/light"
)

func main() {
	mux, err := light.Handler("../../index.html", "../../ico.ico")
	if err != nil {
		panic(err)
	}
	panic(http.ListenAndServe(":80", mux))
}
