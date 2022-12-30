package models

type CreatePost struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Medias      []Media  `json:"medias"`
	Reviews     []Review `json:"reviews"`
}
type Media struct {
	Name string `json:"name"`
	Link string `json:"link"`
	Type string `json:"type"`
}

type Review struct {
	Name        string `json:"name"`
	Rating      int    `json:"rating"`
	Description string `json:"description"`
}
type PostUpdate struct {
	Id          string        `json:"id"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Medias      []MediaUpdate `json:"medias"`
}

type MediaUpdate struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Link string `json:"link"`
	Type string `json:"type"`
}
