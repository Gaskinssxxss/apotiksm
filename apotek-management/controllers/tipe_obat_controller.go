package controllers

import (
	"apotek-management/config"
	"apotek-management/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetAllTipeObat(c *gin.Context) {
	var tipeObats []models.TipeObat
	if err := config.DB.Find(&tipeObats).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error fetching tipe obat"})
		return
	}
	c.JSON(http.StatusOK, tipeObats)
}

func GetTipeObatByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var tipeObat models.TipeObat
	if err := config.DB.First(&tipeObat, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Tipe obat not found"})
		return
	}
	c.JSON(http.StatusOK, tipeObat)
}

func CreateTipeObat(c *gin.Context) {
	var tipeObat models.TipeObat

	if err := c.ShouldBindJSON(&tipeObat); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	if tipeObat.NamaTipe == "" || tipeObat.KodeTipe == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "NamaTipe and KodeTipe are required fields"})
		return
	}

	if err := config.DB.Create(&tipeObat).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create Tipe Obat: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Tipe Obat created successfully",
		"data":    tipeObat,
	})
}

func UpdateTipeObat(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var tipeObat models.TipeObat
	if err := config.DB.First(&tipeObat, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Tipe obat not found"})
		return
	}

	if err := c.ShouldBindJSON(&tipeObat); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input"})
		return
	}

	if err := config.DB.Save(&tipeObat).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error updating tipe obat"})
		return
	}

	c.JSON(http.StatusOK, tipeObat)
}

func DeleteTipeObat(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var tipeObat models.TipeObat
	if err := config.DB.First(&tipeObat, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Tipe obat not found"})
		return
	}

	if err := config.DB.Delete(&tipeObat).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error deleting tipe obat"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tipe obat deleted successfully"})
}

func CreateBatchTipeObat(c *gin.Context) {
	var tipeObats []models.TipeObat
	if err := c.ShouldBindJSON(&tipeObats); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := config.DB.Create(&tipeObats).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, tipeObats)
}

func UpdateBatchTipeObat(c *gin.Context) {
	var tipeObats []models.TipeObat
	if err := c.ShouldBindJSON(&tipeObats); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for _, tipeObat := range tipeObats {
		if err := config.DB.Save(&tipeObat).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, tipeObats)
}

func DeleteBatchTipeObat(c *gin.Context) {
	var tipeObatIDs []uint
	if err := c.ShouldBindJSON(&tipeObatIDs); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := config.DB.Where("id IN ?", tipeObatIDs).Delete(&models.TipeObat{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tipe obat deleted successfully"})
}
