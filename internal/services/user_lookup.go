package services

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hbc/app-hrm-api-go/internal/config"
	"github.com/hbc/app-hrm-api-go/internal/utils"
)

var userLookupLog = utils.NewLogger("UserLookup")

// UserLookupResult contains user information from database
type UserLookupResult struct {
	Email         string
	Name          string
	CompanyDbName string
	CompanyName   string
	ManagerEmail  string
}

// UserLookupService handles user lookup across SAP databases
type UserLookupService struct {
	queryTimeout time.Duration
}

// NewUserLookupService creates a new user lookup service
func NewUserLookupService() *UserLookupService {
	return &UserLookupService{
		queryTimeout: 3 * time.Second,
	}
}

// FindUserAcrossAllDatabases finds user in all SAP databases
func (s *UserLookupService) FindUserAcrossAllDatabases(email string) (*UserLookupResult, error) {
	if email == "" {
		userLookupLog.Warn("Invalid email provided: empty")
		return nil, nil
	}

	normalizedEmail := strings.ToLower(strings.TrimSpace(email))
	userLookupLog.Info("Searching for user: %s across all databases", normalizedEmail)

	dbList := s.getDatabaseList()
	userLookupLog.Info("Searching in %d databases: %s", len(dbList), strings.Join(dbList, ", "))

	// Search databases with timeout
	for _, dbName := range dbList {
		result, err := s.searchInDatabaseWithTimeout(normalizedEmail, dbName)
		if err != nil {
			userLookupLog.Warn("Failed to search in %s: %v", dbName, err)
			continue
		}
		if result != nil {
			userLookupLog.Info("Found user in database: %s", result.CompanyDbName)
			return result, nil
		}
	}

	userLookupLog.Info("User not found in any database: %s", normalizedEmail)
	return nil, nil
}

// searchInDatabaseWithTimeout searches for user in a database with timeout
func (s *UserLookupService) searchInDatabaseWithTimeout(email, dbName string) (*UserLookupResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.queryTimeout)
	defer cancel()

	resultChan := make(chan *UserLookupResult, 1)
	errChan := make(chan error, 1)

	go func() {
		result, err := s.searchInDatabase(email, dbName)
		if err != nil {
			errChan <- err
			return
		}
		resultChan <- result
	}()

	select {
	case <-ctx.Done():
		userLookupLog.Warn("Timeout searching in database: %s", dbName)
		return nil, ctx.Err()
	case err := <-errChan:
		return nil, err
	case result := <-resultChan:
		return result, nil
	}
}

// searchInDatabase searches for user in a specific database
func (s *UserLookupService) searchInDatabase(email, dbName string) (*UserLookupResult, error) {
	db, err := config.ConnectSapDB(dbName)
	if err != nil {
		return nil, err
	}

	query := `
		SELECT 
			e.Email,
			CONCAT(e.lastName, ' ', e.middleName, ' ', e.firstName) AS FullName,
			m.Email AS ManagerEmail
		FROM OHEM e
		LEFT JOIN OHEM m ON e.manager = m.empID
		WHERE e.Active = 'Y' AND LOWER(e.Email) = @email
	`

	row := db.QueryRow(query, email)

	var emailResult, fullName string
	var managerEmail *string

	err = row.Scan(&emailResult, &fullName, &managerEmail)
	if err != nil {
		// No rows found is not an error
		return nil, nil
	}

	manager := ""
	if managerEmail != nil {
		manager = *managerEmail
	}

	return &UserLookupResult{
		Email:         emailResult,
		Name:          strings.TrimSpace(fullName),
		CompanyDbName: dbName,
		CompanyName:   utils.GetCompanyName(dbName),
		ManagerEmail:  manager,
	}, nil
}

// getDatabaseList returns list of databases to search
func (s *UserLookupService) getDatabaseList() []string {
	cfg := config.Get()
	if cfg.SAPDBList == "" {
		userLookupLog.Warn("SAP_DB_LIST not configured, using default databases")
		return []string{"HBCCORP", "KBI"}
	}

	parts := strings.Split(cfg.SAPDBList, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

// GetCompanyName returns company display name
func (s *UserLookupService) GetCompanyName(dbName string) string {
	info, ok := utils.GetCompanyInfo(dbName)
	if !ok {
		return dbName
	}
	return info.Name
}
