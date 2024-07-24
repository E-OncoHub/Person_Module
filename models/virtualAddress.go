package models

import (
	"database/sql"
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
	stmt, err := tx.Prepare("INSERT INTO virtual_address (ID_VIRTUAL_ADDRESS, email, phone_number, date_in) VALUES (VIRTUAL_ADDRESS_SEQ.nextval, :1, :2, :3)")
	if err != nil {
		return err
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
		}
	}(stmt)

	_, err = stmt.Exec(v.Email, v.PhoneNumber, v.DateIn)
	if err != nil {
		return err
	}
	return nil
}
