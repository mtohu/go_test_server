package main

import (
	"errors"
	"flag"
	"fmt"
	"gotest/common"
	"gotest/protocols"
	"gotest/routes"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"os"
	"strings"
	"time"
)

type Argss struct {
	A, B int
}
type Quotient struct {
	Quo, Rem int
}
type Arith int
func (t *Arith) Muliply(args *Argss, reply *int) error {
	*reply = args.A * args.B
	return nil
}

func (t *Arith) Divide(args *Argss, quo *Quotient) error {
	if args.B == 0 {
		return errors.New("divide by zero")
	}
	quo.Quo = args.A * args.B
	quo.Rem = args.A / args.B
	return nil
}
func init() {
	flag.BoolVar(&common.H, "h", false, "this help")
	flag.StringVar(&common.IpServer,"s","127.0.0.1","ip server")
	flag.StringVar(&common.MasterPort,"mp","1235","master port")
	flag.StringVar(&common.RpcPort,"rp","1234","rpc port")
	flag.Usage = common.Usage
}

func show(args ...interface{}){
	for k, v := range args {
		fmt.Println(k,v)
	}
}

func main() {
	flag.Parse()
	if common.H {
		flag.Usage()
	}
	/*var ss string = ""
	if strings.Trim(ss,"") == ""{
		common.Log("aaa",len(ss))
	}*/
	//var countryCapitalMap map[string]string
	//common.Log(len(countryCapitalMap))
	//slice:=make([]interface{},0)
	//slice=append(slice,5)
	//show(slice...)
	//os.Exit(1)
	go rpcListener()
	masterListener()
}
func rpcListener()  {
	byteaddr :=[]byte(common.IpServer+":"+common.RpcPort)
	serveraddr := string(byteaddr)
	common.Log("==================rpc addr=====",serveraddr)
	arith := new(Arith)
	rpcrouter := new(routes.RpcRouter)
	rpc.Register(arith)
	rpc.Register(rpcrouter)

	tcpAddr, err := net.ResolveTCPAddr("tcp", serveraddr)
	if(err !=nil) {
		common.CheckError(err)
	}
	Listener, err := net.ListenTCP("tcp", tcpAddr)
	if(err !=nil) {
		common.CheckError(err)
	}
	for {
		conn, err := Listener.Accept()
		if err != nil {
			common.Log(os.Stderr, "accept err: %s", err.Error())
			common.RpcServerState = false
			continue
		}
		common.RpcServerState = true
		go jsonrpc.ServeConn(conn)
	}
}
func masterListener()  {
	byteaddr :=[]byte(common.IpServer+":"+common.MasterPort)
	serveraddr := string(byteaddr)
	common.Log("==================master addr=====",serveraddr)
	tcpAddr, err := net.ResolveTCPAddr("tcp", serveraddr)
	if(err !=nil) {
		common.CheckError(err)
	}
	Listener, err := net.ListenTCP("tcp", tcpAddr)
	if(err !=nil) {
		common.CheckError(err)
	}
	for{
		conn, err := Listener.Accept()
		if err != nil {
			common.Log(os.Stderr, "accept err: %s", err.Error())
			continue
		}
		if(common.RpcServerState == false){
			go rpcListener()
		}
		connip := conn.RemoteAddr().String()
		common.AddConns(conn)
		common.Log(connip, "已经建立连接")
		go handleConnection(conn, 5)
	}
}
//长连接入口
func handleConnection(conn net.Conn,timeout int) {

	Data := make([]byte, 0)
	buffer := make([]byte, 2048)
	messnager := make(chan []byte,10)
	ackchan:= make(chan []byte,1)
	//心跳计时
	go HeartBeating(conn,ackchan,timeout)
	//任务
	go dispatch(conn,messnager,ackchan)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			common.Log(conn.RemoteAddr().String(), " connection read error: ", err)
			common.KickConn(conn)
			return
		}
		Data = protocols.Unpack(append(Data, buffer[:n]...), messnager)
		if(len(Data)>0){
			ackchan <- Data
		}
		common.Log( "receive data length:",n)
		common.Log(conn.RemoteAddr().String(), "receive data string=====:", len(Data),string(Data))
	}
}
func dispatch(conn net.Conn, msgChannel chan []byte,ackchan chan []byte){
	for{
		select {
		case msgData := <- msgChannel:
			ackchan <- msgData
			common.Log("kkkkk----------",string(msgData))
			if len(msgData) >2 && strings.Compare(string(msgData),"ack-heartbeat") != 0 {
				//go 业务处理
			}
		}
	}
}
//心跳计时，根据判断Client是否在设定时间内发来信息
func HeartBeating(conn net.Conn, readerChannel chan []byte,timeout int) {
	for{
		select {
		case fk := <-readerChannel:
			common.Log(conn.RemoteAddr().String(), "heart:", string(fk))
			conn.SetDeadline(time.Now().Add(time.Duration(timeout) * time.Second))
			go func() {
				time.Sleep(1 * time.Second)
				common.SendToConn(conn,[]byte("ack-heartbeat"))
			}()
		case <-time.After(time.Second*5):
			common.Log("It's really weird to get Nothing!!!")
			close(readerChannel)
			/*err := common.Db.Close()//关闭数据库链接
			if err !=nil {
				common.Log(" db close error")
			}*/
			common.KickConn(conn)
			break

		}

	}
}





