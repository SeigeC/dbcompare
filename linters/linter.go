package linters

import "gorm.io/gorm"

var (
	LintMap map[string]func(newdb *gorm.DB, olddb *gorm.DB) bool
	Message []string
)

func init() {
	LintMap = map[string]func(newdb *gorm.DB, olddb *gorm.DB) bool{
		"checkTable": CheckTable,
		"checkDate":  CheckData,
		"checkSql":   CheckSql,
	}
}


