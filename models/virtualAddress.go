package models

import (
	"database/sql"
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
