package service

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/edgedb/edgedb-go"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type LoginCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type ChangePasswordCredentials struct {
	Username string `json:"username"`
	CurrentPassword string `json:"currentPassword"`
	NewPassword string `json:"newPassword"`
}

type Token struct {
	AccessToken string `json:"accessToken"`
	ExpiresAt time.Time `json:"expiresAt"`
}

type Claims struct {
	ID string `json:"uuid"`
	jwt.RegisteredClaims
}

type UserHeader struct {
	ID string `header:"User-ID"`
}

func Login(c *gin.Context, db *edgedb.Client) {
	var credentials LoginCredentials
	bindError := c.BindJSON(&credentials)
	if bindError != nil {
		log.Println(bindError)
		c.JSON(http.StatusBadRequest, "Error parsing request body.")
		return
	}

	// Build query arguments
	var dbUser []User
	args := map[string]interface{}{
		"username": credentials.Username,
	}
	query := `SELECT user::User {
		id,
		username,
		password,
	} FILTER .username = <str>$username`

	// Run query
	dbError := db.Query(c, query, &dbUser, args)
	if dbError != nil {
		log.Println(dbError)
		c.JSON(http.StatusBadRequest, "Error getting user.")
		return
	}

	// Comparing the password with the hash
	if compareError := bcrypt.CompareHashAndPassword([]byte(dbUser[0].Password), []byte(credentials.Password)); compareError != nil {
		log.Println(compareError)
		c.JSON(http.StatusUnauthorized, "Invalid username and/or password.")
		return
	}

	
	// Issue JWT
	token, jwtError := generateJWT(dbUser[0].ID)
	if jwtError != nil {
		log.Println("jwt error")
		
		log.Println(jwtError)
		c.JSON(http.StatusBadRequest, "Error generating token.")
		return
	}

	c.JSON(http.StatusOK, token)
}

func ChangePassword(c *gin.Context, db *edgedb.Client) {
	// Bind headers
	var header UserHeader
	headerError := c.BindHeader(&header)
	if headerError != nil {
		c.JSON(http.StatusBadRequest, "Error binding header.")
	}

	// Bind request body
	var credentials ChangePasswordCredentials
	bindError := c.BindJSON(&credentials)
	if bindError != nil {
		log.Println(bindError)
		c.JSON(http.StatusBadRequest, "Error parsing request body.")
		return
	}

	// Build query arguments
	var dbUser []User
	getUserArgs := map[string]interface{}{
		"username": credentials.Username,
	}
	getUserQuery := `SELECT user::User {
		id,
		username,
		password,
	} FILTER .username = <str>$username`

	// Run query
	getUserError := db.Query(c, getUserQuery, &dbUser, getUserArgs)
	if getUserError != nil {
		log.Println(getUserError)
		c.JSON(http.StatusBadRequest, "Error getting user.")
		return
	}

	// Comparing the current password with the hash
	if compareError := bcrypt.CompareHashAndPassword([]byte(dbUser[0].Password), []byte(credentials.CurrentPassword)); compareError != nil {
		log.Println(compareError)
		c.JSON(http.StatusUnauthorized, "Invalid username and/or current password.")
		return
	}

	// Generate a hash with the new password
	hash, passError := bcrypt.GenerateFromPassword([]byte(credentials.NewPassword), bcrypt.DefaultCost)
	if passError != nil {
		log.Println(passError)
		c.JSON(http.StatusBadRequest, "Error with password.")
		return
	}

	// Build args
	now := time.Now()
	changePasswordArgs := map[string]interface{}{
		"uuid": dbUser[0].ID,
		"password": string(hash),
		"now": now,
	}

	// Prepare query
	changePasswordQuery := `UPDATE user::User
		FILTER .id = <uuid>$uuid
		SET {
			password := <str>$password,
			updated_at := <datetime>$now
		}
	`

	// Run query
	changePasswordError := db.Query(c, changePasswordQuery, &dbUser, changePasswordArgs)
	if changePasswordError != nil {
		log.Println(changePasswordError)
		c.JSON(http.StatusBadRequest, "Error setting new password.")
		return
	}

	c.JSON(http.StatusNoContent, "")
}

func GetMe(c *gin.Context, db *edgedb.Client) {
	var header UserHeader
	headerError := c.BindHeader(&header)
	if headerError != nil {
		c.JSON(http.StatusBadRequest, "Error binding header.")
	}
	GetUser(c, header.ID, db)
}


func generateJWT(uuid edgedb.UUID) (Token, error){
	secretKey := []byte(os.Getenv("JWT_SECRET"))

	log.Println(secretKey)

	// expirationTime := time.Now().Add(10 * time.Second)
	expirationTime := time.Now().Add(30 * time.Minute)

	claims := &Claims{
		ID: uuid.String(),
		RegisteredClaims: jwt.RegisteredClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Create the JWT string
	tokenString, jwtError := token.SignedString(secretKey)
	if jwtError != nil {
		return Token{} , jwtError
	}

	return Token{ AccessToken: tokenString, ExpiresAt: expirationTime}, nil
}