package controllers

import (
	"apotek-management/config"
	"apotek-management/models"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// func CreateTransaksi(c *gin.Context) {
// 	var transaksi models.Transaksi

// 	if err := c.ShouldBindJSON(&transaksi); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
// 		return
// 	}

// 	for _, detail := range transaksi.Obats {
// 		if detail.Jumlah <= 0 {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": "Jumlah obat harus lebih dari 0"})
// 			return
// 		}
// 	}

// 	tx := config.DB.Begin()

// 	if err := tx.Create(&transaksi).Error; err != nil {
// 		tx.Rollback()
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create transaksi: " + err.Error()})
// 		return
// 	}

// 	for _, detail := range transaksi.Obats {
// 		var stok models.Stok
// 		if err := tx.Where("obat_id = ?", detail.ObatID).First(&stok).Error; err != nil {
// 			tx.Rollback()
// 			c.JSON(http.StatusNotFound, gin.H{"error": "Obat with ID not found"})
// 			return
// 		}

// 		if stok.StokAkhir < detail.Jumlah {
// 			tx.Rollback()
// 			c.JSON(http.StatusBadRequest, gin.H{"error": "Stok tidak mencukupi untuk obat dengan ID " + fmt.Sprint(detail.ObatID)})
// 			return
// 		}

// 		stok.StokAkhir -= detail.Jumlah
// 		stok.JumlahStokTransaksi = detail.Jumlah
// 		stok.TipeTransaksi = "KELUAR"

// 		if err := tx.Save(&stok).Error; err != nil {
// 			tx.Rollback()
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update stok: " + err.Error()})
// 			return
// 		}
// 	}

// 	tx.Commit()

// 	c.JSON(http.StatusCreated, gin.H{"message": "Transaksi berhasil dibuat", "data": transaksi})
// }

func CreateTransaksi(c *gin.Context) {
	var transaksi models.Transaksi

	if err := c.ShouldBindJSON(&transaksi); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	for _, detail := range transaksi.Obats {
		if detail.Jumlah <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Jumlah obat harus lebih dari 0"})
			return
		}
	}

	if err := config.DB.Create(&transaksi).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create transaksi: " + err.Error()})
		return
	}

	var createdTransaksi models.Transaksi
	if err := config.DB.
		Preload("Obats.Obat").
		Preload("Obats.Obat.TipeObat").
		Preload("Obats.Obat.Tags").
		Preload("Obats.Obat.Stok").
		First(&createdTransaksi, transaksi.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load transaksi with relations: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, createdTransaksi)
}

func GetAllTransaksi(c *gin.Context) {
	var transaksiList []models.Transaksi

	if err := config.DB.
		Preload("Obats").
		Preload("Obats.Obat.Tags").
		Preload("Obats.Obat.TipeObat").Preload("Obats.Obat.Stok").
		Find(&transaksiList).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, transaksiList)
}

func GetTransaksiByID(c *gin.Context) {
	id := c.Param("id")
	var transaksi models.Transaksi
	if err := config.DB.Preload("Obats.Obat.Tags").Preload("Obats.Obat.TipeObat").First(&transaksi, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaksi not found"})
		return
	}

	c.JSON(http.StatusOK, transaksi)
}

func UpdateTransaksi(c *gin.Context) {
	id := c.Param("id") // ID transaksi dari parameter URL
	var transaksiBaru models.Transaksi
	var transaksiLama models.Transaksi

	// Ambil data transaksi lama
	if err := config.DB.Preload("Obats").First(&transaksiLama, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaksi not found"})
		return
	}

	// Bind data transaksi baru dari request
	if err := c.ShouldBindJSON(&transaksiBaru); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Gunakan transaksi database untuk menjaga konsistensi
	if err := config.DB.Transaction(func(tx *gorm.DB) error {
		for _, detailBaru := range transaksiBaru.Obats {
			// Ambil detail transaksi lama
			var detailLama models.TransaksiDetail
			if err := tx.Where("id = ?", detailBaru.ID).First(&detailLama).Error; err != nil {
				return errors.New("Transaksi detail not found for id: " + fmt.Sprint(detailBaru.ID))
			}

			// Hitung selisih jumlah
			selisihJumlah := detailBaru.Jumlah - detailLama.Jumlah

			// Ambil stok terakhir untuk obat terkait
			var stok models.Stok
			if err := tx.Where("obat_id = ?", detailBaru.ObatID).Order("created_at desc").First(&stok).Error; err != nil {
				return errors.New("stok not found for obat_id: " + fmt.Sprint(detailBaru.ObatID))
			}

			// Update stok berdasarkan selisih jumlah
			if selisihJumlah > 0 {
				// Jika jumlah bertambah, kurangi stok
				if stok.StokAkhir < selisihJumlah {
					return errors.New("stok tidak mencukupi untuk obat_id: " + fmt.Sprint(detailBaru.ObatID))
				}
				stok.StokAkhir -= selisihJumlah
			} else if selisihJumlah < 0 {
				// Jika jumlah berkurang, tambahkan stok
				stok.StokAkhir += -selisihJumlah
			}

			// Update stok_awal dan jumlah_stok_transaksi
			stok.StokAwal = stok.StokAkhir + selisihJumlah
			stok.JumlahStokTransaksi = detailBaru.Jumlah

			// Simpan perubahan stok
			if err := tx.Save(&stok).Error; err != nil {
				return err
			}

			// Perbarui detail transaksi
			if err := tx.Model(&detailLama).Updates(detailBaru).Error; err != nil {
				return err
			}
		}

		// Perbarui transaksi utama
		if err := tx.Model(&transaksiLama).Updates(transaksiBaru).Error; err != nil {
			return err
		}

		return nil
	}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update transaksi: " + err.Error()})
		return
	}

	// Kembalikan respon sukses
	c.JSON(http.StatusOK, transaksiBaru)
}

func DeleteTransaksi(c *gin.Context) {
	id := c.Param("id")
	var transaksi models.Transaksi
	if err := config.DB.First(&transaksi, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaksi not found"})
		return
	}

	if err := config.DB.Preload("Obat").
		Preload("Obat.TipeObat").
		Preload("Obat.TagObat").Delete(&transaksi).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Transaksi deleted successfully"})
}

func CreateBatchTransaksi(c *gin.Context) {
	var transaksiList []models.Transaksi
	if err := c.ShouldBindJSON(&transaksiList); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	if err := config.DB.Create(&transaksiList).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to create transactions"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Batch transactions created successfully",
		"data":    transaksiList,
	})
}

func UpdateBatchTransaksi(c *gin.Context) {
	var transaksiList []models.Transaksi
	if err := c.ShouldBindJSON(&transaksiList); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	for _, transaksi := range transaksiList {
		if err := config.DB.Save(&transaksi).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to update some transactions"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Batch transactions updated successfully",
		"data":    transaksiList,
	})
}

func DeleteBatchTransaksi(c *gin.Context) {
	var ids []uint
	if err := c.ShouldBindJSON(&ids); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	if err := config.DB.Delete(&models.Transaksi{}, "id IN ?", ids).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to delete transactions"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Batch transactions deleted successfully",
	})
}
