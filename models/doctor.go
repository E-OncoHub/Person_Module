package models

type Doctor struct {
	ID       int    `json:"id_doctor"`
	IDPerson int    `json:"id_person"`
	Parafa   string `json:"parafa"`
}
