package entity

type CardDefault struct {
	ID                   string    `json:"id"`
	Type                 string    `json:"type"`
	Symbol               string    `json:"symbol"`
	Name                 string    `json:"name"`
	Family               string    `json:"family"`
	Heading              string    `json:"heading"`
	Description          string    `json:"description"`
	Cost                 *int      `json:"cost,omitempty"`
	Rule                 *string   `json:"rule,omitempty"`
	IsConditional        *bool     `json:"is_conditional,omitempty"`
	HasBabies            *bool     `json:"has_babies,omitempty"`
	Plus                 *bool     `json:"plus,omitempty"`
	Value                *int      `json:"value,omitempty"`
	Mortgage             *int      `json:"mortgage,omitempty"`
	DownPayment          *int      `json:"down_payment,omitempty"`
	CashFlow             *int      `json:"cash_flow,omitempty"`
	Price                *int      `json:"price,omitempty"`
	Count                *int      `json:"count,omitempty"`
	OnlyYou              *bool     `json:"only_you,omitempty"`
	Range                *[]int    `json:"range,omitempty"`
	SubRule              *[]string `json:"subRule,omitempty"`
	ApplicableToEveryOne *bool     `json:"applicable_to_every_one,omitempty"`
	Percent              *int      `json:"percent,omitempty"`
	Success              *[]int    `json:"success,omitempty"`
	Dices                *[]int    `json:"dices,omitempty"`
	CostPerOne           *float32  `json:"cost_per_one,omitempty"`
	ExtraDices           *int      `json:"extra_dices,omitempty"`
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
	Plus        bool    `json:"plus"`
	Cost        int     `json:"cost"`
	Value       int     `json:"value"`
	Mortgage    *int    `json:"mortgage"`
	DownPayment *int    `json:"down_payment"`
	CashFlow    *int    `json:"cash_flow"`
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
	Count       int    `json:"count"`
	OnlyYou     bool   `json:"only_you"`
	Range       []int  `json:"range"`
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

type CardDice struct {
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
	ApplicableToEveryOne *bool    `json:"applicable_to_every_one"`
}

type CardRiskBusiness struct {
	ID          string     `json:"id"`
	Type        string     `json:"type"`
	Dices       []CardDice `json:"dices"`
	ExtraDices  int        `json:"extra_dices"`
	Symbol      string     `json:"symbol"`
	Heading     string     `json:"heading"`
	Description string     `json:"description"`
	Cost        int        `json:"cost"`
}

type CardRiskStocks struct {
	ID          string     `json:"id"`
	Type        string     `json:"type"`
	Count       int        `json:"count"`
	Cost        int        `json:"cost"`
	Dices       []CardDice `json:"dices"`
	ExtraDices  int        `json:"extra_dices"`
	Symbol      string     `json:"symbol"`
	Heading     string     `json:"heading"`
	Description string     `json:"description"`
	CostPerOne  float64    `json:"cost_per_one"`
}
