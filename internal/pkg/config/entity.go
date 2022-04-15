package config

type Http struct {
	BindAddress string
	Mode string
	Prefix string
}

type Blockchain struct {
	Chains []Chain
}
type Chain struct {
	ChainID uint
	API      API
	Subgraph Subgraph
}
type Subgraph struct {
	AggregatorRouter string
	KyberswapRouter string
}

type API struct {
	Price string
}

type Token struct {
	Address  string
	Name     string
	Symbol   string
	Decimals int
	CkgID    string `json:"cgkId"`
}

var ETH = map[uint]Token{
	1: Token{
		Address: "",
		Name: "ETH",
		Symbol: "ETH",
		Decimals: 18,
	},
	56: Token{
		Address: "",
		Name: "BNB",
		Symbol: "BNB",
		Decimals: 18,
	},
	43114: Token{
		Address: "",
		Name: "AVAX",
		Symbol: "AVAX",
		Decimals: 18,
	},
	25: Token{
		Address: "",
		Name: "CRO",
		Symbol: "CRO",
		Decimals: 18,
	},
	250: Token{
		Address: "",
		Name: "FTM",
		Symbol: "FTM",
		Decimals: 18,
	},
	137: Token{
		Address: "",
		Name: "MATIC",
		Symbol: "MATIC",
		Decimals: 18,
	},
}

