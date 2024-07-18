package models

type PatientDoctor struct {
	IDPatientDoctor  int    `json:"id_patient_doctor"`
	IDPatient        int    `json:"id_patient"`
	IDDoctorHospital int    `json:"id_doctor_hospital"`
	Status           string `json:"status"`
}
