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

type AlterProject struct {
	ProjectId   *int             `json:"projectId"`
	ProjectName *string          `json:"projectName"`
	Description *string          `json:"description"`
	StartDate   *time.Time       `json:"startDate"`
	TargetDate  *time.Time       `json:"targetDate"`
	PicId       *int             `json:"picId"`
	UserRoles   []UserRoleChange `json:"userRoles"`
}

type NewBacklog struct {
	ProjectId   int       `json:"projectId"`
	BacklogName string    `json:"backlogName"`
	Description string    `json:"description"`
	StartDate   time.Time `json:"startDate"`
	TargetDate  time.Time `json:"targetDate"`
	CreatedBy   int       `json:"createdBy"`
	PicId       int       `json:"picId"`
	PriorityId  int       `json:"priorityId"`
}

type AlterBacklog struct {
	BacklogId   int        `json:"backlogId"`
	BacklogName *string    `json:"backlogName"`
	Description *string    `json:"description"`
	StartDate   *time.Time `json:"startDate"`
	TargetDate  *time.Time `json:"targetDate"`
	PicId       *int       `json:"picId"`
	PriorityId  *int       `json:"priorityId"`
}

type NewWork struct {
	BacklogId      int       `json:"backlogId"`
	WorkName       string    `json:"workName"`
	Description    string    `json:"description"`
	StartDate      time.Time `json:"startDate"`
	TargetDate     time.Time `json:"targetDate"`
	PicId          *int      `json:"picId"`
	CurrentState   int       `json:"currentState"`
	CreatedBy      int       `json:"createdBy"`
	PriorityId     int       `json:"priorityId"`
	EstimatedHours int       `json:"estimatedHours"`
	TrackerId      int       `json:"trackerId"`
	ActivityId     int       `json:"activityId"`
	UsersAdded     []int     `json:"usersAdded"`
}

type AlterWork struct {
	WorkId         int        `json:"workId"`
	WorkName       *string    `json:"workName"`
	Description    *string    `json:"description"`
	StartDate      *time.Time `json:"startDate"`
	TargetDate     *time.Time `json:"targetDate"`
	PicId          *int       `json:"picId"`
	CurrentState   *int       `json:"currentState"`
	PriorityId     *int       `json:"priorityId"`
	EstimatedHours *int       `json:"estimatedHours"`
	TrackerId      *int       `json:"trackerId"`
	ActivityId     *int       `json:"activityId"`
	UsersRemoved   []int      `json:"usersRemoved"`
	UsersAdded     []int      `json:"usersAdded"`
}

type UserWorkChange struct {
	WorkId       int   `json:"workId"`
	UsersAdded   []int `json:"usersAdded"`
	UsersRemoved []int `json:"usersRemoved"`
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
	router.PUT("/putAlterProject", putAlterProject)

	// User Project Roles
	router.GET("/getUserProjectRoles", getUserProjectRoles)
	router.PUT("/putUserProjectRole", putUserProjectRole)

	// Backlog
	router.GET("/getProjectBacklogs", getProjectBacklogs)
	router.POST("/postNewBacklog", postNewBacklog)

	// Work
	router.POST("/postNewWork", postNewWork)
	router.GET("/getBacklogWorks", getBacklogWorks)
	router.PUT("/putAlterWork", putAlterWork)
	router.GET("/getUserTodoList", getUserTodoList)

	// User Work Assignment
	router.GET("/getUserWorkAssignment", getUserWorkAssignment)
	router.PUT("/putAlterUserWorkAssignment", putAlterUserWorkAssignment)

	// router.DELETE("/removeUserProjectRole", removeUserProjectRole)

	// Other data
	router.GET("/getUsernames", getUsernames)
	router.GET("/getProjectAssignedUsernames", getProjectAssignedUsernames)
	router.GET("/getStartBundle", getTrackerActivityPriorityStateList)
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
func checkEmpty(c *gin.Context, str string) bool {
	if str == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing query parameters"})
		c.Abort() // Stop processing if a required parameter is missing.
		return true
	}
	return false
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

func getProjectAssignedUsernames(c *gin.Context) {
	var data string
	projectIdInput := c.Query("projectId")
	if checkEmpty(c, projectIdInput) {
		return
	}

	roleIdInput := c.Query("roleId")
	var query string
	var err error

	if roleIdInput == "" {
		query = `SELECT project_manager.get_project_assigned_usernames($1)`
		err = db.QueryRow(query, projectIdInput).Scan(&data)
	} else {
		query = `SELECT project_manager.get_project_assigned_usernames($1, $2)`
		err = db.QueryRow(query, projectIdInput, roleIdInput).Scan(&data)
	}
	if err != nil {
		checkErr(c, http.StatusBadRequest, err, "Failed to get project usernames")
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

func putAlterProject(c *gin.Context) {
	var ap AlterProject
	if err := c.BindJSON(&ap); err != nil {
		checkErr(c, http.StatusBadRequest, err, "Invalid input")
		return
	}
	query := `CALL project_manager.put_alter_project($1,$2,$3,$4,$5)`
	if _, err := db.Exec(query, ap.ProjectId, ap.ProjectName, ap.Description, ap.TargetDate, ap.PicId); err != nil {
		checkErr(c, http.StatusBadRequest, err, "Failed to update project")
		return
	}

	for _, userRole := range ap.UserRoles {
		if len(userRole.UsersAdded) != 0 && len(userRole.UsersRemoved) == 0 {
			userRole.ProjectId = *ap.ProjectId
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
	if checkEmpty(c, projectIdInput) {
		return
	}
	query := `SELECT project_manager.get_user_project_roles($1)`
	if err := db.QueryRow(query, projectIdInput).Scan(&data); err != nil {
		checkErr(c, http.StatusBadRequest, err, "Failed to get user project roles")
		return
	}
	// Return the raw JSON data from the database directly to the client.
	c.Data(http.StatusOK, "application/json", []byte(data))
}

func putUserProjectRole(c *gin.Context) {
	var alterTarget UserRoleChange
	if err := c.BindJSON(&alterTarget); err != nil {
		checkErr(c, http.StatusBadRequest, err, "Invalid input")
		return
	}

	if err := AlterUserProjectRole(c, alterTarget); err != nil {
		checkErr(c, http.StatusBadRequest, err, "Failed to alter user project role")
		return
	}

	c.IndentedJSON(http.StatusOK, "Succesfully altered user project role")
}

func AlterUserProjectRole(c *gin.Context, alterTarget UserRoleChange) error {
	query := `CALL project_manager.alter_user_project_role($1,$2,$3, $4)`
	if _, err := db.Exec(query, alterTarget.ProjectId, alterTarget.RoleId, alterTarget.UsersRemoved, alterTarget.UsersAdded); err != nil {
		return err
	}
	return nil
}

func getProjectBacklogs(c *gin.Context) {
	var data string
	projectIdInput := c.Query("projectId")
	if checkEmpty(c, projectIdInput) {
		return
	}
	query := `SELECT project_manager.get_project_backlogs($1)`
	if err := db.QueryRow(query, projectIdInput).Scan(&data); err != nil {
		checkErr(c, http.StatusBadRequest, err, "Failed to get project backlogs")
		return
	}
	// Return the raw JSON data from the database directly to the client.
	c.Data(http.StatusOK, "application/json", []byte(data))
}

func postNewBacklog(c *gin.Context) {
	var nb NewBacklog
	if err := c.BindJSON(&nb); err != nil {
		checkErr(c, http.StatusBadRequest, err, "Invalid input")
		return
	}

	query := `CALL project_manager.post_new_backlog($1,$2,$3,$4,$5,$6,$7,$8)`
	if _, err := db.Exec(query,
		nb.ProjectId,
		nb.BacklogName,
		nb.Description,
		nb.StartDate,
		nb.TargetDate,
		nb.CreatedBy,
		nb.PicId,
		nb.PriorityId,
	); err != nil {
		checkErr(c, http.StatusBadRequest, err, "Failed to create backlog")
		return
	}

	c.IndentedJSON(http.StatusOK, "Backlog created successfully")
}

func putAlterBacklog(c *gin.Context) {

	var alterTarget AlterBacklog
	if err := c.BindJSON(&alterTarget); err != nil {
		checkErr(c, http.StatusBadRequest, err, "Invalid input")
		return
	}

	query := `CALL project_manager.put_alter_backlog($1, $2, $3, $4, $5, $6, $7)`
	if _, err := db.Exec(query,
		alterTarget.BacklogId,
		alterTarget.BacklogName,
		alterTarget.Description,
		alterTarget.StartDate,
		alterTarget.TargetDate,
		alterTarget.PicId,
		alterTarget.PriorityId,
	); err != nil {
		checkErr(c, http.StatusBadRequest, err, "Failed to update backlog")
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": "Backlog updated successfully"})
}

func getBacklogWorks(c *gin.Context) {
	var data string
	backlogIdInput := c.Query("backlogId")
	if checkEmpty(c, backlogIdInput) {
		return
	}
	query := `SELECT project_manager.get_backlog_works($1)`
	if err := db.QueryRow(query, backlogIdInput).Scan(&data); err != nil {
		checkErr(c, http.StatusBadRequest, err, "Failed to get backlog works")
		return
	}
	// Return the raw JSON data from the database directly to the client.
	c.Data(http.StatusOK, "application/json", []byte(data))
}

func getUserTodoList(c *gin.Context) {
	var data string
	userIdInput := c.Query("userId")
	if checkEmpty(c, userIdInput) {
		return
	}
	query := `SELECT project_manager.get_user_todo_list($1)`
	if err := db.QueryRow(query, userIdInput).Scan(&data); err != nil {
		checkErr(c, http.StatusBadRequest, err, "Failed to get user todo list")
		return
	}
	// Return the raw JSON data from the database directly to the client.
	c.Data(http.StatusOK, "application/json", []byte(data))
}

func getUserWorkAssignment(c *gin.Context) {
	var data string
	workIdInput := c.Query("workId")
	if checkEmpty(c, workIdInput) {
		return
	}
	query := `SELECT project_manager.get_user_work_assignment($1)`
	if err := db.QueryRow(query, workIdInput).Scan(&data); err != nil {
		checkErr(c, http.StatusBadRequest, err, "Failed to get user work assignment")
		return
	}
	// Return the raw JSON data from the database directly to the client.
	c.Data(http.StatusOK, "application/json", []byte(data))
}

func postNewWork(c *gin.Context) {
	var nw NewWork
	if err := c.BindJSON(&nw); err != nil {
		checkErr(c, http.StatusBadRequest, err, "Invalid input")
		return
	}

	_, err := db.Exec(
		`CALL project_manager.post_new_work($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13)`,
		nw.BacklogId,
		nw.WorkName,
		nw.PriorityId,
		nw.PicId,
		nw.Description,
		nw.CurrentState,
		nw.CreatedBy,
		nw.TargetDate,
		nw.StartDate,
		nw.TrackerId,
		nw.ActivityId,
		nw.UsersAdded,
		nw.EstimatedHours,
	)
	if err != nil {
		checkErr(c, http.StatusBadRequest, err, "Failed to create work")
		return
	}

	c.IndentedJSON(http.StatusOK, "Work created successfully")
}

func putAlterWork(c *gin.Context) {
	var alterTarget AlterWork

	// 1. Bind the incoming JSON to the AlterWork struct.
	if err := c.BindJSON(&alterTarget); err != nil {
		checkErr(c, http.StatusBadRequest, err, "Invalid input format")
		return
	}

	// 2. Define the SQL query to call the stored procedure with all 12 parameters.
	query := `CALL project_manager.put_alter_work($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)`

	if _, err := db.Exec(query,
		alterTarget.WorkId,
		alterTarget.WorkName,
		alterTarget.Description,
		alterTarget.StartDate,
		alterTarget.TargetDate,
		alterTarget.CurrentState,
		alterTarget.PicId,
		alterTarget.PriorityId,
		alterTarget.EstimatedHours,
		alterTarget.TrackerId,
		alterTarget.ActivityId,
		alterTarget.UsersRemoved,
		alterTarget.UsersAdded,
	); err != nil {
		checkErr(c, http.StatusInternalServerError, err, "Failed to alter work details")
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": "Successfully altered work assignment"})
}

func putAlterUserWorkAssignment(c *gin.Context) {
	var alterTarget UserWorkChange
	if err := c.BindJSON(&alterTarget); err != nil {
		checkErr(c, http.StatusBadRequest, err, "Invalid input")
		return
	}
	query := `CALL project_manager.alter_user_work_assignment($1,$2,$3)`
	if _, err := db.Exec(query, alterTarget.WorkId, alterTarget.UsersRemoved, alterTarget.UsersAdded); err != nil {
		checkErr(c, http.StatusBadRequest, err, "Failed to alter user work assignment")
		return
	}
	c.IndentedJSON(http.StatusOK, "Succesfully altered user work assignment")
}

func getTrackerActivityPriorityStateList(c *gin.Context) {
	var data string
	query := `SELECT project_manager.get_tracker_activity_priority_state_list()`
	if err := db.QueryRow(query).Scan(&data); err != nil {
		checkErr(c, http.StatusBadRequest, err, "Failed to get start data")
		return
	}
	// Return the raw JSON data from the database directly to the client.
	c.Data(http.StatusOK, "application/json", []byte(data))
}
