package models

import "time"

type TagObat struct {
	ID        uint      `json:"id_tag_obat" gorm:"primaryKey;column:id_tag_obat"`
	NamaTag   string    `json:"nama_tag" gorm:"type:varchar(100);not null"`
	CreatedAt time.Time `json:"created_at" gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `json:"updated_at" gorm:"type:timestamp;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
}

func (TagObat) TableName() string {
	return "tag_obat"
}
