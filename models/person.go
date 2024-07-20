package models

import (
	"eoncohub.com/person_module/db"
	"time"
)

type Person struct {
	ID               int       `json:"id_person"`
	FirstName        string    `json:"f_name"`
	LastName         string    `json:"l_name"`
	CNP              string    `json:"cnp"`
	BornDate         time.Time `json:"born_date"`
	IDAddress        int       `json:"id_address"`
	IDVirtualAddress int       `json:"id_virtual_address"`
}

func GetPersonById(id int) (Person, error) {
	var person Person
	err := db.DB.QueryRow("SELECT id_person, f_name, l_name, cnp, born_date, id_address, id_virtual_address FROM PERSONS WHERE id_person = :1", id).Scan(&person.ID, &person.FirstName, &person.LastName, &person.CNP, &person.BornDate, &person.IDAddress, &person.IDVirtualAddress)
	if err != nil {
		return person, err
	}
	return person, nil
}
