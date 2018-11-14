package main

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

//go 小写 和 大写区分  大写 public  小写private
type home struct {
	Id int
	Name string
}
type Per struct {
	name string
	bh   int
}
type DebugInfo struct {
	Level  string `json:"level"`   // level 解码为 Level
	Msg    string `json:"message"` // message 解码为 Msg
	Author string `json:"-"`       // 忽略Author
}
func (dbgInfo DebugInfo) String() string {
	return fmt.Sprintf("{Level: %s, Msg: %s}", dbgInfo.Level, dbgInfo.Msg)
}
type Point struct{ X, Y int }

func (pt Point)MarshalJSON() ([]byte, error) {
	fmt.Println("00000===")
	return []byte(fmt.Sprintf(`{"XX":%d,"YY":%d}`, pt.X, pt.Y)), nil
}
func main() {
	var ssc,vvc string = "333","444"
	vvb := strings.Join([]string{ssc, vvc},"")
	fmt.Println(vvb)
	fmt.Fprintf(os.Stderr, "an %s\n", "error")
	if data, err := json.Marshal(Point{50, 50}); err == nil {
		fmt.Printf("%s\n", data)
	}
	dbgInfs := []DebugInfo{
		DebugInfo{"debug", `File: "test.txt" Not Found`, "Cynhard"},
		DebugInfo{"", "Logic error", "Gopher"},
	}
	if data, err := json.Marshal(dbgInfs); err == nil {
		fmt.Printf("%s\n", data)
	}
	/*var bh home
	bh.id=1
	bh.name="ddd"*/
	//bh := home{1,"bbb"}
	//vv_test(bh)
	vv_ss("bbbbssss===")
	data := `[{"level":"debug","message":"File Not Found","author":"Cynhard"},` +
		`{"level":"","message":"Logic error","author":"Gopher"}]`
	var dbgInfos []DebugInfo
	json.Unmarshal([]byte(data), &dbgInfos)
	fmt.Println(dbgInfos)

	var dbgInfoss []map[string]string
	json.Unmarshal([]byte(data), &dbgInfoss)

	fmt.Println("======",dbgInfoss)

}
func vv_ss(a interface{}){
	c := a.(string)
	fmt.Print(c)
}
func vv_test(a interface{})  {
	getType := reflect.TypeOf(a)
	fmt.Println("ssss===",getType.Name())
	getValue := reflect.ValueOf(a)
	fmt.Println("get all Fields is:", getValue)
	for i := 0; i < getType.NumField(); i++ {
		field := getType.Field(i)
		valuess := getValue.Field(i).Interface()
		fmt.Printf("%s: %v = %v\n", field.Name, field.Type,valuess)
	}
	//fmt.Printf("cat4 information: nil?:%5v, type=%15v, type.kind=%5v, value=%5v  \n",reflect.TypeOf(a).Kind())
	p := &Per{}
	for i := 0; i < 100; i++ {
		p.bh=i
		p.name=strconv.Itoa(i)
		go pers(p)
	}
}
func pers(p *Per) {
	fmt.Println(p)
}
