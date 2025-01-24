package config

import (
	"apotek-management/models"
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	dsn := "root:password@tcp(127.0.0.1:3306)/apotek?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Gagal koneksi ke database:", err)
	}

	fmt.Println("Database berhasil terkoneksi")

	err = db.AutoMigrate(
		&models.Stok{},
		&models.TagObat{},
		&models.TipeObat{},
		&models.Obat{},
		&models.Transaksi{},
	)
	if err != nil {
		log.Fatal("err migrasi:", err)
	}

	DB = db
}
