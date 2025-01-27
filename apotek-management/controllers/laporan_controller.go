package controllers

import (
	"apotek-management/config"
	"apotek-management/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetLaporanTransaksi(c *gin.Context) {
	var transaksi []models.Transaksi
	err := config.DB.Preload("Obats.Obat").Find(&transaksi).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch transactions"})
		return
	}

	var laporan []map[string]interface{}
	for _, trx := range transaksi {
		for _, detail := range trx.Obats {
			laporan = append(laporan, map[string]interface{}{
				"Kode Transaksi":  trx.KodeTransaksi,
				"Tanggal":         trx.CreatedAt,
				"Obat":            detail.Obat.NamaObat,
				"Jumlah":          detail.Jumlah,
				"Harga Satuan":    detail.Obat.Harga,
				"Total Harga":     detail.Jumlah * int(detail.Obat.Harga),
				"Jenis Transaksi": "Pengeluaran",
				"Keterangan":      "Penjualan",
			})
		}
	}

	c.JSON(http.StatusOK, laporan)
}

func GetLaporanStok(c *gin.Context) {
	var stok []models.Stok
	err := config.DB.Find(&stok).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch stock data"})
		return
	}

	var laporan []map[string]interface{}
	for _, s := range stok {
		laporan = append(laporan, map[string]interface{}{
			"Obat":             s.ObatID,
			"Tipe Transaksi":   s.TipeTransaksi,
			"Stok Awal":        s.StokAwal,
			"Jumlah Transaksi": s.JumlahStokTransaksi,
			"Stok Akhir":       s.StokAkhir,
			"Keterangan":       s.Keterangan,
			"Tanggal":          s.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, laporan)
}

func GetLaporanLabaRugi(c *gin.Context) {
	var transaksi []models.Transaksi
	err := config.DB.Preload("Obats.Obat").Find(&transaksi).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch transactions"})
		return
	}

	var pendapatan, pengeluaran int
	for _, trx := range transaksi {
		for _, detail := range trx.Obats {
			pendapatan += detail.Jumlah * int(detail.Obat.Harga)
			// Asumsi harga beli adalah 70% dari harga jual
			pengeluaran += detail.Jumlah * int(float64(detail.Obat.Harga)*0.7)
		}
	}

	laba := pendapatan - pengeluaran

	c.JSON(http.StatusOK, map[string]interface{}{
		"pendapatan":  pendapatan,
		"pengeluaran": pengeluaran,
		"laba_bersih": laba,
	})
}
