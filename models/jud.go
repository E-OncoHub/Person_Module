package models

import (
	"database/sql"
	"errors"
	go_ora "github.com/sijms/go-ora/v2"
)

type Jud struct {
	ID   int64  `json:"id_jud"`
	Name string `json:"name"`
}

func (j *Jud) CreateJud(tx *sql.Tx) error {
	err := tx.QueryRow("SELECT ID_JUD FROM JUD WHERE NAME = :1", j.Name).Scan(&j.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			stmt, err := tx.Prepare("INSERT INTO JUD (ID_JUD, NAME) VALUES (JUD_SEQ.nextval, :1) RETURNING ID_JUD INTO :2")
			if err != nil {
				return err
			}
			defer stmt.Close()
			_, err = stmt.Exec(j.Name, go_ora.Out{Dest: &j.ID})
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}
	return nil
}
