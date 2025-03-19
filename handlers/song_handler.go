package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"songapp/logger"
	"songapp/models"
	"songapp/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SongHandler struct {
	db          *gorm.DB
	songService *services.SongService
}

func NewSongHandler(db *gorm.DB, songService *services.SongService) *SongHandler {
	logger.Info("Initializing SongHandler")
	return &SongHandler{
		db:          db,
		songService: songService,
	}
}

func (h *SongHandler) GetSongs(c *gin.Context) {
	logger.Debug("Handling GetSongs request")
	var filter models.SongFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		logger.LogError(err, "GetSongs")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	query := h.db.Model(&models.Song{})

	if filter.Group != "" {
		query = query.Where("`group` LIKE ?", "%"+filter.Group+"%")
	}
	if filter.Song != "" {
		query = query.Where("song LIKE ?", "%"+filter.Song+"%")
	}

	var total int64
	query.Count(&total)
	logger.Debug("Found %d total songs", total)

	var songs []models.Song
	offset := (filter.Page - 1) * filter.PageSize
	if err := query.Offset(offset).Limit(filter.PageSize).Find(&songs).Error; err != nil {
		logger.LogError(err, "GetSongs")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	logger.Info("Successfully retrieved %d songs", len(songs))
	c.JSON(http.StatusOK, gin.H{
		"data": songs,
		"meta": gin.H{
			"total":     total,
			"page":      filter.Page,
			"page_size": filter.PageSize,
		},
	})
}

func (h *SongHandler) GetSongInfo(c *gin.Context) {
	logger.Debug("Handling GetSongInfo request")
	group := c.Query("group")
	song := c.Query("song")

	if group == "" || song == "" {
		logger.LogError(fmt.Errorf("missing required parameters"), "GetSongInfo")
		c.JSON(http.StatusBadRequest, gin.H{"error": "group and song parameters are required"})
		return
	}

	songDetail, err := h.songService.GetSongInfo(group, song)
	if err != nil {
		logger.LogError(err, "GetSongInfo")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	logger.Info("Successfully retrieved song info for: %s - %s", group, song)
	c.JSON(http.StatusOK, songDetail)
}


func (h *SongHandler) GetSongLyrics(c *gin.Context) {
	logger.Debug("Handling GetSongLyrics request")
	var filter models.LyricFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		logger.LogError(err, "GetSongLyrics")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var total int64
	h.db.Model(&models.Lyric{}).Where("song_id = ?", filter.SongID).Count(&total)
	logger.Debug("Found %d total lyrics for song %d", total, filter.SongID)

	var lyrics []models.Lyric
	offset := (filter.Page - 1) * filter.PageSize
	if err := h.db.Where("song_id = ?", filter.SongID).
		Order("`order`").
		Offset(offset).
		Limit(filter.PageSize).
		Find(&lyrics).Error; err != nil {
		logger.LogError(err, "GetSongLyrics")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	logger.Info("Successfully retrieved %d lyrics for song %d", len(lyrics), filter.SongID)
	c.JSON(http.StatusOK, gin.H{
		"data": lyrics,
		"meta": gin.H{
			"total":     total,
			"page":      filter.Page,
			"page_size": filter.PageSize,
		},
	})
}

func (h *SongHandler) DeleteSong(c *gin.Context) {
	logger.Debug("Handling DeleteSong request")
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		logger.LogError(err, "DeleteSong")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := h.db.Delete(&models.Song{}, id).Error; err != nil {
		logger.LogError(err, "DeleteSong")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	logger.Info("Successfully deleted song with ID: %d", id)
	c.JSON(http.StatusOK, gin.H{"message": "Song deleted successfully"})
}


func (h *SongHandler) UpdateSong(c *gin.Context) {
	logger.Debug("Handling UpdateSong request")
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		logger.LogError(err, "UpdateSong")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var song models.Song
	if err := c.ShouldBindJSON(&song); err != nil {
		logger.LogError(err, "UpdateSong")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.db.Model(&models.Song{}).Where("id = ?", id).Updates(song).Error; err != nil {
		logger.LogError(err, "UpdateSong")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	logger.Info("Successfully updated song with ID: %d", id)
	c.JSON(http.StatusOK, gin.H{"message": "Song updated successfully"})
}


func (h *SongHandler) CreateSong(c *gin.Context) {
	logger.Debug("Handling CreateSong request")
	var song models.Song
	if err := c.ShouldBindJSON(&song); err != nil {
		logger.LogError(err, "CreateSong")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	songDetail, err := h.songService.GetSongInfo(song.Group, song.Song)
	if err != nil {
		logger.LogError(err, "CreateSong")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get song details"})
		return
	}

	if err := h.db.Create(&song).Error; err != nil {
		logger.LogError(err, "CreateSong")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	lyrics := models.Lyric{
		SongID: song.ID,
		Text:   songDetail.Text,
		Order:  1,
	}
	if err := h.db.Create(&lyrics).Error; err != nil {
		logger.LogError(err, "CreateSong")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save song lyrics"})
		return
	}

	logger.Info("Successfully created new song: %s - %s (ID: %d)", song.Group, song.Song, song.ID)
	response := gin.H{
		"id":          song.ID,
		"group":       song.Group,
		"song":        song.Song,
		"releaseDate": songDetail.ReleaseDate,
		"link":        songDetail.Link,
		"lyrics":      songDetail.Text,
	}

	c.JSON(http.StatusCreated, response)
}
