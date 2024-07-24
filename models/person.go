package models

import (
	"errors"
	"time"
)

type Person struct {
	IDPerson       int64          `json:"id_person"`
	FName          string         `json:"f_name"`
	LName          string         `json:"l_name"`
	CNP            string         `json:"cnp"`
	BornDate       time.Time      `json:"born_date"`
	Address        Address        `json:"address"`
	VirtualAddress VirtualAddress `json:"virtual_address"`
}

func (p *Person) Create() error {

	err := p.VirtualAddress.CreateVirtualAddress()
	if err != nil {
		return errors.New("error creating virtual address")
	}
	err = p.Address.CreateAddress()
	if err != nil {
		return errors.New("error creating address")
	}
	return nil
}
