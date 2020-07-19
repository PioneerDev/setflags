package models

import (
	"time"

	"github.com/gofrs/uuid"
	"github.com/jinzhu/gorm"
)

// Evidence entity
type Evidence struct {
	ID           uuid.UUID `gorm:"type:uuid;primary_key;" json:"id"`
	AttachmentID string    `json:"attachment_id"`
	FlagID       uuid.UUID `json:"flag_id"`
	URL          string    `json:"url"`
	Type         string    `json:"type"`
	Period       int       `json:"period"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// CreateEvidence create evidence
func CreateEvidence(flagID uuid.UUID, attachmentID, mediaType, url string, period int) bool {

	db.Create(&Evidence{
		AttachmentID: attachmentID,
		FlagID:       flagID,
		Type:         mediaType,
		URL:          url,
		Period:       period,
	})
	return true
}

// FindEvidencesByFlag 返回自昨天开始的evidence
func FindEvidencesByFlag(flagID uuid.UUID, currentPage, pageSize int) (evidences []Evidence, total int) {
	// 获取当前时间
	now := time.Now()
	// 回到昨天
	yesterday := now.Add(time.Hour * -24)
	yesterday = time.Date(yesterday.Year(), yesterday.Month(), yesterday.Day(), 0, 0, 0, 0, yesterday.Location())

	skip := (currentPage - 1) * pageSize
	db.Offset(skip).Limit(pageSize).Where("flag_id = ? and created_at >= ?", flagID.String(), yesterday).Order("created_at desc").Find(&evidences)
	db.Model(&Evidence{}).Where("flag_id = ? and created_at >= ?", flagID.String(), yesterday).Count(&total)
	return
}

// GetAllEvidenceByFlagID GetAllEvidenceByFlagID
func GetAllEvidenceByFlagID(flagID uuid.UUID, currentPage, pageSize int) (evidences []Evidence, total int) {
	skip := (currentPage - 1) * pageSize
	db.Offset(skip).Limit(pageSize).Where("flag_id = ?", flagID).Order("period desc, updated_at desc").Find(&evidences)
	db.Model(&Evidence{}).Where("flag_id = ?", flagID).Count(&total)
	return
}

// FindEvidencesByFlagAndPeriod FindEvidencesByFlagAndPeriod
func FindEvidencesByFlagAndPeriod(flagID uuid.UUID, currentPage, pageSize, period int) (evidences []Evidence, total int) {

	skip := (currentPage - 1) * pageSize
	db.Offset(skip).Limit(pageSize).Where("flag_id = ? and period = ?", flagID, period).
		Order("period desc, updated_at desc").Find(&evidences)
	db.Model(&Evidence{}).Where("flag_id = ? and period = ?", flagID, period).Count(&total)
	return
}

// FindEvidenceByFlagIDAndAttachmentID find evidence by flagID
func FindEvidenceByFlagIDAndAttachmentID(flagID, attachmentID uuid.UUID, currentPage, pageSize int) (evidences []Evidence, total int) {
	skip := (currentPage - 1) * pageSize
	db.Offset(skip).
		Limit(pageSize).
		Where("flag_id = ? and attachment_id = ?", flagID.String(), attachmentID.String()).
		Find(&evidences)
	db.Model(&Evidence{}).Where("flag_id = ? and attachment_id = ?", flagID.String(), attachmentID.String()).Count(&total)
	return
}

// BeforeCreate will set a UUID rather than numeric ID.
func (e *Evidence) BeforeCreate(scope *gorm.Scope) error {
	evidenceID, _ := uuid.NewV4()
	scope.SetColumn("ID", evidenceID)
	scope.SetColumn("CreatedAt", time.Now())
	return nil
}

// BeforeUpdate set field updateAt
func (e *Evidence) BeforeUpdate(scope *gorm.Scope) error {
	scope.SetColumn("UpdatedAt", time.Now())
	return nil
}
