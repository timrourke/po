package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"net/http"
	// "errors"
	// "github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/contrib/jwt"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/manyminds/api2go"
	"github.com/manyminds/api2go-adapter/gingonic"
	"github.com/timrourke/po/auth"
	"github.com/timrourke/po/database"
	"github.com/timrourke/po/model"
	"github.com/timrourke/po/resource"
	"github.com/timrourke/po/storage"
	"log"
	"os"
)

func init() {
	var err error

	// Load the dotenv files for the application
	err = godotenv.Load()
	if err != nil {
		log.Fatal("Failed to load the dotenv file.")
	}

	// Attempt to open a sql connection
	database.DB, err = sqlx.Open("mysql",
		fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/%s?parseTime=true",
			os.Getenv("MYSQL_USER"),
			os.Getenv("MYSQL_PASSWORD"),
			os.Getenv("MYSQL_DB")))
	if err != nil {
		log.Panic(err)
	}
	err = database.DB.Ping()
	if err != nil {
		log.Panic(err)
	} else {
		log.Println("Successfully connected to mysql database golang.")
	}
}

func main() {
	r := gin.Default()

	// Attempt to register sessions using redis
	// store, _ := sessions.NewRedisStore(10, "tcp", "localhost:6379", os.Getenv("REDIS_PASSWORD"))
	// r.Use(sessions.Sessions("session", store))

	api := api2go.NewAPIWithRouting(
		"api",
		api2go.NewStaticResolver("/"),
		gingonic.New(r),
	)

	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	r.POST("/signup", func(c *gin.Context) {
		tokenString, newUser, err := auth.HandleSignup(c)
		if err != nil {
			log.Println("Error while trying to sign user up", err)
			c.Abort() // Just to be safe
			return
		}
		c.JSON(200, gin.H{
			"token": tokenString,
			"user":  newUser,
		})
	})

	r.POST("/login", func(c *gin.Context) {
		tokenString, err := auth.HandleLogin(c)
		if err != nil {
			log.Println("Error while trying to log user in", err)
			c.Abort() // Just to be safe
			return
		}
		c.JSON(200, gin.H{"token": tokenString})
	})

	r.GET("/", func(c *gin.Context) {
		http.Redirect(c.Writer, c.Request, "/app", http.StatusMovedPermanently)
	})

	r.Use(static.Serve("/app", static.LocalFile("./frontend", true)))

	r.GET("/app/*wildcard", func(c *gin.Context) {
		path := fmt.Sprintf("%v", c.Request.URL)
		path = "./frontend/index.html"
		fmt.Println("trying to load path: ", path)
		http.ServeFile(c.Writer, c.Request, path)
	})

	// All routes going forward are protected with JWT authorization
	r.Use(jwt.Auth(os.Getenv("JWT_SECRET")))

	nounStorage := storage.NounStorage{}
	verbStorage := storage.VerbStorage{}
	tensePresIndStorage := storage.TensePresIndStorage{}

	api.AddResource(model.Noun{}, resource.NounResource{NounStorage: nounStorage})
	api.AddResource(model.Verb{}, resource.VerbResource{VerbStorage: verbStorage})
	api.AddResource(model.TensePresentIndicative{}, resource.TensePresIndResource{TensePresIndStorage: tensePresIndStorage})

	r.Run(":8080")
}
