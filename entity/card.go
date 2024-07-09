package entity

var LotteryTypes = struct {
	Money    string
	CashFlow string
}{
	Money:    "money",
	CashFlow: "cashflow",
}

var BusinessTypes = struct {
	Startup string
	Limited string
}{
	Startup: "startup",
	Limited: "limited",
}

var OtherAssetTypes = struct {
	Piece string
	Whole string
}{
	Piece: "piece",
	Whole: "whole",
}

var RealEstateTypes = struct {
	Building string
	Single   string
}{
	Building: "building",
	Single:   "single",
}

type Card struct {
	ID                   string        `json:"id"`
	Type                 string        `json:"type"`
	Symbol               string        `json:"symbol"`
	Name                 string        `json:"name"`
	Family               string        `json:"family"`
	Heading              string        `json:"heading"`
	Description          string        `json:"description"`
	Cost                 int           `json:"cost,omitempty"`
	Rule                 string        `json:"rule,omitempty"`
	IsConditional        bool          `json:"is_conditional,omitempty"`
	HasBabies            bool          `json:"has_babies,omitempty"`
	Plus                 bool          `json:"plus,omitempty"`
	Value                int           `json:"value,omitempty"`
	Mortgage             int           `json:"mortgage,omitempty"`
	DownPayment          int           `json:"down_payment,omitempty"`
	CashFlow             int           `json:"cash_flow,omitempty"`
	Price                int           `json:"price,omitempty"`
	AssetType            string        `json:"asset_type,omitempty"`
	Count                int           `json:"count,omitempty"`
	History              []CardHistory `json:"history,omitempty"`
	Increase             int           `json:"increase,omitempty"`
	Decrease             int           `json:"decrease,omitempty"`
	OnlyYou              bool          `json:"only_you,omitempty"`
	Range                []int         `json:"range,omitempty"`
	SubRule              []string      `json:"sub_rule,omitempty"`
	Lottery              string        `json:"lottery,omitempty"`
	Failure              []int         `json:"failure,omitempty"`
	ApplicableToEveryOne bool          `json:"applicable_to_every_one,omitempty"`
	Percent              int           `json:"percent,omitempty"`
	Success              []int         `json:"success,omitempty"`
	CostPerOne           int           `json:"cost_per_one,omitempty"`
	ExtraDices           int           `json:"extra_dices,omitempty"`
	Limit                int           `json:"limit,omitempty"`
	Outcome              struct {
		Failure int `json:"failure"`
		Success int `json:"success"`
	} `json:"outcome,omitempty"`
}

func (c *Card) MultiUserFlow() bool {
	multiFlows := map[string]map[string]bool{
		"deal": {
			"stock": true,
		},
		"market": {
			"other":      true,
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

func (c *CardStocks) SetCardHistory(history CardHistory) {
	history.SumCost()

	if c.History != nil {
		data := c.History

		var check bool

		for index, cH := range data {
			if cH.Price == history.Price {
				data[index].Count += history.Count
				data[index].SumCost()
				check = true

				break
			}
		}

		if !check {
			data = append(data, history)
		}

		c.History = data
	} else {
		histories := append([]CardHistory{}, history)
		c.History = histories
	}
}

func (c *CardBusiness) SetCardHistory(history CardHistory) {
	history.SumCost()

	if len(c.History) > 0 {
		data := c.History

		var check bool

		for index, cH := range data {
			if cH.Price == history.Price {
				data[index].Count += history.Count
				data[index].SumCost()
				check = true

				break
			}
		}

		if !check {
			data = append(data, history)
		}

		c.History = data
	} else {
		histories := append([]CardHistory{}, history)
		c.History = histories
	}
}

type CardHistory struct {
	Cost  int `json:"cost"`
	Price int `json:"price"`
	Count int `json:"count"`
}

func (h *CardHistory) SumCost() {
	h.Cost = h.Price * h.Count
}

type CardDoodad struct {
	ID            string `json:"id"`
	Type          string `json:"type"`
	Heading       string `json:"heading"`
	Symbol        string `json:"symbol,omitempty"`
	Description   string `json:"description"`
	Cost          int    `json:"cost"`
	Rule          string `json:"rule,omitempty"`
	IsConditional bool   `json:"is_conditional"`
	HasBabies     bool   `json:"has_babies"`
}

type CardRealEstate struct {
	ID          string `json:"id"`
	Type        string `json:"type"`
	Symbol      string `json:"symbol"`
	Heading     string `json:"heading"`
	Description string `json:"description"`
	Rule        string `json:"rule,omitempty"`
	IsOwner     bool   `json:"is_owner,omitempty"`
	AssetType   string `json:"asset_type"`
	Cost        int    `json:"cost"`
	Mortgage    int    `json:"mortgage,omitempty"`
	DownPayment int    `json:"down_payment,omitempty"`
	CashFlow    int    `json:"cash_flow,omitempty"`
	Percent     int    `json:"percent,omitempty"`
	Count       int    `json:"count,omitempty"`
	WholeCost   int    `json:"-"`
}

type CardMarketRealEstate struct {
	ID          string   `json:"id"`
	Type        string   `json:"type"`
	Heading     string   `json:"heading"`
	Symbol      string   `json:"symbol"`
	Description string   `json:"description"`
	Rule        string   `json:"rule,omitempty"`
	Plus        bool     `json:"plus,omitempty"`
	AssetType   string   `json:"asset_type"`
	SubRule     []string `json:"sub_rule,omitempty"`
	Cost        int      `json:"cost,omitempty"`
	Range       []int    `json:"range,omitempty"` //for 2nd, 4th, 8th ... flats building
}

type CardMarketOtherAssets struct {
	ID          string   `json:"id"`
	Type        string   `json:"type"`
	Heading     string   `json:"heading"`
	Symbol      string   `json:"symbol"`
	Description string   `json:"description"`
	Rule        string   `json:"rule,omitempty"`
	AssetType   string   `json:"asset_type"`
	SubRule     []string `json:"sub_rule,omitempty"`
	Cost        int      `json:"cost,omitempty"`
}

type CardMarketBusiness struct {
	ID          string   `json:"id"`
	Type        string   `json:"type"`
	Heading     string   `json:"heading"`
	Symbol      string   `json:"symbol"`
	Description string   `json:"description"`
	Rule        string   `json:"rule"`
	SubRule     []string `json:"sub_rule"`
	AssetType   string   `json:"asset_type"`
	Plus        bool     `json:"plus"`
	Cost        int      `json:"cost,omitempty"`
	CashFlow    int      `json:"cash_flow,omitempty"`
}

func (c *CardMarketBusiness) Fill(card Card) {
	c.ID = card.ID
	c.Type = card.Type
	c.Heading = card.Heading
	c.Symbol = card.Symbol
	c.Description = card.Description
	c.Rule = card.Rule
	c.AssetType = card.AssetType
	c.Plus = card.Plus
	c.SubRule = card.SubRule
	c.Cost = card.Cost
	c.CashFlow = card.CashFlow
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
	ID          string        `json:"id"`
	Type        string        `json:"type"`
	Symbol      string        `json:"symbol"`
	Heading     string        `json:"heading"`
	Description string        `json:"description"`
	Rule        string        `json:"rule,omitempty"`
	Cost        int           `json:"cost"`
	Limit       int           `json:"limit,omitempty"`
	IsOwner     bool          `json:"is_owner,omitempty"`
	AssetType   string        `json:"asset_type,omitempty"`
	Count       int           `json:"count,omitempty"`
	ExtraDices  int           `json:"extra_dices,omitempty"`
	History     []CardHistory `json:"history,omitempty"`
	Mortgage    int           `json:"mortgage,omitempty"`
	DownPayment int           `json:"down_payment,omitempty"`
	CashFlow    int           `json:"cash_flow,omitempty"`
	Percent     int           `json:"percent,omitempty"`
	WholeCost   int           `json:"-"`
}

type CardStocks struct {
	ID          string        `json:"id"`
	Type        string        `json:"type"`
	Symbol      string        `json:"symbol"`
	Heading     string        `json:"heading"`
	Description string        `json:"description"`
	Rule        string        `json:"rule"`
	Price       int           `json:"price"`
	Count       int           `json:"count,omitempty"`
	History     []CardHistory `json:"history,omitempty"`
	Increase    int           `json:"increase,omitempty"`
	Decrease    int           `json:"decrease,omitempty"`
	OnlyYou     bool          `json:"only_you,omitempty"`
	Range       []int         `json:"range,omitempty"`
}

func (c *CardStocks) Fill(card Card) {
	c.ID = card.ID
	c.Type = card.Type
	c.Heading = card.Heading
	c.Symbol = card.Symbol
	c.Description = card.Description
	c.Rule = card.Rule
	c.Price = card.Price
	c.Increase = card.Increase
	c.Decrease = card.Decrease
	c.Count = card.Count
	c.OnlyYou = card.OnlyYou
	c.Range = card.Range
}

type CardOtherAssets struct {
	ID          string `json:"id"`
	Type        string `json:"type"`
	Cost        int    `json:"cost"`
	CostPerOne  int    `json:"cost_per_one,omitempty"`
	Count       int    `json:"count,omitempty"`
	AssetType   string `json:"asset_type"`
	IsOwner     bool   `json:"is_owner"`
	Symbol      string `json:"symbol"`
	Heading     string `json:"heading"`
	Description string `json:"description"`
	WholeCost   int    `json:"-"`
}

type CardLottery struct {
	ID          string   `json:"id"`
	Type        string   `json:"type"`
	Symbol      string   `json:"symbol"`
	Heading     string   `json:"heading"`
	Description string   `json:"description"`
	Cost        int      `json:"cost"`
	Lottery     string   `json:"lottery"`
	Rule        string   `json:"rule"`
	SubRule     []string `json:"sub_rule"`
	Failure     []int    `json:"failure"`
	Success     []int    `json:"success"`
	Outcome     struct {
		Failure int `json:"failure"`
		Success int `json:"success"`
	} `json:"outcome"`
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
	ID                   string   `json:"id"`
	Type                 string   `json:"type"`
	Cost                 int      `json:"cost,omitempty"`
	Count                int      `json:"count,omitempty"`
	Symbol               string   `json:"symbol"`
	Heading              string   `json:"heading"`
	Description          string   `json:"description"`
	Percent              int      `json:"percent,omitempty"`
	Rule                 string   `json:"rule,omitempty"`
	Price                int      `json:"price,omitempty"`
	OnlyYou              bool     `json:"only_you,omitempty"`
	Range                []int    `json:"range,omitempty"`
	SubRule              []string `json:"subRule,omitempty"`
	ApplicableToEveryOne bool     `json:"applicable_to_every_one,omitempty"`
}

type CardMarket struct {
	ID                   string   `json:"id"`
	Type                 string   `json:"type"`
	Heading              string   `json:"heading"`
	Symbol               string   `json:"symbol"`
	Description          string   `json:"description"`
	Rule                 string   `json:"rule"`
	Plus                 bool     `json:"plus,omitempty"`
	SubRule              []string `json:"sub_rule"`
	Success              []int    `json:"success,omitempty"`
	Cost                 int      `json:"cost,omitempty"`
	CashFlow             int      `json:"cash_flow,omitempty"`
	CostPerOne           int      `json:"cost_per_one,omitempty"`
	Value                int      `json:"value,omitempty"`
	ApplicableToEveryOne bool     `json:"applicable_to_every_one,omitempty"`
}

func (c *CardMarket) Fill(card Card) {
	c.ID = card.ID
	c.Type = card.Type
	c.Heading = card.Heading
	c.Symbol = card.Symbol
	c.Description = card.Description
	c.Rule = card.Rule
	c.SubRule = card.SubRule
	c.Cost = card.Cost
	c.CashFlow = card.CashFlow
	c.ApplicableToEveryOne = card.ApplicableToEveryOne
	c.Success = card.Success
	c.Plus = card.Plus
}
