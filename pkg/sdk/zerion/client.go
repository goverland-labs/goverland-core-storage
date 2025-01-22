package zerion

import (
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
	"io"
	"net/http"
)

type (
	Client struct {
		client  *http.Client
		apiURL  string
		authKey string
	}
)

func NewClient(apiURL, authKey string, client *http.Client) *Client {
	return &Client{
		client:  client,
		apiURL:  apiURL,
		authKey: authKey,
	}
}

func (c *Client) GetFungibleList(ids string, address string) (*FungibleList, error) {
	req, err := c.buildRequest(
		http.MethodGet,
		"fungibles/",
		"fungibles-list",
		map[string]string{
			"currency":                       "usd",
			"filter[fungible_ids]":           ids,
			"filter[implementation_address]": address,
		},
	)
	log.Info().Msgf("request %v", req)
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request do: %w", err)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read body: %w", err)
	}

	var list FungibleList
	if err = json.Unmarshal(body, &list); err != nil {
		return nil, fmt.Errorf("unmarshal body: %w", err)
	}

	return &list, nil
}

func (c *Client) buildRequest(method, subURL, alias string, params map[string]string) (*http.Request, error) {
	req, err := http.NewRequest(
		method,
		fmt.Sprintf("%s/%s", c.apiURL, subURL),
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}

	q := req.URL.Query()
	for k, v := range params {
		if v != "" {
			q.Add(k, v)
		}
	}
	req.URL.RawQuery = q.Encode()

	req.Header.Add("alias", alias)

	return c.withAuth(req), nil
}

func (c *Client) withAuth(req *http.Request) *http.Request {
	if c.authKey == "" {
		return req
	}

	req.Header.Add("Authorization", fmt.Sprintf("Basic %s", c.authKey))

	return req
}
