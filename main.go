package main

import (
	"fmt"
	rpio "github.com/stianeikeland/go-rpio"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"time"
)

const PinOffset = 14

var html = `
<html style="height:100%;">
<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <meta name="description" content="the controller of everything.">
    <title>Sean's Home</title>
    <link rel="stylesheet" href="http://yui.yahooapis.com/pure/0.6.0/pure-min.css">
    <link rel="stylesheet" href="http://yui.yahooapis.com/pure/0.6.0/grids-responsive-min.css">
    <link rel="stylesheet" href="css/blog.css">
</head>
<body style="height:100%;">
<div style="margin:0;height=100%;">
<style scoped>
.button-xlarge{
	font-size: 500%;
	height: 100%;
}
.button-green {
            background: rgb(28, 184, 65); /* this is a green */
}
.button-blue {
	background: rgb(66, 184, 221);
}
#open {
	width: 100%;
	height: 50%;
	font-size: 500%;
}
#close {
	width: 100%;
	height: 50%;
	font-size: 500%;
}
</style>
<a class="pure-button button-green" href="/open" align="center" id="open">Turn On</a>
<a class="pure-button button-blue" href="/close" align="center" id="close">Turn Off</a>
</div>
</body>
</html>`

var fav string

func ReadFile() error {
	file, err := os.Open("./ico.ico")
	if err != nil {
		return err
	}
	defer file.Close()
	body, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}
	fav = string(body)
	return nil
}

func open() error {
	err := rpio.Open()
	if err != nil {
		return err
	}
	pin := rpio.Pin(PinOffset)
	pin.Output()
	pin.High() // Low()
	rpio.Close()
	return nil
}

func close() error {
	err := rpio.Open()
	if err != nil {
		return err
	}
	pin := rpio.Pin(PinOffset)
	pin.Output()
	pin.Low() // close
	rpio.Close()
	return nil
}

func Close(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, html)
	go close()
}

func Open(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, html)
	go open()
}

func Main(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, html)
}

func Fav(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, fav)
}

func File(w http.ResponseWriter, req *http.Request) {
	u := req.RequestURI[5:]
	if u == "/" {
		io.WriteString(w, "Forbidden")
		return
	}
	fn := "/home/pi/txt" + u + ".mp4"
	http.ServeFile(w, req, fn)
}

func C() error {
	c, err := net.Dial("tcp4", "dil.fish:11400")
	if err != nil {
		fmt.Println("dial error", err)
		return err
	}
	c.Close()
	return nil
}

func Report() {
	for {
		err := C()
		if err != nil {
			fmt.Println("c error is", time.Now(), err)
		} else {
			fmt.Println("ip reported")
		}
		time.Sleep(time.Second * 5)
	}
}

func main() {
	err := ReadFile()
	if err != nil {
		panic(err)
	}
	go Report()
	http.HandleFunc("/", Main)
	http.HandleFunc("/file/", File)
	http.HandleFunc("/open", Open)
	http.HandleFunc("/close", Close)
	http.HandleFunc("/favicon.ico", Fav)
	err = http.ListenAndServe(":80", nil)
	if err != nil {
		log.Fatal("ListenAndServ", err)
	}
}
