# Go To-Do App ~~with MS SQL~~ - SQL Backends with and without ORM

This README documents my journey learning Go and building a simple **To-Do** application using **Gin (Go web framework)** and **Microsoft SQL Server**. This guide covers setup, installation, improvements, and how to run the project.

A progression of working with different SQL backends and techniques in Go:

| Branch         | Description                            |
|----------------|----------------------------------------|
| `main`         | Plain Go + MSSQL (no ORM)              |
| `gorm-mssql`   | Same MSSQL backend, using GORM ORM     |
| `gorm-mysql`   | GORM ORM with MySQL backend            |

---

## 1️⃣ Environment Setup

### **Install Go**
Download and install Go from [golang.org](https://go.dev/dl/).

```sh
# Verify installation
go version
```

### **Install Visual Studio Code (VS Code)**
Download and install [VS Code](https://code.visualstudio.com/). Install the **Go extension** from the Extensions Marketplace.

### Hello World
- As with learning any programming language, hello world to ensure my environment is set up correctly.
- I first create a new folder named `sample-go-app` and in the root directory, a file named `main.go`.
- Open a terminal, **Terminal > New Terminal**, then run the command `go mod init sample-go-app` to initialize the sample Go app.
- Copied and pasted the following into `main.go`
```go
package main

import "fmt"

func main() {
    fmt.Println("Hello, World!")
}
```
- In the terminal, I ran `go run .`.
- I can get a list of other go commands by calling `go help`.
- I proceed to follow the tutorial on [`go.dev/doc/tutorial/getting-started`](https://go.dev/doc/tutorial/getting-started) to call code in an external package and modified the code to the following:
```go
package main

import (
	"fmt"

	"rsc.io/quote"
)

func main() {
    fmt.Println(quote.Go())
}
```

### **Install Microsoft SQL Server (Local)**
- I was exploring Microsoft SQL Server a couple weeks back (at the time of writing this), that is why I had Microsoft SQL Server installed on my local machine.
- If you're following this guide, you can use with whatever storage services.
- Download and install **Microsoft SQL Server** (Express Edition) from [Microsoft](https://www.microsoft.com/en-us/sql-server/sql-server-downloads).
- Use **SQL Server Management Studio (SSMS)** for database management.
- Ensure SQL Server is running on `localhost`, port `1433`.
- To create a **New Login** using **SSMS GUI**
    - Right-click **Logins** > **New Login...** under **Security**
    - Enter:
        - Login name
        - Authentication type (SQL Server Authentication)
        - Password
        - Default database (optional)
    - In **User Mapping** tab:
        - Check the databases the user should access (TodoApp)
        - Assign apropriate roles (e.g. `db_datareader`, `db_datawriter`, etc.)
    - Test the connection:
    ```sh
    sqlcmd -S localhost -d TodoApp -U your_username -P your_password
    ```

## 2️⃣ Project Setup

### **Initialize the Go Project**
```sh
mkdir go-todo-app
cd go-todo-app

# Initialize Go module
go mod init go-todo-app
```

### **Install Dependencies**
```sh
# Install the Gin web framework
go get -u github.com/gin-gonic/gin

# Install MS SQL Server driver
go get -u github.com/denisenkom/go-mssqldb
```

## 3️⃣ Database Setup

### **Create the Database and Table**
Run the following SQL script in SSMS:

```sql
CREATE DATABASE TodoApp;
GO

USE TodoApp;
GO

CREATE TABLE Tasks (
    ID INT IDENTITY(1,1) PRIMARY KEY,
    Text NVARCHAR(255) NOT NULL,
    Done BIT NOT NULL DEFAULT 0
);
```

## 4️⃣ Implementing the To-Do App
- If you want to follow through, as below. Otherwise, skip to the final version that implements a To-Do App with MS SQL, [here](#final).

### Initial `main.go`
```go
package main

import (
    "net/http"
    "github.com/gin-gonic/gin"
)

type Task struct {
    ID   int    `json:"id"`
    Text string `json:"text"`
    Done bool   `json:"done"`
}

var tasks []Task
var nextID = 1

func main() {
    r := gin.Default()

    r.LoadHTMLGlob("templates/*")

    r.GET("/", func(c *gin.Context) {
        c.HTML(http.StatusOK, "index.html", gin.H{"tasks": tasks})
    })

    r.POST("/add", func(c *gin.Context) {
        text := c.PostForm("text")
        if text != "" {
            tasks = append(tasks, Task{ID: nextID, Text: text, Done: false})
            nextID++
        }
        c.Redirect(http.StatusSeeOther, "/")
    })

    r.POST("/done/:id", func(c *gin.Context) {
        id := c.Param("id")
        for i, task := range tasks {
            if id == strconv.Itoa(task.ID) {
                tasks[i].Done = true
                break
            }
        }
        c.Redirect(http.StatusSeeOther, "/")
    })

    r.Run(":8080") // Start server
}
```

### Initial `templates/index.html` file:
```html
<!DOCTYPE html>
<html>
<head>
    <title>Go To-Do App</title>
</head>
<body>
    <h1>To-Do List</h1>
    <form action="/add" method="POST">
        <input type="text" name="text" required>
        <button type="submit">Add Task</button>
    </form>
    <ul>
        {{range .tasks}}
            <li>
                {{if .Done}}✅{{else}}❌{{end}}
                {{.Text}}
                <form action="/done/{{.ID}}" method="POST" style="display:inline;">
                    <button type="submit">Mark Done</button>
                </form>
            </li>
        {{end}}
    </ul>
</body>
</html>
```

### Added bootstrap
```html
<!DOCTYPE html>
<!-- Step 2: Implement bootstrap -->
<html>
    <head>
        <meta charset="UTF-8">
        <meta name="viewport" content="width=device-width, initial-scale=1">
        <title>Go To-Do App</title>
        <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.0/dist/css/bootstrap.min.css" rel="stylesheet">
    </head>

    <body class="bg-light">
        <div class="container mt-5">
            <h1 class="text-center">Go To-Do List</h1>

            <form action="/add" method="POST" class="input-group my-4">
                <input type="text" name="text" class="form-control" placeholder="Enter a task..." required>
                <button type="submit" class="btn btn-primary">Add Task</button>
            </form>

            <ul class="list-group">
                {{range .tasks}}
                <li class="list-group-item d-flex justify-content-between align-items-center">
                    <span class="{{if .Done}}text-decoration-line-through text-muted{{end}}">
                        {{.Text}}
                    </span>
                    <div>
                        <form action="/done/{{.ID}}" method="POST" style="display:inline;">
                            <button type="submit" class="btn btn-sm btn-success">✔</button>
                        </form>
                        <form action="/delete/{{.ID}}" method="POST" style="display:inline;">
                            <button type="submit" class="btn btn-sm btn-danger">✖</button>
                        </form>
                    </div>
                </li>
                {{end}}
            </ul>
        </div>

        <script src="https://cdn.jsdelivr.net/npm/bootstrap@5.3.0/dist/js/bootstrap.bundle.min.js"></script>

    </body>
</html>
```

### <a id="final" />**Evolved `main.go`**
The following is the final, cleaned version. My `main.go` comes with comments and steps I took, showing progression.

```go
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
    ID   int
    Text string
    Done bool
}

func main() {
    // Database connection
    server := "localhost"
    database := "TodoApp"
    user:="sa"
    password:="YourPassword"

    connString := fmt.Sprintf("server=%s;port=%d;database=%s;user id=%s;password=%s;",
        server, port, database, user, password)

    var err error
    db, err = sql.Open("sqlserver", connString)
    if err != nil {
        panic("Failed to connect to database: " + err.Error())
    }
    defer db.Close()

    err = db.Ping()
    if err != nil {
        panic("Database connection failed: " + err.Error())
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
            _, err := db.Exec("INSERT INTO Tasks (Text, Done) VALUES (@p1, @p2)", text, false)
            if err != nil {
                c.String(http.StatusInternalServerError, "Error adding task: "+err.Error())
                return
            }
        }
        c.Redirect(http.StatusSeeOther, "/")
    })

	// Mark a task as done (Update)
    r.POST("/done/:id", func(c *gin.Context) {
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
```

## 5️⃣ Running the Application

### **Start the Server**
```sh
go run main.go
```

### **Access the App**
Open a browser and visit:
```
http://localhost:8080
```

## 6️⃣ Clean Up Resources
- I drop my `TodoApp` database using:
```sql
USE master;
ALTER DATABASE TodoApp SET SINGLE_USER WITH ROLLBACK IMMEDIATE;
DROP DATABASE TodoApp;
```
- Uploaded this to GitHub and deleted the local copy. (Well, you can skip these.)

## 7️⃣ Future Improvements
- **Dockerize the application** to simplify deployment.
- **Use connection pooling** for better SQL performance.
- **Convert to a REST API** for frontend frameworks like React or Vue.
- **Add authentication** to restrict access.

## 8️⃣ Summary
This project helped me learn Go while building a CRUD-based To-Do app with SQL Server. The journey included:
✅ Setting up the environment  
✅ Connecting Go to MS SQL Server  
✅ Creating RESTful endpoints  
✅ Rendering tasks in an HTML frontend  
✅ Implementing database CRUD operations

## 9️⃣ Finding out more
1. Understanding `func getTasks() []Task`
This function signature means:
- `getTasks` is a function in Go.
- It returns a slice of `Task` objects (`[]Task`).
- `tasks` is a **slice** (similar to an array, but dynamic).
- `[]Task` is the return type (like `List<Task>` in C#).
2. Is `getTasks` **Private** or **Public**?
Yes, `getTasks` is private in Go because its name starts with a lowercase letter.
- **Go's Visibility Rules**
    - **Functions, variables, and structs starting with lowercase** are **private** (only accessible within the same package).
    - **Functions starting with uppercase** are **public** (accessible outside the package).
3. Does Go Support `async` and `await` Like C#?
No, Go **does not** have `async` and `await` like C#.
- Instead, Go uses **goroutines** (lightweight threads) and channels for concurrency.
4. What's that underscore for in `_ "github.com/denisenkom/go-mssqldb"`?
The underscore (_) tells Go to register the package for side effects (required for database drivers).
5. Can't seem to connect to MS SQL.
- Check if SQL Server is Running
Run this in Command Prompt (CMD):
```sh
sqlcmd -S .\SQLEXPRESS -Q "SELECT name FROM sys.databases"
```
- Try SQL Server Authentication instead:
    - Open **SQL Server Management Studio (SSMS)**.
    - Expand **Security** → **Logins**.
    - Right-click **New Logins..** → Set login name → Set **SQL Server Authentication** → Set password.
- Ensure SQL Server allows SQL Authentication:
    - Open SSMS.
    - Right-click your server → Properties → Security.
    - Check SQL Server and Windows Authentication Mode.
    - Restart SQL Server.
- Check if SQL Server is listening on port 1433
```sh
netstat -an | findstr 1433
```
- If SQL Server is listening, you should see output like:
```
TCP    0.0.0.0:1433      0.0.0.0:0      LISTENING
TCP    [::]:1433         [::]:0         LISTENING
```
- Enable Port 1433 in SQL Server Configuration Manager
    - Open** SQL Server Configuration Manager**.
    - Go to **SQL Server Network Configuration** → Click **Protocols for SQLEXPRESS**.
    - Open **TCP/IP**, then scroll down to **IPAll**:
    - Set **TCP Port** to `1433`.
    - Clear ((TCP Dynamic Ports)) (leave it blank).
    - Click **OK**, then restart the SQL Server service.
6. In C# console app, I can use environment variables for user and password, how do I do that with Go?
- Set Environment Variables
    - Windows (Command Prompt)
    ```sh
    set MSSQL_USER=myusername
    set MSSQL_PASSWORD=mypassword
    ```
    - Windows (PowerShell)
    ```sh
    $env:MSSQL_USER="myusername"
    $env:MSSQL_PASSWORD="mypassword"
    ```
    - Linux/macOS (Bash)
    ```sh
    export MSSQL_USER=myusername
    export MSSQL_PASSWORD=mypassword
    ```
- Use Environment Variables in Go
    - Modify connection string:
        ```go
        user := os.Getenv("MSSQL_USER")
        password := os.Getenv("MSSQL_PASSWORD")
        ```
- This way, your username and password are not hardcoded in your source code, making your application more secure and flexible.
- We'll use `go get github.com/joho/godotenv`, a popular package that loads `.env` files into your environment.
```sh
go get github.com/joho/godotenv
```
- Don't forget to call it in your `main.go`:
```go
// Load .env file
err := godotenv.Load()
if err != nil {
    log.Fatal("Error loading .env file")
}
```

## 🔟 References
- [Configure Go with Visual Studio Code](https://learn.microsoft.com/en-us/azure/developer/go/configure-visual-studio-code) 
- [Go Documentation](https://go.dev/doc/)
- [Gin Web Framework](https://github.com/gin-gonic/gin)
- [Microsoft SQL Server](https://www.microsoft.com/en-us/sql-server/)
- [GoDotEnv](https://pkg.go.dev/github.com/joho/godotenv)
- [GORM Documentation](https://gorm.io/)
