package models

import "time"

type Jud struct {
	ID   int    `json:"id_jud"`
	Name string `json:"name"`
}

type Loc struct {
	ID    int    `json:"id_loc"`
	Name  string `json:"name"`
	IDJud int    `json:"id_jud"`
}

type Address struct {
	ID      int    `json:"id_address"`
	IDLoc   int    `json:"id_loc"`
	Address string `json:"address"`
}

type Hospital struct {
	ID        int    `json:"id_hospital"`
	Name      string `json:"name"`
	IDAddress int    `json:"id_address"`
}

type Person struct {
	ID               int       `json:"id_person"`
	FirstName        string    `json:"f_name"`
	LastName         string    `json:"l_name"`
	CNP              string    `json:"cnp"`
	BornDate         time.Time `json:"born_date"`
	IDAddress        int       `json:"id_address"`
	IDVirtualAddress int       `json:"id_virtual_address"`
}

type Doctor struct {
	ID       int    `json:"id_doctor"`
	IDPerson int    `json:"id_person"`
	Parafa   string `json:"parafa"`
}

type Patient struct {
	ID       int `json:"id_patient"`
	IDPerson int `json:"id_person"`
}

type DoctorHospital struct {
	IDDoctorHospital int `json:"id_doctor_hospital"`
	IDDoctor         int `json:"id_doctor"`
	IDHospital       int `json:"id_hospital"`
}

type PatientDoctor struct {
	IDPatientDoctor  int    `json:"id_patient_doctor"`
	IDPatient        int    `json:"id_patient"`
	IDDoctorHospital int    `json:"id_doctor_hospital"`
	Status           string `json:"status"`
}
