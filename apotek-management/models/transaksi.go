package models

import (
	"time"
)

type Transaksi struct {
	ID            uint              `json:"id_transaksi" gorm:"primaryKey"`
	KodeTransaksi string            `json:"kode_transaksi" gorm:"type:varchar(20);unique;not null"`
	TotalHarga    int               `json:"total_harga" gorm:"not null"`
	Status        string            `json:"status" gorm:"type:varchar(50);not null"`
	Obats         []TransaksiDetail `json:"obats" gorm:"foreignKey:TransaksiID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	CreatedAt     time.Time         `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt     time.Time         `json:"updated_at" gorm:"autoUpdateTime"`
}
