package zerion

import (
	"encoding/json"
	"fmt"
	"time"
)

type (
	MarketData struct {
		TotalSupply           float64 `json:"total_supply"`
		CirculatingSupply     float64 `json:"circulating_supply"`
		MarketCap             float64 `json:"market_cap"`
		FullyDilutedValuation float64 `json:"fully_diluted_valuation"`
		Price                 float64 `json:"price"`
	}

	Attributes struct {
		Name            string            `json:"name"`
		Symbol          string            `json:"symbol"`
		MarketData      MarketData        `json:"market_data"`
		Implementations []Implementations `json:"implementations"`
	}

	Implementations struct {
		ChainID  string `json:"chain_id"`
		Address  string `json:"address"`
		Decimals int    `json:"decimals"`
	}

	FungibleData struct {
		ID         string     `json:"id"`
		Attributes Attributes `json:"attributes"`
	}

	Fungible struct {
		FungibleData FungibleData `json:"data"`
	}

	FungibleList struct {
		List []FungibleData `json:"data"`
	}

	ChartData struct {
		ID              string          `json:"id"`
		ChartAttributes ChartAttributes `json:"attributes"`
	}

	Chart struct {
		ChartData ChartData `json:"data"`
	}

	ChartAttributes struct {
		BeginAt string  `json:"begin_at"`
		EndAt   string  `json:"end_at"`
		Stats   Stats   `json:"stats"`
		Points  []Point `json:"points"`
	}

	Stats struct {
		First float64 `json:"first"`
		Min   float64 `json:"min"`
		Avg   float64 `json:"avg"`
		Max   float64 `json:"max"`
		Last  float64 `json:"last"`
	}

	Point struct {
		Time  time.Time
		Price float64
	}

	Chains struct {
		Data []ChainData `json:"data"`
	}

	ChainData struct {
		ID         string          `json:"id"`
		Type       string          `json:"type"`
		Attributes ChainAttributes `json:"attributes"`
	}

	ChainAttributes struct {
		ExternalID string    `json:"external_id"`
		Name       string    `json:"name"`
		Icon       ChainIcon `json:"icon,omitempty"`
	}

	ChainIcon struct {
		URL string `json:"url,omitempty"`
	}
)

func (p *Point) UnmarshalJSON(data []byte) error {
	var v []interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		return fmt.Errorf("failed to parse data: %w", err)
	}
	if len(v) != 2 {
		return fmt.Errorf("incorrect elements in point object %v", v)
	}
	p.Time = time.Unix(int64(v[0].(float64)), 0)
	p.Price = v[1].(float64)
	return nil
}
