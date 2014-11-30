package faviconfetch

import (
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
)

// download favicon
func GetFavicon(uri string, faviconUri string) []byte {
	if faviconUri == "" {
		return nil
	}
	var favicon string
	// remove any new line characters from url
	re := regexp.MustCompile("%20")

	faviconUri = re.ReplaceAllString(faviconUri, "")
	client := &http.Client{}
	// get request for favicon
	req, err := SetHTTPHeaders(faviconUri)
	if err != nil {
		return nil
	}
	resp, err := client.Do(req)
	if os.Getenv("DEBUG") != "" {
		DumpResponse(resp)
	}
	if err != nil {
		log.Println(err)
		return nil
	}
	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("contents error: %s", err)
		return nil
	}
	matchHtml, err := regexp.MatchString("(?i)<html>", string(contents))
	if !matchHtml && string(contents) != "" {
		contentType := http.DetectContentType(contents)
		// return plain text as nil
		matchText, _ := regexp.MatchString("text", contentType)
		if matchText {
			return nil
		}
		// unzip any compression
		matchZip, err := regexp.MatchString("(?i)zip", contentType)
		if err != nil {
			log.Print("Regex error")
			return nil
		}
		if !matchZip {
			favicon := contents
			return favicon
		} else {
			faviconGzipReader, err := gzip.NewReader(req.Body)
			faviconContents, err := gzip.NewReader(faviconGzipReader)
			favicon, err := ioutil.ReadAll(faviconContents)
			return favicon
			if err != nil {
				log.Print("Could not read favicon")
				return nil
			}
		}
	} else {
		log.Print("Favicon was HTML")
		return nil
	}
	if favicon == "" {
		log.Print("Favicon was empty")
		return nil
	}
	new_uri := faviconUri + "/favicon.ico"
	if faviconUri != new_uri {
		return GetFavicon(uri, new_uri)
	}
	log.Println("Failed to download favicon")
	// failed to download favicon, give up.
	return nil
}

func DumpResponse(resp *http.Response) {
	fmt.Println("\n")
	fmt.Printf("Url: %s\n", resp.Request.URL)
	fmt.Printf("Status: %s\n", resp.Status)
	fmt.Printf("Headers: %s\n", resp.Header)
	//fmt.Printf("Body: %s\n", resp.Body)

}
