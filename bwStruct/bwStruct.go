package bwStruct

type Block struct {
	Id  int        `json:"id"`
	Loc [3]float64 `json:"loc"`
}

type Data struct {
	World   []Block `json:"world"`
	Version int     `json:"version"`
	Input   string  `json:"input"`
	Error   string  `json:"error"`
}
