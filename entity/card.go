package entity

type Card struct {
	ID            string        `json:"id"`
	Type          string        `json:"type"`
	Symbol        string        `json:"symbol"`
	Name          string        `json:"name"`
	Family        string        `json:"family"`
	Heading       string        `json:"heading"`
	Description   string        `json:"description"`
	Cost          int           `json:"cost,omitempty"`
	Rule          string        `json:"rule,omitempty"`
	IsConditional bool          `json:"is_conditional,omitempty"`
	HasBabies     bool          `json:"has_babies,omitempty"`
	Plus          bool          `json:"plus,omitempty"`
	Value         int           `json:"value,omitempty"`
	Mortgage      int           `json:"mortgage,omitempty"`
	DownPayment   int           `json:"down_payment,omitempty"`
	CashFlow      int           `json:"cash_flow,omitempty"`
	Price         int           `json:"price,omitempty"`
	AssetType     string        `json:"asset_type,omitempty"`
	Count         int           `json:"count,omitempty"`
	History       []CardHistory `json:"history,omitempty"`
	Increase      int           `json:"increase,omitempty"`
	Decrease      int           `json:"decrease,omitempty"`
	OnlyYou       bool          `json:"only_you,omitempty"`
	Range         []int         `json:"range,omitempty"`
	SubRule       []string      `json:"sub_rule,omitempty"`
	Lottery       string        `json:"lottery,omitempty"`
	Failure       []int         `json:"failure,omitempty"`
	Percent       int           `json:"percent,omitempty"`
	PlayerId      int           `json:"player_id,omitempty"`
	Success       []int         `json:"success,omitempty"`
	CostPerOne    int           `json:"cost_per_one,omitempty"`
	ExtraDices    int           `json:"extra_dices,omitempty"`
	Limit         int           `json:"limit,omitempty"`
	Outcome       interface{}   `json:"outcome,omitempty"`
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
	OnlyYou     bool     `json:"only_you,omitempty"`
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
	AssetType   string        `json:"asset_type"`
	Symbol      string        `json:"symbol"`
	Heading     string        `json:"heading"`
	Description string        `json:"description"`
	Rule        string        `json:"rule,omitempty"`
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
	c.AssetType = card.AssetType
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
	Cost        int    `json:"cost,omitempty"`
	CostPerOne  int    `json:"cost_per_one,omitempty"`
	Count       int    `json:"count,omitempty"`
	AssetType   string `json:"asset_type"`
	IsOwner     bool   `json:"is_owner,omitempty"`
	Symbol      string `json:"symbol"`
	Heading     string `json:"heading"`
	Description string `json:"description"`
	WholeCost   int    `json:"-"`
}

type CardLotteryOutcome struct {
	Failure int `json:"failure,omitempty"`
	Success int `json:"success,omitempty"`
}

type CardLottery struct {
	ID          string             `json:"id"`
	Type        string             `json:"type"`
	Symbol      string             `json:"symbol"`
	Heading     string             `json:"heading"`
	Description string             `json:"description"`
	Cost        int                `json:"cost"`
	AssetType   string             `json:"asset_type,omitempty"`
	Rule        string             `json:"rule,omitempty"`
	SubRule     []string           `json:"sub_rule,omitempty"`
	Failure     []int              `json:"failure,omitempty"`
	Success     []int              `json:"success,omitempty"`
	Outcome     CardLotteryOutcome `json:"outcome,omitempty"`
}

type CardDream struct {
	ID          string `json:"id"`
	Heading     string `json:"heading"`
	Description string `json:"description"`
	Type        string `json:"type"`
	AssetType   string `json:"asset_type,omitempty"`
	Cost        int    `json:"cost"`
	PlayerId    int    `json:"player_id,omitempty"`
}

type CardCharity struct {
	ID          string `json:"id"`
	Heading     string `json:"heading"`
	Description string `json:"description"`
	Cost        int    `json:"cost"`
	ExtraDices  int    `json:"extra_dices"`
	Limit       int    `json:"limit"`
	Percent     int    `json:"percent"`
	Type        string `json:"type"`
	Symbol      string `json:"symbol"`
}

func (c *CardCharity) Fill(card Card) {
	c.ID = card.ID
	c.Type = card.Type
	c.Heading = card.Heading
	c.Symbol = card.Symbol
	c.Description = card.Description
	c.Cost = card.Cost
	c.Limit = card.Limit
	c.Percent = card.Percent
	c.ExtraDices = card.ExtraDices
}

type CardPayTax struct {
	ID          string `json:"id"`
	Heading     string `json:"heading"`
	Description string `json:"description"`
	Type        string `json:"type"`
	Percent     int    `json:"percent,omitempty"`
}

func (c *CardPayTax) Fill(card Card) {
	c.ID = card.ID
	c.Type = card.Type
	c.Heading = card.Heading
	c.Description = card.Description
	c.Percent = card.Percent
}

type CardDownsized struct {
	ID          string `json:"id"`
	Heading     string `json:"heading"`
	Description string `json:"description"`
	Type        string `json:"type"`
	Percent     int    `json:"percent,omitempty"`
}

type CardMarket struct {
	ID          string   `json:"id"`
	Type        string   `json:"type"`
	Heading     string   `json:"heading"`
	Symbol      string   `json:"symbol"`
	Description string   `json:"description"`
	Rule        string   `json:"rule"`
	SubRule     []string `json:"sub_rule"`
	Success     []int    `json:"success,omitempty"`
	Cost        int      `json:"cost,omitempty"`
	CashFlow    int      `json:"cash_flow,omitempty"`
	CostPerOne  int      `json:"cost_per_one,omitempty"`
	AssetType   string   `json:"asset_type,omitempty"`
	OnlyYou     bool     `json:"only_you,omitempty"`
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
	c.CostPerOne = card.CostPerOne
	c.OnlyYou = card.OnlyYou
	c.AssetType = card.AssetType
	c.Success = card.Success
}
