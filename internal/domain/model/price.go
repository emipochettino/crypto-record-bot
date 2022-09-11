package model

type PriceWithChange struct {
	Usd          float32 `json:"usd"`
	Usd24HChange float32 `json:"usd_24h_change"`
}

func (p PriceWithChange) GetChangeSymbol() string {
	symbol := ""
	if p.Usd24HChange >= 15 {
		symbol = "🚀"
	} else if p.Usd24HChange > 0 {
		symbol = "😎"
	} else if p.Usd24HChange < 0 {
		symbol = "😓"
	}

	return symbol
}

type SimplePrice struct {
	Id string
	PriceWithChange
}
