package graphql

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
)

// Client is a client for interacting with a GraphQL API.
type Client struct {
	endpoint         string
	httpClient       *http.Client
	useMultipartForm bool
}

// NewClient makes a new Client capable of making GraphQL requests.
func NewClient(endpoint string) *Client {
	c := &Client{
		endpoint: endpoint,
	}

	if c.httpClient == nil {
		c.httpClient = http.DefaultClient
	}
	return c
}

func (c *Client) Run(ctx context.Context, req *Request, resp interface{}) error {
	var requestBody bytes.Buffer
	requestBodyObj := struct {
		Query     string                 `json:"query"`
		Variables map[string]interface{} `json:"variables"`
	}{
		Query:     req.q,
		Variables: req.vars,
	}
	if err := json.NewEncoder(&requestBody).Encode(requestBodyObj); err != nil {
		return err
	}

	r, err := http.NewRequest(http.MethodPost, c.endpoint, &requestBody)
	if err != nil {
		return err
	}
	r.Header.Set("Content-Type", "application/json; charset=utf-8")
	r.Header.Set("Accept", "application/json; charset=utf-8")
	for key, values := range req.Header {
		for _, value := range values {
			r.Header.Add(key, value)
		}
	}

	r = r.WithContext(ctx)
	res, err := c.httpClient.Do(r)
	if err != nil {
		return err
	}

	defer res.Body.Close()
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, res.Body); err != nil {
		return err
	}

	if err := json.NewDecoder(&buf).Decode(&resp); err != nil {
		return err
	}

	return nil
}

type Request struct {
	q      string
	vars   map[string]interface{}
	Header http.Header
}

func NewRequest(q string) *Request {
	req := &Request{
		q:      q,
		Header: make(map[string][]string),
	}
	return req
}
