package models

import (
	"database/sql"
	"errors"
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
	// Check if VirtualAddress exists
	var exists bool
	err := tx.QueryRow("SELECT 1 FROM VIRTUAL_ADDRESS WHERE ID_VIRTUAL_ADDRESS = :1", v.ID).Scan(&exists)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// If not exists, create new VirtualAddress
			return v.CreateVirtualAddress(tx)
		}
		return err
	}

	// If exists, update the VirtualAddress
	stmt, err := tx.Prepare(`
        UPDATE VIRTUAL_ADDRESS 
        SET EMAIL = :1, PHONE_NUMBER = :2, DATE_OUT = :3 
        WHERE ID_VIRTUAL_ADDRESS = :4
    `)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(v.Email, v.PhoneNumber, v.DateOut, v.ID)
	if err != nil {
		return err
	}

	// If DateOut is set, create a new active VirtualAddress
	if !v.DateOut.IsZero() {
		newVA := VirtualAddress{
			Email:       v.Email,
			PhoneNumber: v.PhoneNumber,
			DateIn:      time.Now(),
		}
		err = newVA.CreateVirtualAddress(tx)
		if err != nil {
			return err
		}
		// Update the ID to the new VirtualAddress
		v.ID = newVA.ID
	}

	return nil
}
