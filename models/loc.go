package models

import (
	"database/sql"
	"errors"
	go_ora "github.com/sijms/go-ora/v2"
)

type Loc struct {
	ID   int64  `json:"id_loc"`
	Name string `json:"name"`
	Jud  Jud    `json:"jud"`
}

func (l *Loc) CreateLoc(tx *sql.Tx) error {
	err := tx.QueryRow("SELECT ID_LOC FROM LOC WHERE upper(NAME) = :1", l.Name).Scan(&l.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			stmt, err := tx.Prepare("INSERT INTO LOC (ID_LOC, NAME, ID_JUD) VALUES (LOC_SEQ.nextval, :1, :2) RETURNING ID_LOC INTO :3")
			if err != nil {
				return err
			}
			defer stmt.Close()

			_, err = stmt.Exec(l.Name, l.Jud.ID, go_ora.Out{Dest: &l.ID})
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}
	return nil
}
