package linters

import (
	"fmt"
	"os"
	"os/exec"

	"dbcompare/common"

	"gorm.io/gorm"
)

var (
	newsql = "newtest.sql"
	oldsql = "oldtest.sql"
	parse  = "parse.py"
)

func dumpStr(config common.DBconfig) []string {
	arr := []string{
		"--host=" + config.IP,
		"--port=" + config.Port,
		"-u" + config.User,
		"-p" + config.Password,
		"-t",
		config.Name,
	}
	return arr
}
func init() {
	getwd, _ := os.Getwd()
	newsql = fmt.Sprintf("%s/%s", getwd, newsql)
	oldsql = fmt.Sprintf("%s/%s", getwd, oldsql)
	parse = fmt.Sprintf("%s/linters/%s", getwd, parse)
}

func getSqlFile(config common.DBconfig, file string) error {
	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer f.Close()
	command := exec.Command("mysqldump", dumpStr(config)...)
	out, err := command.CombinedOutput()
	if err != nil {
		return err
	}
	f.WriteString(string(out))
	return nil
}

func CheckSql(newdb *gorm.DB, olddb *gorm.DB) bool {

	_, err := exec.LookPath("mysqldump")
	if err != nil {
		Message = append(Message, "mysqldump 不存在")
		return false
	}

	_, err = exec.LookPath("diff-so-fancy")
	if err != nil {
		Message = append(Message, "diff-so-fancy 不存在")
		return false
	}

	if err = getSqlFile(common.Config.NewDB, newsql); err != nil {
		return false
	}

	if err = getSqlFile(common.Config.OleDb, oldsql); err != nil {
		return false
	}
	// 格式化 sql 文件
	// TODO: 找到一个用 go 实现的 sql format 包
	command := exec.Command("python3", parse, oldsql, newsql)
	err = command.Run()
	if err != nil {
		panic(err)
	}

	command = exec.Command("bash", "-c",
		fmt.Sprintf("diff -u %s %s | diff-so-fancy", oldsql, newsql))
	output, err := command.CombinedOutput()
	if err != nil {
		panic(err)
	}
	Message = append(Message, string(output))
	return true
}
