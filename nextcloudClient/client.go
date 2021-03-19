package nextcloudClient

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type Client struct {
	HostURL    string
	HTTPClient *http.Client
	username   string
	password   string
}

func NewClient(host, username, password string) *Client {
	c := Client{
		HostURL:    host + "/ocs/v1.php",
		username:   username,
		password:   password,
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
	}
	return &c
}

func (c *Client) addHeadersForBody(req *http.Request, contentLength int) {
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Content-Length", strconv.Itoa(contentLength))
}

func (c *Client) doRequest(req *http.Request) ([]byte, error) {
	req.Header.Set("OCS-APIRequest", "true")
	req.SetBasicAuth(c.username, c.password)

	response, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status: %d, body: %s", response.StatusCode, body)
	}
	return body, err
}

func doSimpleRequest(c *Client, method string, endpoint string, bodyData *url.Values) (bool, error) {

	req, err := http.NewRequest(method, endpoint, strings.NewReader(bodyData.Encode()))
	if err != nil {
		return false, err
	}
	if bodyData != nil {
		c.addHeadersForBody(req, len(bodyData.Encode()))
	} else {
		c.addHeadersForBody(req, 0)
	}
	body, err := c.doRequest(req)
	if err != nil {
		return false, err
	}

	response := SimpleResponse{}
	if err := xml.Unmarshal(body, &response); err != nil {
		return false, err
	}

	if response.RequestMeta.StatusCode != StatusSuccess {
		return false, errors.New(fmt.Sprintf("Api returned a status code %d indicating failure. Message: %s", response.RequestMeta.StatusCode, response.RequestMeta.Status))
	}

	return true, nil
}
