package models

import "database/sql"

type Address struct {
	IDAddress int64  `json:"id_address"`
	Loc       Loc    `json:"loc"`
	Address   string `json:"address"`
}

func (a *Address) CreateAddress(tx *sql.Tx) error {
	err := a.Loc.Jud.CreateJud(tx)
	if err != nil {
		return err
	}
	err = a.Loc.CreateLoc(tx)
	if err != nil {
		return err
	}

	_, err = tx.Exec("INSERT INTO ADDRESS (id_address, id_loc, address) VALUES (ADDRESS_SEQ.nextval, :1, :2)", a.Loc.ID, a.Address)
	if err != nil {
		return err
	}

	return nil
}
