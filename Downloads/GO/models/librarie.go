package models

type Librarie struct {
	Id             int    `json:"id"`
	Owner_name     string `json:"owner_name"`
	Owner_password string `json:"owner_password"`
	Is_Premium     bool   `json:"is_premium"`
	Creation_year  int    `json:"creation_year"`
}
