package models

type Address struct {
	IDAddress int64  `json:"id_address"`
	Loc       Loc    `json:"loc"`
	Address   string `json:"address"`
}

func (a *Address) CreateAddress() error {
	err := a.Loc.Jud.CreateJud()
	if err != nil {
		return err
	}
	err = a.Loc.CreateLoc()
	if err != nil {
		return err
	}
	return nil
}
