package models

type NotifParam struct {
	Name        string
	Activated   bool
	Param_id    int
	Speed_limit struct {
		MaxSpeed int `json:"max_speed"`
	} `json:"param"`
	Creator_id int
}
