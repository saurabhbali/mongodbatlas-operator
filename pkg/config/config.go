package config

import (
	"net/http"
	"os"
	"strconv"
	"time"

	dac "github.com/akshaykarle/go-http-digest-auth-client"
	//ma "github.com/akshaykarle/go-mongodbatlas/mongodbatlas"
        "github.com/dghubble/sling"
)

const apiURL = "https://cloud.mongodb.com/api/atlas/v1.0/"

// APIError represents a MongDB Atlas API Error response
// https://docs.atlas.mongodb.com/api/#errors
type APIError struct {
	Detail    string `json:"detail"`
	Code      int    `json:"error"`
	ErrorCode string `json:"errorCode"`
	Reason    string `json:"reason"`
}

/////////
//func (e APIError) Error() string {
//	if e == (APIError{}) {
//		return ""
//	}
//	return fmt.Sprintf("MongoDB Atlas: %d %v", e.Code, e.Detail)
//}

// relevantError returns any non-nil http-related error (creating the request,
// getting the response, decoding) if any. If the decoded apiError is non-nil
// the apiError is returned. Otherwise, no errors occurred, returns nil.
func relevantError(httpError error, apiError APIError) error {
	if httpError != nil {
		return httpError
	}
	if apiError == (APIError{}) {
		return nil
	}
	return apiError
}

// Project represents a projecting connection information in MongoDB.
//type Project struct {
//	ID           string `json:"id,omitempty"`
//	Name         string `json:"name,omitempty"`
//	OrgID        string `json:"orgId,omitempty"`
//	Created      string `json:"created,omitempty"`
//	ClusterCount int    `json:"clusterCount,omitempty"`
//}

// ProjectService provides methods for accessing MongoDB Atlas Projects API endpoints.
type ProjectService struct {
	sling *sling.Sling
}

// newProjectService returns a new ProjectService.
func newProjectService(sling *sling.Sling) *ProjectService {
	return &ProjectService{
		sling: sling.Path("groups/"),
	}
}

// GetByName information about the project associated to group name
// https://docs.atlas.mongodb.com/reference/api/project-get-one-by-name/
//func (c *ProjectService) GetByName(name string) (*Project, *http.Response, error) {
//	project := new(Project)
//	apiError := new(APIError)
//	path := fmt.Sprintf("byName/%s", name)
//	resp, err := c.sling.New().Get(path).Receive(project, apiError)
//	return project, resp, relevantError(err, *apiError)
//}

// DatabaseUserService provides methods for accessing MongoDB Atlas DatabaseUsers API endpoints.
type DatabaseUserService struct {
	sling *sling.Sling
}

// newDatabaseUserService returns a new DatabaseUserService.
func newDatabaseUserService(sling *sling.Sling) *DatabaseUserService {
	return &DatabaseUserService{
		sling: sling.Path("groups/"),
	}
}

// Role allows the user to perform particular actions on the specified database.
// A role on the admin database can include privileges that apply to the other databases as well.
type Role struct {
	DatabaseName   string `json:"databaseName,omitempty"`
	CollectionName string `json:"collectionName,omitempty"`
	RoleName       string `json:"roleName,omitempty"`
}

type Scope struct {
        Name   string `json:"name,omitempty"`
        Type   string `json:"type,omitempty"`
}

// DatabaseUser represents MongoDB users in your cluster.
type DatabaseUser struct {
	GroupID         string `json:"groupId,omitempty"`
	Username        string `json:"username,omitempty"`
	Password        string `json:"password,omitempty"`
	DatabaseName    string `json:"databaseName,omitempty"`
	DeleteAfterDate string `json:"deleteAfterDate,omitempty"`
	Roles           []Role `json:"roles,omitempty"`
	Scopes          []Scope `json:"scopes,omitempty"`
}

// Client is a MongoDB Atlas client for making MongoDB API requests.
type Client struct {
	sling               *sling.Sling
	Projects            *ProjectService
	DatabaseUsers       *DatabaseUserService
}

// NewClient returns a new Client.
func NewClient(httpClient *http.Client) *Client {
	base := sling.New().Client(httpClient).Base(apiURL)

	return &Client{
		sling:               base,
		Projects:            newProjectService(base.New()),
		DatabaseUsers:       newDatabaseUserService(base.New()),
	}
}

// AtlasConfig stores Programmatic API Keys for authentication to Atlas API
type AtlasConfig struct {
	AtlasPublicKey  string
	AtlasPrivateKey string
}

// NewMongoDBAtlasClient returns a REST API client for MongoDB Atlas
//func (c *AtlasConfig) newMongoDBAtlasClient() *ma.Client {
func (c *AtlasConfig) newMongoDBAtlasClient() *Client {
	t := dac.NewTransport(c.AtlasPublicKey, c.AtlasPrivateKey)
	httpClient := &http.Client{Transport: &t}
	//client := ma.NewClient(httpClient)
	client := NewClient(httpClient)
	return client
}

// GetAtlasClient returns a MongoDB Atlas client
//func GetAtlasClient() *ma.Client {
func GetAtlasClient() *Client {
	// create MongoDB Atlas client
	privateKey, ok := os.LookupEnv("ATLAS_PRIVATE_KEY")
	if ok != true {
		panic("Error fetching private key: Env variable ATLAS_PRIVATE_KEY not set.")
	}
	publicKey, ok := os.LookupEnv("ATLAS_PUBLIC_KEY")
	if ok != true {
		panic("Error fetching public key: Env variable ATLAS_PUBLIC_KEY not set.")
	}
	atlasConfig := AtlasConfig{
		AtlasPublicKey:  publicKey,
		AtlasPrivateKey: privateKey,
	}
	return atlasConfig.newMongoDBAtlasClient()
}

// ReconciliationConfig let us customize reconcilitation
type ReconciliationConfig struct {
	Time time.Duration
}

// GetReconcilitationConfig gives us default values
func GetReconcilitationConfig() *ReconciliationConfig {
	// default reconciliation loop time is 2 minutes
	timeString := getenv("RECONCILIATION_TIME", "120")
	timeInt, _ := strconv.Atoi(timeString)
	reconciliationTime := time.Second * time.Duration(timeInt)
	return &ReconciliationConfig{
		Time: reconciliationTime,
	}
}

func getenv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}

// Get a databaseUser in the specified group.
// https://docs.atlas.mongodb.com/reference/api/database-users-get-single-user/
func (c *DatabaseUserService) Get(gid string, username string) (*DatabaseUser, *http.Response, error) {
	databaseUser := new(DatabaseUser)
	apiError := new(APIError)
	path := fmt.Sprintf("%s/databaseUsers/admin/%s", gid, username)
	resp, err := c.sling.New().Get(path).Receive(databaseUser, apiError)
	return databaseUser, resp, relevantError(err, *apiError)
}

// Create a databaseUser in the specified group.
// https://docs.atlas.mongodb.com/reference/api/databaseUsers-create-one/
func (c *DatabaseUserService) Create(gid string, databaseUserParams *DatabaseUser) (*DatabaseUser, *http.Response, error) {
	databaseUser := new(DatabaseUser)
	apiError := new(APIError)
	path := fmt.Sprintf("%s/databaseUsers", gid)
	resp, err := c.sling.New().Post(path).BodyJSON(databaseUserParams).Receive(databaseUser, apiError)
	return databaseUser, resp, relevantError(err, *apiError)
}

// Update a databaseUser in the specified group.
// https://docs.atlas.mongodb.com/reference/api/databaseUsers-modify-one/
func (c *DatabaseUserService) Update(gid string, username string, databaseUserParams *DatabaseUser) (*DatabaseUser, *http.Response, error) {
	databaseUser := new(DatabaseUser)
	apiError := new(APIError)
	path := fmt.Sprintf("%s/databaseUsers/admin/%s", gid, username)
	resp, err := c.sling.New().Patch(path).BodyJSON(databaseUserParams).Receive(databaseUser, apiError)
	return databaseUser, resp, relevantError(err, *apiError)
}

// Delete a databaseUser in the specified group.
// https://docs.atlas.mongodb.com/reference/api/databaseUsers-delete-one/
func (c *DatabaseUserService) Delete(gid string, username string) (*http.Response, error) {
	databaseUser := new(DatabaseUser)
	apiError := new(APIError)
	path := fmt.Sprintf("%s/databaseUsers/admin/%s", gid, username)
	resp, err := c.sling.New().Delete(path).Receive(databaseUser, apiError)
	return resp, relevantError(err, *apiError)
}
