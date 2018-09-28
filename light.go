package light

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
var TestMode = false


func SetTestMode(set bool) {
    TestMode = set
}


func (rh *rootHandler) onOff() {
    for {
        set, ok := <-rh.cStatus
        if ok == false {
            break
        }
        if TestMode == true {
            continue
        }
        err := rpio.Open()
        if err != nil {
            fmt.Println("open rpio", err)
            continue
        }
        pin := rpio.Pin(PinOffset)
        pin.Output()
        if set == true {
            pin.High()
        } else {
            pin.Low()
        }
        rpio.Close()
    }
}

func (rh *rootHandler) on() {
    Status = true
    rh.cStatus <-true
}

func (rh *rootHandler) off() {
    Status = false
    rh.cStatus <-false
}

func (static *staticHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    uri := r.URL.Path
	if uri == "/favicon.ico" {
		w.Write(static.Fav)
		return
	}
	w.Write(static.Page)
}

func (rh *rootHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println("request is", r.URL.Path, r.Method)
    uri := r.URL.Path
	if uri == "/api/on" {
		rh.on()
        w.Write([]byte("ok"))
		return
	}
	if uri == "/api/off" {
		rh.off()
        w.Write([]byte("ok"))
		return
	}
	if uri == "/api/status" {
		var status string
		if Status == true {
			status = "on"
		} else {
			status = "off"
		}
		w.Write([]byte(status))
	}
}

func ReadFile(fn string) ([]byte, error) {
	fmt.Println("read file", fn)
	file, err := os.Open(fn)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return ioutil.ReadAll(file)
}

func Handler(index, ico string) (http.Handler, error) {
	var rh rootHandler
	var sh staticHandler
    rh.cStatus = make(chan bool, 1)
    go rh.onOff()
    mux := http.NewServeMux()
	page, err := ReadFile(index)
	if err != nil {
		return nil, err
	}
	fav, err := ReadFile(ico)
	if err != nil {
		return nil, err
	}
	sh.Page = page
	sh.Fav = fav
	mux.Handle("/api/", &rh)
	mux.Handle("/", &sh)
    return mux, nil
}

type rootHandler struct{
    cStatus chan bool
}

type staticHandler struct {
	Page []byte
	Fav  []byte
}
