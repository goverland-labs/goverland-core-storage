package zerion

type (
	MarketData struct {
		Price float64 `json:"price"`
	}

	Attributes struct {
		Symbol     string     `json:"symbol"`
		MarketData MarketData `json:"market_data"`
	}

	FungibleData struct {
		ID         string     `json:"id"`
		Attributes Attributes `json:"attributes"`
	}

	FungibleList struct {
		List []FungibleData `json:"data"`
	}
)
