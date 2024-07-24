package models

import (
	"database/sql"
	"eoncohub.com/person_module/db"
	"errors"
)

type Loc struct {
	ID   int64  `json:"id_loc"`
	Name string `json:"name"`
	Jud  Jud    `json:"jud"`
}

func (l *Loc) CreateLoc(tx *sql.Tx) error {
	err := tx.QueryRow("SELECT ID_LOC FROM LOC WHERE NAME = :1", l.Name).Scan(&l.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			_, err = db.DB.Exec("INSERT INTO LOC (ID_LOC, NAME, ID_JUD) VALUES (LOC_SEQ.nextval, :1, :2)", l.Name, l.Jud.ID)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}
	return nil
}
