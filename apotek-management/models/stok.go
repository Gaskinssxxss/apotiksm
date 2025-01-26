package models

import "time"

type Stok struct {
	ID                  uint      `json:"id_stok" gorm:"primaryKey;column:id_stok"`
	StokAwal            int       `json:"stok_awal"`
	StokAkhir           int       `json:"stok_akhir"`
	JumlahStokTransaksi int       `json:"jumlah_stok_transaksi"`
	TipeTransaksi       string    `json:"tipe_transaksi" gorm:"type:enum('MASUK', 'KELUAR');not null"`
	Keterangan          string    `json:"keterangan" gorm:"type:text"`
	ObatID              uint      `json:"obat_id" gorm:"not null"`
	CreatedAt           time.Time `json:"created_at" gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
	UpdatedAt           time.Time `json:"updated_at" gorm:"type:timestamp;default:CURRENT_TIMESTAMP;autoUpdateTime"`
}
