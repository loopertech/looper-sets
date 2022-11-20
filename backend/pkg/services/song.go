package service

import "github.com/edgedb/edgedb-go"

type Song struct {
	ID edgedb.UUID `edgedb:"id" json:"id"` 
	Title string `edgedb:"title" json:"title"`
	Artist string `edgedb:"artist" json:"artist"`
	Genre string `edgedb:"genre" json:"genre"`
	LoopingType string `edgedb:"looping_type" json:"loopingType"`
	BPM edgedb.OptionalInt16 `edgedb:"bpm" json:"bpm"`
	Key edgedb.OptionalStr `edgedb:"key" json:"key"`
	Layers interface{} `edgedb:"layers" json:"layers"`
	Text edgedb.OptionalStr `edgedb:"text" json:"text"`
	VideoURL edgedb.OptionalStr `edgedb:"video_url" json:"video"`
	SongURL edgedb.OptionalStr `edgedb:"song_url" json:"songUrl"`
	MusicURL edgedb.OptionalStr `edgedb:"music_url" json:"musicUrl"`
	Submitter User `edgedb:"user" json:"user"`
	CreatedAt edgedb.OptionalDateTime `edgedb:"created_at" json:"createdAt"`
	UpdatedAt edgedb.OptionalDateTime `edgedb:"updated_at" json:"updatedAt"`
}