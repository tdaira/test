package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/base64"
	"github.com/PuerkitoBio/goquery"
	"github.com/tdaira/gocrawl"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

type ExampleExtender struct {
	gocrawl.DefaultExtender
	ValidURLRegex *regexp.Regexp
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
	x.writeFileWithDir(outDir, shaStr, body)

	return nil, true
}

func (x *ExampleExtender) Filter(ctx *gocrawl.URLContext, isVisited bool) bool {
	return !isVisited && x.ValidURLRegex.MatchString(ctx.NormalizedURL().String())
}

func (x *ExampleExtender) writeFileWithDir(dir string, fileName string, body []byte) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.Mkdir(dir, os.ModePerm)
	}
	file, err := os.Create(dir + "/" + fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	file.Write(body)
}

func main() {
	opts := gocrawl.NewOptions(&ExampleExtender{
		ValidURLRegex: regexp.MustCompile(`(^http://news\.yahoo\.co\.jp/flash$)|(^http://headlines\.yahoo\.co\.jp/hl\?.*)`)})

	opts.RobotUserAgent = "MirageBot"
	opts.UserAgent = "Mozilla/5.0 (compatible; MirageBot/1.0; +http://miragebot.com)"

	opts.CrawlDelay = 1 * time.Second
	opts.LogFlags = gocrawl.LogAll
	opts.SameHostOnly = false

	opts.MaxVisits = 20

	c := gocrawl.NewCrawlerWithOptions(opts)
	c.Run([]string{"https://news.yahoo.co.jp/flash"})
}
