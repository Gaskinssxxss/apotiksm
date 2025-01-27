package routes

import (
	"apotek-management/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	api := router.Group("/api")
	{
		// Tag Obat
		api.GET("/tag_obat", controllers.GetAllTagObat)
		api.GET("/tag_obat/:id", controllers.GetTagObatByID)
		api.POST("/tag_obat", controllers.CreateTagObat)
		api.PUT("/tag_obat/:id", controllers.UpdateTagObat)
		api.DELETE("/tag_obat/:id", controllers.DeleteTagObat)

		// Tipe Obat
		api.GET("/tipe_obat", controllers.GetAllTipeObat)
		api.GET("/tipe_obat/:id", controllers.GetTipeObatByID)
		api.POST("/tipe_obat", controllers.CreateTipeObat)
		api.PUT("/tipe_obat/:id", controllers.UpdateTipeObat)
		api.DELETE("/tipe_obat/:id", controllers.DeleteTipeObat)

		api.POST("/tipe_obat/batch_create", controllers.CreateBatchTipeObat)
		api.PUT("/tipe_obat/batch_update", controllers.UpdateBatchTipeObat)
		api.DELETE("/tipe_obat/batch_delete", controllers.DeleteBatchTipeObat)

		// Stok
		api.GET("/stok", controllers.GetAllStok)
		api.GET("stok/:id", controllers.GetStokByID)
		api.POST("/stok", controllers.CreateStok)
		api.PUT("/stok/:id", controllers.UpdateStok)
		api.DELETE("/stok/:id", controllers.DeleteStok)

		// Transaksi
		api.POST("/transaksi", controllers.CreateTransaksi)
		api.GET("/transaksi", controllers.GetAllTransaksi)
		api.GET("/transaksi/:id", controllers.GetTransaksiByID)
		api.PUT("/transaksi/:id", controllers.UpdateTransaksi)
		api.DELETE("/transaksi/:id", controllers.DeleteTransaksi)

		api.POST("/transaksi/batch_create", controllers.CreateBatchTransaksi)
		api.PUT("/transaksi/batch_update", controllers.UpdateBatchTransaksi)
		api.DELETE("/transaksi/batch_delete", controllers.DeleteBatchTransaksi)

		// Obat
		api.GET("/obat", controllers.GetAllObat)
		api.GET("/obat/:id", controllers.GetObatByID)
		api.POST("/obat", controllers.CreateObat)
		api.POST("/obat/batch_create", controllers.CreateBatchObat)

		api.PUT("/obat/:id", controllers.UpdateObat)
		api.PUT("/obat/batch_update", controllers.UpdateBatchObat)
		api.DELETE("/obat/:id", controllers.DeleteObat)
		api.DELETE("/obat/batch_delete", controllers.DeleteBatchObat)

		//laporan
		api.GET("/laporan/laporan-transaksi", controllers.GetLaporanTransaksi)
		api.GET("/laporan/laporan-stok", controllers.GetLaporanStok)
		api.GET("/laporan/laporan-labarugi", controllers.GetLaporanLabaRugi)
	}

	// Welcome route
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Welcome to the Apotek Management API!",
		})
	})
}
