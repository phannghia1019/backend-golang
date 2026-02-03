package config

import (
	"database/sql"
	"fmt"
	"sync"

	"github.com/hbc/app-hrm-api-go/internal/utils"
	_ "github.com/denisenkom/go-mssqldb"
)

var (
	sapConnections = make(map[string]*sql.DB)
	sapMutex       sync.RWMutex
	sapLog         = utils.NewLogger("SAPDB")
)

// ConnectSapDB connects to a specific SAP database
func ConnectSapDB(dbName string) (*sql.DB, error) {
	// Check if connection already exists
	sapMutex.RLock()
	if conn, ok := sapConnections[dbName]; ok {
		sapMutex.RUnlock()
		return conn, nil
	}
	sapMutex.RUnlock()

	cfg := Get()

	// Build connection string
	connString := fmt.Sprintf(
		"server=%s;port=%s;user id=%s;password=%s;database=%s;encrypt=disable",
		cfg.SAPDBHost,
		cfg.SAPDBPort,
		cfg.SAPDBUser,
		cfg.SAPDBPass,
		dbName,
	)

	db, err := sql.Open("sqlserver", connString)
	if err != nil {
		sapLog.Error("Failed to open connection to %s: %v", dbName, err)
		return nil, err
	}

	// Verify connection
	if err := db.Ping(); err != nil {
		sapLog.Error("Failed to ping %s: %v", dbName, err)
		db.Close()
		return nil, err
	}

	// Store connection
	sapMutex.Lock()
	sapConnections[dbName] = db
	sapMutex.Unlock()

	sapLog.Info("Connected to SAP database: %s", dbName)
	return db, nil
}

// GetSapDB returns an existing SAP database connection
func GetSapDB(dbName string) (*sql.DB, bool) {
	sapMutex.RLock()
	defer sapMutex.RUnlock()
	conn, ok := sapConnections[dbName]
	return conn, ok
}

// CloseSapConnections closes all SAP database connections
func CloseSapConnections() {
	sapMutex.Lock()
	defer sapMutex.Unlock()

	for name, conn := range sapConnections {
		if err := conn.Close(); err != nil {
			sapLog.Error("Failed to close connection to %s: %v", name, err)
		}
	}
	sapConnections = make(map[string]*sql.DB)
}
