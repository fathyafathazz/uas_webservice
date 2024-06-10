package artist

type Artist struct {
	ArtistId int `json:"id"`
	Name string `json:"name"`
	Genre string `json:"genre"`
}