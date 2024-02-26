package models

type MsgAll struct {
	Id  int      `json:"id"`
	T   int      `json:"t"`
	Ign int      `json:"ignition"`
	Pos Position `json:"pos"`
}

type Position struct {
	S int     `json:"s"`
	X float64 `json:"x"`
	Y float64 `json:"y"`
}
