package models

import "time"

type Obat struct {
	ID         uint      `json:"id_obat" gorm:"primaryKey;column:id_obat"`
	KodeObat   string    `json:"kode_obat" gorm:"type:varchar(100);not null"`
	NamaObat   string    `json:"nama_obat" gorm:"type:varchar(100);not null"`
	Dosis      string    `json:"dosis_obat" gorm:"column:dosis_obat;type:varchar(255);not null"`
	Gambar     string    `json:"gambar_obat" gorm:"column:gambar_obat;type:varchar(255);not null"`
	Deskripsi  string    `json:"deskripsi" gorm:"type:text"`
	Harga      uint64    `json:"harga_obat" gorm:"type:bigint unsigned;not null"`
	TipeObatID uint      `json:"id_tipe_obat" gorm:"not null"`
	TipeObat   TipeObat  `json:"tipe_obat" gorm:"foreignKey:TipeObatID;references:ID"`
	Tags       []TagObat `json:"tags" gorm:"many2many:obat_tags;joinForeignKey:ObatID;joinReferences:TagObatID"`
	Stok       []Stok    `json:"stok" gorm:"foreignKey:ObatID"` // Relasi one-to-many dengan Stok
	CreatedAt  time.Time `json:"created_at" gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
	UpdatedAt  time.Time `json:"updated_at" gorm:"type:timestamp;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
}

type ObatTag struct {
	ObatID    uint      `json:"obat_id" gorm:"primaryKey"`
	TagObatID uint      `json:"tag_obat_id" gorm:"primaryKey"`
	CreatedAt time.Time `json:"created_at" gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
}
