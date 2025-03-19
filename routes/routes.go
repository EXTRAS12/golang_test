package routes

import (
	"songapp/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, songHandler *handlers.SongHandler) {
	api := r.Group("/api")
	{
		songs := api.Group("/songs")
		{
			songs.GET("", songHandler.GetSongs)
			songs.POST("", songHandler.CreateSong)
			songs.PUT("/:id", songHandler.UpdateSong)
			songs.DELETE("/:id", songHandler.DeleteSong)
			songs.GET("/:id/lyrics", songHandler.GetSongLyrics)
		}
	}
}
