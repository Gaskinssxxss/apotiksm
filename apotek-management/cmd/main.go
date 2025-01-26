package main

import (
	"apotek-management/config"
	"apotek-management/routes"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	router.Static("/uploads", "./uploads")

	routes.SetupRoutes(router)

	return router
}

func main() {
	config.ConnectDB()
	r := setupRouter()
	if err := r.Run(":3000"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}

// package main

// import (
// 	"apotek-management/config"
// 	"apotek-management/routes"
// 	"log"

// 	"github.com/gin-gonic/gin"
// )

// func setupRouter() *gin.Engine {
// 	router := gin.Default()
// 	routes.SetupRoutes(router)

// 	return router
// }

// func main() {
// 	config.ConnectDB()
// 	r := setupRouter()
// 	if err := r.Run(":3000"); err != nil {
// 		log.Fatalf("Failed to run server: %v", err)
// 	}
// }
