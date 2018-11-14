package protocols

import (
	"bytes"
	"encoding/binary"
	"strings"
)

//通讯协议处理，主要处理封包和解包的过程
const (
	ConstHeader         = "www.uf101.com"   //包头
	ConstHeaderLength   = 13                //包头长度
	ConstPackageEof     = "/r/n/r/n"        //设置EOF
	ConstPackageLength  = 8                 //EOF长度
	ConstSaveDataLength = 4                 //存储包体长度
	ConstOpenEofCheck   = true              //开启EOFF检测
)
//封包
func Packet(message []byte) []byte {
	if(ConstOpenEofCheck == false){
		return append(append([]byte(ConstHeader), IntToBytes(len(message))...), message...)
	}
	return append(append(append([]byte(ConstHeader), IntToBytes(len(message)+ConstPackageLength)...), message...),ConstPackageEof...)
}
//解包
func Unpack(buffer []byte, readerChannel chan []byte) []byte {
	length := len(buffer)

	var i int
	for i = 0; i < length; i = i + 1 {
		if length < i+ConstHeaderLength+ConstSaveDataLength {
			break
		}

		if ConstOpenEofCheck == true {
			if length < i+ConstHeaderLength+ConstPackageLength+ConstPackageLength {
				break
			}
		}
		if string(buffer[i:i+ConstHeaderLength]) == ConstHeader {
			messageLength := BytesToInt(buffer[i+ConstHeaderLength : i+ConstHeaderLength+ConstSaveDataLength])
			if length < i+ConstHeaderLength+ConstSaveDataLength+messageLength {
				break
			}
			if ConstOpenEofCheck == false {
				data := buffer[i+ConstHeaderLength+ConstSaveDataLength : i+ConstHeaderLength+ConstSaveDataLength+messageLength]
				readerChannel <- data
			} else if ConstOpenEofCheck == true {
				eof := buffer[i+ConstHeaderLength+ConstSaveDataLength+messageLength-ConstPackageLength : i+ConstHeaderLength+ConstSaveDataLength+messageLength]
				if(strings.Compare(string(eof),ConstPackageEof) !=0){
					break;
				}
				data := buffer[i+ConstHeaderLength+ConstSaveDataLength : i+ConstHeaderLength+ConstSaveDataLength+messageLength-ConstPackageLength]
				readerChannel <- data
			}

			i += ConstHeaderLength + ConstSaveDataLength + messageLength  - 1
		}
	}

	if i == length {
		return make([]byte, 0)
	}
	return buffer[i:]
}
//整形转换成字节
func IntToBytes(n int) []byte {
	x := int32(n)

	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, x)
	return bytesBuffer.Bytes()
}

//字节转换成整形
func BytesToInt(b []byte) int {
	bytesBuffer := bytes.NewBuffer(b)

	var x int32
	binary.Read(bytesBuffer, binary.BigEndian, &x)

	return int(x)
}
