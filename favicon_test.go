package faviconfetch

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http/httptest"
	"testing"
)

var (
	server    *httptest.Server
	serverUrl string
	reader    io.Reader
)

func init() {
	server = httptest.NewServer(HandlerOne())
	serverUrl = server.URL

}

// Test for vanilla /favicon.ico
func TestOne(t *testing.T) {

	uri := serverUrl
	fmt.Println(uri)
	faviconChannel := make(chan []byte)
	go func() {
		faviconChannel <- Fetch(uri)
	}()
	favicon := <-faviconChannel
	// get actual favicon to compare them to
	f, _ := ioutil.ReadFile("./favicons/ddgfavicon.ico")
	// compare favicon downloaded with favicon served by html template
	if bytes.Compare(favicon, f) != 0 {
		t.Errorf("Did not download correct favicon")
	}
	server.Close()
}

// Test for favicon with base url
func TestTwo(t *testing.T) {
	server = httptest.NewServer(HandlerTwo())
	serverUrl = server.URL
	uri := serverUrl
	faviconChannel := make(chan []byte)
	go func() {
		faviconChannel <- Fetch(uri)
	}()
	favicon := <-faviconChannel
	// get actual favicon to compare them to
	f, _ := ioutil.ReadFile("./favicons/ddgfavicon.ico")
	// compare favicon downloaded with favicon served by html template
	if bytes.Compare(favicon, f) != 0 {
		t.Errorf("Did not download correct favicon")
	}
	server.Close()

}

// Test for .gzip favicon
