package main

import (
	"bufio"
	"crypto/tls"
	"io"
	"net/http"
	"os"
	"regexp"

	"github.com/gin-gonic/gin"
)

var reg []*regexp.Regexp

func Proxy(c *gin.Context) {
	s := c.Request.URL.String()
	if s == "/" {
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

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	req, err := http.NewRequest(c.Request.Method, s, c.Request.Body)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	for header, values := range c.Request.Header {
		for _, value := range values {
			req.Header.Add(header, value)
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		c.Status(http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	for header, values := range resp.Header {
		for _, value := range values {
			c.Writer.Header().Add(header, value)
		}
	}
	c.Status(resp.StatusCode)
	io.Copy(c.Writer, resp.Body)
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
