package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"songapp/logger"
	"songapp/models"
	"time"
)

type SongService struct {
	baseURL string
}

func NewSongService(baseURL string) *SongService {
	logger.Info("Initializing SongService with baseURL: %s", baseURL)
	return &SongService{
		baseURL: baseURL,
	}
}

func (s *SongService) GetSongInfo(group, songName string) (*models.SongDetail, error) {
	logger.Debug("Getting song info for group: %s, song: %s", group, songName)

	params := url.Values{}
	params.Add("group", group)
	params.Add("song", songName)

	requestURL := fmt.Sprintf("%s/songs/info?%s", s.baseURL, params.Encode())
	logger.Debug("Making request to: %s", requestURL)

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(requestURL)
	if err != nil {
		logger.LogError(err, "GetSongInfo")
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logger.LogError(fmt.Errorf("unexpected status code: %d", resp.StatusCode), "GetSongInfo")
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var songDetail models.SongDetail
	if err := json.NewDecoder(resp.Body).Decode(&songDetail); err != nil {
		logger.LogError(err, "GetSongInfo")
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if songDetail.Text == "" {
		logger.LogError(fmt.Errorf("empty song text"), "GetSongInfo")
		return nil, fmt.Errorf("empty song text")
	}

	logger.Debug("Successfully retrieved song info for: %s - %s", group, songName)
	return &songDetail, nil
}
