package models

import (
	"eoncohub.com/person_module/db"
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
	tx, err := db.DB.Begin()
	if err != nil {
		return err
	}

	// Create virtual address within the same transaction
	err = (&p.VirtualAddress).CreateVirtualAddress(tx)
	if err != nil {
		tx.Rollback()
		return errors.New("error creating virtual address: " + err.Error())
	}

	// Create address within the same transaction
	err = p.Address.CreateAddress(tx)
	if err != nil {
		tx.Rollback()
		return errors.New("error creating address: " + err.Error())
	}

	// Insert person within the same transaction
	_, err = tx.Exec("INSERT INTO PERSONS (id_person, f_name, l_name, cnp, born_date, id_address, id_virtual_address) VALUES (PERSONS_SEQ.nextval, :1, :2, :3, :4, :5, :6)", p.FName, p.LName, p.CNP, p.BornDate, p.Address.IDAddress, p.VirtualAddress.ID)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}
