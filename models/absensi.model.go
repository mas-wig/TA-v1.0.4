package models

import (
	"github.com/google/uuid"
)

type Presensi struct {
	ID             string    `gorm:"type:varchar(36);default:uuid();primary_key"`
	NamaSiswa      string    `gorm:"type:varchar(255);not null"`
	IDSiswa        uuid.UUID `gorm:"not null"`
	Hari           string    `gorm:"type:varchar(20);not null"`
	Lokasi         string    `gorm:"type:varchar(255);not null"`
	TanggalWaktu   string    `gorm:"type:varchar(40);not null"`
	Kehadiran      string    `gorm:"type:varchar(20);not null"`
	InformasiMedis string    `gorm:"type:text;not null"`
}

type CreatePresensi struct {
	Hari           string `form:"hari" binding:"required"`
	Lokasi         string `form:"lokasi" binding:"required"`
	Kehadiran      string `form:"kehadiran" binding:"required"`
	TanggalWaktu   string `form:"date" binding:"required"`
	InformasiMedis string `form:"catatankesehatan" binding:"required"`
}
