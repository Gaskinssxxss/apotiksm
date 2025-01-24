package models

import "time"

type Stok struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	Jumlah        int       `json:"jumlah" gorm:"not null"`
	TipeTransaksi string    `json:"tipe_transaksi" gorm:"type:enum('MASUK', 'KELUAR');not null"`
	Keterangan    string    `json:"keterangan" gorm:"type:text"`
	ObatID        uint      `json:"obat_id" gorm:"unique;not null"`
	CreatedAt     time.Time `json:"created_at" gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
	UpdatedAt     time.Time `json:"updated_at" gorm:"type:timestamp;default:CURRENT_TIMESTAMP;autoUpdateTime"`
}
