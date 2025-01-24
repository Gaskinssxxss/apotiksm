package controllers

import (
	"apotek-management/config"
	"apotek-management/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

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
	id := c.Param("id")
	var transaksi models.Transaksi

	if err := config.DB.First(&transaksi, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaksi not found"})
		return
	}

	if err := c.ShouldBindJSON(&transaksi); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := config.DB.Save(&transaksi).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update transaksi: " + err.Error()})
		return
	}

	if err := config.DB.
		Preload("Obat").
		Preload("Obat.TipeObat").
		Preload("Obat.TagObat").
		First(&transaksi, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load updated transaksi: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, transaksi)
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
