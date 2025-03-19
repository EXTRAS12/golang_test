package models


type SongDetail struct {
	ReleaseDate string `json:"releaseDate" example:"1975-10-31"`
	Text        string `json:"text" example:"Is this the real life? Is this just fantasy?"`
	Link        string `json:"link" example:"https://www.youtube.com/watch?v=fJ9rUzIMcZQ"`
}
