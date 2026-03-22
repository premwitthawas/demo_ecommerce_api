package pkg_elasticsearch

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch/v9"
)

func NewElasticsearch(urlWithComma, username, password string) (*elasticsearch.Client, error) {
	urls := strings.Split(urlWithComma, ",")
	cfg := elasticsearch.Config{
		Addresses: urls,
		Username:  username,
		Password:  password,
		Transport: &http.Transport{
			MaxIdleConns:          10,
			ResponseHeaderTimeout: time.Second * 5,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return nil, err
	}
	res, err := es.Info()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to elasticsearch: %w", err)
	}
	defer res.Body.Close()
	return es, nil
}
