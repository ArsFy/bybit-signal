package main

type Signal struct {
	Symbol      string `json:"symbol"`
	Side        string `json:"side"`
	Qty         string `json:"qty"`
	Price       string `json:"price"`
	TriggerTime string `json:"trigger_time"`
	StrategyID  string `json:"strategy_id"`
}
