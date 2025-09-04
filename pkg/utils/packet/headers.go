package packet

type EthernetHdr struct {
	Dst          [6]byte
	Src          [6]byte
	ProtocolType uint16
}
type ArpHdr struct {
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

type Ipv4Hdr struct {
	Version_HdrLen uint8   //4位版本号 and 4位ip数据包头部字节数
	TOS            uint8   //服务类型
	TotalLen       uint16  //总长度
	Identification uint16  //标识符
	FragmentOffset uint16  //3位分片标识符和13位分片偏移
	TTL            uint8   //最大生存时间
	Protocol       uint8   //协议
	CheckSum       uint16  //校验和
	Src            [4]byte //
	Dst            [4]byte //
	Options        []byte  //可选项40bit
}

// func (ip4 *Ipv4Hdr)
type Ipv6Hdr struct {
	Version_TrafficClass_flowLable uint32 //4位版本号，8位流量类别，20位流标签
	payloadLen                     uint16 //有效载荷长度，超过65535时设置为0，然后用扩展报头中的超大有效载荷表示
	NextHdr                        uint8  //下一层的头部
	HopLimit                       uint8  //类似TTL
	Src                            [8]uint16
	Dst                            [8]uint16
	ExtensionHdrs                  []byte
}

type TcpHdr struct {
	Src            uint16
	Dst            uint16
	Seq            uint32
	Ack            uint32
	HdrLen_reserve uint8
	Sign           uint8
	WindowSize     uint16
	CheckSum       uint16
	UrgentPointer  uint16
	Options        []uint32
}
type UdpHdr struct {
	Src      uint16
	Dst      uint16
	Len      uint16
	CheckSum uint16
}

type IcmpPkt struct {
	Type       uint8  //8:echo ,0:reply
	Code       uint8  //default 0
	CheckSum   uint16 //
	Identifier uint16 //标识码
	Sequence   uint16 //序号
	Data       []byte //可选数据40bit
}
