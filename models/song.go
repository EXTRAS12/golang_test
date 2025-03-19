package models


type Song struct {
	ID     uint    `json:"id" gorm:"primaryKey" example:"1"`
	Group  string  `json:"group" example:"Queen"`
	Song   string  `json:"song" example:"Bohemian Rhapsody"`
	Lyrics []Lyric `json:"lyrics,omitempty" gorm:"foreignKey:SongID"`
}


type Lyric struct {
	ID     uint   `json:"id" gorm:"primaryKey" example:"1"`
	SongID uint   `json:"song_id" example:"1"`
	Text   string `json:"text" example:"Is this the real life? Is this just fantasy?"`
	Order  int    `json:"order" example:"1"`
}


type SongFilter struct {
	Group    string `form:"group" example:"Queen"`
	Song     string `form:"song" example:"Bohemian Rhapsody"`
	Page     int    `form:"page,default=1" example:"1"`
	PageSize int    `form:"page_size,default=10" example:"10"`
}


type LyricFilter struct {
	SongID   uint `form:"song_id" example:"1"`
	Page     int  `form:"page,default=1" example:"1"`
	PageSize int  `form:"page_size,default=5" example:"5"`
}
