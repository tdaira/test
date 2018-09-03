package crawler

import (
	"cloud.google.com/go/storage"
	"context"
	"encoding/json"
	"io"
	"log"
)

type SiteInfo struct {
	Seeds []string `json:"seeds"`
	ValidUrlRegex string `json:"valid_url_regex"`
	CrawlDelay int `json:"crawl_delay"`
	SameHostOnly bool `json:"same_host_only"`
	MaxVisits int `json:"max_visits"`
}

func (s SiteInfo) ToByteArray() ([]byte, error) {
	jsonBytes, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}
	return jsonBytes, nil
}

func GetSiteList() ([]SiteInfo, error) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, err
	}
	bucketName := "mirage-crawler"
	objectName := "site-list/list.jsonl"
	rc, err := client.Bucket(bucketName).Object(objectName).NewReader(ctx)
	if err != nil {
		return nil, err
	}
	defer rc.Close()

	dec := json.NewDecoder(rc)
	var siteList []SiteInfo
	for {
		var s SiteInfo
		if err := dec.Decode(&s); err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}
		siteList = append(siteList, s)
	}

	return siteList, err
}

func StringToSiteInfo(str string) (*SiteInfo, error) {
	jsonBytes := ([]byte)(str)
	siteInfo := new(SiteInfo)

	err := json.Unmarshal(jsonBytes, siteInfo);
	if err != nil {
		return nil, err
	}
	return siteInfo, nil
}

