package entity

type PlayerDream struct {
	ID    int    `json:"id" form:"id"`
	Name  string `json:"name" form:"name"`
	Price int    `json:"price" form:"price"`
}

type PlayerInfo struct {
	ID                uint64            `json:"id"`
	Dream             PlayerDream       `json:"dream"`
	FullName          string            `json:"fullName"`
	Language          string            `json:"language"`
	GoalPassiveIncome int               `json:"goalPassiveIncome"`
	Data              PlayerInfoData    `json:"data"`
	Conditions        BigRaceConditions `json:"conditions"`
}

type PlayerInfoData struct {
	PassiveIncome    [][]string                      `json:"passiveIncome,omitempty"`
	CommonIncome     [][]string                      `json:"commonIncome,omitempty"`
	SumExpenses      [][]string                      `json:"sumExpenses,omitempty"`
	AssetsStock      []PlayerInfoDataAssetStock      `json:"assetsStock,omitempty"`
	AssetsRealEstate []PlayerInfoDataAssetRealEstate `json:"assetsRealEstate,omitempty"`
	AssetsBusiness   []PlayerInfoDataAssetBusiness   `json:"assetsBusiness,omitempty"`
	AssetsOther      []PlayerInfoDataAssetOther      `json:"assetsOther,omitempty"`
	AssetsDreams     []PlayerInfoDataAssetDream      `json:"assetsDreams,omitempty"`
	CashFlow         [][]string                      `json:"cashFlow,omitempty"`
	Credits          [][]string                      `json:"credits,omitempty"`
	Expenses         PlayerInfoDataExpenses          `json:"expenses,omitempty"`
	Assets           PlayerInfoDataAssets            `json:"assets,omitempty"`
	Liabilities      PlayerInfoDataLiabilities       `json:"liabilities,omitempty"`
}

// AssetStock represents stock assets
type PlayerInfoDataAssetStock struct {
	Heading string `json:"heading"`
	Count   string `json:"count"`
	Price   string `json:"price"`
	Cost    string `json:"cost"`
}

// AssetStock represents stock assets
type PlayerInfoDataAssetDream struct {
	Heading string `json:"heading"`
	Cost    string `json:"cost"`
}

// AssetRealEstate represents real estate assets
type PlayerInfoDataAssetRealEstate struct {
	Heading     string `json:"heading"`
	DownPayment string `json:"downPayment"`
	Mortgage    string `json:"mortgage"`
	Cost        string `json:"cost"`
	CashFlow    string `json:"cashFlow"`
}

// AssetBusiness represents business assets
type PlayerInfoDataAssetBusiness struct {
	Heading     string `json:"heading"`
	DownPayment string `json:"downPayment"`
	Cost        string `json:"cost"`
	Mortgage    string `json:"mortgage"`
	Result      string `json:"result,omitempty"`
	CashFlow    string `json:"cashFlow"`
}

// AssetOther represents other assets
type PlayerInfoDataAssetOther struct {
	Heading  string `json:"heading"`
	Cost     string `json:"cost"`
	CashFlow string `json:"cashFlow"`
}

// Expenses represents the expenses section
type PlayerInfoDataExpenses struct {
	Taxes               string `json:"taxes"`
	HomeMortgagePayment string `json:"homeMortgagePayment"`
	SchoolLoanPayment   string `json:"schoolLoanPayment"`
	CarLoanPayment      string `json:"carLoanPayment"`
	CreditCardPayment   string `json:"creditCardPayment"`
	BankLoanPayment     string `json:"bankLoanPayment"`
	OtherExpenses       string `json:"otherExpenses"`
	PerChildExpense     string `json:"perChildExpense"`
}

// Assets represents the assets section
type PlayerInfoDataAssets struct {
	Savings  string `json:"savings"`
	Deposits string `json:"deposits"`
}

// Liabilities represents the liabilities section
type PlayerInfoDataLiabilities struct {
	HomeMortgage   string `json:"homeMortgage"`
	CarLoans       string `json:"carLoans"`
	BankLoan       string `json:"bankLoan"`
	CreditCardDebt string `json:"creditCardDebt"`
}
