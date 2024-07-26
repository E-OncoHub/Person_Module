package models

import (
	"database/sql"
	"errors"
	go_ora "github.com/sijms/go-ora/v2"
)

type Address struct {
	IDAddress int64  `json:"id_address"`
	Loc       Loc    `json:"loc"`
	Address   string `json:"address"`
}

func (a *Address) CreateAddress(tx *sql.Tx) error {
	err := (a.Loc.Jud).CreateJud(tx)
	if err != nil {
		return err
	}
	err = a.Loc.CreateLoc(tx)
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare("INSERT INTO ADDRESS (ID_ADDRESS, ID_LOC, ADDRESS) VALUES (ADDRESS_SEQ.nextval, :1, :2) RETURNING ID_ADDRESS INTO :3")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(a.Loc.ID, a.Address, go_ora.Out{Dest: &a.IDAddress})
	if err != nil {
		return err
	}

	return nil
}
func (a *Address) UpdateAddress(tx *sql.Tx) error {
	err := a.Loc.UpdateLoc(tx)
	if err != nil {
		return err
	}

	// Check if Address exists
	var exists bool
	err = tx.QueryRow("SELECT 1 FROM ADDRESS WHERE ID_ADDRESS = :1", a.IDAddress).Scan(&exists)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// If not exists, create new Address
			return a.CreateAddress(tx)
		}
		return err
	}

	// If exists, update the Address
	stmt, err := tx.Prepare("UPDATE ADDRESS SET ID_LOC = :1, ADDRESS = :2 WHERE ID_ADDRESS = :3")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(a.Loc.ID, a.Address, a.IDAddress)
	if err != nil {
		return err
	}

	return nil
}
