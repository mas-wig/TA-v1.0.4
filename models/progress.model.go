package models

import (
	"mime/multipart"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EncodeProgressLatihan struct {
	CreatedAt           time.Time `gorm:"not null"`
	UpdatedAt           time.Time `gorm:"not null"`
	ID                  string    `gorm:"type:char(36);primary_key"`
	TingkatKemajuan     string    `gorm:"type:varchar(255)"`
	CatatanPerkembangan string    `gorm:"type:text"`
	PenilaianTeknik     string    `gorm:"type:varchar(255)"`
	PenilaianDisiplin   string    `gorm:"type:varchar(255)"`
	CatatanEvaluasi     string    `gorm:"type:text"`
	Media               string    `gorm:"type:varchar(255)"`
	SiswaID             uuid.UUID `gorm:"type:varchar(90)"`
}

func (note *EncodeProgressLatihan) BeforeCreate(tx *gorm.DB) (err error) {
	note.ID = uuid.New().String()
	return nil
}

type DecodeProgressLatihan struct {
	CreatedAt           time.Time `gorm:"not null"`
	UpdatedAt           time.Time `gorm:"not null"`
	ID                  string    `gorm:"type:char(36);primary_key"`
	TingkatKemajuan     string    `gorm:"type:varchar(255)"`
	CatatanPerkembangan string    `gorm:"type:text"`
	PenilaianTeknik     string    `gorm:"type:varchar(255)"`
	PenilaianDisiplin   string    `gorm:"type:varchar(255)"`
	CatatanEvaluasi     string    `gorm:"type:text"`
	Media               string    `gorm:"type:varchar(255)"`
	SiswaID             uuid.UUID `gorm:"type:varchar(90)"`
}

func (note *DecodeProgressLatihan) BeforeCreate(tx *gorm.DB) (err error) {
	note.ID = uuid.New().String()
	return nil
}

type CreateProgress struct {
	Media               *multipart.FileHeader `form:"media" binding:"required"`
	TingkatKemajuan     string                `form:"tingkatkemajuan" binding:"required"`
	CatatanPerkembangan string                `form:"progressnote" binding:"required"`
	PenilaianTeknik     string                `form:"nilaiteknik" binding:"required"`
	PenilaianDisiplin   string                `form:"nilaidisiplin" binding:"required"`
	CatatanEvaluasi     string                `form:"evaluasinote" binding:"required"`
	SiswaID             uuid.UUID
}
