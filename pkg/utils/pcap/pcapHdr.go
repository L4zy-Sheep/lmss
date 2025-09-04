package pcap

import "unsafe"

const (
	PCAPHDRLEN          = int(unsafe.Sizeof(PcapHdr{}))
	PKTHDRLEN           = int(unsafe.Sizeof(PktHdr{}))
	BIGENDIAN    uint32 = 0xd4c3b2a1
	LITTLEENDIAN uint32 = 0xa1b2c3d4
)

type PcapHdr struct {
	Magic    uint32 //0xA1B2C3D4:小端序，0xD4C3B2A1:大端序
	Major    uint16 //主版本，一般为0x0200
	Minor    uint16 //次版本，一般为0x0400
	ThisZone uint32 //标准时间，一般为GMT，这个值全0
	SigFigs  uint32 //时间戳的精度，一般全0
	SnapLen  uint32 //捕获的最大数据包的长度
	LinkType uint32 //链路类型，1：以太网
}

type PktHdr struct {
	TimeStampHigh uint32 //时间戳高位，精确到second
	TimeStampLow  uint32 //时间戳地位，精确到microsecond
	CapLen        uint32 //当前数据区的长度，由此得到下一个数据区
	Len           uint32 //实际数据包长度，一般不大于CapLen，多数情况下和CapLen一样，暂不确定实际区别
}
