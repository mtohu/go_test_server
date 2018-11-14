package datasource

import (
	"database/sql"
	"fmt"
	"gopkg.in/yaml.v2"
	"gotest/common"
	"io/ioutil"
	"path/filepath"
	_ "github.com/go-sql-driver/mysql"
)

func Getinstance() (*sql.DB,error){

	var c DBconfigs
	c,err := db_config("./configs/databases.yml")
	if(err !=nil){
		common.Log("sql db databases yml error",err)
		return nil,err
	}
    var dns string
	dns = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s",c.UserName,c.PassWord,c.Host,c.Port,c.DBName,c.Charset)
	common.Db, err = sql.Open("mysql", dns)
	if(err !=nil){
		common.Log("open mysql error",err)
		return nil,err
	}
	//defer db.Close()
	common.Db.SetMaxOpenConns(c.OpenNum)
	common.Db.SetMaxIdleConns(c.IdleNum)
	err =common.Db.Ping()
	if(err !=nil){
		common.Log("mysql ping error",err)
		return nil,err
	}
	return common.Db,nil
}
func SwtichDb(dbname string) (*sql.DB,error)  {
	var c DBconfigs
	var db *sql.DB
	var filedb = append(append([]byte("./configs/"),dbname...),"_databases.yml"...)
	c,err := db_config(string(filedb))
	if(err !=nil){
		common.Log("sql db databases yml error",err)
		return nil,err
	}
	var dns string
	dns = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s",c.UserName,c.PassWord,c.Host,c.Port,c.DBName,c.Charset)
	db, err = sql.Open("mysql", dns)
	if(err !=nil){
		common.Log("open mysql error",err)
		return nil,err
	}
	//defer db.Close()
	db.SetMaxOpenConns(c.OpenNum)
	db.SetMaxIdleConns(c.IdleNum)
	err =db.Ping()
	if(err !=nil){
		common.Log("mysql ping error",err)
		return nil,err
	}
	return db,nil
}
func db_config(filename string) (DBconfigs,error){
	c :=DefaultDbconfig()
	yamlAbsPath, err := filepath.Abs(filename)
	if err != nil {
		return c,err
	}
	// read the raw contents of the file
	data, err := ioutil.ReadFile(yamlAbsPath)
	if err != nil {
		return c, err
	}
	// put the file's contents as yaml to the default configuration(c)
	if err := yaml.Unmarshal(data, &c); err != nil {
		return c, err
	}
	return c, nil

}
