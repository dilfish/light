package light

import (
	"testing"
    "github.com/appleboy/gofight"
    "github.com/stretchr/testify/assert"
    "net/http"
)

func TestReadFile(t *testing.T) {
	bt, err := ReadFile("testdata/ReadFile.txt")
	if err != nil || string(bt) != "abcd\n" {
		t.Error("expect nil and abcd, got", string(bt), err)
	}
}


func DoTest() {
    SetTestMode(true)
}


type Case struct {
    t *testing.T
    status string
    fav []byte
    root []byte
}

func (c *Case) Status (r gofight.HTTPResponse, rq gofight.HTTPRequest) {
    assert.Equal(c.t, c.status, r.Body.String())
    assert.Equal(c.t, http.StatusOK, r.Code)
}


func (c *Case) OnOff (r gofight.HTTPResponse, rq gofight.HTTPRequest) {
    assert.Equal(c.t, "ok", r.Body.String())
    assert.Equal(c.t, http.StatusOK, r.Code)
}


func (c *Case) Root (r gofight.HTTPResponse, rq gofight.HTTPRequest) {
    assert.Equal(c.t, string(c.root), r.Body.String())
    assert.Equal(c.t, http.StatusOK, r.Code)
}

func (c *Case) Ico (r gofight.HTTPResponse, rq gofight.HTTPRequest) {
    assert.Equal(c.t, string(c.fav), r.Body.String())
    assert.Equal(c.t, http.StatusOK, r.Code)
}


func TestHandler(t *testing.T) {
    DoTest()
    h, err := Handler("a", "b")
    if err == nil {
        t.Error("open empty 1")
    }
    h, err = Handler("index.html", "b")
    if err == nil {
        t.Error("open empty 2")
    }
    h, err = Handler("index.html", "ico.ico")
    if err != nil {
        t.Error("get handler", err)
    }
    r := gofight.New()
    var c Case
    c.t = t
    c.status = "off"
    r.GET("/api/status").SetDebug(true).Run(h, c.Status)
    r.GET("/api/on").SetDebug(true).Run(h, c.OnOff)
    c.status = "on"
    r.GET("/api/status").SetDebug(true).Run(h, c.Status)
    r.GET("/api/off").SetDebug(true).Run(h, c.OnOff)
    root, err := ReadFile("index.html")
    if err != nil {
        t.Error("read index", err)
    }
    fav, err := ReadFile("ico.ico")
    if err != nil {
        t.Error("read ico", err)
    }
    c.fav = fav
    c.root = root
    r.GET("/").SetDebug(true).Run(h, c.Root)
    r.GET("/favicon.ico").SetDebug(true).Run(h, c.Ico)
}
