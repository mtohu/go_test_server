package main

import (
	"flag"
	"gotest/common"
	"gotest/protocols"
	"net"
	"os"
	"strings"
	"time"
)

func init() {
	flag.BoolVar(&common.H, "h", false, "this help")
	flag.StringVar(&common.IpServer,"s","127.0.0.1","ip server")
	flag.StringVar(&common.MasterPort,"mp","1238","master port")
	flag.StringVar(&common.RpcPort,"rp","1239","rpc port")
	flag.Usage = common.Usage
}

func main()  {
	flag.Parse()
	if common.H {
		flag.Usage()
	}
	//go rpcListener()
	forwardListener()
}
func forwardListener()  {
	byteaddr :=[]byte(common.IpServer+":"+common.MasterPort)
	serveraddr := string(byteaddr)
	common.Log("==================forward master addr=====",serveraddr)
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
			common.Log(os.Stderr, "forward accept err: %s", err.Error())
			continue
		}
		connip := conn.RemoteAddr().String()
		common.Conns[connip] = conn
		common.Log(connip, "forward 已经建立连接")
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
	go heartBeating(conn,ackchan,timeout)
	//任务
	go dispatch(conn,messnager,ackchan)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			common.Log(conn.RemoteAddr().String(), "forward connection read error: ", err)
			common.KickConn(conn)
			return
		}
		Data = protocols.Unpack(append(Data, buffer[:n]...), messnager)
		if(len(Data)>0){
			ackchan <- Data
		}
		common.Log( "forward receive data length:",n)
		common.Log(conn.RemoteAddr().String(), "forward receive data string=====:", len(Data),string(Data))
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
func heartBeating(conn net.Conn, readerChannel chan []byte,timeout int) {
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
			common.KickConn(conn)
			break

		}

	}
}
