package models

type Address struct {
	ID      int    `json:"id_address"`
	IDLoc   int    `json:"id_loc"`
	Address string `json:"address"`
}
