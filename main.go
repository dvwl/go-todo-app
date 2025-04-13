package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

var db *gorm.DB

type Task struct {
    ID   int    `gorm:"primaryKey" json:"id"`
    Text string `json:"text"`
    Done bool   `json:"done"`
}

func main() {
    // Load .env file
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }

    // Database connection
    server := os.Getenv("MSSQL_SERVER")
    port := 1433
    database := os.Getenv("MSSQL_DB")
    user:= os.Getenv("MSSQL_USER")
    password:= os.Getenv("MSSQL_PASSWORD")

    connString := fmt.Sprintf("server=%s;port=%d;database=%s;user id=%s;password=%s;",
        server, port, database, user, password)
    
    db, err = gorm.Open(sqlserver.Open(connString), &gorm.Config{})
    if err != nil {
        log.Fatal("Failed to connect to database: ", err)
    }

    appEnv := os.Getenv("APP_ENV")

    if appEnv == "development" {
        err = db.AutoMigrate(&Task{})
        if err != nil {
            log.Fatalf("AutoMigrate failed: %v", err)
        }
    }

    fmt.Println("Connected to MS SQL Server!")
    
    r := gin.Default()
    r.LoadHTMLGlob("templates/*")

    // Routes
    // Render home page
    r.GET("/", func(c *gin.Context) {
        tasks := getTasks()
        c.HTML(http.StatusOK, "index.html", gin.H{"tasks": tasks})
    })

    // Add a new task (Create)
    r.POST("/add", func(c *gin.Context) {
        text := c.PostForm("text")
        if text != "" {
            task := Task{Text: text, Done: false}
            if result := db.Create(&task); result.Error != nil {
                c.String(http.StatusInternalServerError, "Error adding task: " + result.Error.Error())
                return
            }
        }
        c.Redirect(http.StatusSeeOther, "/")
    })

    // Mark a task as done (Update)
    r.POST("/done/:id", func(c *gin.Context) {
        id, err := strconv.Atoi(c.Param("id"))
        if err == nil {
            result := db.Model(&Task{}).Where("id = ?", id).Update("done", true)
            if result.Error != nil {
                c.String(http.StatusInternalServerError, "Error marking task as done: " + result.Error.Error())
                return
            }
        }
        c.Redirect(http.StatusSeeOther, "/")
    })

    // Delete a task
    r.POST("/delete/:id", func(c *gin.Context) {
        id, err := strconv.Atoi(c.Param("id"))
        if err == nil {
            result := db.Delete(&Task{}, id)
            if result.Error != nil {
                c.String(http.StatusInternalServerError, "Error deleting task: " + result.Error.Error())
                return
            }
        }
        c.Redirect(http.StatusSeeOther, "/")
    })

    r.Run(":8080") // Start server
}

// Fetch all tasks from the database
func getTasks() []Task {
    var tasks []Task
    result := db.Find(&tasks)
    if result.Error != nil {
        fmt.Println("Error fetching tasks:", result.Error)
    }
    return tasks
}
