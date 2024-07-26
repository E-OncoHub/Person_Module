package models

import (
	"database/sql"
	"fmt"
	go_ora "github.com/sijms/go-ora/v2"
	"time"
)

type VirtualAddress struct {
	ID          int64     `json:"id_virtual_address"`
	Email       string    `json:"email"`
	PhoneNumber string    `json:"phone_number"`
	DateIn      time.Time `json:"date_in"`
	DateOut     time.Time `json:"date_out"`
}

func (v *VirtualAddress) CreateVirtualAddress(tx *sql.Tx) error {
	var id int64
	v.DateIn = time.Now()
	stmt, err := tx.Prepare("INSERT INTO virtual_address (ID_VIRTUAL_ADDRESS, email, phone_number, date_in) VALUES (VIRTUAL_ADDRESS_SEQ.nextval, :1, :2, :3) RETURNING ID_VIRTUAL_ADDRESS INTO :4")
	if err != nil {
		return err
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			panic(err)
		}
	}(stmt)
	_, err = stmt.Exec(v.Email, v.PhoneNumber, v.DateIn, go_ora.Out{Dest: &id})
	if err != nil {
		return err
	}
	v.ID = id
	return nil
}
func (v *VirtualAddress) UpdateVirtualAddress(tx *sql.Tx) error {
	// First, expire the old record
	_, err := tx.Exec(`
        UPDATE VIRTUAL_ADDRESS 
        SET DATE_OUT = SYSDATE 
        WHERE ID_VIRTUAL_ADDRESS = :1 AND DATE_OUT IS NULL
    `, v.ID)
	if err != nil {
		return fmt.Errorf("error expiring old virtual address: %w", err)
	}

	// Now, create a new record
	stmt, err := tx.Prepare(`
        INSERT INTO VIRTUAL_ADDRESS (ID_VIRTUAL_ADDRESS, EMAIL, PHONE_NUMBER, DATE_IN) 
        VALUES (VIRTUAL_ADDRESS_SEQ.nextval, :1, :2, SYSDATE) 
        RETURNING ID_VIRTUAL_ADDRESS INTO :3
    `)
	if err != nil {
		return fmt.Errorf("error preparing statement for new virtual address: %w", err)
	}
	defer stmt.Close()

	var newID int64
	_, err = stmt.Exec(v.Email, v.PhoneNumber, go_ora.Out{Dest: &newID})
	if err != nil {
		return fmt.Errorf("error creating new virtual address: %w", err)
	}

	// Update the ID in the struct to reflect the new record
	v.ID = newID
	v.DateIn = time.Now()
	v.DateOut = time.Time{} // Reset DateOut as it's a new record

	return nil
}
