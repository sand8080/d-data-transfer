package fetch

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

const (
	ReportFormatCSV string = "csv"
	ReportFormatXML string = "XML"
)

type Client struct {
	http.Client
	url     string
	authKey string
}

func NewClient(url, authKey string) *Client {
	cli := Client{url: url, authKey: authKey}
	cli.Timeout = time.Second * 10
	return &cli
}

func (c Client) GetReports(from, to time.Time) ([]byte, error) {

	form := url.Values{
		"datefrom": {from.Format("2006-01-02")},
		"dateto":   {to.Format("2006-01-02")},
		"reportoption": {"csv"},
	}

	body := bytes.NewBufferString(form.Encode())
	fmt.Println(form.Encode())
	req, err := http.NewRequest(http.MethodGet, c.url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Auth", c.authKey)
	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	return ioutil.ReadAll(resp.Body)
	//return nil, nil
}
