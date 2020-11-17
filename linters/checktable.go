package linters

import (
	"fmt"

	"dbcompare/common"

	"gorm.io/gorm"
)

func CheckTable(newdb *gorm.DB, olddb *gorm.DB) bool {
	var find = func(tables []string, tb string) bool {
		for _, ot := range common.Oldtables {
			if ot == tb {
				return true
			}
		}
		return false
	}
	for _, nt := range common.Newtables {
		if !find(common.Oldtables, nt) {
			Message = append(Message, fmt.Sprint("新数据库中不存在 ? 表", nt))
		}
	}
	// TODO: 验证表结构
	return true
}
