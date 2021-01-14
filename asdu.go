package iec104

import (
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"time"
)

//ASDU 应用服务数据单元
type ASDU struct {
	TypeID        byte    //类型标识
	Sequence      bool    //是否连续
	Length        byte    //可变结构限定词
	Cause         uint16  //传输原因
	PublicAddress uint16  //公共地址
	Ts            float64 //毫秒级时间戳
}

//数据类型
const (
	//MSpNa1 不带游标的单点遥信，3个字节的地址，1个字节的值
	MSpNa1 = 1
	//MDpNa1 不带时标的双点遥信，每个遥信占1个字节
	MDpNa1 = 3
	//MMeNc1 带品质描述的测量值，每个遥测值占3个字节
	MMeNa1 = 9
	//MMeNc1 带品质描述的浮点值，每个遥测值占5个字节
	MMeNc1 = 13
	//MItNa1 电度总量,每个遥脉值占5个字节
	MItNa1 = 15
	//MSpTb1 带游标的单点遥信，3个字节的地址，1个字节的值，7个字节短时标
	MSpTb1 = 30
	//MEiNA1 初始化结束
	MEiNA1 = 70
	//CIcNa1 总召唤
	CIcNa1 = 100
	//CCiNa1 电度总召唤
	CCiNa1 = 101
)

// ParseASDU 解析asdu
func (asdu *ASDU) ParseASDU(asduBytes []byte) (signals []*Signal, err error) {
	signals = make([]*Signal, 0, 0)
	if asduBytes == nil || len(asduBytes) < 4 {
		err = fmt.Errorf("asdu[%X]非法", asduBytes)
		return
	}
	asdu.TypeID = asduBytes[0]
	//数据是否连续
	asdu.Sequence, asdu.Length = asdu.ParseVariable(asduBytes[1])
	var firstAddress uint32

	asdu.Cause = binary.LittleEndian.Uint16([]byte{asduBytes[2], asduBytes[3]})
	asdu.PublicAddress = binary.LittleEndian.Uint16([]byte{asduBytes[4], asduBytes[5]})

	if asdu.Sequence {
		firstAddress = binary.LittleEndian.Uint32([]byte{asduBytes[6], asduBytes[7], asduBytes[8], 0x00})
	}
	for i := 0; i < int(asdu.Length); i++ {
		s := new(Signal)
		s.TypeID = uint(asdu.TypeID)
		if asdu.Sequence {
			s.Address = firstAddress
			firstAddress++
		}
		switch asdu.TypeID {
		case MSpNa1, MDpNa1:
			size := 4
			if asdu.Sequence {
				s.Value = float64(asduBytes[9+i])
			} else {
				s.Address = binary.LittleEndian.Uint32([]byte{asduBytes[6+i*size], asduBytes[6+i*size+1], asduBytes[6+i*size+2], 0x00})
				s.Value = float64(asduBytes[6+i*size+3])
			}
		case MMeNa1:
			size := 6
			if asdu.Sequence {
				size := 3
				s.Value = float64(binary.LittleEndian.Uint16([]byte{asduBytes[9+i*size], asduBytes[9+i*size+1]}))
				s.Quality = asduBytes[9+i*size+2]
			} else {
				s.Address = binary.LittleEndian.Uint32([]byte{asduBytes[6+i*size], asduBytes[6+i*size+1], asduBytes[6+i*size+2], 0x00})
				s.Value = float64(binary.LittleEndian.Uint16([]byte{asduBytes[6+i*size+3], asduBytes[6+i*size+4]}))
				s.Quality = asduBytes[6+i*size+5]
			}
		case MMeNc1:
			size := 8
			if asdu.Sequence {
				size := 5
				s.Value = float64(math.Float32frombits(binary.LittleEndian.Uint32([]byte{asduBytes[9+i*size], asduBytes[9+i*size+1],
					asduBytes[9+i*size+2], asduBytes[9+i*size+3]})))
				s.Quality = asduBytes[9+i*size+4]
			} else {
				s.Address = binary.LittleEndian.Uint32([]byte{asduBytes[6+i*size], asduBytes[6+i*size+1], asduBytes[6+i*size+2], 0x00})
				s.Value = float64(math.Float32frombits(binary.LittleEndian.Uint32([]byte{asduBytes[6+i*size+3], asduBytes[9+i*size+4],
					asduBytes[9+i*size+5], asduBytes[9+i*size+6]})))
				s.Quality = asduBytes[6+i*size+7]
			}
		case MItNa1:
			size := 8
			if asdu.Sequence {
				size := 5
				s.Value = float64(binary.LittleEndian.Uint32([]byte{asduBytes[9+i*size], asduBytes[9+i*size+1],
					asduBytes[9+i*size+2], asduBytes[9+i*size+3]}))
			} else {
				s.Address = binary.LittleEndian.Uint32([]byte{asduBytes[6+i*size], asduBytes[6+i*size+1], asduBytes[6+i*size+2], 0x00})
				s.Value = float64(binary.LittleEndian.Uint32([]byte{asduBytes[6+i*size+3], asduBytes[9+i*size+4],
					asduBytes[9+i*size+5], asduBytes[9+i*size+6]}))
			}
		case MSpTb1:
			size := 11
			s.Address = binary.LittleEndian.Uint32([]byte{asduBytes[6+i*size], asduBytes[6+i*size+1], asduBytes[6+i*size+2], 0x00})
			s.Value = float64(asduBytes[6+i*size+3])
			s.Ts = asdu.ParseTime(asduBytes[6+i*size+4 : 6+i*size+11])
		case CIcNa1, CCiNa1, MEiNA1:
		default:
			log.Fatalln("暂不支持的数据类型:", asdu.TypeID)
		}
		signals = append(signals, s)
	}
	return
}

// ParseVariable 解析asdu可变结构限定词
func (asdu *ASDU) ParseVariable(b byte) (sq bool, length byte) {
	//最高位是否为1
	sq = b&128>>7 == 1
	if sq {
		length = b - 1<<7
		return
	}
	length = b
	return
}

// ParseTime 解析asdu中7个字节时表,转为带毫秒的时间戳
func (asdu *ASDU) ParseTime(asduBytes []byte) float64 {
	if len(asduBytes) != 7 {
		return 0
	}
	milliseconds := binary.LittleEndian.Uint16([]byte{asduBytes[0], asduBytes[1]})
	nanosecond := (int(milliseconds) % 1000) * 1000000
	second := int(milliseconds / 1000)
	minute := int(asduBytes[2])
	hour := int(asduBytes[3])
	day := int(asduBytes[4])
	month := int(asduBytes[5])
	year := int(asduBytes[6]) + 2000
	return float64(time.Date(year, time.Month(month), day, hour, minute, second, nanosecond, time.Local).Unix()) + float64(nanosecond)/1000000000.0
}
