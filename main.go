package main

import (
	"log"
	"os"

	"songapp/config"
	"songapp/handlers"
	"songapp/models"
	"songapp/services"

	_ "songapp/docs"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	config.ConnectDB()

	if err := config.DB.AutoMigrate(&models.Song{}, &models.Lyric{}); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	songService := services.NewSongService("http://localhost:8081/api") // api для getinfo
	songHandler := handlers.NewSongHandler(config.DB, songService)

	router := gin.Default()

	api := router.Group("/api")
	{

		api.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

		songs := api.Group("/songs")
		{
			songs.GET("", songHandler.GetSongs)
			songs.POST("", songHandler.CreateSong)
			songs.PUT("/:id", songHandler.UpdateSong)
			songs.DELETE("/:id", songHandler.DeleteSong)
			songs.GET("/:id/lyrics", songHandler.GetSongLyrics)
			songs.GET("/info", songHandler.GetSongInfo)
		}
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
