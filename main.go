package main

import (
  "fmt"
  "net/http"
  _ "github.com/go-sql-driver/mysql"
  "github.com/jmoiron/sqlx"
  
  "github.com/gin-gonic/gin"
  "github.com/gin-gonic/contrib/static"
  "github.com/manyminds/api2go"
  "github.com/manyminds/api2go-adapter/gingonic"
  "github.com/timrourke/po/model"
  "github.com/timrourke/po/resource"
  "github.com/timrourke/po/storage"
  "log"
)

var db *sqlx.DB

func init() {
  var err error
  db, err = sqlx.Open("mysql", "golang:golang@tcp(127.0.0.1:3306)/golang?parseTime=true")
  if err != nil {
    log.Panic(err)
  }
  err = db.Ping()
  if err != nil {
    log.Panic(err)
  } else {
    log.Println("Successfully connected to mysql database golang.")
  }
}

func main() {
  r := gin.Default()

  api := api2go.NewAPIWithRouting(
    "api",
    api2go.NewStaticResolver("/"),
    gingonic.New(r),
  )

  nounStorage := storage.NewNounStorage(db)

  api.AddResource(model.Noun{}, resource.NounResource{NounStorage: nounStorage})

  r.GET("/ping", func(c *gin.Context) {
    c.String(200, "pong")
  })

  r.GET("/", func (c *gin.Context) {
    http.Redirect(c.Writer, c.Request, "/app", http.StatusMovedPermanently)
  })
  
  r.Use(static.Serve("/app", static.LocalFile("./frontend", true)))

  r.GET("/app/*wildcard", func (c *gin.Context) {
    path := fmt.Sprintf("%v", c.Request.URL)
    path = "./frontend/index.html"
    fmt.Println("trying to load path: ", path)
    http.ServeFile(c.Writer, c.Request, path)
  })
  
  r.Run(":8080")
}