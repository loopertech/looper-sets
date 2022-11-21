package service

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/edgedb/edgedb-go"
	"github.com/gin-gonic/gin"
)

type Song struct {
	ID edgedb.UUID `edgedb:"id" json:"id"` 
	Title string `edgedb:"title" json:"title"`
	Artist string `edgedb:"artist" json:"artist"`
	Genre string `edgedb:"genre" json:"genre"`
	LoopingType string `edgedb:"looping_type" json:"loopingType"`
	BPM edgedb.OptionalInt16 `edgedb:"bpm" json:"bpm"`
	Key edgedb.OptionalStr `edgedb:"key" json:"key"`
	Layers interface{} `edgedb:"layers" json:"layers"`
	Text interface{} `edgedb:"text" json:"text"`
	VideoURL edgedb.OptionalStr `edgedb:"video_url" json:"video"`
	SongURL edgedb.OptionalStr `edgedb:"song_url" json:"songUrl"`
	MusicURL edgedb.OptionalStr `edgedb:"music_url" json:"musicUrl"`
	Submitter User `edgedb:"submitter" json:"submitter"`
	CreatedAt edgedb.OptionalDateTime `edgedb:"created_at" json:"createdAt"`
	UpdatedAt edgedb.OptionalDateTime `edgedb:"updated_at" json:"updatedAt"`
}

func CreateSong(c *gin.Context, db *edgedb.Client) {
	var header UserHeader
	headerError := c.BindHeader(&header)
	if headerError != nil {
		c.JSON(http.StatusBadRequest, "Error binding header.")
	}
	// Parse Header ID
	headerID, headerParseError := edgedb.ParseUUID(header.ID)
	if headerParseError != nil {
		log.Println(headerParseError)
		c.JSON(http.StatusBadRequest, "Error parsing UUID.")
		return
	}


	// Parse body
	var songBody Song
	bindError := c.BindJSON(&songBody)
	if bindError != nil {
		log.Println(bindError)
		c.JSON(http.StatusBadRequest, "Error parsing request body.")
		return
	}

	// TODO: add genre check

	layers, layersError := json.Marshal(songBody.Layers)
	if layersError != nil {
		log.Println(layersError)
		c.JSON(http.StatusBadRequest, "Error marshalling layers.")
		return
	}

	text, textError := json.Marshal(songBody.Text)
	if textError != nil {
		log.Println(textError)
		c.JSON(http.StatusBadRequest, "Error marshalling text.")
		return
	}
	
	// Build query arguments
	now := time.Now()
	args := map[string]interface{}{
		"title": songBody.Title,
		"artist": songBody.Artist,
		"genre": songBody.Genre,
		"type": songBody.LoopingType,
		"bpm": songBody.BPM,
		"key": songBody.Key,
		"layers": layers,
		"text": text,
		"video_url": songBody.VideoURL,
		"song_url": songBody.SongURL,
		"music_url": songBody.MusicURL,
		"submitter_uuid": headerID,
		"now": now,
	}
	
	// Prepare query
	query := `WITH song := (
		INSERT song::Song {
			title :=  <str>$title,
			artist :=  <str>$artist,
			genre :=  <str>$genre,
			looping_type := <song::LoopingType>$type,
			bpm :=  <optional int16>$bpm,
			key :=  <optional str>$key,
		  layers :=  <json>$layers,
			text :=  <json>$text,
			video_url :=  <optional str>$video_url,
			song_url :=  <optional str>$song_url,
			music_url :=  <optional str>$music_url,
			submitter :=  (SELECT user::User FILTER .id = <uuid>$submitter_uuid LIMIT 1),
			created_at := <datetime>$now,
		}
	) SELECT song {
		id,
		title,
		artist,
		genre,
		looping_type,
		bpm,
		key,
		layers,
		text,
		video_url,
		song_url,
		music_url,
		submitter,
		created_at,
		updated_at
	}`
	
	// Run query
	var song []Song
	queryError := db.Query(c, query, &song, args)
	if queryError != nil {
		log.Println(queryError)
		c.JSON(http.StatusBadRequest, "Error running query.")
		return
	}

	c.JSON(http.StatusCreated, song)
}

// // Get all users
// func GetUsers(c *gin.Context, db *edgedb.Client) {
// 	var dbUsers []User

// 	// Prepare query
// 	query := `SELECT user::User {
// 		id,
// 		email,
// 		username,
// 		display_name,
// 		website,
// 		social_media,
// 		submitted_songs,
// 		saved_songs,
// 		set_lists,
// 		created_at,
// 		updated_at
// 	}`

// 	// Run query
// 	err := db.Query(c, query, &dbUsers)
// 	if err != nil {
// 		log.Println(err)
// 		c.JSON(http.StatusBadRequest, "Error getting users.")
// 		return
// 	}

// 	// Prepare JSON reponse body
// 	var users []UserJson
// 	for _, user := range dbUsers {
// 		users = append(users, UserJson{
// 			ID: user.ID,
// 			Email: user.Email,
// 			Username: user.Username,
// 			DisplayName: user.DisplayName,
// 			Website: user.Website,
// 			SocialMedia: user.SocialMedia,
// 			SubmittedSongs: user.SubmittedSongs,
// 			SavedSongs: user.SavedSongs,
// 			SetLists: user.SetLists,
// 			CreatedAt: user.CreatedAt,
// 			UpdatedAt: user.UpdatedAt,
// 		})
// 	}

// 	c.JSON(http.StatusOK, users)
// }

// // Get a user
// func GetUser(c *gin.Context, id string, db *edgedb.Client) {
// 	var dbUser []User

// 	// Parse UUID
// 	if id == "" {
// 		id = c.Param("uuid")
// 	}
// 	uuid, parseError := edgedb.ParseUUID(id)
// 	if parseError != nil {
// 		log.Println(parseError)
// 		c.JSON(http.StatusBadRequest, "Error parsing UUID.")
// 		return
// 	}

// 	// Build query arguments
// 	args := map[string]interface{}{
// 		"uuid": uuid,
// 	}

// 	query := `SELECT user::User {
// 		id,
// 		email,
// 		username,
// 		display_name,
// 		website,
// 		social_media,
// 		submitted_songs,
// 		saved_songs,
// 		set_lists,
// 		created_at,
// 		updated_at
// 	} FILTER .id = <uuid>$uuid`

// 	// Run query
// 	dbError := db.Query(c, query, &dbUser, args)
// 	if dbError != nil {
// 		log.Println(dbError)
// 		c.JSON(http.StatusBadRequest, "Error getting user.")
// 		return
// 	}

// 	// Check if users found
// 	// If not, return 404
// 	if len(dbUser) == 0 {
// 		c.JSON(http.StatusNotFound, "User not found.")
// 		return
// 	} else {
// 		// Prepare JSON response body
// 		user := UserJson{
// 			ID: dbUser[0].ID,
// 			Email: dbUser[0].Email,
// 			Username: dbUser[0].Username,
// 			DisplayName: dbUser[0].DisplayName,
// 			Website: dbUser[0].Website,
// 			SocialMedia: dbUser[0].SocialMedia,
// 			SubmittedSongs: dbUser[0].SubmittedSongs,
// 			SavedSongs: dbUser[0].SavedSongs,
// 			SetLists: dbUser[0].SetLists,
// 			CreatedAt: dbUser[0].CreatedAt,
// 			UpdatedAt: dbUser[0].UpdatedAt,
// 		}
// 		c.JSON(http.StatusOK, user)
// 		return
// 	}

// }

// // Edit a user
// func EditUser(c *gin.Context, db *edgedb.Client) {
// 	var header UserHeader
// 	headerError := c.BindHeader(&header)
// 	if headerError != nil {
// 		c.JSON(http.StatusBadRequest, "Error binding header.")
// 	}
// 	// Parse Header ID
// 	headerID, headerParseError := edgedb.ParseUUID(header.ID)
// 	if headerParseError != nil {
// 		log.Println(headerParseError)
// 		c.JSON(http.StatusBadRequest, "Error parsing UUID.")
// 		return
// 	}

// 	// Parse UUID
// 	uuid, parseError := edgedb.ParseUUID(c.Param("uuid"))
// 	if parseError != nil {
// 		log.Println(parseError)
// 		c.JSON(http.StatusBadRequest, "Error parsing UUID.")
// 		return
// 	}

// 	// Check if header ID matches user ID
// 	if headerID != uuid {
// 		c.Status(http.StatusUnauthorized)
// 		return
// 	}
	
// 	// Parse request body
// 	var userBody User
// 	bindError := c.BindJSON(&userBody)
// 	if bindError != nil {
// 		log.Println(bindError)
// 		c.JSON(http.StatusBadRequest, "Error parsing request body.")
// 		return
// 	}
	
// 	// Build query arguments
// 	var dbUser []User
// 	now := time.Now()
// 	args := map[string]interface{}{
// 		"uuid": uuid,
// 		"email": userBody.Email,
// 		"username": userBody.Username,
// 		"display_name": userBody.DisplayName,
// 		"website": userBody.Website,
// 		"social_media": userBody.SocialMedia,
// 		"now": now,
// 	}

// 	// Prepare first part of query
// 	query := `WITH user := (
// 		UPDATE user::User
// 		FILTER .id = <uuid>$uuid
// 		SET { ` 
	
// 	// Check which fields are being edited and append them to the query
// 	if userBody.Email != "" {
// 		query += `email := <str>$email, `
// 	}
// 	if userBody.Username != "" {
// 		query += `username := <str>$username, `
// 	}
// 	if userBody.DisplayName != "" {
// 		query += `display_name := <str>$display_name, `
// 	}
// 	website, _ := userBody.Website.Get()
// 	if website != "" {
// 		query += `website := <optional str>$website, `
// 	}
// 	social, _ := userBody.SocialMedia.Get()
// 	if social != "" {
// 		query += `social_media := <optional str>$social_media, `
// 	}
	
// 	// Complete query string
// 	query += `updated_at := <datetime>$now })
// 		SELECT user {
// 			id,
// 			email,
// 			username,
// 			display_name,
// 			website,
// 			social_media,
// 			created_at,
// 			updated_at
// 		}
// 	`

// 	// Run query
// 	dbError := db.Query(c, query, &dbUser, args)
// 	if dbError != nil {
// 		log.Println(dbError)
// 		c.JSON(http.StatusBadRequest, "Error editing user.")
// 		return
// 	}

// 	// Prepare JSON response body
// 	user := UserJson{
// 		ID: dbUser[0].ID,
// 		Email: dbUser[0].Email,
// 		Username: dbUser[0].Username,
// 		DisplayName: dbUser[0].DisplayName,
// 		Website: dbUser[0].Website,
// 		SocialMedia: dbUser[0].SocialMedia,
// 		SubmittedSongs: dbUser[0].SubmittedSongs,
// 		SavedSongs: dbUser[0].SavedSongs,
// 		SetLists: dbUser[0].SetLists,
// 		CreatedAt: dbUser[0].CreatedAt,
// 		UpdatedAt: dbUser[0].UpdatedAt,
// 	}
// 	c.JSON(http.StatusOK, user)
// }

// // Delete a user
// func DeleteUser(c *gin.Context, db *edgedb.Client) {
// 	var header UserHeader
// 	headerError := c.BindHeader(&header)
// 	if headerError != nil {
// 		c.JSON(http.StatusBadRequest, "Error binding header.")
// 	}
// 	// Parse Header ID
// 	headerID, headerParseError := edgedb.ParseUUID(header.ID)
// 	if headerParseError != nil {
// 		log.Println(headerParseError)
// 		c.JSON(http.StatusBadRequest, "Error parsing UUID.")
// 		return
// 	}

// 	// Parse UUID
// 	uuid, parseError := edgedb.ParseUUID(c.Param("uuid"))
// 	if parseError != nil {
// 		log.Println(parseError)
// 		c.JSON(http.StatusBadRequest, "Error parsing UUID.")
// 		return
// 	}

// 	// Check if header ID matches user ID
// 	if headerID != uuid {
// 		c.Status(http.StatusUnauthorized)
// 		return
// 	}
	
// 	// Build query arguments
// 	args := map[string]interface{}{
// 		"uuid": uuid,
// 	}

// 	// Prepare query
// 	query := `DELETE user::User FILTER .id = <uuid>$uuid`
	
// 	// Run query
// 	var dbUser []User
// 	dbError := db.Query(c, query, &dbUser, args)
// 	if dbError != nil {
// 		log.Println(dbError)
// 		c.JSON(http.StatusBadRequest, "Error deleting user.")
// 		return
// 	}

// 	// Check if users found
// 	// If not, return 404
// 	// Else return no content
// 	if len(dbUser) == 0 {
// 		c.JSON(http.StatusNotFound, "User not found.")
// 		return
// 	} else {
// 		c.JSON(http.StatusNoContent, nil)
// 		return
// 	}
// }