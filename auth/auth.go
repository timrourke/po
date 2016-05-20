package auth

import (
	"database/sql"
	"fmt"
	jwt_lib "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/timrourke/po/database"
	"github.com/timrourke/po/model"
	"golang.org/x/crypto/bcrypt"
	"log"
	"os"
	"time"
)

type Credentials struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Authenticates user and returns signed JWT
func HandleLogin(c *gin.Context) (string, error) {
	var credentials Credentials
	c.Bind(&credentials)

	// Look up the user
	user := model.User{}
	err := database.DB.Get(&user, "SELECT * FROM user WHERE email = ? LIMIT 1", credentials.Email)

	// Fail immediately if any error occurs
	if err != nil {
		log.Println("Error retrieving user:", err)
		c.JSON(401, gin.H{
			"status":  "Not authorized",
			"message": "The email address or password was incorrect.",
		})
		c.Abort()
		return "", err
	}

	// Compare password with bcrypt
	err = bcrypt.CompareHashAndPassword([]byte(*user.Password), []byte(credentials.Password))
	if err != nil {
		log.Println("Error authenticating user:", err)
		c.JSON(401, gin.H{
			"status":  "Not authorized",
			"message": "The email address or password was incorrect.",
		})
		c.Abort()
		return "", err
	}

	// TODO: consider storing JWTs in Redis for token revocation
	// uuid := uuid.NewV4()
	// fmt.Println("UUID: ", uuid)

	// Build the token
	token := jwt_lib.New(jwt_lib.GetSigningMethod("HS256"))
	token.Claims["userid"] = user.GetID()
	token.Claims["exp"] = time.Now().Add(time.Hour*72).Unix() * 1000
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		c.JSON(500, gin.H{
			"status":  "Error",
			"message": "Could not log user in",
		})
		c.Abort()
		return "", err
	}
	return tokenString, nil
}

// Create a new user and log them in
func HandleSignup(c *gin.Context) (string, model.User, error) {
	var credentials Credentials
	c.Bind(&credentials)

	// Make sure user doesn't already exist
	existingUser := model.User{}
	err := database.DB.Get(&existingUser, "SELECT * FROM user WHERE email = ? LIMIT 1", credentials.Email)
	if err != nil && err != sql.ErrNoRows {
		log.Println("Error checking for existing user in signup", err)
		c.JSON(500, gin.H{
			"status":  "Error",
			"message": "User signup failed",
		})
		c.Abort()
		return "", model.User{}, err
	}
	if *existingUser.Email == credentials.Email {
		log.Println("User already exists, aborting signup")
		c.JSON(409, gin.H{
			"status":  "Conflict",
			"message": fmt.Sprintf("A user with the email address %v is already taken", credentials.Email),
		})
		c.Abort()
		return "", model.User{}, fmt.Errorf("A user with the email address %v is already taken", credentials.Email)
	}

	// Bcrypt the password for the new user
	bcryptedPassword, err := bcrypt.GenerateFromPassword([]byte(credentials.Password), 11)
	if err != nil {
		log.Println("Error bcrypting password for new user", err)
		c.JSON(500, gin.H{
			"status":  "Error",
			"message": "User signup failed",
		})
		c.Abort()
		return "", model.User{}, err
	}

	// Create new user
	passwordString := string(bcryptedPassword[:])
	user := model.User{
		Username: &credentials.Username,
		Email:    &credentials.Email,
		Password: &passwordString,
	}
	result, err := database.DB.Exec("INSERT INTO user (username, email, password) VALUES (?, ?, ?)",
		user.Username,
		user.Email,
		user.Password)
	if err != nil {
		log.Println("Error saving new user", err)
		c.JSON(500, gin.H{
			"status":  "Error",
			"message": "User signup failed",
		})
		c.Abort()
		return "", model.User{}, err
	}
	insertId, err := result.LastInsertId()
	if err != nil {
		log.Println("Error saving new user and retrieving insert ID", err)
		c.JSON(500, gin.H{
			"status":  "Error",
			"message": "User signup failed",
		})
		c.Abort()
		return "", model.User{}, err
	}
	user.SetID(fmt.Sprintf("%d", insertId))

	// Build the token
	token := jwt_lib.New(jwt_lib.GetSigningMethod("HS256"))
	token.Claims["userid"] = user.GetID()
	token.Claims["exp"] = time.Now().Add(time.Hour*72).Unix() * 1000
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		c.JSON(500, gin.H{
			"status":  "Error",
			"message": "Could not create token for new user",
		})
		c.Abort()
		return "", model.User{}, err
	}

	return tokenString, user, nil
}
