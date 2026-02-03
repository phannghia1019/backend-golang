package utils

// CompanyInfo holds information about a company
type CompanyInfo struct {
	Name   string
	Value  string
	DbName string
}

// CompanyInfoMap stores all company information
var CompanyInfoMap = map[string]CompanyInfo{
	"HBCCORP": {Name: "Hòa Bình", Value: "HBCCORP", DbName: "HBCCORP"},
	"KBI":     {Name: "Kim Bình", Value: "KBI", DbName: "KBI"},
	"HBC_DN":  {Name: "Hồng Hải Bình", Value: "HBC_DN", DbName: "HBC_DN"},
	"HBC_HN":  {Name: "Tâm Bình", Value: "HBC_HN", DbName: "HBC_HN"},
	"HBC_TN":  {Name: "Gia Bình", Value: "HBC_TN", DbName: "HBC_TN"},
	"HBC_LA":  {Name: "An Kiên Bình", Value: "HBC_LA", DbName: "HBC_LA"},
	"AKBC":    {Name: "An Khánh Bình", Value: "AKBC", DbName: "AKBC"},
	"LBC_TT":  {Name: "Long Bình", Value: "LBC_TT", DbName: "LBC_TT"},
}

// GetCompanyInfo returns company info by database name
func GetCompanyInfo(dbName string) (CompanyInfo, bool) {
	info, ok := CompanyInfoMap[dbName]
	return info, ok
}

// GetCompanyName returns company display name
func GetCompanyName(dbName string) string {
	if info, ok := CompanyInfoMap[dbName]; ok {
		return info.Name
	}
	return dbName
}

// IsValidDbName checks if database name is valid
func IsValidDbName(dbName string) bool {
	_, ok := CompanyInfoMap[dbName]
	return ok
}

// GetSupportedDbNames returns all supported database names
func GetSupportedDbNames() []string {
	names := make([]string, 0, len(CompanyInfoMap))
	for k := range CompanyInfoMap {
		names = append(names, k)
	}
	return names
}
