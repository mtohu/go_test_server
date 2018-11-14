package common

import (
	"gotest/protocols"
	"net"
)

var (
	ServerConns  = make(map[string]net.Conn)//map[conn.RemoteAddr().String()]conn
	ForwardUser  = make(map[string]string)//map[uid]conn.RemoteAddr().String()
	ForwardGroup = make(map[string]string)//map[groupid]conn.RemoteAddr().String()
	FH          bool
	ForwarIpServer   string
	ForwarMasterPort string
)

func ForwardSendToGroup(groupid string,msg []byte) {
	if addrs,ok :=ForwardGroup[groupid];ok {
		if conn,ok :=ServerConns[addrs];ok {
			conn.Write(protocols.Packet(msg))
		}
	}
}

func ForwardSendToUid(uid string,msg []byte) {
	if addrs,ok :=ForwardUser[uid];ok {
		if conn,ok :=ServerConns[addrs];ok {
			conn.Write(protocols.Packet(msg))
		}
	}
}

func ForwardSendToUids(uids []string,msg []byte) {
	for _, uid := range uids {
		if addrs, ok := ForwardUser[uid]; ok {
			if conn, ok := ServerConns[addrs]; ok {
				conn.Write(protocols.Packet(msg))
			}
		}
	}
}
