package main

import (
	"bufio"
	"crypto/tls"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"regexp"

	"github.com/gin-gonic/gin"
)

var reg []*regexp.Regexp

func Proxy(c *gin.Context) {
	s := c.Param("path")
	if s == "" {
		c.String(http.StatusOK, "Server running.")
		return
	} else if s == "/reload" {
		err := LoadFile()
		if err != nil {
			c.String(http.StatusOK, err.Error())
		} else {
			c.String(http.StatusOK, "Whitelist reloaded.")
		}
		return
	}
	s = s[1:]

	flag := false
	for _, v := range reg {
		if v.FindStringIndex(s) != nil {
			flag = true
			break
		}
	}
	if !flag {
		c.String(http.StatusForbidden, "URL not in whitelist.")
		return
	}

	u, err := url.Parse(s)
	if err != nil {
		c.AbortWithStatus(500)
		return
	}
	r := c.Request
	r.URL = u
	r.Host = r.URL.Host
	r.RequestURI = u.RawQuery

	director := func(req *http.Request) {
		req = r
	}
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	proxy := &httputil.ReverseProxy{
		Transport: tr,
		Director:  director,
	}

	//data, _ := httputil.DumpRequest(r, true)
	//fmt.Printf("%s", string(data))

	proxy.ServeHTTP(c.Writer, r)

	if c.Writer.Status() == 301 || c.Writer.Status() == 302 {
		c.Redirect(http.StatusMovedPermanently, "/"+c.Writer.Header().Get("Location"))
	}
}

func LoadFile() error {
	var ret []*regexp.Regexp

	file, err := os.Open("./whitelist.txt")
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		ret = append(ret, regexp.MustCompile(scanner.Text()))
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	reg = ret
	return nil
}

func main() {
	err := LoadFile()
	if err != nil {
		panic(err)
	}
	r := gin.Default()
	r.Any("/*path", Proxy)
	_ = r.Run()
}
