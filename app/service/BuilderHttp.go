package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type ClientBuilder struct {
	HttpClient          *http.Client
	EndpointUrl         string
	UserAgent           string
	ApiUser             string
	ApiPassword         string
	AuthorizationHeader string
}

// ClientOptions a client options
type ClientOptions struct {
	UserAgent           string
	EndpointUrl         string
	Timeout             time.Duration
	ApiUser             string
	ApiPassword         string
	AuthorizationHeader string
}

func (c *ClientBuilder) SetAuthorizationHeader(bearerToken string) {
	c.AuthorizationHeader = bearerToken
}

// SetCustomTransport set new custom transport
func (c *ClientBuilder) SetCustomTransport(customHTTPTransport http.RoundTripper) {
	if c.HttpClient != nil {
		c.HttpClient.Transport = customHTTPTransport
	}
}

func (c *ClientBuilder) DoPostJson(path string, query map[string]string, v interface{}) (res *http.Response, err error) {
	body := new(bytes.Buffer)
	if err := json.NewEncoder(body).Encode(v); err != nil {
		return nil, err
	}

	res, err = c.Do(http.MethodPost, path, query, body, "application/json")
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *ClientBuilder) DoPutJson(path string, query map[string]string, v interface{}) error {
	body := new(bytes.Buffer)
	if err := json.NewEncoder(body).Encode(v); err != nil {
		return err
	}

	//nolint:bodyclose
	_, err := c.Do(http.MethodPut, path, query, body, "application/json")
	return err
}

func (c *ClientBuilder) DoPutStream(path string, query map[string]string, body []byte) error {
	bodyReader := bytes.NewReader(body)
	//nolint:bodyclose
	_, err := c.Do(http.MethodPut, path, query, bodyReader, "application/octet-stream")
	fmt.Printf("err: %v", err)
	return err
}

func (c *ClientBuilder) DoDelete(path string, query map[string]string) error {
	//nolint:bodyclose
	_, err := c.Do(http.MethodDelete, path, query, nil, "")
	return err
}

func (c *ClientBuilder) DoPost(path string, query map[string]string) (res *http.Response, err error) {
	return c.Do(http.MethodPost, path, query, nil, "")
}

func (c *ClientBuilder) Do(method, path string, query map[string]string, body io.Reader, contentType string) (res *http.Response, err error) {
	url, err := c.BuildUrl(path, query)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", c.UserAgent)
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}

	if c.AuthorizationHeader != "" {
		req.Header.Set("Authorization", c.AuthorizationHeader)
	} else {
		req.SetBasicAuth(c.ApiUser, c.ApiPassword)
	}

	res, err = c.HttpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to do request: %w", err)
	}
	if err := c.CheckResponse(res); err != nil {
		return nil, fmt.Errorf("failed to check response: %w", err)
	}

	return
}

func (c *ClientBuilder) DoRaw(method string, urlEndpoint string, query *map[string]string, data *interface{}, headers *map[string]interface{}) (res *http.Response, err error) {
	url, err := c.BuildUrlWithEndpoint(urlEndpoint, *query)
	if err != nil {
		return nil, err
	}
	body := new(bytes.Buffer)
	if err := json.NewEncoder(body).Encode(data); err != nil {
		return nil, err
	}
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	if headers != nil {
		for k, v := range *headers {
			req.Header.Set(k, fmt.Sprintf("%v", v))
		}
	}

	res, err = c.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	if err := c.CheckResponse(res); err != nil {
		return nil, err
	}

	return
}

func (c *ClientBuilder) DoGet(path string, query map[string]string) (res *http.Response, err error) {
	return c.Do(http.MethodGet, path, query, nil, "")
}

func (c *ClientBuilder) CheckResponse(res *http.Response) error {
	if res.StatusCode >= 200 && res.StatusCode <= 299 {
		return nil
	}

	defer res.Body.Close()

	if strings.Contains(res.Header.Get("Content-Type"), "application/json") {
		if res.StatusCode == 404 {
			return ErrorNotFound
		}

		jsonErr := &Error{}
		err := json.NewDecoder(res.Body).Decode(jsonErr)
		if err != nil {
			return fmt.Errorf("response error with status code %d: failed unmarshal error response: %w", res.StatusCode, err)
		}

		return fmt.Errorf("raw error : %v",jsonErr)
	}

	errText, err := ioutil.ReadAll(res.Body)
	if err == nil {
		return fmt.Errorf("response error with status code %d: %s", res.StatusCode, string(errText))
	}

	return fmt.Errorf("response error with status code %d", res.StatusCode)
}

func (c *ClientBuilder) ReadJsonResponse(res *http.Response, v interface{}) error {
	defer res.Body.Close()
	err := json.NewDecoder(res.Body).Decode(v)
	if err != nil {
		return err
	}
	return nil
}

func (c *ClientBuilder) BuildUrl(path string, query map[string]string) (string, error) {
	if len(query) == 0 {
		return c.EndpointUrl + path, nil
	}
	url, err := url.Parse(c.EndpointUrl + path)
	if err != nil {
		return "", err
	}

	q := url.Query()
	for k, v := range query {
		q.Set(k, v)
	}

	url.RawQuery = q.Encode()
	return url.String(), nil
}

func (c *ClientBuilder) BuildUrlWithEndpoint(urlEndpoint string, query map[string]string) (string, error) {
	if len(query) == 0 {
		return urlEndpoint, nil
	}
	url, err := url.Parse(urlEndpoint)
	if err != nil {
		return "", err
	}

	q := url.Query()
	for k, v := range query {
		q.Set(k, v)
	}

	url.RawQuery = q.Encode()
	return url.String(), nil
}

var ErrorNotFound = &Error{
	Type:    "NotFound",
	Message: "Not found",
}

// Error a custom error type
type Error struct {
	ErrorData 	 interface{} `json:"error"`
	Type    string `json:"type"`
	Message string `json:"message"`
}

// Error error message
func (e *Error) Error() string {
	if e.Message == "" {
		return "An error occurred " + fmt.Sprintf("%v", e.ErrorData)
	}
	return e.Message
}
