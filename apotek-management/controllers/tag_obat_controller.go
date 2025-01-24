package controllers

import (
	"apotek-management/config"
	"apotek-management/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetAllTagObat(c *gin.Context) {
	var tagObats []models.TagObat
	if err := config.DB.Find(&tagObats).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to fetch tag obat"})
		return
	}
	c.JSON(http.StatusOK, tagObats)
}

func GetTagObatByID(c *gin.Context) {
	id := c.Param("id")
	var tagObat models.TagObat
	if err := config.DB.First(&tagObat, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tag obat not found"})
		return
	}
	c.JSON(http.StatusOK, tagObat)
}

func CreateTagObat(c *gin.Context) {
	var tagObat models.TagObat
	if err := c.ShouldBindJSON(&tagObat); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	if err := config.DB.Create(&tagObat).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to create tag obat"})
		return
	}

	c.JSON(http.StatusOK, tagObat)
}

func UpdateTagObat(c *gin.Context) {
	id := c.Param("id")
	var tagObat models.TagObat
	if err := config.DB.First(&tagObat, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tag obat not found"})
		return
	}

	if err := c.ShouldBindJSON(&tagObat); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	if err := config.DB.Save(&tagObat).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to update tag obat"})
		return
	}

	c.JSON(http.StatusOK, tagObat)
}

func DeleteTagObat(c *gin.Context) {
	id := c.Param("id")
	var tagObat models.TagObat
	if err := config.DB.First(&tagObat, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tag obat not found"})
		return
	}

	if err := config.DB.Delete(&tagObat).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to delete tag obat"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tag obat deleted successfully"})
}

func CreateBatchTagObat(c *gin.Context) {
	var tagObats []models.TagObat
	if err := c.ShouldBindJSON(&tagObats); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	if err := config.DB.Create(&tagObats).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to create tag obat"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Batch create successful", "data": tagObats})
}

func UpdateBatchTagObat(c *gin.Context) {
	var tagObats []models.TagObat
	if err := c.ShouldBindJSON(&tagObats); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	for _, tagObat := range tagObats {
		if err := config.DB.Save(&tagObat).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to update tag obat"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Batch update successful", "data": tagObats})
}

func DeleteBatchTagObat(c *gin.Context) {
	var ids []uint
	if err := c.ShouldBindJSON(&ids); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	if err := config.DB.Delete(&models.TagObat{}, ids).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to delete tag obat"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Batch delete successful"})
}
