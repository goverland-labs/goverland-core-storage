package zerion

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/rs/zerolog/log"
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
	resp, err := c.getResponse(
		http.MethodGet,
		"fungibles/",
		"fungibles-list",
		map[string]string{
			"currency":                       "usd",
			"filter[fungible_ids]":           ids,
			"filter[implementation_address]": address,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("get response: %w", err)
	}

	var list FungibleList
	if err = json.Unmarshal(resp, &list); err != nil {
		return nil, fmt.Errorf("unmarshal body: %w", err)
	}

	return &list, nil
}

func (c *Client) GetChains() ([]ChainData, error) {
	resp, err := c.getResponse(
		http.MethodGet,
		"chains",
		"chains",
		map[string]string{},
	)
	if err != nil {
		return nil, fmt.Errorf("get response: %w", err)
	}

	var chains Chains
	if err = json.Unmarshal(resp, &chains); err != nil {
		return nil, fmt.Errorf("unmarshal body: %w", err)
	}

	return chains.Data, nil
}

func (c *Client) GetFungibleData(id string) (*FungibleData, error) {
	resp, err := c.getResponse(
		http.MethodGet,
		fmt.Sprintf("fungibles/%s/", id),
		"fungible-data",
		map[string]string{
			"currency": "usd",
		},
	)
	if err != nil {
		return nil, fmt.Errorf("get response: %w", err)
	}

	var fungible Fungible
	if err = json.Unmarshal(resp, &fungible); err != nil {
		return nil, fmt.Errorf("unmarshal body: %w", err)
	}

	return &fungible.FungibleData, nil
}

func (c *Client) GetFungibleChart(id string, period string) (*ChartData, error) {
	resp, err := c.getResponse(
		http.MethodGet,
		fmt.Sprintf("fungibles/%s/charts/%s/", id, period),
		"fungible-chart",
		map[string]string{
			"currency": "usd",
		},
	)
	if err != nil {
		return nil, fmt.Errorf("get response: %w", err)
	}

	var chart Chart
	if err = json.Unmarshal(resp, &chart); err != nil {
		return nil, fmt.Errorf("unmarshal body: %w", err)
	}

	return &chart.ChartData, nil
}

func (c *Client) getResponse(method, subURL, alias string, params map[string]string) ([]byte, error) {
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

	request := c.withAuth(req)

	log.Info().Msgf("request %v", request)

	resp, err := c.client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("request do: %w", err)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read body: %w", err)
	}

	return body, err
}

func (c *Client) withAuth(req *http.Request) *http.Request {
	if c.authKey == "" {
		return req
	}

	req.Header.Add("Authorization", fmt.Sprintf("Basic %s", c.authKey))

	return req
}
