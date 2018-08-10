package main

import (
	"github.com/tdaira/gocrawl"
	"net/http"
	"time"
	"github.com/PuerkitoBio/goquery"
	"crypto/sha1"
	"os"
	"bytes"
	"encoding/base64"
	"strings"
)

type ExampleExtender struct {
	gocrawl.DefaultExtender
}

func (x *ExampleExtender) Visit(ctx *gocrawl.URLContext, res *http.Response, doc *goquery.Document) (interface{}, bool) {
	bufBody := new(bytes.Buffer)
	bufBody.ReadFrom(res.Body)
	body := bufBody.Bytes()
	url := ctx.NormalizedURL()

	// Save data as a file with a name generated from sha1.
	sha := sha1.New()
	sha.Write([]byte(url.String()))
	shaStr := base64.URLEncoding.EncodeToString(sha.Sum(nil))
	outDir := "./data/" + strings.Replace(url.Host, ".", "_", -1)
	if _, err := os.Stat(outDir); os.IsNotExist(err) {
		os.Mkdir(outDir, os.ModePerm)
	}
	file, err := os.Create(outDir + "/" + shaStr)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	file.Write(body)

	return nil, true
}

func main() {
	opts := gocrawl.NewOptions(new(ExampleExtender))

	opts.RobotUserAgent = "Example"
	opts.UserAgent = "Mozilla/5.0 (compatible; Example/1.0; +http://example.com)"

	opts.CrawlDelay = 1 * time.Second
	opts.LogFlags = gocrawl.LogEnqueued

	opts.MaxVisits = 20

	c := gocrawl.NewCrawlerWithOptions(opts)
	c.Run([]string{"https://news.yahoo.co.jp/", "https://headlines.yahoo.co.jp/"})
}
