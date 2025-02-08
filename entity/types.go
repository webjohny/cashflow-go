package entity

var LotteryTypes = struct {
	Cash     string
	CashFlow string
}{
	Cash:     "cash",
	CashFlow: "cashflow",
}

var BusinessTypes = struct {
	Startup  string
	Business string
	Limited  string
}{
	Startup:  "startup",
	Business: "business",
	Limited:  "limited",
}

var OtherAssetTypes = struct {
	Piece            string
	Whole            string
	HealthyInsurance string
}{
	Piece:            "piece",
	Whole:            "whole",
	HealthyInsurance: "healthyInsurance",
}

var RealEstateTypes = struct {
	Building string
	Single   string
}{
	Building: "building",
	Single:   "single",
}

var StockTypes = struct {
	Manipulation string
}{
	Manipulation: "manipulation",
}

var SmallDealTypes = struct {
	Lottery string
}{
	Lottery: "lottery",
}

var BigBusinessTypes = struct {
	RiskBusiness string
	RiskStocks   string
}{
	RiskBusiness: "riskBusiness",
	RiskStocks:   "riskStock",
}

var BankruptTypes = struct {
	Absolute string
	Relative string
}{
	Absolute: "absolute",
	Relative: "relative",
}

var MarketTypes = struct {
	AnyRealEstate  string
	EachRealEstate string
	AnyBusiness    string
	EachBusiness   string
	AnyStartup     string
	EachStartup    string
}{
	AnyRealEstate:  "anyRealEstate",
	EachRealEstate: "eachRealEstate",
	AnyBusiness:    "anyBusiness",
	EachBusiness:   "eachBusiness",
	AnyStartup:     "anyStartup",
	EachStartup:    "eachStartup",
}

var UserRequestTypes = struct {
	Baby   string
	Salary string
}{
	Baby:   "baby",
	Salary: "salary",
}
