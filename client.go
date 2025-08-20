package godruid

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	DefaultEndPoint = "/druid/v2"
)

type Client struct {
	Url      string
	EndPoint string
	Timeout  time.Duration

	Debug        bool
	LastRequest  string
	LastResponse string
}

func (c *Client) Query(query Query) (err error) {
	return c.QueryWithContext(context.Background(), query)
}

func (c *Client) QueryWithContext(ctx context.Context, query Query) (err error) {
	query.setup()
	var reqJson []byte
	if c.Debug {
		reqJson, err = json.MarshalIndent(query, "", "  ")
	} else {
		reqJson, err = json.Marshal(query)
	}
	if err != nil {
		return
	}
	result, err := c.QueryRawWithContext(ctx, reqJson)
	if err != nil {
		return
	}

	return query.onResponse(result)
}

func (c *Client) QueryRaw(req []byte) (result []byte, err error) {
	return c.QueryRawWithContext(context.Background(), req)
}

func (c *Client) QueryRawWithContext(ctx context.Context, req []byte) (result []byte, err error) {
	if c.EndPoint == "" {
		c.EndPoint = DefaultEndPoint
	}
	endPoint := c.EndPoint
	if c.Debug {
		endPoint += "?pretty"
	}
	c.LastRequest = string(req)

	// By default, use 60 second timeout unless specified otherwise
	// by the caller
	clientTimeout := 60 * time.Second
	if c.Timeout != 0 {
		clientTimeout = c.Timeout
	}

	httpClient := &http.Client{
		Timeout: clientTimeout,
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", c.Url+endPoint, bytes.NewBuffer(req))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := httpClient.Do(httpReq)
	if err != nil {
		return
	}
	defer func() {
		resp.Body.Close()
	}()

	result, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	if c.Debug {
		c.LastResponse = string(result)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %s", resp.Status, string(result))
	}

	return
}
