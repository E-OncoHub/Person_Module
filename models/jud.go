package models

import (
	"database/sql"
	"eoncohub.com/person_module/db"
	"errors"
)

type Jud struct {
	ID   int64  `json:"id_jud"`
	Name string `json:"name"`
}

func (j *Jud) CreateJud() error {
	err := db.DB.QueryRow("SELECT ID_JUD FROM JUD WHERE NAME = :1", j.Name).Scan(&j.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			_, err = db.DB.Exec("INSERT INTO JUD (ID_JUD, NAME) VALUES (JUD_SEQ.nextval, :1)", j.Name)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}
	return nil
}
