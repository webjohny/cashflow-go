package entity

type Card struct {
	ID                   string           `json:"id"`
	Type                 string           `json:"type"`
	Symbol               string           `json:"symbol"`
	Name                 string           `json:"name"`
	Family               string           `json:"family"`
	Heading              string           `json:"heading"`
	Description          string           `json:"description"`
	Cost                 *int             `json:"cost,omitempty"`
	Rule                 *string          `json:"rule,omitempty"`
	IsConditional        *bool            `json:"is_conditional,omitempty"`
	HasBabies            *bool            `json:"has_babies,omitempty"`
	Plus                 *bool            `json:"plus,omitempty"`
	Value                *int             `json:"value,omitempty"`
	Mortgage             *int             `json:"mortgage,omitempty"`
	DownPayment          *int             `json:"down_payment,omitempty"`
	CashFlow             *int             `json:"cash_flow,omitempty"`
	Price                *int             `json:"price,omitempty"`
	Count                *int             `json:"count,omitempty"`
	Increase             *int             `json:"increase,omitempty"`
	Decrease             *int             `json:"decrease,omitempty"`
	OnlyYou              *bool            `json:"only_you,omitempty"`
	Range                *[]int           `json:"range,omitempty"`
	SubRule              *[]string        `json:"sub_rule,omitempty"`
	ApplicableToEveryOne *bool            `json:"applicable_to_every_one,omitempty"`
	Percent              *int             `json:"percent,omitempty"`
	Success              *[]int           `json:"success,omitempty"`
	Dices                *[]CardRiskDices `json:"dices,omitempty"`
	CostPerOne           *float32         `json:"cost_per_one,omitempty"`
	ExtraDices           *int             `json:"extra_dices,omitempty"`
}

func (c *Card) MultiUserFlow() bool {
	multiFlows := map[string]map[string]bool{
		"deal": {
			"stock": true,
		},
		"market": {
			"goldCoins":  true,
			"realEstate": true,
		},
	}

	if _, ok := multiFlows[c.Family]; ok {
		if _, okType := multiFlows[c.Family][c.Type]; okType {
			return multiFlows[c.Family][c.Type]
		}
	}

	return false
}

type CardDoodad struct {
	ID            string `json:"id"`
	Type          string `json:"type"`
	Heading       string `json:"heading"`
	Description   string `json:"description"`
	Cost          int    `json:"cost"`
	Rule          string `json:"rule"`
	IsConditional bool   `json:"is_conditional"`
	HasBabies     bool   `json:"has_babies"`
}

type CardRealEstate struct {
	ID          string  `json:"id"`
	Type        string  `json:"type"`
	Symbol      string  `json:"symbol"`
	Heading     string  `json:"heading"`
	Description string  `json:"description"`
	Rule        *string `json:"rule"`
	Cost        int     `json:"cost"`
	Mortgage    *int    `json:"mortgage"`
	DownPayment *int    `json:"down_payment"`
	CashFlow    *int    `json:"cash_flow"`
}

type CardMarketRealEstate struct {
	ID          string   `json:"id"`
	Type        string   `json:"type"`
	Heading     string   `json:"heading"`
	Symbol      string   `json:"symbol"`
	Description string   `json:"description"`
	Rule        string   `json:"rule"`
	Plus        bool     `json:"plus"`
	SubRule     []string `json:"sub_rule"`
	Value       int      `json:"value"`
}

type CardMarketDamage struct {
	ID                   string   `json:"id"`
	Type                 string   `json:"type"`
	Heading              string   `json:"heading"`
	Symbol               string   `json:"symbol"`
	Description          string   `json:"description"`
	Rule                 string   `json:"rule"`
	SubRule              []string `json:"sub_rule"`
	Cost                 int      `json:"cost"`
	ApplicableToEveryOne bool     `json:"applicable_to_every_one"`
}

type CardBusiness struct {
	ID          string  `json:"id"`
	Type        string  `json:"type"`
	Symbol      string  `json:"symbol"`
	Heading     string  `json:"heading"`
	Description string  `json:"description"`
	Rule        *string `json:"rule"`
	Cost        int     `json:"cost"`
	CashFlow    *int    `json:"cash_flow"`
}

type CardStocks struct {
	ID          string `json:"id"`
	Type        string `json:"type"`
	Symbol      string `json:"symbol"`
	Heading     string `json:"heading"`
	Description string `json:"description"`
	Rule        string `json:"rule"`
	Price       int    `json:"price"`
	Count       *int   `json:"count"`
	Increase    *int   `json:"increase"`
	Decrease    *int   `json:"decrease"`
	OnlyYou     *bool  `json:"only_you"`
	Range       *[]int `json:"range"`
}

type CardPreciousMetals struct {
	ID          string `json:"id"`
	Type        string `json:"type"`
	Cost        int    `json:"cost"`
	Count       int    `json:"count"`
	Symbol      string `json:"symbol"`
	Heading     string `json:"heading"`
	Description string `json:"description"`
}

type CardDream struct {
	ID          string `json:"id"`
	Heading     string `json:"heading"`
	Description string `json:"description"`
	Type        string `json:"type"`
	Cost        int    `json:"cost"`
}

type CardCharity struct {
	ID          string `json:"id"`
	Heading     string `json:"heading"`
	Description string `json:"description"`
	Cost        int    `json:"cost"`
	Type        string `json:"type"`
	Symbol      string `json:"symbol"`
}

type CardPayTax struct {
	ID          string `json:"id"`
	Heading     string `json:"heading"`
	Description string `json:"description"`
	Type        string `json:"type"`
	Percent     int    `json:"percent"`
}

type CardDownsized struct {
	ID          string `json:"id"`
	Heading     string `json:"heading"`
	Description string `json:"description"`
	Type        string `json:"type"`
	Percent     int    `json:"percent"`
}

type CardSmallDeal struct {
	ID                   string    `json:"id"`
	Type                 string    `json:"type"`
	Cost                 *int      `json:"cost"`
	Count                *int      `json:"count"`
	Symbol               string    `json:"symbol"`
	Heading              string    `json:"heading"`
	Description          string    `json:"description"`
	Percent              *int      `json:"percent"`
	Rule                 *string   `json:"rule"`
	Price                *int      `json:"price"`
	OnlyYou              *bool     `json:"only_you"`
	Range                *[]int    `json:"range"`
	SubRule              *[]string `json:"subRule"`
	ApplicableToEveryOne *bool     `json:"applicable_to_every_one"`
}

type CardRiskDices struct {
	Dices      []int    `json:"dices"`
	CashFlow   *int     `json:"cash_flow"`
	CostPerOne *float32 `json:"cost_per_one"`
}

type CardMarket struct {
	ID                   string   `json:"id"`
	Type                 string   `json:"type"`
	Heading              string   `json:"heading"`
	Symbol               string   `json:"symbol"`
	Description          string   `json:"description"`
	Rule                 string   `json:"rule"`
	Plus                 *bool    `json:"plus"`
	SubRule              []string `json:"sub_rule"`
	Success              *[]int   `json:"success"`
	Cost                 *int     `json:"cost"`
	Value                *int     `json:"value"`
	ApplicableToEveryOne *bool    `json:"applicable_to_every_one"`
}

type CardRiskBusiness struct {
	ID          string          `json:"id"`
	Type        string          `json:"type"`
	Dices       []CardRiskDices `json:"dices"`
	ExtraDices  int             `json:"extra_dices"`
	Symbol      string          `json:"symbol"`
	Heading     string          `json:"heading"`
	Description string          `json:"description"`
	Cost        int             `json:"cost"`
}

type CardRiskStocks struct {
	ID          string          `json:"id"`
	Type        string          `json:"type"`
	Count       int             `json:"count"`
	Cost        int             `json:"cost"`
	Dices       []CardRiskDices `json:"dices"`
	ExtraDices  int             `json:"extra_dices"`
	Symbol      string          `json:"symbol"`
	Heading     string          `json:"heading"`
	Description string          `json:"description"`
	CostPerOne  float32         `json:"cost_per_one"`
}
