package service

import "github.com/edgedb/edgedb-go"

type SetList struct {
	ID edgedb.UUID `edgedb:"id" json:"id"` 
	Title string `edgedb:"title" json:"title"`
	Genres string `edgedb:"genres" json:"genres"`
	Text edgedb.OptionalStr `edgedb:"text" json:"text"`
	Visibility string `edgedb:"visibility" json:"visibility"`
	CreatedBy User `edgedb:"created_by" json:"createdBy"`
	Songs []Song `edgedb:"songs" json:"songs"`
	SongOrder string `edgedb:"song_order" json:"songOrder"`
	CreatedAt edgedb.OptionalDateTime `edgedb:"created_at" json:"createdAt"`
	UpdatedAt edgedb.OptionalDateTime `edgedb:"updated_at" json:"updatedAt"`
}