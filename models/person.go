package models

import (
	"database/sql"
	"errors"
	"fmt"
	go_ora "github.com/sijms/go-ora/v2"
	"time"

	"eoncohub.com/person_module/db"
	"eoncohub.com/person_module/utils"
)

var ErrPersonNotFound = errors.New("person not found or already expired")

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
	defer tx.Rollback() // This will be a no-op if the tx has been committed later

	// Create virtual address within the same transaction
	err = (&p.VirtualAddress).CreateVirtualAddress(tx)
	if err != nil {
		return errors.New("error creating virtual address: " + err.Error())
	}

	// Create address within the same transaction
	err = (&p.Address).CreateAddress(tx)
	if err != nil {
		return errors.New("error creating address: " + err.Error())
	}

	// Prepare the statement
	stmt, err := tx.Prepare(`
        INSERT INTO PERSONS (id_person, f_name, l_name, cnp, born_date, id_address, id_virtual_address) 
        VALUES (PERSONS_SEQ.nextval, :1, :2, :3, :4, :5, :6)
        RETURNING id_person INTO :7`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Register the output parameter
	var newID int64
	_, err = stmt.Exec(
		p.FName, p.LName, p.CNP, p.BornDate, p.Address.IDAddress, p.VirtualAddress.ID,
		go_ora.Out{Dest: &newID, In: false},
	)
	if err != nil {
		return errors.New("error inserting person: " + err.Error())
	}

	// Set the new ID to the Person struct
	p.IDPerson = newID

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return errors.New("error committing transaction: " + err.Error())
	}

	return nil
}

func (p *Person) Update() error {
	tx, err := db.DB.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	if !utils.IsEmptyStruct(p.VirtualAddress) {
		var currentVirtualAddressID int64
		err = tx.QueryRow("SELECT ID_VIRTUAL_ADDRESS FROM PERSONS WHERE ID_PERSON = :1", p.IDPerson).Scan(&currentVirtualAddressID)
		if err != nil {
			return fmt.Errorf("error getting current virtual address ID: %w", err)
		}
		p.VirtualAddress.ID = currentVirtualAddressID
		err = p.VirtualAddress.UpdateVirtualAddress(tx)
		if err != nil {
			return fmt.Errorf("error updating virtual address: %w", err)
		}

		_, err = tx.Exec("UPDATE PERSONS SET ID_VIRTUAL_ADDRESS = :1 WHERE ID_PERSON = :2", p.VirtualAddress.ID, p.IDPerson)
		if err != nil {
			return fmt.Errorf("error updating person's virtual address reference: %w", err)
		}
	}

	if !utils.IsEmptyStruct(p.Address) {
		var currentAddressID int64
		err = tx.QueryRow("SELECT ID_ADDRESS FROM PERSONS WHERE ID_PERSON = :1", p.IDPerson).Scan(&currentAddressID)
		if err != nil {
			return fmt.Errorf("error getting current address ID: %w", err)
		}
		p.Address.IDAddress = currentAddressID
		err = p.Address.UpdateAddress(tx)
		if err != nil {
			return fmt.Errorf("error updating address: %w", err)
		}

		// Update the PERSONS table with the new address ID
		_, err = tx.Exec("UPDATE PERSONS SET ID_ADDRESS = :1 WHERE ID_PERSON = :2", p.Address.IDAddress, p.IDPerson)
		if err != nil {
			return fmt.Errorf("error updating person's address reference: %w", err)
		}
	}

	// Update person fields
	updateQuery := "UPDATE PERSONS SET "
	updateParams := []interface{}{}
	paramCount := 1

	if p.FName != "" {
		updateQuery += fmt.Sprintf("f_name = :%d, ", paramCount)
		updateParams = append(updateParams, p.FName)
		paramCount++
	}
	if p.LName != "" {
		updateQuery += fmt.Sprintf("l_name = :%d, ", paramCount)
		updateParams = append(updateParams, p.LName)
		paramCount++
	}
	if p.CNP != "" {
		updateQuery += fmt.Sprintf("cnp = :%d, ", paramCount)
		updateParams = append(updateParams, p.CNP)
		paramCount++
	}
	if !p.BornDate.IsZero() {
		updateQuery += fmt.Sprintf("born_date = :%d, ", paramCount)
		updateParams = append(updateParams, p.BornDate)
		paramCount++
	}

	// Remove trailing comma and space
	updateQuery = updateQuery[:len(updateQuery)-2]
	updateQuery += fmt.Sprintf(" WHERE id_person = :%d", paramCount)
	updateParams = append(updateParams, p.IDPerson)

	// Execute the update only if there are fields to update
	if len(updateParams) > 1 { // More than just the ID
		_, err = tx.Exec(updateQuery, updateParams...)
		if err != nil {
			return fmt.Errorf("error updating person: %w", err)
		}
	}
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("error committing transaction: %w", err)
	}
	return nil
}

func DeletePerson(personID int64) error {
	tx, err := db.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	// Get and delete virtual address
	var virtualAddressID int64
	err = tx.QueryRow("SELECT ID_VIRTUAL_ADDRESS FROM PERSONS WHERE ID_PERSON = :1", personID).Scan(&virtualAddressID)
	if err != nil {
		if err == sql.ErrNoRows {
			return ErrPersonNotFound
		}
		return err
	}

	// Get and delete address
	var addressID int64
	err = tx.QueryRow("SELECT ID_ADDRESS FROM PERSONS WHERE ID_PERSON = :1", personID).Scan(&addressID)
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	// Delete the person record
	result, err := tx.Exec("DELETE FROM PERSONS WHERE ID_PERSON = :1", personID)
	if err != nil {
		return err
	}

	if addressID != 0 {
		_, err = tx.Exec("DELETE FROM ADDRESS WHERE ID_ADDRESS = :1", addressID)
		if err != nil {
			return err
		}
	}

	_, err = tx.Exec("DELETE FROM VIRTUAL_ADDRESS WHERE ID_VIRTUAL_ADDRESS = :1", virtualAddressID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrPersonNotFound
	}

	return tx.Commit()
}
