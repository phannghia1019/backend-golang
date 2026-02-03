package services

import (
	"github.com/hbc/app-hrm-api-go/internal/utils"
)

var companyMappingLog = utils.NewLogger("CompanyMapping")

// CompanyMappingService handles email to company mapping
type CompanyMappingService struct {
	userLookup     *UserLookupService
	defaultCompany string
}

// NewCompanyMappingService creates a new company mapping service
func NewCompanyMappingService() *CompanyMappingService {
	return &CompanyMappingService{
		userLookup:     NewUserLookupService(),
		defaultCompany: "HBCCORP",
	}
}

// GetCompanyDbNameByEmail maps email to company database name
func (s *CompanyMappingService) GetCompanyDbNameByEmail(email string) (string, error) {
	if email == "" {
		companyMappingLog.Warn("Invalid email provided: empty, using default company")
		return s.defaultCompany, nil
	}

	userInfo, err := s.userLookup.FindUserAcrossAllDatabases(email)
	if err != nil {
		companyMappingLog.Error("Error looking up %s: %v", email, err)
		return s.defaultCompany, err
	}

	if userInfo != nil && userInfo.CompanyDbName != "" {
		companyMappingLog.Info("Found %s in database: %s", email, userInfo.CompanyDbName)
		return userInfo.CompanyDbName, nil
	}

	companyMappingLog.Warn("Email %s not found in any SAP database, using default company", email)
	return s.defaultCompany, nil
}

// GetCompanyInfo returns company info by database name
func (s *CompanyMappingService) GetCompanyInfo(dbName string) (utils.CompanyInfo, bool) {
	return utils.GetCompanyInfo(dbName)
}

// GetAllCompanies returns all supported companies
func (s *CompanyMappingService) GetAllCompanies() map[string]utils.CompanyInfo {
	return utils.CompanyInfoMap
}

// IsValidDbName validates if database name is supported
func (s *CompanyMappingService) IsValidDbName(dbName string) bool {
	return utils.IsValidDbName(dbName)
}

// GetCompanyName returns company display name
func (s *CompanyMappingService) GetCompanyName(dbName string) string {
	return utils.GetCompanyName(dbName)
}
