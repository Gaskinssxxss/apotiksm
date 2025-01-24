package controllers

import (
	"apotek-management/config"
	"apotek-management/models"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CreateObat(c *gin.Context) {
	namaObat := c.PostForm("nama_obat")
	dosisObat := c.PostForm("dosis_obat")
	deskripsi := c.PostForm("deskripsi")
	idTipeObat, err := strconv.Atoi(c.PostForm("id_tipe_obat"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id_tipe_obat"})
		return
	}
	hargaObat, err := strconv.ParseUint(c.PostForm("harga_obat"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid harga_obat"})
		return
	}

	tagIDs := c.PostFormArray("tags[]")
	var tags []models.TagObat
	for _, tagID := range tagIDs {
		var tag models.TagObat
		if err := config.DB.First(&tag, tagID).Error; err == nil {
			tags = append(tags, tag)
		}
	}

	file, err := c.FormFile("gambar")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Gambar file is required"})
		return
	}
	filePath := "uploads/obat/" + file.Filename
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save image"})
		return
	}

	obat := models.Obat{
		NamaObat:   namaObat,
		Dosis:      dosisObat,
		Deskripsi:  deskripsi,
		TipeObatID: uint(idTipeObat),
		Harga:      hargaObat,
		Gambar:     filePath,
		Tags:       tags,
	}
	if err := config.DB.Create(&obat).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, obat)
}

func CreateBatchObat(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid multipart form data"})
		return
	}

	files := form.File["gambar"]
	data := form.Value["data"]

	if len(files) != len(data) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Mismatch between number of images and data"})
		return
	}

	var obatList []models.Obat
	for i, jsonData := range data {
		var obat models.Obat
		if err := json.Unmarshal([]byte(jsonData), &obat); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON in data"})
			return
		}

		file := files[i]
		filePath := "uploads/obat/" + file.Filename
		if err := c.SaveUploadedFile(file, filePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save image: " + err.Error()})
			return
		}

		obat.Gambar = filePath
		obatList = append(obatList, obat)
	}

	if err := config.DB.Create(&obatList).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, obatList)
}

func GetAllObat(c *gin.Context) {
	var obatList []models.Obat

	if err := config.DB.Preload("TipeObat").Preload("Tags").Preload("Stok").Find(&obatList).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for _, obat := range obatList {
		log.Printf("Obat: %+v\n", obat)
	}

	c.JSON(http.StatusOK, obatList)
}

func GetObatByID(c *gin.Context) {
	id := c.Param("id")
	var obat models.Obat

	if err := config.DB.
		Preload("TipeObat").
		Preload("Tags").Preload("Stok").
		First(&obat, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Obat not found"})
		return
	}

	c.JSON(http.StatusOK, obat)
}

func UpdateObat(c *gin.Context) {
	id := c.Param("id")
	var existingObat models.Obat

	if err := config.DB.First(&existingObat, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Obat not found"})
		return
	}

	var updatedObat struct {
		KodeObat   string           `json:"kode_obat"`
		NamaObat   string           `json:"nama_obat"`
		Dosis      string           `json:"dosis_obat"`
		Deskripsi  string           `json:"deskripsi"`
		HargaObat  uint64           `json:"harga_obat"`
		TipeObatID uint             `json:"id_tipe_obat"`
		Tags       []models.TagObat `json:"tags"`
	}

	if err := c.ShouldBindJSON(&updatedObat); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	file, err := c.FormFile("gambar")
	if err == nil {
		if existingObat.Gambar != "" {
			_ = os.Remove(existingObat.Gambar)
		}

		filePath := "uploads/obat/" + file.Filename
		if err := c.SaveUploadedFile(file, filePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save image: " + err.Error()})
			return
		}

		existingObat.Gambar = filePath
	}

	existingObat.KodeObat = updatedObat.KodeObat
	existingObat.NamaObat = updatedObat.NamaObat
	existingObat.Dosis = updatedObat.Dosis
	existingObat.Deskripsi = updatedObat.Deskripsi
	existingObat.Harga = updatedObat.HargaObat
	existingObat.TipeObatID = updatedObat.TipeObatID

	if err := config.DB.Save(&existingObat).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update Obat: " + err.Error()})
		return
	}

	if len(updatedObat.Tags) > 0 {
		var tagIDs []uint
		for _, tag := range updatedObat.Tags {
			tagIDs = append(tagIDs, tag.ID)
		}

		var tags []models.TagObat
		if err := config.DB.Where("id_tag_obat IN ?", tagIDs).Find(&tags).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update tags: " + err.Error()})
			return
		}
		if err := config.DB.Model(&existingObat).Association("Tags").Replace(&tags); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update tags association: " + err.Error()})
			return
		}
	}

	if err := config.DB.Preload("TipeObat").Preload("Tags").First(&existingObat, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load updated Obat: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, existingObat)
}

func UpdateBatchObat(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid multipart form data"})
		return
	}

	files := form.File["gambar"]
	data := form.Value["data"]

	if len(files) > 0 && len(files) != len(data) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Mismatch between number of images and data"})
		return
	}

	var obatList []models.Obat
	for i, jsonData := range data {
		var obat models.Obat

		if err := json.Unmarshal([]byte(jsonData), &obat); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON in data: " + err.Error()})
			return
		}

		var existingObat models.Obat
		if err := config.DB.First(&existingObat, obat.ID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Obat not found with ID: " + fmt.Sprint(obat.ID)})
			return
		}

		if len(files) > 0 {
			file := files[i]

			if existingObat.Gambar != "" {
				_ = os.Remove(existingObat.Gambar)
			}

			filePath := "uploads/obat/" + file.Filename
			if err := c.SaveUploadedFile(file, filePath); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save image for ID: " + fmt.Sprint(obat.ID)})
				return
			}

			obat.Gambar = filePath
		} else {
			obat.Gambar = existingObat.Gambar
		}

		existingObat.KodeObat = obat.KodeObat
		existingObat.NamaObat = obat.NamaObat
		existingObat.Deskripsi = obat.Deskripsi
		existingObat.Harga = obat.Harga
		existingObat.TipeObatID = obat.TipeObatID

		if err := config.DB.Save(&existingObat).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update Obat: " + fmt.Sprint(obat.ID)})
			return
		}

		obatList = append(obatList, existingObat)
	}

	c.JSON(http.StatusOK, obatList)
}

func DeleteObat(c *gin.Context) {
	id := c.Param("id")
	var obat models.Obat

	if err := config.DB.First(&obat, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Obat not found"})
		return
	}

	if err := config.DB.Model(&obat).Association("Tags").Clear(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete related tags: " + err.Error()})
		return
	}

	if err := config.DB.Where("obat_id = ?", obat.ID).Delete(&models.Stok{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete related stock: " + err.Error()})
		return
	}

	if err := config.DB.Delete(&obat).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete obat: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Obat and related data deleted successfully"})
}

func DeleteBatchObat(c *gin.Context) {
	var ids []uint

	if err := c.ShouldBindJSON(&ids); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	if err := config.DB.Delete(&models.Obat{}, ids).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Obat(s) deleted successfully"})
}
