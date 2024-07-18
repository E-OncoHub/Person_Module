package models

import "time"

type Person struct {
	ID               int       `json:"id_person"`
	FirstName        string    `json:"f_name"`
	LastName         string    `json:"l_name"`
	CNP              string    `json:"cnp"`
	BornDate         time.Time `json:"born_date"`
	IDAddress        int       `json:"id_address"`
	IDVirtualAddress int       `json:"id_virtual_address"`
}
