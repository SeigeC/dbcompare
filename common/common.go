package common

import (
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v3"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	Newtables tables
	Oldtables tables
	Config    config
)

func (t *tables) getTables(conn *gorm.DB, table string) {
	conn.Table("TABLES").Where("TABLE_SCHEMA = ?", table).Pluck("table_name", t)
}

type tables []string

func init() {
	Config = loadConfigFile()
	conn := GetTables(Config.NewDB)
	Newtables.getTables(conn, Config.NewDB.Name)
	conn = GetTables(Config.OleDb)
	Oldtables.getTables(conn, Config.OleDb.Name)
}
func GetDB(config DBconfig) (*gorm.DB, error) {
	str := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.User,
		config.Password,
		config.IP,
		config.Port,
		config.Name,
	)
	return gorm.Open(mysql.Open(str), &gorm.Config{})
}

func GetTables(config DBconfig) *gorm.DB {
	str := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.User,
		config.Password,
		config.IP,
		config.Port,
		"information_schema",
	)
	db, err := gorm.Open(mysql.Open(str), &gorm.Config{})
	if err != nil {
		panic("获取数据库表失败")
	}
	return db
}

type config struct {
	NewDB   DBconfig `yaml:"newDB"`
	OleDb   DBconfig `yaml:"oldDB"`
	Linters []string `yaml:"linter"`
}

type DBconfig struct {
	Name     string `yaml:"name"`
	IP       string `yaml:"ip"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

func loadConfigFile() config {
	wd, _ := os.Getwd()
	f, err := ioutil.ReadFile(wd + "/config.yml")
	if err != nil {
		panic(err)
	}
	return decodeFile(f)
}

func decodeFile(b []byte) config {
	var config config
	err := yaml.Unmarshal(b, &config)
	if err != nil {
		panic(err)
	}
	return config
}
