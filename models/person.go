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

func GetPerson(id int64) (Person, error) {
	var p Person
	err := db.DB.QueryRow(
		`SELECT ID_PERSON,
	       F_NAME,
	       L_NAME,
	       CNP,
	       BORN_DATE,
	       EMAIL,
	       PHONE_NUMBER,
	       ad.ADDRESS,
	       l.NAME,
	       j.NAME
	from PERSONS
	         join VIRTUAL_ADDRESS on PERSONS.ID_VIRTUAL_ADDRESS = VIRTUAL_ADDRESS.ID_VIRTUAL_ADDRESS
	         join ADDRESS Ad on PERSONS.ID_ADDRESS = Ad.ID_ADDRESS
	         join Loc l on Ad.ID_LOC = l.ID_LOC
	         join Jud j on l.ID_JUD = j.ID_JUD
	where DATE_OUT is null
	 and ID_PERSON = :1
	`, id).Scan(&p.IDPerson, &p.FName, &p.LName, &p.CNP, &p.BornDate, &p.VirtualAddress.Email, &p.VirtualAddress.PhoneNumber, &p.Address.Address, &p.Address.Loc.Name, &p.Address.Loc.Jud.Name)
	if err != nil {
		return Person{}, err
	}
	return p, nil
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
	err = (&p.Address).CreateAddress(tx)
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

func (p *Person) Update() error {
	tx, err := db.DB.Begin()
	if err != nil {
		return err
	}

	// Update virtual address
	err = p.VirtualAddress.UpdateVirtualAddress(tx)
	if err != nil {
		tx.Rollback()
		return errors.New("error updating virtual address: " + err.Error())
	}

	// Update address
	err = p.Address.UpdateAddress(tx)
	if err != nil {
		tx.Rollback()
		return errors.New("error updating address: " + err.Error())
	}

	_, err = tx.Exec(`
        UPDATE PERSONS 
        SET f_name = :1, l_name = :2, cnp = :3, born_date = :4 
        WHERE id_person = :5
    `, p.FName, p.LName, p.CNP, p.BornDate, p.IDPerson)
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
