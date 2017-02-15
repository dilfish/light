package main

import (
	"fmt"
	rpio "github.com/stianeikeland/go-rpio"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"time"
)

const PinOffset = 14

var Status = false

func on() error {
	err := rpio.Open()
	if err != nil {
		fmt.Println("open error", err)
		return err
	}
	pin := rpio.Pin(PinOffset)
	pin.Output()
	pin.High() // Low()
	rpio.Close()
	Status = true
	return nil
}

func off() error {
	err := rpio.Open()
	if err != nil {
		return err
	}
	pin := rpio.Pin(PinOffset)
	pin.Output()
	pin.Low() // close
	rpio.Close()
	Status = false
	return nil
}

func (static *staticHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.RequestURI == "/favicon.ico" {
		w.Write(static.Fav)
		return
	}
	w.Write(static.Page)
}

func (rh *rootHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println("request is", r.RequestURI, r.Method)
	if r.RequestURI == "/api/on" {
		go on()
		return
	}
	if r.RequestURI == "/api/off" {
		go off()
		return
	}
	if r.RequestURI == "/api/status" {
		var status string
		if Status == true {
			status = "on"
		} else {
			status = "off"
		}
		w.Write([]byte(status))
	}
}

func readFile(fn string) ([]byte, error) {
	fmt.Println("read file", fn)
	file, err := os.Open(fn)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return ioutil.ReadAll(file)
}

func Handler() error {
	var rh rootHandler
	var sh staticHandler
	page, err := readFile("index.html")
	if err != nil {
		return err
	}
	fav, err := readFile("ico.ico")
	if err != nil {
		return err
	}
	sh.Page = page
	sh.Fav = fav
	http.Handle("/api/", &rh)
	http.Handle("/", &sh)
	fmt.Println("listen on 80")
	return http.ListenAndServe(":80", nil)
}

type rootHandler struct{}
type staticHandler struct {
	Page []byte
	Fav  []byte
}

func ReportIP() {
	for {
		_, err := http.PostForm("https://dil.fish/util/homeip", url.Values{"key": {"Value"}, "id": {"123"}})
		if err != nil {
			fmt.Println("set home ip", err)
		}
		time.Sleep(time.Second)
	}
}

func main() {
	go ReportIP()
	panic(Handler())
}
