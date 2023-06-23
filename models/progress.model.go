package models

import (
	"mime/multipart"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EncodeProgressLatihan struct {
	ID                  string    `gorm:"type:char(36);primary_key"`
	SiswaID             uuid.UUID `gorm:"type:varchar(90)"`
	TingkatKemajuan     string    `gorm:"type:varchar(255)"`
	CatatanPerkembangan string    `gorm:"type:text"`
	PenilaianTeknik     string    `gorm:"type:varchar(255)"`
	PenilaianDisiplin   string    `gorm:"type:varchar(255)"`
	CatatanEvaluasi     string    `gorm:"type:text"`
	Media               string    `gorm:"type:varchar(255)"`
	CreatedAt           time.Time `gorm:"not null"`
	UpdatedAt           time.Time `gorm:"not null"`
}

func (note *EncodeProgressLatihan) BeforeCreate(tx *gorm.DB) (err error) {
	note.ID = uuid.New().String()
	return nil
}

type DecodeProgressLatihan struct {
	ID                  string    `gorm:"type:char(36);primary_key"`
	SiswaID             uuid.UUID `gorm:"type:varchar(90)"`
	TingkatKemajuan     string    `gorm:"type:varchar(255)"`
	CatatanPerkembangan string    `gorm:"type:text"`
	PenilaianTeknik     string    `gorm:"type:varchar(255)"`
	PenilaianDisiplin   string    `gorm:"type:varchar(255)"`
	CatatanEvaluasi     string    `gorm:"type:text"`
	Media               string    `gorm:"type:varchar(255)"`
	CreatedAt           time.Time `gorm:"not null"`
	UpdatedAt           time.Time `gorm:"not null"`
}

func (note *DecodeProgressLatihan) BeforeCreate(tx *gorm.DB) (err error) {
	note.ID = uuid.New().String()
	return nil
}

type CreateProgress struct {
	SiswaID             uuid.UUID
	TingkatKemajuan     string                `form:"tingkatkemajuan" binding:"required"`
	CatatanPerkembangan string                `form:"progressnote" binding:"required"`
	PenilaianTeknik     string                `form:"nilaiteknik" binding:"required"`
	PenilaianDisiplin   string                `form:"nilaidisiplin" binding:"required"`
	CatatanEvaluasi     string                `form:"evaluasinote" binding:"required"`
	Media               *multipart.FileHeader `form:"media" binding:"required"`
}
