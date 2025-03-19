package models

type Song struct {
	ID     uint    `json:"id" gorm:"primaryKey"`
	Group  string  `json:"group"`
	Song   string  `json:"song"`
	Lyrics []Lyric `json:"lyrics,omitempty" gorm:"foreignKey:SongID"`
}

type Lyric struct {
	ID     uint   `json:"id" gorm:"primaryKey"`
	SongID uint   `json:"song_id"`
	Text   string `json:"text"`
	Order  int    `json:"order"`
}

type SongFilter struct {
	Group    string `form:"group"`
	Song     string `form:"song"`
	Page     int    `form:"page,default=1"`
	PageSize int    `form:"page_size,default=10"`
}

type LyricFilter struct {
	SongID   uint `form:"song_id"`
	Page     int  `form:"page,default=1"`
	PageSize int  `form:"page_size,default=5"`
}
