package models

type DoctorHospital struct {
	IDDoctorHospital int `json:"id_doctor_hospital"`
	IDDoctor         int `json:"id_doctor"`
	IDHospital       int `json:"id_hospital"`
}
