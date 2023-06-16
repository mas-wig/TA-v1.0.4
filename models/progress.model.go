package models

import (
	"mime/multipart"
	"time"

	"github.com/google/uuid"
)

type EncodeProgressLatihan struct {
	ID                  string    `gorm:"type:varchar(36);default:uuid();primary_key"`
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

type CreateProgress struct {
	SiswaID             uuid.UUID
	TingkatKemajuan     string                `form:"tingkatkemajuan" binding:"required"`
	CatatanPerkembangan string                `form:"progressnote" binding:"required"`
	PenilaianTeknik     string                `form:"nilaiteknik" binding:"required"`
	PenilaianDisiplin   string                `form:"nilaidisiplin" binding:"required"`
	CatatanEvaluasi     string                `form:"evaluasinote" binding:"required"`
	Media               *multipart.FileHeader `form:"media" binding:"required"`
}
