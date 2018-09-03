package crawler

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
)

type MirageBotExtender struct {
	gocrawl.DefaultExtender
	ValidURLRegex *regexp.Regexp
}

func (x *MirageBotExtender) Visit(ctx *gocrawl.URLContext, res *http.Response, doc *goquery.Document) (interface{}, bool) {
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

func (x *MirageBotExtender) Filter(ctx *gocrawl.URLContext, isVisited bool) bool {
	return !isVisited && x.ValidURLRegex.MatchString(ctx.NormalizedURL().String())
}

func (x *MirageBotExtender) writeFileWithDir(dir string, fileName string, body []byte) {
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
