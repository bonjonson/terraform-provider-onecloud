// client/client.go
package client

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
    "time"
)

type Client struct {
    BaseURL    string
    HTTPClient *http.Client
    Token      string
}

func NewClient(baseURL, token string) (*Client, error) {
    if token == "" {
        return nil, fmt.Errorf("API token is required")
    }
    return &Client{
        BaseURL: baseURL,
        Token:   token,
        HTTPClient: &http.Client{
            Timeout: 30 * time.Second,
        },
    }, nil
}

func (c *Client) doRequest(method, path string, body interface{}) (*http.Response, error) {
    var reqBody []byte
    var err error
    if body != nil {
        reqBody, err = json.Marshal(body)
        if err != nil {
            return nil, err
        }
    }
    url := c.BaseURL + path
    req, err := http.NewRequest(method, url, bytes.NewBuffer(reqBody))
    if err != nil {
        return nil, err
    }
    req.Header.Set("Authorization", "Bearer "+c.Token)
    req.Header.Set("Content-Type", "application/json")
    return c.HTTPClient.Do(req)
}
