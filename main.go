package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/gin-gonic/gin"
)

var db *sql.DB

type Task struct {
    ID   int    `json:"id"`
    Text string `json:"text"`
    Done bool   `json:"done"`
}

// var tasks []Task
// var nextID = 1

func main() {    
	// Database connection
    server := "localhost"
	port := 1433
	database := "TodoApp"
    user:="<your user name>"
    password:="<your password>"

    connString := fmt.Sprintf("server=%s;port=%d;database=%s;user id=%s;password=%s;",
        server, port, database, user, password)

	var err error
    db, err = sql.Open("sqlserver", connString)
    if err != nil {
        panic("Failed to connect to database: " + err.Error())
    }
    defer db.Close()

    // Ensure the connection is working
    err = db.Ping()
    if err != nil {
        panic("Database connection failed: " + err.Error())
    }

    fmt.Println("Connected to MS SQL Server!")
	
    r := gin.Default()
    r.LoadHTMLGlob("templates/*")

	// Routes
	// Render home page
	// Step 4: modify for MS SQL
    r.GET("/", func(c *gin.Context) {
        // c.HTML(http.StatusOK, "index.html", gin.H{"tasks": tasks})
		tasks := getTasks()
        c.HTML(http.StatusOK, "index.html", gin.H{"tasks": tasks})
	})

	// Add a new task (Create)
    r.POST("/add", func(c *gin.Context) {
        text := c.PostForm("text")
        if text != "" {
            // tasks = append(tasks, Task{ID: nextID, Text: text, Done: false})
            // nextID++
			_, err := db.Exec("INSERT INTO Tasks (Text, Done) VALUES (@p1, @p2)", text, false)
            if err != nil {
                c.String(http.StatusInternalServerError, "Error adding task: "+err.Error())
                return
            }
        }
        c.Redirect(http.StatusSeeOther, "/")
    })

	// Step 3: implement delete
	// Mark a task as done (Update)
    r.POST("/done/:id", func(c *gin.Context) {
        // id := c.Param("id")
        // for i, task := range tasks {
        //     if id == strconv.Itoa(task.ID) {
        //         tasks[i].Done = true
        //         break
        //     }
        // }
		id, err := strconv.Atoi(c.Param("id"))
        if err == nil {
            _, err := db.Exec("UPDATE Tasks SET Done = 1 WHERE ID = @p1", id)
            if err != nil {
                c.String(http.StatusInternalServerError, "Error marking task as done: "+err.Error())
                return
            }
        }
        c.Redirect(http.StatusSeeOther, "/")
    })

	// Delete a task
	r.POST("/delete/:id", func(c *gin.Context) {
        // id, err := strconv.Atoi(c.Param("id"))
        // if err == nil {
        //     for i, task := range tasks {
        //         if task.ID == id {
        //             tasks = append(tasks[:i], tasks[i+1:]...) // Remove task
        //             break
        //         }
        //     }
        // }
		id, err := strconv.Atoi(c.Param("id"))
        if err == nil {
            _, err := db.Exec("DELETE FROM Tasks WHERE ID = @p1", id)
            if err != nil {
                c.String(http.StatusInternalServerError, "Error deleting task: "+err.Error())
                return
            }
        }
        c.Redirect(http.StatusSeeOther, "/")
    })

    r.Run(":8080") // Start server
}

// Fetch all tasks from the database
func getTasks() []Task {
    rows, err := db.Query("SELECT ID, Text, Done FROM Tasks")
    if err != nil {
        fmt.Println("Error fetching tasks:", err)
        return nil
    }
    defer rows.Close()

    var tasks []Task
    for rows.Next() {
        var t Task
        if err := rows.Scan(&t.ID, &t.Text, &t.Done); err != nil {
            fmt.Println("Error scanning row:", err)
            continue
        }
        tasks = append(tasks, t)
    }
    return tasks
}
