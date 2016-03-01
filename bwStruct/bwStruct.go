package bwStruct

type BWBlock struct {
	Id  int        `json:"id"`
	Loc [3]float64 `json:"loc"`
}

type BWData struct {
	World   []BWBlock `json:"world"`
	Version int       `json:"version"`
	Input   string    `json:"input"`
	Error   string    `json:"error"`
}
