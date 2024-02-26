package models

type GeoZone struct {
	Id       int
	Name     string
	Points   [][]float64
	Position []float64
	Radius   float64
	Type     int
}
type Pos struct {
	X, Y      float64
	PolygonId int
}
type Circle struct {
	Name   string
	X, Y   float64
	Radius float64
	Cars   map[int]bool
}
type Quad struct {
	Points   []Pos
	Quads    []Quad
	Position Pos
	H        float64
	W        float64
}

type Polygon struct {
	Name     string
	Vertices []Pos
	Cars     map[int]bool
}

type User struct {
	Polygons []Polygon
	MainQuad Quad
	Circles  []Circle
}

func (u *User) GetMainQuad() *Quad {
	return &u.MainQuad
}
