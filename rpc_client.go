package main

import (
	"fmt"
	"gotest/common"
	"gotest/protocols"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"os"
	"time"
)

type RepRpc struct {
	C string
	A string
	Args interface{}
}
type Args struct {
	A, B int
}

type quo struct {
	Quo, Rem int
}

var clientarr=make(map[int]*rpc.Client)

func main() {
	service := "127.0.0.1:1234"
	service2 := "127.0.0.1:1235"
	tcpaddr, err := net.ResolveTCPAddr("tcp4", service2);
	tcpconn, err := net.DialTCP("tcp",nil, tcpaddr)
	if err !=nil {
		common.Log("dial error2==:", err)
		os.Exit(1)
	}
	fmt.Println("=======start=====")
	go rpc_listen(service)
	tcpconn.Write(protocols.Packet([]byte("ack-heartbeat")))
	clientConnection(tcpconn,10)
}
func rpc_listen(ipaddr string){
	client, err := jsonrpc.Dial("tcp", ipaddr)
	if err != nil {
		common.Log("dial error:", err)
		os.Exit(1)
	}
	clientarr[0]=client
	args := Args{1, 2}
	var reply int
	err = client.Call("Arith.Muliply", args, &reply)
	if err != nil {
		common.Log("Arith.Muliply call error:", err)
		os.Exit(1)
	}
	common.Log("the arith.mutiply is :", reply)
	var quto quo
	err = client.Call("Arith.Divide", args, &quto)
	if err != nil {
		common.Log("Arith.Divide call error:", err)
		os.Exit(1)
	}
	common.Log("the arith.devide is :", quto.Quo, quto.Rem)
	rargs := RepRpc{"cc","an","bbb"}
	var reps common.ResData
	err = client.Call("RpcRouter.RpcAccept", rargs, &reps)
	if err != nil {
		common.Log("RpcRouter.RpcAccept call error:", err)
		os.Exit(1)
	}
	common.Log("the RpcRouter.RpcAccept is :", reps.Code, reps.Msg,reps.Data)
	//call :=<-client.Go("Rpcroute.Rpcaccept", rargs, &reps,nil).Done
	//call.Reply
	//client.Close()


}
//长连接入口
func clientConnection(conn *net.TCPConn,timeout int) {

	Data := make([]byte, 0)
	messnager := make(chan []byte,2)
	buffer := make([]byte, 2048)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			common.Log(conn.RemoteAddr().String(), " connection error: ", err)
			return
		}
		args := Args{1, 2}
		var reply int
		err = clientarr[0].Call("Arith.Muliply", args, &reply)
		if err != nil {
			common.Log("=======222=====Arith.Muliply call error========kkkkkk:", err)
		}
		common.Log("the==== arith.mutiply is===== :", reply)
		//心跳计时
		go HeartBeat(conn,messnager,timeout)
		Data = protocols.Unpack(append(Data, buffer[:n]...), messnager)
		//检测每次Client是否有数据传来
		//go GravelChannels(Data,messnager)
		common.Log( "receive data length:",n)
		common.Log("receive data string========:", string(Data))
	}
}
//心跳计时，根据GravelChannel判断Client是否在设定时间内发来信息
func HeartBeat(conn *net.TCPConn, readerChannel chan []byte,timeout int) {
	select {
	case fk := <-readerChannel:
		common.Log(conn.RemoteAddr().String(), "receive data string:", string(fk))
		conn.SetDeadline(time.Now().Add(time.Duration(timeout) * time.Second))
		go func() {
			time.Sleep(2 * time.Second)
			conn.Write(protocols.Packet([]byte("ack-heartbeat")))
		}()
		break
	case <-time.After(time.Second*5):
		common.Log("It's really weird to get Nothing!!!")
		conn.Close()
	}

}

func GravelChannels(n []byte,mess chan byte){
	for _ , v := range n{
		mess <- v
	}
	close(mess)
}
