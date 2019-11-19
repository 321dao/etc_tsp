package term

import (
	"fmt"
	"net"
	"time"

	"audotsp/utils"
)

type Terminal struct {
	imei      string
	iccid     string
	vin       string
	pepsver   uint32
	loginTime time.Time
	seqNum    uint16
	phoneNum  []byte
	Conn      *net.Conn
}

type UpdateConf struct {
	Url        string
	DialName   string
	DialUser   string
	DialPasswd string
	Ip         string
	TcpPort    uint16
	UdpPort    uint16
	ManulId    []byte
	HwVersion  string
	FwVersion  string
	ConTime    uint16
}

const (
	protoHeader byte = 0x7e

	register    uint16 = 0x0100
	registerAck uint16 = 0x8100
	unregister  uint16 = 0x0003
	login       uint16 = 0x0102
	heartbeat   uint16 = 0x0002
	gpsinfo     uint16 = 0x0200
	platAck     uint16 = 0x8001
)

func init() {
	fmt.Println("hello module init function")
}

func (uc UpdateConf) toByteArray() []byte {
	retByte := make([]byte, 0)
	retByte = append(retByte, uc.Url...)
	retByte = append(retByte, uc.DialName...)
	retByte = append(retByte, uc.DialUser...)
	retByte = append(retByte, uc.DialPasswd...)
	retByte = append(retByte, uc.Ip...)

	retByte = append(retByte, utils.Word2Bytes(uc.TcpPort)...)
	retByte = append(retByte, utils.Word2Bytes(uc.UdpPort)...)
	retByte = append(retByte, utils.Word2Bytes(uc.UdpPort)...)

	return retByte
}

func (t Terminal) MakeFrame(cmd uint16, phone []byte, seq uint16, apdu []byte) []byte {
	data := make([]byte, 0)
	tempbytes := utils.Word2Bytes(cmd)
	data = append(data, tempbytes...)
	datalen := uint16(len(apdu))
	tempbytes = utils.Word2Bytes(datalen)
	data = append(data, tempbytes...)

	data = append(data, phone...)

	tempbytes = utils.Word2Bytes(seq)
	data = append(data, tempbytes...)

	data = append(data, apdu...)

	csdata := byte(t.checkSum(data[:]))
	data = append(data, csdata)

	//转义

	//添加头尾
	var tmpdata []byte = []byte{0x7e}
	data = append(tmpdata, data...)
	data = append(data, 0x7e)

	return data
}

func (t Terminal) makeApduRegisterAck(res uint8, authkey string) []byte {
	data := make([]byte, 0)
	tempbytes := utils.Word2Bytes(t.seqNum)
	data = append(data, tempbytes...)

	data = append(data, res)

	for _, item := range authkey {
		data = append(data, byte(item))
	}

	return data
}

func (t Terminal) makeApduCommonAck(cmdid uint16, res byte) []byte {
	data := make([]byte, 0)
	tempbytes := utils.Word2Bytes(t.seqNum)
	data = append(data, tempbytes...)

	tempbytes = utils.Word2Bytes(cmdid)
	data = append(data, tempbytes...)

	data = append(data, res)

	fmt.Println("apdu:", data)
	return data
}

func (t Terminal) DataFilter(data []byte) int {
	//--------------------------------------------------
	//int iRet = 0;
	// static int curLen=0;
	fmt.Printf("len = %d,data[0]=0x%X.\n", len(data), data[0])
	if data[0] == protoHeader {
		fmt.Println("find start.")
		var endindex int = -1
		for i := 1; i < len(data); i++ {
			if data[i] == protoHeader {
				fmt.Println("find end.")
				endindex = i
				break
			}
		}

		if endindex > 0 {
			data = data[:endindex+1]
		}

		return len(data)
	} else {
		return -2
	}
}

func (t Terminal) FrameHandle(data []byte) []byte {
	//bodylen := data[1] - 5
	//cmdid := data[5]
	cmdid := utils.Bytes2Word(data[1:3])
	t.phoneNum = data[5:11]
	t.seqNum = utils.Bytes2Word(data[11:13])
	fmt.Println("cmdid:", cmdid)
	len := len(data)
	return t.apduHandle(cmdid, data[13:len-2])
}

func (t Terminal) apduHandle(cmdType uint16, apdu []byte) []byte {
	switch cmdType {
	case register:
		fmt.Println("rcv register.")
		apduack := t.makeApduRegisterAck(0, "AACAB")
		sendBuf := t.MakeFrame(registerAck, t.phoneNum, t.seqNum, apduack)
		return sendBuf
	case login:
		fmt.Println("rcv login.")
		apduack := t.makeApduCommonAck(cmdType, 0)
		sendBuf := t.MakeFrame(platAck, t.phoneNum, t.seqNum, apduack)

		return sendBuf
		//return []byte{}
	case heartbeat:
		fmt.Println("rcv heartbeat.")
		apduack := t.makeApduCommonAck(cmdType, 0)
		sendBuf := t.MakeFrame(platAck, t.phoneNum, t.seqNum, apduack)

		return sendBuf
	case gpsinfo:
		fmt.Println("rcv gpsinfo.")
		apduack := t.makeApduCommonAck(cmdType, 0)
		sendBuf := t.MakeFrame(platAck, t.phoneNum, t.seqNum, apduack)

		return sendBuf
	}

	return nil
}

func (t Terminal) paramHandle(id byte, data []byte) int {
	return 0
}

func (t Terminal) checkSum(data []byte) byte {
	var sum byte = 0
	for _, itemdata := range data {
		sum ^= itemdata
	}
	return sum
}

func (t Terminal) GetImei() string {
	return t.imei
}

func (t Terminal) GetIccid() string {
	return t.iccid
}
