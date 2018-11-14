package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"net/rpc"
	"os"
	"strings"
	"time"
)

type Argsss struct {
	A, B int
}

func checkError(err error) {
	if err != nil {
		fmt.Fprint(os.Stderr, "Usage: %s", err.Error())
		os.Exit(1)
	}
}

type Quotients struct {
	Quo, Rem int
}

type Ariths int
func (t *Ariths) Muliply(args *Argsss, reply *int) error {
	*reply = args.A * args.B
	return nil
}

func (t *Ariths) Divide(args *Argsss, quo *Quotients) error {
	if args.B == 0 {
		return errors.New("divide by zero")
	}
	quo.Quo = args.A * args.B
	quo.Rem = args.A / args.B
	return nil
}

var connss = make(map[string]net.Conn)
func main(){
	ariths := new(Ariths)
	rpc.Register(ariths)

	tcpAddr, err := net.ResolveTCPAddr("tcp", ":12345")
	checkError(err)

	Listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)

	for {
		conn, err := Listener.Accept()
		if err != nil {
			fmt.Fprint(os.Stderr, "accept err: %s", err.Error())
			continue
		}
		connip := conn.RemoteAddr().String()
		connlocalip :=conn.LocalAddr().String()
		fmt.Println(connlocalip+"===="+connip)
		mstringip := strings.Join([]string{connip,connlocalip},"-")
		connss[mstringip] = conn
		log.Println(connip,"已经建立连接")
		handleConnection(conn,5)
		//jsonrpc.ServeConn(conn)
	}

}
//长连接入口
func handleConnections(conn net.Conn,timeout int) {

	buffer := make([]byte, 2048)
	for {
		n, err := conn.Read(buffer)

		if err != nil {
			Log(conn.RemoteAddr().String(), " connection error: ", err)
			return
		}
		Data :=(buffer[:n])
		messnager := make(chan byte)
		//postda :=make(chan byte)
		//心跳计时
		go HeartBeatings(conn,messnager,timeout)
		//检测每次Client是否有数据传来
		go GravelChannel(Data,messnager)
		Log( "receive data length:",n)
		Log(conn.RemoteAddr().String(), "receive data string:", string(Data))

	}
}
//心跳计时，根据GravelChannel判断Client是否在设定时间内发来信息
func HeartBeatings(conn net.Conn, readerChannel chan byte,timeout int) {
	select {
	case fk := <-readerChannel:
		Log(conn.RemoteAddr().String(), "receive data string:", string(fk))
		conn.SetDeadline(time.Now().Add(time.Duration(timeout) * time.Second))
		break
	case <-time.After(time.Second*5):
		Log("It's really weird to get Nothing!!!")
		conn.Close()
	}

}

func GravelChannel(n []byte,mess chan byte){
	for _ , v := range n{
		mess <- v
	}
	close(mess)
}

func Log(v ...interface{}) {
	log.Println(v...)
}

