package main

// package handler

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
)

// User represents a user for authentication purposes.
type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// type NewProject struct {
// 	ProjectName string    `json:"projectName"`
// 	Description string    `json:"description"`
// 	CreatorID   int       `json:"creatorId"`
// 	TargetDate  time.Time `json:"targetDate"`
// 	PicId       int       `json:"picId"`
// }

// type AlterUserProjectRole struct {
// 	RoleId    int   `json:"roleId"`
// 	ProjectId int   `json:"projectId"`
// 	UserIds   []int `json:"userIds"`
// }

type UserRoleChange struct {
	RoleId       int   `json:"roleId"`
	ProjectId    int   `json:"projectId"`
	UsersAdded   []int `json:"usersAdded"`
	UsersRemoved []int `json:"usersRemoved"`
}

type NewProject struct {
	ProjectName string           `json:"projectName"`
	Description string           `json:"description"`
	CreatedBy   int              `json:"createdBy"`
	StartDate   time.Time        `json:"startDate"`
	TargetDate  time.Time        `json:"targetDate"`
	PicId       int              `json:"picId"`
	UserRoles   []UserRoleChange `json:"userRoles"`
}

// Global variables for the database connection and the Gin engine.
var (
	db  *sql.DB
	app *gin.Engine
)

// init is a special Go function that runs once when the package is initialized.
// For a Vercel serverless function, this serves as the cold-start entry point.
func init() {
	// Establish the database connection pool.
	if err := godotenv.Load(); err != nil {
		log.Println("Error loading .env file")
	}
	db = openDB()
	// Create a new Gin router with default middleware.
	app = gin.Default()

	// Configure CORS (Cross-Origin Resource Sharing) middleware to allow requests from specified frontend origins.
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:4200"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}
	app.Use(cors.New(config))

	// Group all routes under the "/api" prefix for versioning and organization.
	apiGroup := app.Group("/api")
	// Register all application-specific routes.
	registerRoutes(apiGroup)
}

// registerRoutes defines all the API endpoints for the application.
func registerRoutes(router *gin.RouterGroup) {
	// Authentication
	router.POST("/login", checkUserCredentials)

	// Project
	router.POST("/postNewProject", postNewProject)
	router.GET("/getProjects", getProjects)

	// User Project Roles
	router.GET("/getUserProjectRoles", getUserProjectRoles)
	router.PUT("/putUserProjectRole", putUserProjectRole)
	// router.DELETE("/removeUserProjectRole", removeUserProjectRole)

	// Other data
	router.GET("/getUsernames", getUsernames)
}

// Handler is the entry point for Vercel Serverless Functions.
func Handler(w http.ResponseWriter, r *http.Request) {
	app.ServeHTTP(w, r)
}

// main is the entry point for local development. It is ignored by Vercel.
func main() {
	port := "9090"
	log.Printf("INFO: Starting local server on http://localhost:%s\n", port)
	http.ListenAndServe(":"+port, http.HandlerFunc(Handler))
}

// openDB establishes a connection to the PostgreSQL database.
// It uses the DATABASE_URL environment variable for establishing the connection
func openDB() *sql.DB {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		// Fallback for local development if the environment variable is not set.
		databaseURL = "postgres://postgres:12345678@localhost:5432/gudang_garam?sslmode=disable"
		log.Println("INFO: DATABASE_URL not set, using local fallback.")
	}

	// Open a connection using the pgx driver.
	db, err := sql.Open("pgx", databaseURL)
	if err != nil {
		// If the connection string is invalid, the application cannot run.
		log.Fatalf("FATAL: Error opening database: %v", err)
	}
	// Ping the database to verify that the connection is alive.
	if err = db.Ping(); err != nil {
		// If the database is unreachable, the application cannot run.
		log.Fatalf("FATAL: Error pinging database: %v", err)
	}
	log.Println("INFO: Database connection successful.")
	return db
}

// checkErr is a centralized error handling utility.
// It logs the technical error for debugging and sends a standardized, user-friendly
// JSON error response to the client, preventing further execution.
func checkErr(c *gin.Context, errType int, err error, errMsg string) {
	if err != nil {
		log.Printf("ERROR: %v", err) // Log the detailed error for server-side debugging.
		// Send a JSON response with the appropriate HTTP status code.
		if errType == http.StatusInternalServerError {
			c.JSON(http.StatusInternalServerError, gin.H{"error": errMsg})
		} else if errType == http.StatusBadRequest {
			c.JSON(http.StatusBadRequest, gin.H{"error": errMsg})
		}
		c.Abort() // Stop processing the request.
	}
}

// checkEmpty validates that a required query parameter is not empty.
// This prevents nil pointer errors and ensures handlers receive necessary data.
func checkEmpty(c *gin.Context, str string) {
	if str == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing query parameters"})
		c.Abort() // Stop processing if a required parameter is missing.
	}
}

func checkUserCredentials(c *gin.Context) {
	var newUser User
	var data string
	// Attempt to bind the request body to the User struct.
	if err := c.BindJSON(&newUser); err != nil {
		checkErr(c, http.StatusBadRequest, err, "Invalid input")
		return
	}
	log.Printf("INFO: Login attempt for user: %s", newUser.Username)

	// Call the corresponding database function to authenticate the user.
	query := `SELECT project_manager.get_user_id_by_credentials($1, $2)`
	if err := db.QueryRow(query, newUser.Username, newUser.Password).Scan(&data); err != nil {
		checkErr(c, http.StatusBadRequest, err, "Failed to get user ID")
		return
	}
	// Return the raw JSON data from the database directly to the client.
	c.Data(http.StatusOK, "application/json", []byte(data))
	// c.IndentedJSON(http.StatusOK, "ok")
}

func getUsernames(c *gin.Context) {
	var data string

	query := `SELECT project_manager.get_usernames()`
	if err := db.QueryRow(query).Scan(&data); err != nil {
		checkErr(c, http.StatusBadRequest, err, "Failed to get usernames")
		return
	}
	// Return the raw JSON data from the database directly to the client.
	c.Data(http.StatusOK, "application/json", []byte(data))
}

func getProjects(c *gin.Context) {
	var data string

	// Call the function to get the projects data
	query := `SELECT project_manager.get_projects()`
	if err := db.QueryRow(query).Scan(&data); err != nil {
		checkErr(c, http.StatusBadRequest, err, "Failed to get projects")
		return
	}
	// Return the raw JSON data from the database directly to the client.
	c.Data(http.StatusOK, "application/json", []byte(data))
}

func postNewProject(c *gin.Context) {
	var np NewProject
	if err := c.BindJSON(&np); err != nil {
		checkErr(c, http.StatusBadRequest, err, "Invalid input")
		return
	}

	var projectIdTemp int
	query := `SELECT project_manager.post_new_project($1,$2,$3,$4,$5)`
	if err := db.QueryRow(query, np.ProjectName, np.Description, np.CreatedBy, np.TargetDate, np.PicId).Scan(&projectIdTemp); err != nil {
		checkErr(c, http.StatusBadRequest, err, "Failed to create project")
		return
	}
	log.Printf("INFO: Project created with ID: %d", projectIdTemp)
	for _, userRole := range np.UserRoles {
		if len(userRole.UsersAdded) != 0 && len(userRole.UsersRemoved) == 0 {
			userRole.ProjectId = projectIdTemp
			if err := AlterUserProjectRole(c, userRole); err != nil {
				checkErr(c, http.StatusBadRequest, err, "Project created successfully but Failed to set user project role")
				return
			}
		}
	}

	c.IndentedJSON(http.StatusOK, "Project created successfully")
}

func getUserProjectRoles(c *gin.Context) {
	var data string
	projectIdInput := c.Query("projectId")
	query := `SELECT project_manager.get_user_project_roles($1)`
	if err := db.QueryRow(query, projectIdInput).Scan(&data); err != nil {
		checkErr(c, http.StatusBadRequest, err, "Failed to get user project roles")
		return
	}
	// Return the raw JSON data from the database directly to the client.
	c.Data(http.StatusOK, "application/json", []byte(data))
}

// func removeUserProjectRole(c *gin.Context) {
// 	var alterTarget AlterUserProjectRole
// 	if err := c.BindJSON(&alterTarget); err != nil {
// 		checkErr(c, http.StatusBadRequest, err, "Invalid input")
// 		return
// 	}
// 	log.Println("INFO: Removing user project role for project ID:", alterTarget.ProjectId, "and user IDs:", alterTarget.UserIds)
// 	query := `CALL project_manager.remove_user_project_role($1,$2)`
// 	if _, err := db.Exec(query, alterTarget.ProjectId, alterTarget.UserIds); err != nil {
// 		checkErr(c, http.StatusBadRequest, err, "Failed to remove user project role")
// 		return
// 	}
// 	c.IndentedJSON(http.StatusOK, "User project role removed successfully")
// }

func putUserProjectRole(c *gin.Context) {
	var alterTarget UserRoleChange
	if err := c.BindJSON(&alterTarget); err != nil {
		checkErr(c, http.StatusBadRequest, err, "Invalid input")
		return
	}

	AlterUserProjectRole(c, alterTarget)
}

func AlterUserProjectRole(c *gin.Context, alterTarget UserRoleChange) error {
	query := `CALL project_manager.alter_user_project_role($1,$2,$3, $4)`
	if _, err := db.Exec(query, alterTarget.ProjectId, alterTarget.RoleId, alterTarget.UsersRemoved, alterTarget.UsersAdded); err != nil {
		return err
	}
	return nil
}
