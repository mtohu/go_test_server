package repositories

import (
	"database/sql"
	"gotest/common"
	"gotest/datasource"
	"strings"
)

type BaseRepository struct {
    TableName string
    FieldStr  string
	AliasName string
	WhereStr  string
	OrderStr  string
	GroupStr  string
	HavingStr string
	JoinStr   string
    LimitStr  string
    UpdateStr string
    InsertCols string
    InsertValues string
    SaveData   map[string]interface{}
    Tx         *sql.Tx
    Db         *sql.DB
    IsTran     int
    ActionType int //1=select 2=insert 3=update 4=delete 5=count 6=max 7=sum
}



func getSql(b *BaseRepository) string {
	var sql []string
	if b.ActionType == 0 {
		return ""
	}
	switch b.ActionType {
	case  1,5,6,7:
		sql=append(append(append(append(sql,"select"),b.FieldStr),"from"),b.TableName)
	case  2:
		sql=append(append(sql,"insert into"),b.TableName)
	case  3:
		sql=append(append(sql,"update"),b.TableName)
		if len(strings.Trim(b.AliasName," ")) > 0{
			sql=append(sql,b.AliasName)
		}
		sql=append(sql,"set")
	case  4:
		sql=append(append(append(sql,"delete"),"from"),b.TableName)

	}
	if len(strings.Trim(b.AliasName," ")) > 0 && b.ActionType !=2 {
		sql=append(sql,b.AliasName)
	}
	if len(strings.Trim(b.JoinStr," ")) > 0{
		sql=append(sql,b.JoinStr)
	}
	if len(strings.Trim(b.UpdateStr," ")) > 0{
		sql=append(sql,b.UpdateStr)
	}
	if len(strings.Trim(b.WhereStr," ")) > 0{
		sql=append(append(sql,"where"),b.WhereStr)
	}
	if len(strings.Trim(b.GroupStr," ")) > 0{
		sql=append(append(sql,"group by"),b.GroupStr)
	}
	if len(strings.Trim(b.HavingStr," ")) > 0{
		sql=append(append(sql,"having"),b.HavingStr)
	}
	if len(strings.Trim(b.OrderStr," ")) > 0{
		sql=append(append(sql,"order by"),b.OrderStr)
	}
	if len(strings.Trim(b.LimitStr," ")) > 0{
		sql=append(append(sql,"limit"),b.LimitStr)
	}
	return strings.Join(sql," ")
}
func (b *BaseRepository) GroupData(actionType int) []interface{} {
	if len(b.SaveData) <=0 {
		return nil
	}
	vargs:=make([]interface{},0)
	vcols:=make([]string,0)
    if actionType == 3{
		for cols,vals :=range b.SaveData {
			vcols = append(vcols,cols+"=?")
			vargs=append(vargs,vals)
		}
		b.UpdateStr=strings.Join(vcols,",")
	}else if actionType == 2 {
		vvargs :=make([]string,0)
		for cols,vals :=range b.SaveData {
			vcols = append(vcols,cols)
			vargs=append(vargs,vals)
			vvargs=append(vvargs,"?")
		}
		buf := make([]byte,0)
		buf=append(append(append(append(buf,"( "...),strings.Join(vcols,",")...)," )"...)," values "...)
		buf=append(append(append(buf,"( "...),strings.Join(vvargs,",")...)," ) "...)
		b.InsertCols=string(buf)
	}
	return vargs
}
func (b *BaseRepository) Query (sql string,args ...interface{}) (*sql.Rows,error){
	 result,err :=b.Db.Query(sql,args)
	 return result,err
}
func (b *BaseRepository) Execute(sql string,args ...interface{}) (int,int,error) {
     result,err :=b.Db.Exec(sql,args)
	 i,errr:=result.LastInsertId()
	 if errr !=nil{
	 	i = 0
	 }
     r,errr:=result.RowsAffected()
     if errr !=nil{
     	r = 0
	 }
     return int(i),int(r),err
}
func (b *BaseRepository) Begin() error {
	var err error
    b.Tx,err=common.Db.Begin()
    b.IsTran=1
    return err
}
func (b *BaseRepository) Rollback() error {
	err := b.Tx.Rollback()
	return err
}
func (b *BaseRepository) Commit() error {
	err := b.Tx.Commit()
	return err
}
func (b *BaseRepository) Counts(a string) int {
	b.ActionType=5
	b.FieldStr=string([]byte("count("+a+")"))
	b.LimitStr="1"
	sqls := getSql(b)
	common.Log(sqls)
	return 0
}

func (b *BaseRepository) FetchSql(actionType int) string {
	b.ActionType =actionType
	if actionType == 3 || actionType == 2{
		b.GroupData(actionType)
	}
	sqls := getSql(b)
	return sqls
}

func (b *BaseRepository) DataMap(data map[string]interface{}) IBaseRepository {
	b.SaveData=data
	return b
}

func (b *BaseRepository) Update() int {
	if len(b.SaveData) <=0 {
		return 0
	}
    b.ActionType=3
	vargs :=b.GroupData(b.ActionType)
	sqls := getSql(b)
	var rid int
	var i int64 = 0
	if b.IsTran == 1{
		var err error
		stmt, err:= b.Tx.Prepare(sqls)
		if err == nil {
			result, _ := stmt.Exec(vargs)
			_,err=result.LastInsertId()
			i,err=result.RowsAffected()
			if err !=nil{
				i = 0
			}
		}
		rid = int(i)
	}else{
		_,rid,_ =b.Execute(sqls,vargs)
	}
	return rid
}


func (b *BaseRepository) Insert() int {
	if len(b.SaveData) <=0 {
		return 0
	}
	b.ActionType=2
	vargs :=b.GroupData(b.ActionType)
	sqls := getSql(b)
	var lid int
	var i int64 = 0
	if b.IsTran == 1{
		var err error
		stmt, err:= b.Tx.Prepare(sqls)
		if err == nil {
			result, _ := stmt.Exec(vargs)
			i,err=result.LastInsertId()
			if err !=nil{
				i = 0
			}
		}
		lid = int(i)
	}else{
		_,lid,_ =b.Execute(sqls,vargs)
	}
	return lid
}

func (b *BaseRepository) Select() map[int]map[string]string {
	b.ActionType = 1
	var err error
	sqls := getSql(b)
	rows,err:=b.Query(sqls)
	if err !=nil{
		return nil
	}
	cols, _ := rows.Columns()
	//这里表示一行所有列的值，用[]byte表示
	vals := make([][]byte, len(cols))
	//这里表示一行填充数据
	scans := make([]interface{}, len(cols))
	//这里scans引用vals，把数据填充到[]byte里
	for k, _ := range vals {
		scans[k] = &vals[k]
	}
	i := 0
	result := make(map[int]map[string]string)
	for rows.Next(){
		//填充数据
		rows.Scan(scans...)
		//每行数据
		row := make(map[string]string)
		//把vals中的数据复制到row中
		for k, v := range vals {
			key := cols[k]
			//这里把[]byte数据转成string
			row[key] = string(v)
		}
		//放入结果集
		result[i] = row
		i++
	}
	return result
}

func (b *BaseRepository) FindOne(args ...interface{}) error{
	b.ActionType = 1
	sqls := getSql(b)
	err := b.Db.QueryRow(sqls).Scan(args)
	return err
}

func (b *BaseRepository) Delete() int {
	b.ActionType = 4
	sqls :=getSql(b)
	var rid int
	var i int64 = 0
	if b.IsTran == 1{
		var err error
		stmt, err:= b.Tx.Prepare(sqls)
		if err == nil {
			result, _ := stmt.Exec()
			_,err=result.LastInsertId()
			i,err=result.RowsAffected()
			if err !=nil{
				i = 0
			}
		}
		rid = int(i)
	}else{
		_,rid,_=b.Execute(sqls)
	}
	return rid
}

func (b *BaseRepository) Pages(p int, offsets int) IBaseRepository{
	var start int = 0
	if (p-1) > 0{
		p = p-1
	}else{
		p = 0
	}
	start = p * offsets
	b.Limit(start,offsets)
	return b
}

func (b *BaseRepository) Fields(f string) IBaseRepository{
	b.FieldStr = f
	return b
}

func (b *BaseRepository) Where(w string) IBaseRepository{
	b.WhereStr = w
	return b
}

func (b *BaseRepository) Order(o string) IBaseRepository{
	b.OrderStr = o
	return b
}

func (b *BaseRepository) Group(g string) IBaseRepository{
	b.GroupStr = g
	return b
}

func (b *BaseRepository) Having(h string) IBaseRepository{
	b.HavingStr = h
	return b
}

func (b *BaseRepository) Join(j string) IBaseRepository{
	b.JoinStr=j
	return b
}

func (b *BaseRepository) Limit(start int, end int) IBaseRepository{
	if start >0 && end >0 {
		b.LimitStr=strings.Join([]string{string(start),string(end)},",")
	}else if start >0 {
		b.LimitStr=string(start)
	}else {
		b.LimitStr=string(end)
	}
	return b
}

func (b *BaseRepository) connect(db string) IBaseRepository {
	var err error
	b.Db,err=datasource.SwtichDb(db)
	if (err !=nil) {
		common.Log(err)
	}
	return b
}

func (b *BaseRepository) Table(name string) IBaseRepository{
	b.TableName=name
	return b
}
func (b *BaseRepository) Alias(a string) IBaseRepository{
	b.AliasName=a
	return b
}

type IBaseRepository interface {
	Table (name string) IBaseRepository
	Alias (a string)  IBaseRepository
	Fields (f string)  IBaseRepository
	Where (w string)  IBaseRepository
	Order (o string)  IBaseRepository
	Group (g string)  IBaseRepository
	Having(h string)  IBaseRepository
	Join  (j string)  IBaseRepository
	Limit (start int,end int) IBaseRepository
	Counts (a string) int //只能在最后调用
	FetchSql (actionType int) string  //只能在最后调用
	Select   () map[int]map[string]string
	FindOne(args ...interface{}) error
	Update   () int
	Insert   () int
	DataMap  (s  map[string]interface{}) IBaseRepository
	Pages     (p int,offsets int) IBaseRepository //跟Limit各选一种
	Query (sql string,args ...interface{}) (*sql.Rows,error)
	Execute (sql string,args ...interface{}) (int,int,error)
	Delete() int
	Begin() error
	Rollback() error
	Commit() error
	connect(db string) IBaseRepository
}

func NewBaseRepository() IBaseRepository {
	var b BaseRepository
	b.TableName=""
	b.FieldStr=""
	b.AliasName=""
	b.WhereStr=""
	b.OrderStr=""
	b.GroupStr=""
	b.HavingStr=""
	b.JoinStr=""
	b.LimitStr=""
	b.ActionType=0
	b.UpdateStr=""
	b.InsertCols=""
	b.InsertValues=""
	b.Db = common.Db
	b.IsTran = 0
	return &b
}

