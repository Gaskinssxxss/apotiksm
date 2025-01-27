package models

import (
	"time"
)

type TransaksiDetail struct {
	ID          uint      `json:"id_transaksi_detail" gorm:"primaryKey;autoIncrement"`
	TransaksiID uint      `json:"id_transaksi" gorm:"not null"`
	ObatID      uint      `json:"id_obat" gorm:"not null"`
	Jumlah      int       `json:"jumlah" gorm:"not null;check:jumlah > 0"`
	Obat        Obat      `json:"obat" gorm:"foreignKey:ObatID;references:ID;`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
