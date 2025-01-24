package models

import "time"

type TipeObat struct {
	ID        uint      `json:"id_tipe_obat" gorm:"primaryKey;column:id_tipe_obat;type:bigint unsigned"`
	NamaTipe  string    `json:"nama_tipe" gorm:"type:varchar(100);not null"`
	KodeTipe  string    `json:"kode_tipe" gorm:"type:varchar(100);not null"`
	CreatedAt time.Time `json:"created_at" gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `json:"updated_at" gorm:"type:timestamp;default:CURRENT_TIMESTAMP;autoUpdateTime"`
}

func (TipeObat) TableName() string {
	return "tipe_obats"
}
