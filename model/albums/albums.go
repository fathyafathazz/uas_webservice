package album

type Album struct {
	AlbumId       int     `json:"id"`
	Title         string  `json:"title"`
	Artist_Id     int  `json:"artist_id"`
	Price         float32 `json:"price"`
	Year int `json:"year"`
}