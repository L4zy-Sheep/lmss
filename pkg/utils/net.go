package utils

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"lmss/pkg/log"
	"math/rand/v2"
	"net"
	"time"
	"unsafe"
)

const (
	ARPPROTOCOL  = 0x0806
	IPPROTOCOL   = 0x06
	IPV4PROTOCOL = 0x0800
)

type EthernetPkt struct {
	DstAddr   [6]byte
	SrcAddr   [6]byte
	EtherType uint16
}

type ArpPkt struct {
	HardWareType  uint16 //1:以太网
	ProtocolType  uint16
	HardWareSize  uint8
	ProtocolSize  uint8
	Opcode        uint16 //1:arp_req,2:arp_res,3:rarp_req,4:rarp_res
	SenderMacAddr [6]byte
	SenderIpAddr  [4]byte
	TargetMacAddr [6]byte //arp广播时设置为0
	TargetIpAddr  [4]byte
}
type IPPkt struct {
	version_HdrLen uint8
	//...
}

type IcmpPkt struct {
	Type       uint8  //8:echo ,0:reply
	Code       uint8  //default 0
	Checksum   uint16 //
	Identifier uint16 //标识码
	Sequence   uint16 //序号
	//Data       []byte //可选数据
}

func Icmp(host string, timeout uint8) {
	var buf = new(bytes.Buffer)
	res := make([]byte, 32)
	var conn net.Conn
	icmpPkt := IcmpPkt{8, 0, 0, uint16(rand.Uint()), 1}
	err := binary.Write(buf, binary.BigEndian, icmpPkt)
	if err != nil {
		log.Error(err.Error())
		goto end
	}
	icmpPkt.Checksum = Checksum(buf.Bytes())
	conn, err = net.Dial("ip4:icmp", host)
	if err != nil {
		log.Error(err.Error())
		goto end
	}
	defer conn.Close()
	conn.SetReadDeadline(time.Now().Add(time.Duration(timeout) * time.Millisecond))
	err = binary.Write(conn, binary.BigEndian, icmpPkt)
	if err != nil {
		log.Error(err.Error())
	}
	_, err = conn.Read(res)
	if err != nil {
		goto end
	}
	if res[21] == 0 && binary.BigEndian.Uint16(res[24:26]) == icmpPkt.Identifier {
		log.Success(host + " alive")
	}
end:
}

type TcpPkt struct {
	SrcPort       uint16
	DstPort       uint16
	Seq           uint32
	AckNum        uint32
	HdrLen_Sign   uint16 //前四位	urg,ack,psh,rst,syn,fin
	WindowSize    uint16
	Checksum      uint16
	UrgentPointer uint16
}

func (tcp *TcpPkt) Fin() *TcpPkt {
	tcp.HdrLen_Sign |= 1 << 15
	return tcp
}
func (tcp *TcpPkt) Syn() *TcpPkt {
	tcp.HdrLen_Sign |= 1 << 14
	return tcp
}
func (tcp *TcpPkt) Ack() *TcpPkt {
	tcp.HdrLen_Sign |= 1 << 11
	return tcp
}

func Syn(host string, port int) {
	var buf = new(bytes.Buffer)
	syn := TcpPkt{
		SrcPort:    uint16(rand.Uint()),
		DstPort:    uint16(port),
		Seq:        1,
		AckNum:     0,
		WindowSize: 0xff,
	}
	syn.Syn()
	buf := binary.Write(buf, binary.BigEndian, syn)
	syn.Checksum = Checksum(buf.Bytes())
	conn, err := net.Dial("ip4", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		goto end
	}
	defer conn.Close()
	err = binary.Write(conn, binary.BigEndian, syn)
	if err != nil {
		goto end
	}
	_, err = conn.Read(make([]byte, 1))

}

func Arp(host string) bool {
	//log.Info(fmt.Sprintf("%d - %d - %d", syscall.Geteuid(), syscall.Getuid(), syscall.Getegid()))
	//if syscall.Geteuid() != 0 {
	//	panic("arp scan requires root privilege")
	//}
	/*	var macAddr [6]byte
		var ipAddr [4]byte
		var targetIpaddr [4]byte
		copy(targetIpaddr[:], net.ParseIP(host).To4())
		interfaces, err := net.Interfaces()
		if err != nil {
			panic("get interfaces failed")
		}
		for _, iface := range interfaces {
			if iface.Flags&net.FlagUp != 0 && iface.Flags&net.FlagLoopback == 0 && iface.Flags&net.FlagBroadcast != 0 {
				copy(macAddr[:], iface.HardwareAddr)
				addrs, err := iface.Addrs()
				if err != nil {
					panic("get interfaces failed")
				}
				for _, addr := range addrs {
					if ipNet := addr.(*net.IPNet); ipNet.IP.To4() != nil {
						copy(ipAddr[:], ipNet.IP.To4())
						break
					}
				}
				sockAddr := &syscall.SockaddrInet4{
					Addr: ipAddr,
					Port: 0,
				}
				etherPkt := EthernetPkt{
					[6]byte{0, 0, 0, 0, 0, 0},
					macAddr,
					ARPPROTOCOL,
				}
				arpPkt := ArpPkt{
					1,
					IPV4PROTOCOL,
					6, //mac地址字节数
					4, //ip地址字节数,ipv4=4
					1,
					macAddr,
					ipAddr,
					[6]byte{0, 0, 0, 0, 0, 0},
					targetIpaddr,
				}
				buf := new(bytes.Buffer)
				binary.Write(buf, binary.BigEndian, etherPkt)
				binary.Write(buf, binary.BigEndian, arpPkt)
				syscall.WSAStartup(2, &syscall.WSAData{})
				sock, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_RAW, syscall.IPPROTO_IP)
				if err != nil {
					panic(err)
				}
				err = syscall.WSASend(sock, syscall.WSABuf{uint32(len(buf.Bytes())), &(buf.Bytes()[0])}, 0, sockAddr)
				if err != nil {
					panic(err)
				}
			}
		}*/
	return false
}

func Checksum(data []byte) uint16 {
	var sum uint32
	for i := 0; i < len(data); i += 2 {
		sum += uint32(data[i])<<8 + uint32(data[i+1])
		for sum > 0xffff {
			sum = (sum >> 16) + (sum & 0xffff)
		}
	}
	return ^uint16(sum)

}
