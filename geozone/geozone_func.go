package geozone

import (
	"Organize/models"
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"math"
)

var Users = make(map[int]*models.User)

func StartSave(db *pgxpool.Pool) {
	rows_creators, err := db.Query(context.Background(), "SELECT id from auth_user")
	if err != nil {
		log.Fatal(err)
	}
	var creator_id int
	for rows_creators.Next() {
		rows_creators.Scan(&creator_id)
		UpdateGeoPolygon(creator_id, db)
		UpdateGeoCircle(creator_id, db)
	}
}

func InitGeo(geo_id int, db *pgxpool.Pool) {
	rows_creators, err := db.Query(context.Background(), "Select user_id from main_geozones_creator where geozones_id = $1", geo_id)
	if err != nil {
		log.Fatal(err)
	}
	var creator_id int
	for rows_creators.Next() {
		rows_creators.Scan(&creator_id)
		UpdateGeoPolygon(creator_id, db)
		UpdateGeoCircle(creator_id, db)
	}
}

func build(MainQuad *models.Quad) {
	for i := 1; i <= 4; i++ {
		var SubQuad models.Quad
		SubQuad.H = MainQuad.H / 2
		SubQuad.W = MainQuad.W / 2
		switch i {
		case 1:
			SubQuad.Position = models.Pos{(MainQuad.Position.X*2 + MainQuad.W) / 2, (MainQuad.Position.Y*2 + MainQuad.H) / 2, 0}
			//log.Println(SubQuad.Position)
			//log.Println(SubQuad.H, SubQuad.W)
			MainQuad.Quads = append(MainQuad.Quads, SubQuad)
		case 2:
			SubQuad.Position = models.Pos{(MainQuad.Position.X*2 - MainQuad.W) / 2, (MainQuad.Position.Y*2 + MainQuad.H) / 2, 0}
			//log.Println(SubQuad.Position)
			MainQuad.Quads = append(MainQuad.Quads, SubQuad)
		case 3:
			SubQuad.Position = models.Pos{(MainQuad.Position.X*2 - MainQuad.W) / 2, (MainQuad.Position.Y*2 - MainQuad.H) / 2, 0}
			//log.Println(SubQuad.Position)
			MainQuad.Quads = append(MainQuad.Quads, SubQuad)
		case 4:
			SubQuad.Position = models.Pos{(MainQuad.Position.X*2 + MainQuad.W) / 2, (MainQuad.Position.Y*2 - MainQuad.H) / 2, 0}
			//log.Println(SubQuad.Position)
			MainQuad.Quads = append(MainQuad.Quads, SubQuad)
		}
	}
	for _, point := range MainQuad.Points {
		if point.X > MainQuad.Position.X && point.Y > MainQuad.Position.Y {
			MainQuad.Quads[0].Points = append(MainQuad.Quads[0].Points, point)
			//log.Println(point, " into 0")
		}
		if point.X < MainQuad.Position.X && point.Y > MainQuad.Position.Y {
			MainQuad.Quads[1].Points = append(MainQuad.Quads[1].Points, point)
			//log.Println(point, " into 1")
		}
		if point.X < MainQuad.Position.X && point.Y < MainQuad.Position.Y {
			MainQuad.Quads[2].Points = append(MainQuad.Quads[2].Points, point)
			//log.Println(point, " into 2")
		}
		if point.X > MainQuad.Position.X && point.Y < MainQuad.Position.Y {
			MainQuad.Quads[3].Points = append(MainQuad.Quads[3].Points, point)
			//log.Println(point, " into 3")
		}
	}

	for i, _ := range MainQuad.Quads {
		//log.Println(MainQuad.Quads[i], i)
		if len(MainQuad.Quads[i].Points) > 10 {
			build(&MainQuad.Quads[i])
		}

	}
}
func findNear(car models.Pos, MainQuad *models.Quad) []models.Pos {
	var found []models.Pos
	if len(MainQuad.Points) == 10 {
		//fmt.Println("MainQuad poins len is 1")
		return MainQuad.Points
	}
	if car.X > MainQuad.Position.X && car.Y > MainQuad.Position.Y {
		if len(MainQuad.Quads[0].Points) <= 9 {
			//fmt.Println("MainQuad poins len is 0")
			return MainQuad.Points
		}
		found = findNear(car, &MainQuad.Quads[0])
	}
	if car.X < MainQuad.Position.X && car.Y > MainQuad.Position.Y {
		if len(MainQuad.Quads[1].Points) <= 9 {
			//fmt.Println("MainQuad poins len is 0")
			return MainQuad.Points
		}
		found = findNear(car, &MainQuad.Quads[1])
	}
	if car.X < MainQuad.Position.X && car.Y < MainQuad.Position.Y {
		if len(MainQuad.Quads[2].Points) <= 9 {
			//fmt.Println("MainQuad poins len is 0")
			return MainQuad.Points
		}
		found = findNear(car, &MainQuad.Quads[2])
	}
	if car.X > MainQuad.Position.X && car.Y < MainQuad.Position.Y {
		if len(MainQuad.Quads[3].Points) <= 9 {
			//fmt.Println("MainQuad poins len is 0")
			return MainQuad.Points
		}
		found = findNear(car, &MainQuad.Quads[3])
	}
	return found
}
func Contains(point models.Pos, Pol_Id int, Id int) (bool, bool) {
	if len(Users[Id].Polygons) < 1 {
		return false, true
	}
	p := Users[Id].Polygons[Pol_Id]

	numVertices := len(p.Vertices)
	if numVertices < 3 {
		return false, true // Полигон должен иметь хотя бы три вершины
	}

	var inside bool
	for i, j := 0, numVertices-1; i < numVertices; i++ {
		if (p.Vertices[i].Y > point.Y) != (p.Vertices[j].Y > point.Y) &&
			point.X < (p.Vertices[j].X-p.Vertices[i].X)*(point.Y-p.Vertices[i].Y)/(p.Vertices[j].Y-p.Vertices[i].Y)+p.Vertices[i].X {
			inside = !inside
		}
		j = i
	}

	return inside, false
}
func UpdateGeoPolygon(Id int, db *pgxpool.Pool) {
	var User models.User
	rows_geozones, err := db.Query(context.Background(), "SELECT geozones_id from main_geozones_creator where user_id = $1", Id)
	if err != nil {
		log.Fatal(err)
	}
	if Users[Id] == nil {
		Users[Id] = &models.User{} // Инициализация, если nil
	}
	Users[Id].Polygons = User.Polygons
	var geo_id int
	var d int
	d = 0
	for rows_geozones.Next() {
		rows_geozones.Scan(&geo_id)
		var geozone models.GeoZone
		db.QueryRow(context.Background(), "SELECT id,name,coord,type_id,radius from main_geozones where id = $1", geo_id).Scan(&geozone.Id, &geozone.Name, &geozone.Points, &geozone.Type, &geozone.Radius)
		if geozone.Type != 2 {
			continue
		}
		var polygon models.Polygon
		polygon.Name = geozone.Name
		for i, _ := range geozone.Points {
			var Vertice models.Pos
			Vertice.X = geozone.Points[i][0]
			Vertice.Y = geozone.Points[i][1]
			Vertice.PolygonId = d
			polygon.Vertices = append(polygon.Vertices, Vertice)
			User.MainQuad.Points = append(User.MainQuad.Points, Vertice)
		}
		if d < len(Users[Id].Polygons) {
			polygon.Cars = Users[Id].Polygons[d].Cars
		} else {
			polygon.Cars = make(map[int]bool)
		}
		User.Polygons = append(User.Polygons, polygon)
		d++
	}
	Users[Id].Polygons = User.Polygons
	User.MainQuad.H = 90
	User.MainQuad.W = 180
	Users[Id].MainQuad = User.MainQuad
	build(&Users[Id].MainQuad)
	log.Println(Id, 1)
}
func FindCircleDis(point models.Pos, Id int) ([]models.Circle, []float64) {
	var minDis []float64
	var circle []models.Circle

	for i, _ := range Users[Id].Circles {
		dis := math.Sqrt(math.Pow(point.X-Users[Id].Circles[i].X, 2)+math.Pow(point.Y-Users[Id].Circles[i].Y, 2)) - Users[Id].Circles[i].Radius
		minDis = append(minDis, dis)
		circle = append(circle, Users[Id].Circles[i])
	}
	return circle, minDis
}
func UpdateGeoCircle(Id int, db *pgxpool.Pool) {
	var User models.User
	rows_geozones, err := db.Query(context.Background(), "SELECT geozones_id from main_geozones_creator where user_id = $1", Id)
	if err != nil {
		log.Fatal(err)
	}
	var geo_id int
	var d int
	d = 0
	for rows_geozones.Next() {
		rows_geozones.Scan(&geo_id)
		var geozone models.GeoZone
		db.QueryRow(context.Background(), "SELECT id,name,coord,type_id,radius from main_geozones where id = $1", geo_id).Scan(&geozone.Id, &geozone.Name, &geozone.Points, &geozone.Type, &geozone.Radius)
		if geozone.Type == 2 {
			continue
		}
		if geozone.Points == nil {
			continue
		}
		var circle models.Circle
		circle.X = geozone.Points[0][0]
		circle.Y = geozone.Points[0][1]
		circle.Name = geozone.Name
		circle.Radius = geozone.Radius
		User.Circles = append(User.Circles, circle)
		if d < len(Users[Id].Circles) {
			circle.Cars = Users[Id].Circles[d].Cars
		} else {
			circle.Cars = make(map[int]bool)
		}
		d++
	}
	if Users[Id] == nil {
		Users[Id] = &models.User{} // Инициализация, если nil
	}
	Users[Id].Circles = User.Circles
	log.Println(Id, 2)
}

func FindPolygon(Id int, car models.Pos, car_id int) (bool, []bool, []string) {
	var o bool
	var inout []bool
	var Name []string
	o = false
	points := findNear(car, &Users[Id].MainQuad)
	//var Mindis float64
	//var point models.Pos
	//Mindis = math.Inf(1)
	for _, point := range points {
		inside, werr := Contains(car, point.PolygonId, Id)
		if werr {
			return false, nil, nil
		}
		if inside != Users[Id].Polygons[point.PolygonId].Cars[car_id] {
			if Users[Id].Polygons[point.PolygonId].Cars[car_id] {
				Users[Id].Polygons[point.PolygonId].Cars[car_id] = false
				Name = append(Name, Users[Id].Polygons[point.PolygonId].Name)
				o = true
				inout = append(inout, false)
			} else {
				Users[Id].Polygons[point.PolygonId].Cars[car_id] = true
				Name = append(Name, Users[Id].Polygons[point.PolygonId].Name)
				o = true
				inout = append(inout, true)
			}
		}
	}
	return o, inout, Name
}
func FindCircle(Id int, car models.Pos, car_id int) (bool, []bool, []string) {
	var o bool
	var inout []boo
	var Name []string
	o = false
	points, dis := FindCircleDis(car, Id)
	for i, point := range points {
		if dis[i] <= 0 && !point.Cars[car_id] {
			point.Cars[car_id] = true
			Name = append(Name, point.Name)
			o = true
			inout = append(inout, true)
		}

		if dis[i] > 0 && point.Cars[car_id] {
			point.Cars[car_id] = false
			Name = append(Name, point.Name)
			o = true
			inout = append(inout, false)
		}
	}

	return o, inout, Name
}
