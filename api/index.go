// package main

package handler

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v5/stdlib"
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

type NewSubModule struct {
	ProjectId     int       `json:"projectId"`
	SubModuleName string    `json:"subModuleName"`
	Description   string    `json:"description"`
	StartDate     time.Time `json:"startDate"`
	TargetDate    time.Time `json:"targetDate"`
	CreatedBy     int       `json:"createdBy"`
	PicId         int       `json:"picId"`
	PriorityId    int       `json:"priorityId"`
}

type AlterSubModule struct {
	SubModuleId   int        `json:"subModuleId"`
	SubModuleName *string    `json:"subModuleName"`
	Description   *string    `json:"description"`
	StartDate     *time.Time `json:"startDate"`
	TargetDate    *time.Time `json:"targetDate"`
	PicId         *int       `json:"picId"`
	PriorityId    *int       `json:"priorityId"`
}

type NewWork struct {
	SubModuleId    int       `json:"subModuleId"`
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

type NewBug struct {
	WorkName       string    `json:"workName"`
	Description    string    `json:"description"`
	StartDate      time.Time `json:"startDate"`
	TargetDate     time.Time `json:"targetDate"`
	PicId          *int      `json:"picId"`
	CurrentState   int       `json:"currentState"`
	CreatedBy      int       `json:"createdBy"`
	PriorityId     int       `json:"priorityId"`
	EstimatedHours int       `json:"estimatedHours"`
	UsersAdded     []int     `json:"usersAdded"`
	WorkAffected   int       `json:"workAffected"`
	DefectCause    int       `json:"defectCause"`
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
type AlterBug struct {
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
	WorkAffected   *int       `json:"workAffected"`
	DefectCause    *int       `json:"defectCause"`
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
	// if err := godotenv.Load(); err != nil {
	// 	log.Println("Error loading .env file")
	// }
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
	router.GET("/getAllProjects", getAllProjects)
	router.GET("/getUserProjects", getUserProjects)
	router.PUT("/putAlterProject", putAlterProject)
	router.DELETE("/dropProject", dropProject)
	router.GET("/getGanttDataOfProject", getGanttDataOfProject)

	// User Project Roles
	router.GET("/getUserProjectRoles", getUserProjectRoles)
	router.PUT("/putUserProjectRole", putUserProjectRole)

	// Module
	router.GET("/getModulesOfProject", getModulesOfProject)
	router.GET("/getModuleDetails", getModuleDetails)

	// subModule
	router.GET("/getProjectSubModules", getProjectSubModules)
	router.POST("/postNewSubModule", postNewSubModule)
	router.PUT("/putAlterSubModule", putAlterSubModule)
	router.DELETE("/dropSubModule", dropSubModule)

	// Work
	router.POST("/postNewWork", postNewWork)
	router.GET("/getSubModuleWorks", getSubModuleWorks)
	router.GET("/getWorkDetails", getWorkDetails)
	router.PUT("/putAlterWork", putAlterWork)
	router.DELETE("/dropWork", dropWork)
	router.GET("/getUserTodoList", getUserTodoList)
	router.GET("/getWorkNameListOfProjectDev", getWorkNameListOfProjectDev)

	// Bug
	router.POST("/postNewBug", postNewBug)
	router.GET("/getProjectBugs", getProjectBugs)
	router.PUT("/putAlterBug", putAlterBug)
	router.GET("/getBugDetails", getBugDetails)

	// User Work Assignment
	router.GET("/getUserWorkAssignment", getUserWorkAssignment)
	router.PUT("/putAlterUserWorkAssignment", putAlterUserWorkAssignment)

	// router.DELETE("/removeUserProjectRole", removeUserProjectRole)

	// Other data
	router.GET("/getUsernames", getUsernames)
	router.GET("/getProjectAssignedUsernames", getProjectAssignedUsernames)
	router.GET("/getStartBundle", getTrackerActivityPriorityStateList)
	router.GET("/getProjectAndWorkNames", getProjectAndWorkNames)
	router.GET("/getDefectCauseList", getDefectCauseList)
}

// Handler is the entry point for Vercel Serverless Functions.
func Handler(w http.ResponseWriter, r *http.Request) {
	app.ServeHTTP(w, r)
}

// main is the entry point for local development. It is ignored by Vercel.
// func main() {
// 	port := "9090"
// 	log.Printf("INFO: Starting local server on http://localhost:%s\n", port)
// 	http.ListenAndServe(":"+port, http.HandlerFunc(Handler))
// }

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

func getProjectAndWorkNames(c *gin.Context) {
	var data string
	userIdInput := c.Query("userId")
	if checkEmpty(c, userIdInput) {
		return
	}

	query := `SELECT project_manager.get_project_and_work_names($1)`
	if err := db.QueryRow(query, userIdInput).Scan(&data); err != nil {
		checkErr(c, http.StatusBadRequest, err, "Failed to get project and work names")
		return
	}
	// Return the raw JSON data from the database directly to the client.
	c.Data(http.StatusOK, "application/json", []byte(data))
}

func getWorkNameListOfProjectDev(c *gin.Context) {
	var data string
	projectIdInput := c.Query("projectId")
	if checkEmpty(c, projectIdInput) {
		return
	}

	query := `SELECT project_manager.get_work_name_list_of_project_dev($1)`
	if err := db.QueryRow(query, projectIdInput).Scan(&data); err != nil {
		checkErr(c, http.StatusBadRequest, err, "Failed to get work name list of project")
		return
	}
	// Return the raw JSON data from the database directly to the client.
	c.Data(http.StatusOK, "application/json", []byte(data))
}

func getModulesOfProject(c *gin.Context) {
	var data string
	projectIdInput := c.Query("projectId")
	if checkEmpty(c, projectIdInput) {
		return
	}

	query := `SELECT project_manager.get_modules_of_project($1)`
	if err := db.QueryRow(query, projectIdInput).Scan(&data); err != nil {
		checkErr(c, http.StatusBadRequest, err, "Failed to get modules of project")
		return
	}
	// Return the raw JSON data from the database directly to the client.
	c.Data(http.StatusOK, "application/json", []byte(data))
}

func getModuleDetails(c *gin.Context) {
	var data string
	moduleIdInput := c.Query("moduleId")
	if checkEmpty(c, moduleIdInput) {
		return
	}

	query := `SELECT project_manager.get_module_details($1)`
	if err := db.QueryRow(query, moduleIdInput).Scan(&data); err != nil {
		checkErr(c, http.StatusBadRequest, err, "Failed to get module details")
		return
	}
	// Return the raw JSON data from the database directly to the client.
	c.Data(http.StatusOK, "application/json", []byte(data))
}

func getAllProjects(c *gin.Context) {
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

func getUserProjects(c *gin.Context) {
	var data string
	userIdInput := c.Query("userId")
	if checkEmpty(c, userIdInput) {
		return
	}

	// Call the function to get the projects data
	query := `SELECT project_manager.get_projects($1)`
	if err := db.QueryRow(query, userIdInput).Scan(&data); err != nil {
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

func dropProject(c *gin.Context) {
	var projectIdInput = c.Query("projectId")
	if checkEmpty(c, projectIdInput) {
		return
	}
	query := `CALL project_manager.drop_project($1)`
	if _, err := db.Exec(query, projectIdInput); err != nil {
		checkErr(c, http.StatusBadRequest, err, "Failed to drop project")
		return
	}
	c.IndentedJSON(http.StatusOK, "Project dropped successfully")
}

func getGanttDataOfProject(c *gin.Context) {
	var data string
	var projectIdInput = c.Query("projectId")
	if checkEmpty(c, projectIdInput) {
		return
	}

	// Call the function to get the projects data
	query := `SELECT project_manager.get_gantt_data_of_project($1)`
	if err := db.QueryRow(query, projectIdInput).Scan(&data); err != nil {
		checkErr(c, http.StatusBadRequest, err, "Failed to get gantt data")
		return
	}
	// Return the raw JSON data from the database directly to the client.
	c.Data(http.StatusOK, "application/json", []byte(data))
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

func getProjectSubModules(c *gin.Context) {
	var data string
	projectIdInput := c.Query("projectId")
	if checkEmpty(c, projectIdInput) {
		return
	}
	query := `SELECT project_manager.get_project_sub_modules($1)`
	if err := db.QueryRow(query, projectIdInput).Scan(&data); err != nil {
		checkErr(c, http.StatusBadRequest, err, "Failed to get project sub-modules")
		return
	}
	// Return the raw JSON data from the database directly to the client.
	c.Data(http.StatusOK, "application/json", []byte(data))
}

func postNewSubModule(c *gin.Context) {
	var nb NewSubModule
	if err := c.BindJSON(&nb); err != nil {
		checkErr(c, http.StatusBadRequest, err, "Invalid input")
		return
	}

	query := `CALL project_manager.post_new_sub_module($1,$2,$3,$4,$5,$6,$7,$8)`
	if _, err := db.Exec(query,
		nb.ProjectId,
		nb.SubModuleName,
		nb.Description,
		nb.StartDate,
		nb.TargetDate,
		nb.CreatedBy,
		nb.PicId,
		nb.PriorityId,
	); err != nil {
		checkErr(c, http.StatusBadRequest, err, "Failed to create sub-module")
		return
	}

	c.IndentedJSON(http.StatusOK, "Sub-module created successfully")
}

func putAlterSubModule(c *gin.Context) {

	var alterTarget AlterSubModule
	if err := c.BindJSON(&alterTarget); err != nil {
		checkErr(c, http.StatusBadRequest, err, "Invalid input")
		return
	}

	query := `CALL project_manager.put_alter_sub_module($1, $2, $3, $4, $5, $6, $7)`
	if _, err := db.Exec(query,
		alterTarget.SubModuleId,
		alterTarget.SubModuleName,
		alterTarget.Description,
		alterTarget.StartDate,
		alterTarget.TargetDate,
		alterTarget.PicId,
		alterTarget.PriorityId,
	); err != nil {
		checkErr(c, http.StatusBadRequest, err, "Failed to update subModule")
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": "subModule updated successfully"})
}

func dropSubModule(c *gin.Context) {
	var subModuleIdInput = c.Query("subModuleId")
	if checkEmpty(c, subModuleIdInput) {
		return
	}
	query := `CALL project_manager.drop_sub_module($1)`
	if _, err := db.Exec(query, subModuleIdInput); err != nil {
		checkErr(c, http.StatusBadRequest, err, "Failed to drop subModule")
		return
	}

	c.IndentedJSON(http.StatusOK, "subModule dropped successfully")
}

func getSubModuleWorks(c *gin.Context) {
	var data string
	subModuleIdInput := c.Query("subModuleId")
	if checkEmpty(c, subModuleIdInput) {
		return
	}
	query := `SELECT project_manager.get_sub_module_works($1)`
	if err := db.QueryRow(query, subModuleIdInput).Scan(&data); err != nil {
		checkErr(c, http.StatusBadRequest, err, "Failed to get sub-module works")
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

	var newWorkId int
	if err := db.QueryRow(
		`SELECT project_manager.post_new_work($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13)`,
		nw.WorkName,
		nw.PriorityId,
		nw.PicId,
		nw.Description,
		nw.CurrentState,
		nw.CreatedBy,
		nw.TargetDate,
		nw.StartDate,
		nw.UsersAdded,
		nw.EstimatedHours,
		nw.SubModuleId,
		nw.TrackerId,
		nw.ActivityId,
	).Scan(&newWorkId); err != nil {
		checkErr(c, http.StatusBadRequest, err, "Failed to create work")
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"message": "Work created successfully", "workId": newWorkId})
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

func dropWork(c *gin.Context) {
	var workIdInput = c.Query("workId")
	if checkEmpty(c, workIdInput) {
		return
	}
	query := `CALL project_manager.drop_work($1)`
	if _, err := db.Exec(query, workIdInput); err != nil {
		checkErr(c, http.StatusBadRequest, err, "Failed to drop work")
		return
	}
	c.IndentedJSON(http.StatusOK, "Work dropped successfully")
}

func getWorkDetails(c *gin.Context) {
	var data string
	workIdInput := c.Query("workId")
	if checkEmpty(c, workIdInput) {
		return
	}

	query := `SELECT project_manager.get_work_details($1)`
	if err := db.QueryRow(query, workIdInput).Scan(&data); err != nil {
		checkErr(c, http.StatusBadRequest, err, "Failed to get work details")
		return
	}
	// Return the raw JSON data from the database directly to the client.
	c.Data(http.StatusOK, "application/json", []byte(data))
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

func getProjectBugs(c *gin.Context) {
	var data string
	projectIdInput := c.Query("projectId")
	if checkEmpty(c, projectIdInput) {
		return
	}
	query := `SELECT project_manager.get_project_bugs($1)`
	if err := db.QueryRow(query, projectIdInput).Scan(&data); err != nil {
		checkErr(c, http.StatusBadRequest, err, "Failed to get bug list")
		return
	}
	// Return the raw JSON data from the database directly to the client.
	c.Data(http.StatusOK, "application/json", []byte(data))
}

func postNewBug(c *gin.Context) {
	var nb NewBug
	if err := c.BindJSON(&nb); err != nil {
		checkErr(c, http.StatusBadRequest, err, "Invalid input")
		return
	}
	query := `CALL project_manager.post_new_bug($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)`
	if _, err := db.Exec(
		query,
		nb.WorkName,
		nb.PriorityId,
		nb.PicId,
		nb.Description,
		nb.CurrentState,
		nb.CreatedBy,
		nb.TargetDate,
		nb.StartDate,
		nb.UsersAdded,
		nb.EstimatedHours,
		nb.DefectCause,
		nb.WorkAffected,
	); err != nil {
		checkErr(c, http.StatusBadRequest, err, "Failed to create bug")
		return
	}
	c.IndentedJSON(http.StatusOK, "Bug created successfully")
}

func putAlterBug(c *gin.Context) {
	var alterTarget AlterBug

	if err := c.BindJSON(&alterTarget); err != nil {
		checkErr(c, http.StatusBadRequest, err, "Invalid input format")
		return
	}

	query := `CALL project_manager.put_alter_bug($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)`
	log.Printf("%+v\n", alterTarget)
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
		alterTarget.DefectCause,
		alterTarget.WorkAffected,
		alterTarget.UsersRemoved,
		alterTarget.UsersAdded,
	); err != nil {
		checkErr(c, http.StatusInternalServerError, err, "Failed to alter bug details")
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": "Successfully altered bug"})
}

func getBugDetails(c *gin.Context) {
	var data string
	bugIdInput := c.Query("bugId")
	if checkEmpty(c, bugIdInput) {
		return
	}

	query := `SELECT project_manager.get_bug_details($1)`
	if err := db.QueryRow(query, bugIdInput).Scan(&data); err != nil {
		checkErr(c, http.StatusBadRequest, err, "Failed to get bug details")
		return
	}
	// Return the raw JSON data from the database directly to the client.
	c.Data(http.StatusOK, "application/json", []byte(data))
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

func getDefectCauseList(c *gin.Context) {
	var data string
	query := `SELECT project_manager.get_defect_cause_list()`
	if err := db.QueryRow(query).Scan(&data); err != nil {
		checkErr(c, http.StatusBadRequest, err, "Failed to get start data")
		return
	}
	// Return the raw JSON data from the database directly to the client.
	c.Data(http.StatusOK, "application/json", []byte(data))
}
