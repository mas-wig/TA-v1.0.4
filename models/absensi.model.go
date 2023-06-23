package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EncodePresensi struct {
	ID             string    `gorm:"type:char(36);primary_key"`
	NamaSiswa      string    `gorm:"type:varchar(255);not null"`
	IDSiswa        uuid.UUID `gorm:"type:varchar(255);not null"`
	Hari           string    `gorm:"type:varchar(255);not null"`
	Lokasi         string    `gorm:"type:varchar(255);not null"`
	TanggalWaktu   string    `gorm:"type:varchar(255);not null"`
	Kehadiran      string    `gorm:"type:varchar(255);not null"`
	InformasiMedis string    `gorm:"type:text;not null"`
}

func (note *EncodePresensi) BeforeCreate(tx *gorm.DB) (err error) {
	note.ID = uuid.New().String()
	return nil
}

type DecodePresensi struct {
	ID             string    `gorm:"type:char(36);primary_key"`
	NamaSiswa      string    `gorm:"type:varchar(255);not null"`
	IDSiswa        uuid.UUID `gorm:"type:varchar(255);not null"`
	Hari           string    `gorm:"type:varchar(255);not null"`
	Lokasi         string    `gorm:"type:varchar(255);not null"`
	TanggalWaktu   string    `gorm:"type:varchar(255);not null"`
	Kehadiran      string    `gorm:"type:varchar(255);not null"`
	InformasiMedis string    `gorm:"type:text;not null"`
}

func (note *DecodePresensi) BeforeCreate(tx *gorm.DB) (err error) {
	note.ID = uuid.New().String()
	return nil
}

type CreatePresensi struct {
	Hari           string `form:"hari" binding:"required"`
	Lokasi         string `form:"lokasi" binding:"required"`
	Kehadiran      string `form:"kehadiran" binding:"required"`
	TanggalWaktu   string `form:"date" binding:"required"`
	InformasiMedis string `form:"catatankesehatan" binding:"required"`
}

type DecodeKey struct {
	Key string `form:"decodekey" binding:"required"`
}
