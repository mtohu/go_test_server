package datasource

import (
	"fmt"
	"github.com/go-xorm/core"
	"github.com/go-xorm/xorm"
	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/yaml.v2"
	"gotest/common"
	"io/ioutil"
	"path/filepath"
)

func Instance() (*xorm.Engine,error){
	var c DBconfigs
	c,err := xdb_config("./configs/databases.yml")
	if(err !=nil){
		common.Log("xorm databases yml error",err)
		return nil,err
	}
	var dns string
	dns = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s",c.UserName,c.PassWord,c.Host,c.Port,c.DBName,c.Charset)
	//创建orm引擎
	common.Dbengine, err = xorm.NewEngine("mysql", dns)
	if err!=nil{
		common.Log("xorm new engine error",err)
		return nil,err
	}
	//defer engine.Close()
	//连接测试
	if err := common.Dbengine.Ping(); err!=nil{
		common.Log("xorm ping error",err)
		return nil,err
	}

	//日志打印SQL
	common.Dbengine.ShowSQL(c.ShowSQL)

	//设置连接池的空闲数大小
	common.Dbengine.SetMaxIdleConns(c.IdleNum)
	//设置最大打开连接数
	common.Dbengine.SetMaxOpenConns(c.OpenNum)
	//名称映射规则主要负责结构体名称到表名和结构体field到表字段的名称映射
	common.Dbengine.SetTableMapper(core.SnakeMapper{})

    return common.Dbengine,nil
}
func xdb_config(filename string) (DBconfigs,error){
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
