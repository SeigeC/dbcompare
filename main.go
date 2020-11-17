package main

import (
	"fmt"

	"dbcompare/common"
	"dbcompare/linters"

	"gorm.io/gorm"
)

func main() {
	newDB, err := common.GetDB(common.Config.NewDB)
	if err != nil {
		panic("新数据库连接失败")
	}
	oldDB, err := common.GetDB(common.Config.OleDb)
	if err != nil {
		panic("原始数据库连接失败")
	}

	var lints []func(newdb *gorm.DB, olddb *gorm.DB) bool
	for _, v := range common.Config.Linters {
		linter, ok := linters.LintMap[v]
		if !ok {
			panic("插件不存在" + v)
		}
		lints = append(lints, linter)
	}
	if lints == nil {
		fmt.Println("请配置使用的插件")
		return
	}
	for _, f := range lints {
		f(newDB, oldDB)
	}
	for _, v := range linters.Message {
		fmt.Println(v)
	}
}
