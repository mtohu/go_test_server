package common

import (
	"database/sql"
	"flag"
	"github.com/go-xorm/xorm"
	"gotest/protocols"
	"net"
	"reflect"
	"strings"
)

var (
	Conns      = make(map[string]net.Conn)//map[conn.RemoteAddr().String()]
	UserConns  = make(map[string]net.Conn)//map[uid]
	AddrsUserOrGroup = make(map[string]map[string]interface{})//map[conn.RemoteAddr().String()]["uid"]map[conn.RemoteAddr().String()]["group"][0]groupid
	GroupConns = make(map[string]map[string]net.Conn)//map[groupid]map[uid/conn.RemoteAddr().String()]
	H          bool
	IpServer   string
	MasterPort string
	RpcPort    string
	RpcServerState bool = true
    Db *sql.DB
    Dbengine *xorm.Engine
)
type ResData struct {
	Code int
	Msg string
	Data map[string]interface{}
}
func Usage() {
	flag.PrintDefaults()
}
func AddConns(conn net.Conn)  {
	Conns[conn.RemoteAddr().String()]=conn
	AddrsUserOrGroup[conn.RemoteAddr().String()]["uid"]="0"
	AddrsUserOrGroup[conn.RemoteAddr().String()]["group"]=make(map[int]string)
}
//广播所有除了当前连接
func SendToMsg(conn net.Conn,msg []byte){
	for caddr := range Conns {
		if strings.Compare(caddr,conn.RemoteAddr().String()) !=0 {
			if conn,ok :=Conns[caddr];ok {
				conn.Write(protocols.Packet(msg))
			}
		}
	}
}
//广播所有
func SendToAll(msg []byte){
	for caddr := range Conns {
		if conn,ok :=Conns[caddr];ok {
			conn.Write(protocols.Packet(msg))
		}
	}
}
//发给当前连接
func SendToConn(conn net.Conn,msg []byte)  {
	conn.Write(protocols.Packet(msg))
}
//发送给某个用户信息
func SendToUid(uid string,msg []byte){
	if conn,ok :=UserConns[uid];ok {
		conn.Write(protocols.Packet(msg))
	}
}
//发送给多个用户信息
func SendToUids(uids []string,msg []byte){
	for _, uid := range uids {
		if conn,ok :=UserConns[uid];ok {
			conn.Write(protocols.Packet(msg))
		}
	}
}

//发送给组成员
func SendToGroup(groupid string,msg []byte)  {
	if uids,ok :=GroupConns[groupid];ok {
		for _,conn := range uids{
			conn.Write(protocols.Packet(msg))
		}
	}
}
//关闭某个用户
func KickUid(uid string) {
	if conn,ok :=UserConns[uid];ok {
		KickConn(conn)
		conn.Close()
		delete(UserConns, uid)
	}
}
//关闭某个连接
func KickConn(conn net.Conn) {
	if ugar,ok :=AddrsUserOrGroup[conn.RemoteAddr().String()];ok {
		if ugar["uid"].(string) != "0" {
			delete(UserConns,ugar["uid"].(string))
		}
		if reflect.TypeOf(ugar["group"]).Kind() == reflect.Map {
			for _,groupid := range ugar["group"].(map[int]string){
				UnBindUidGroupConns(groupid, conn.RemoteAddr().String())
				UnBindUidGroupConns(groupid, ugar["uid"].(string))
			}
		}
		delete(AddrsUserOrGroup,conn.RemoteAddr().String())
	}
	conn.Close()
	delete(Conns,conn.RemoteAddr().String())

}
//绑定某个用户
func BindUid(conn net.Conn,uid string){
	UserConns[uid]=conn
}
//解绑某个用户
func UnBindUid(uid string){
	delete(UserConns,uid)
}
//绑定用户到小组里/uid可以为临时
func BindGroupConns(groupid string,uid string,conn net.Conn) {
	GroupConns[groupid][uid] = conn
}
//解绑用户在组里面
func UnBindUidGroupConns(groupid string,uid string)  {
	if uids,ok :=GroupConns[groupid];ok {
		if _,ok :=uids[uid];ok {
			delete(GroupConns[groupid],uid)
		}
	}

}
//解绑整个组
func UnBindGroup(groupid string){
	delete(GroupConns,groupid)
}
//用户是否在线
func UidIsOnline(uid string) bool {
	if _,ok :=UserConns[uid];ok {
		return true
	}
	return false
}
//异常恢复
func Grecover() {
	if r := recover(); r != nil {
		Log("Recovered in f", r)
	}
}
