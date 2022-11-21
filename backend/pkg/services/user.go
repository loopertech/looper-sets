package service

import (
	"log"
	"net/http"
	"time"

	"github.com/edgedb/edgedb-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID edgedb.UUID `edgedb:"id" json:"id"` 
	Email string `edgedb:"email" json:"email"`
	Password string `edgedb:"password" json:"password"`
	Username string `edgedb:"username" json:"username"`
	DisplayName string `edgedb:"display_name" json:"displayName"`
	Type string `edgedb:"user_type" json:"type"`
	Website edgedb.OptionalStr `edgedb:"website" json:"website"`
	SocialMedia edgedb.OptionalStr `edgedb:"social_media" json:"socialMedia"`
	SubmittedSongs []Song `edgedb:"submitted_songs" json:"submittedSongs"`
	SavedSongs []Song `edgedb:"saved_songs" json:"savedSongs"`
	SetLists []SetList `edgedb:"set_lists" json:"setLists"`
	CreatedAt edgedb.OptionalDateTime `edgedb:"created_at" json:"createdAt"`
	UpdatedAt edgedb.OptionalDateTime `edgedb:"updated_at" json:"updatedAt"`
}

type UserJson struct {
	ID edgedb.UUID `json:"id"` 
	Email string `json:"email"`
	Username string `json:"username"`
	DisplayName string `json:"displayName"`
	Website edgedb.OptionalStr `json:"website"`
	SocialMedia edgedb.OptionalStr `json:"socialMedia"`
	SubmittedSongs []Song `json:"submittedSongs"`
	SavedSongs []Song `json:"savedSongs"`
	SetLists []SetList `json:"setLists"`
	CreatedAt edgedb.OptionalDateTime `json:"createdAt"`
	UpdatedAt edgedb.OptionalDateTime `json:"updatedAt"`
}

// Create a user
func CreateUser(c *gin.Context, db *edgedb.Client) {
	// Parse body
	var userBody User
	bindError := c.BindJSON(&userBody)
	if bindError != nil {
		log.Println(bindError)
		c.JSON(http.StatusBadRequest, "Error parsing request body.")
		return
	}

	hash, passError := bcrypt.GenerateFromPassword([]byte(userBody.Password), bcrypt.DefaultCost)
	if passError != nil {
		log.Println(passError)
		c.JSON(http.StatusBadRequest, "Error with password.")
		return
	}
	
	// Build query arguments
	now := time.Now()
	args := map[string]interface{}{
		"email": userBody.Email,
		"password": string(hash),
		"username": userBody.Username,
		"display_name": userBody.DisplayName,
		"type": "Looper",
		"now": now,
	}
	
	// Prepare query
	query := `WITH user := (
		INSERT user::User {
			email := <str>$email,
			password := <str>$password,
			username := <str>$username,
			display_name := <str>$display_name,
			user_type := <str>$type,
			created_at := <datetime>$now,
		}
	) SELECT user {
		id,
		email,
		username,
		display_name,
		website,
		social_media,
		created_at,
		updated_at
	}`
	
	// Run query
	var user []User
	queryError := db.Query(c, query, &user, args)
	if queryError != nil {
		log.Println(queryError)
		c.JSON(http.StatusBadRequest, "Error running query.")
		return
	}

	// Prepare JSON response body
	createdUser := UserJson{
		ID: user[0].ID,
		Email: user[0].Email,
		Username: user[0].Username,
		DisplayName: user[0].DisplayName,
	}

	c.JSON(http.StatusCreated, createdUser)
}

// Get all users
func GetUsers(c *gin.Context, db *edgedb.Client) {
	var dbUsers []User

	// Prepare query
	query := `SELECT user::User {
		id,
		email,
		username,
		display_name,
		website,
		social_media,
		submitted_songs,
		saved_songs,
		set_lists,
		created_at,
		updated_at
	}`

	// Run query
	err := db.Query(c, query, &dbUsers)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, "Error getting users.")
		return
	}

	// Prepare JSON reponse body
	var users []UserJson
	for _, user := range dbUsers {
		users = append(users, UserJson{
			ID: user.ID,
			Email: user.Email,
			Username: user.Username,
			DisplayName: user.DisplayName,
			Website: user.Website,
			SocialMedia: user.SocialMedia,
			SubmittedSongs: user.SubmittedSongs,
			SavedSongs: user.SavedSongs,
			SetLists: user.SetLists,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		})
	}

	c.JSON(http.StatusOK, users)
}

// Get a user
func GetUser(c *gin.Context, id string, db *edgedb.Client) {
	var dbUser []User

	// Parse UUID
	if id == "" {
		id = c.Param("uuid")
	}
	uuid, parseError := edgedb.ParseUUID(id)
	if parseError != nil {
		log.Println(parseError)
		c.JSON(http.StatusBadRequest, "Error parsing UUID.")
		return
	}

	// Build query arguments
	args := map[string]interface{}{
		"uuid": uuid,
	}

	query := `SELECT user::User {
		id,
		email,
		username,
		display_name,
		website,
		social_media,
		submitted_songs,
		saved_songs,
		set_lists,
		created_at,
		updated_at
	} FILTER .id = <uuid>$uuid`

	// Run query
	dbError := db.Query(c, query, &dbUser, args)
	if dbError != nil {
		log.Println(dbError)
		c.JSON(http.StatusBadRequest, "Error getting user.")
		return
	}

	// Check if users found
	// If not, return 404
	if len(dbUser) == 0 {
		c.JSON(http.StatusNotFound, "User not found.")
		return
	} else {
		// Prepare JSON response body
		user := UserJson{
			ID: dbUser[0].ID,
			Email: dbUser[0].Email,
			Username: dbUser[0].Username,
			DisplayName: dbUser[0].DisplayName,
			Website: dbUser[0].Website,
			SocialMedia: dbUser[0].SocialMedia,
			SubmittedSongs: dbUser[0].SubmittedSongs,
			SavedSongs: dbUser[0].SavedSongs,
			SetLists: dbUser[0].SetLists,
			CreatedAt: dbUser[0].CreatedAt,
			UpdatedAt: dbUser[0].UpdatedAt,
		}
		c.JSON(http.StatusOK, user)
		return
	}

}

// Edit a user
func EditUser(c *gin.Context, db *edgedb.Client) {
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

	// Parse UUID
	uuid, parseError := edgedb.ParseUUID(c.Param("uuid"))
	if parseError != nil {
		log.Println(parseError)
		c.JSON(http.StatusBadRequest, "Error parsing UUID.")
		return
	}

	// Check if header ID matches user ID
	if headerID != uuid {
		c.Status(http.StatusUnauthorized)
		return
	}
	
	// Parse request body
	var userBody User
	bindError := c.BindJSON(&userBody)
	if bindError != nil {
		log.Println(bindError)
		c.JSON(http.StatusBadRequest, "Error parsing request body.")
		return
	}
	
	// Build query arguments
	var dbUser []User
	now := time.Now()
	args := map[string]interface{}{
		"uuid": uuid,
		"email": userBody.Email,
		"username": userBody.Username,
		"display_name": userBody.DisplayName,
		"website": userBody.Website,
		"social_media": userBody.SocialMedia,
		"now": now,
	}

	// Prepare first part of query
	query := `WITH user := (
		UPDATE user::User
		FILTER .id = <uuid>$uuid
		SET { ` 
	
	// Check which fields are being edited and append them to the query
	if userBody.Email != "" {
		query += `email := <str>$email, `
	}
	if userBody.Username != "" {
		query += `username := <str>$username, `
	}
	if userBody.DisplayName != "" {
		query += `display_name := <str>$display_name, `
	}
	website, _ := userBody.Website.Get()
	if website != "" {
		query += `website := <optional str>$website, `
	}
	social, _ := userBody.SocialMedia.Get()
	if social != "" {
		query += `social_media := <optional str>$social_media, `
	}
	
	// Complete query string
	query += `updated_at := <datetime>$now })
		SELECT user {
			id,
			email,
			username,
			display_name,
			website,
			social_media,
			created_at,
			updated_at
		}
	`

	// Run query
	dbError := db.Query(c, query, &dbUser, args)
	if dbError != nil {
		log.Println(dbError)
		c.JSON(http.StatusBadRequest, "Error editing user.")
		return
	}

	// Prepare JSON response body
	user := UserJson{
		ID: dbUser[0].ID,
		Email: dbUser[0].Email,
		Username: dbUser[0].Username,
		DisplayName: dbUser[0].DisplayName,
		Website: dbUser[0].Website,
		SocialMedia: dbUser[0].SocialMedia,
		SubmittedSongs: dbUser[0].SubmittedSongs,
		SavedSongs: dbUser[0].SavedSongs,
		SetLists: dbUser[0].SetLists,
		CreatedAt: dbUser[0].CreatedAt,
		UpdatedAt: dbUser[0].UpdatedAt,
	}
	c.JSON(http.StatusOK, user)
}

// Delete a user
func DeleteUser(c *gin.Context, db *edgedb.Client) {
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

	// Parse UUID
	uuid, parseError := edgedb.ParseUUID(c.Param("uuid"))
	if parseError != nil {
		log.Println(parseError)
		c.JSON(http.StatusBadRequest, "Error parsing UUID.")
		return
	}

	// Check if header ID matches user ID
	if headerID != uuid {
		c.Status(http.StatusUnauthorized)
		return
	}
	
	// Build query arguments
	args := map[string]interface{}{
		"uuid": uuid,
	}

	// Prepare query
	query := `DELETE user::User FILTER .id = <uuid>$uuid`
	
	// Run query
	var dbUser []User
	dbError := db.Query(c, query, &dbUser, args)
	if dbError != nil {
		log.Println(dbError)
		c.JSON(http.StatusBadRequest, "Error deleting user.")
		return
	}

	// Check if users found
	// If not, return 404
	// Else return no content
	if len(dbUser) == 0 {
		c.JSON(http.StatusNotFound, "User not found.")
		return
	} else {
		c.JSON(http.StatusNoContent, nil)
		return
	}
}