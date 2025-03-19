package handlers

import (
	"net/http"
	"strconv"

	"songapp/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SongHandler struct {
	db *gorm.DB
}

func NewSongHandler(db *gorm.DB) *SongHandler {
	return &SongHandler{db: db}
}

// GetSongs получает список песен с фильтрацией и пагинацией
func (h *SongHandler) GetSongs(c *gin.Context) {
	var filter models.SongFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
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

	var songs []models.Song
	offset := (filter.Page - 1) * filter.PageSize
	if err := query.Offset(offset).Limit(filter.PageSize).Find(&songs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": songs,
		"meta": gin.H{
			"total":     total,
			"page":      filter.Page,
			"page_size": filter.PageSize,
		},
	})
}

// GetSongLyrics получает текст песни с пагинацией по куплетам
func (h *SongHandler) GetSongLyrics(c *gin.Context) {
	var filter models.LyricFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var total int64
	h.db.Model(&models.Lyric{}).Where("song_id = ?", filter.SongID).Count(&total)

	var lyrics []models.Lyric
	offset := (filter.Page - 1) * filter.PageSize
	if err := h.db.Where("song_id = ?", filter.SongID).
		Order("`order`").
		Offset(offset).
		Limit(filter.PageSize).
		Find(&lyrics).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": lyrics,
		"meta": gin.H{
			"total":     total,
			"page":      filter.Page,
			"page_size": filter.PageSize,
		},
	})
}

// DeleteSong удаляет песню
func (h *SongHandler) DeleteSong(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := h.db.Delete(&models.Song{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Song deleted successfully"})
}

// UpdateSong обновляет данные песни
func (h *SongHandler) UpdateSong(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var song models.Song
	if err := c.ShouldBindJSON(&song); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.db.Model(&models.Song{}).Where("id = ?", id).Updates(song).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Song updated successfully"})
}

// CreateSong создает новую песню
func (h *SongHandler) CreateSong(c *gin.Context) {
	var song models.Song
	if err := c.ShouldBindJSON(&song); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.db.Create(&song).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, song)
}
