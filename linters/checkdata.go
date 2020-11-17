package linters

import (
	"fmt"
	"reflect"

	"dbcompare/common"

	"gorm.io/gorm"
)

// TODO: 读取主键名称
func CheckData(newdb *gorm.DB, olddb *gorm.DB) bool {
	for _, table := range common.Newtables {
		var newResults []map[string]interface{}
		var oldResults []map[string]interface{}
		newdb.Table(table).Find(&newResults)
		olddb.Table(table).Find(&oldResults)
		i, j := 0, 0
		for i < len(newResults) && j < len(oldResults) {
			nid, ok := newResults[i]["key"].(int32)
			if !ok {
				Message = append(Message, fmt.Sprintf("主键无法对比: %s %v",
					table, newResults[i]["key"]))
				continue
			}
			oid, ok := oldResults[j]["key"].(int32)
			if !ok {
				Message = append(Message, fmt.Sprintf("主键无法对比: %s %v",
					table, oldResults[j]["key"]))
				continue
			}
			if nid == oid {
				if !reflect.DeepEqual(newResults[i], oldResults[j]) {
					Message = append(Message, fmt.Sprintf("数据不同: %s %v %v",
						table, newResults[i], oldResults[j]))
				}
				i++
				j++
				continue
			}
			if nid > oid {
				Message = append(Message, fmt.Sprintf("缺失数据: %s %v",
					table, oldResults[j]))
				j++
				continue
			}
			if nid < oid {
				Message = append(Message, fmt.Sprintf("多出数据: %s %v",
					table, oldResults[i]))
				i++
			}
		}
		if i < len(newResults) {
			for _, v := range newResults[i:] {
				Message = append(Message, fmt.Sprintf("多出数据: %s %v",
					table, v))
			}
		}
		if j < len(oldResults) {
			for _, v := range oldResults[j:] {
				Message = append(Message, fmt.Sprintf("缺失数据: %s %v",
					table, v))
			}
		}
	}

	return true
}
