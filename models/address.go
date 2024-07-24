package models

import (
	"database/sql"
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
