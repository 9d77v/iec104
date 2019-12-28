package iec104

import (
	"encoding/binary"
	"fmt"
)

//StartFrame 起始符
const startFrame = 0x68

var (
	startDtAct = [4]byte{0x07, 0x00, 0x00, 0x00} //启动激活帧
	startDtCon = [4]byte{0x0b, 0x00, 0x00, 0x00} //启动确认帧
	testFrAct  = [4]byte{0x43, 0x00, 0x00, 0x00} //测试激活帧
	testFrCon  = [4]byte{0x83, 0x00, 0x00, 0x00} //测试确认帧
	stopDtAct  = [4]byte{0x13, 0x00, 0x00, 0x00} //停止激活帧
	stopDtCon  = [4]byte{0x23, 0x00, 0x00, 0x00} //停止确认帧
)

const (
	iFrame byte = 0
	sFrame byte = 1
	uFrame byte = 3
)

//APCI ..
type APCI struct {
	ApduLen int
	Ctr1    byte
	Ctr2    byte
	Ctr3    byte
	Ctr4    byte
}

//IFrame I帧
type IFrame struct {
	Send int16
	Recv int16
}

//SFrame S帧
type SFrame struct {
	Recv int16
}

//UFrame U帧
type UFrame struct {
	cmd [4]byte //激活确认命令
}

func convert4BytesToSlice(b [4]byte) []byte {
	return []byte{b[0], b[1], b[2], b[3]}
}

//convertIntToBytes 转换int为bytes,大端序
func convertIntToBytes(i int) []byte {
	bytes := make([]byte, 4, 4)
	binary.BigEndian.PutUint32(bytes, uint32(i))
	return bytes
}

//parseBigEndianUint16 转换大端Uint16
func parseBigEndianUInt16(i uint16) []byte {
	bytes := make([]byte, 2, 2)
	binary.BigEndian.PutUint16(bytes, i)
	return bytes
}

//parseLittleEndianUint16 转换小端Uint16
func parseLittleEndianUInt16(i uint16) []byte {
	bytes := make([]byte, 2, 2)
	binary.LittleEndian.PutUint16(bytes, i)
	return bytes
}

//convertBytes 转换发送数据
func convertBytes(data []byte) []byte {
	sendData := make([]byte, 0, 0)
	iBytes := parseBigEndianUInt16(uint16(len(data)))
	sendData = append(sendData, startFrame)
	sendData = append(sendData, iBytes[1])
	sendData = append(sendData, data...)
	return sendData
}

//ParseCtr 解析控制域
func (apci *APCI) ParseCtr() (byte, interface{}, error) {
	switch {
	case apci.Ctr1&1 == iFrame:
		//I帧
		t, f := apci.parseIFrame()
		return t, f, nil
	case apci.Ctr1&3 == sFrame:
		//S帧
		t, f := apci.parseSFrame()
		return t, f, nil
	case apci.Ctr1&3 == uFrame:
		//U帧
		t, f := apci.parseUFrame()
		return t, f, nil
	default:
		return 0xFF, nil, fmt.Errorf("未知APCI帧类型")
	}
}

//parseIFrame 解析I帧
func (apci *APCI) parseIFrame() (byte, IFrame) {
	send := int16(apci.Ctr1)>>1 + int16(apci.Ctr2)<<7
	recv := int16(apci.Ctr3)>>1 + int16(apci.Ctr4)<<7
	return iFrame, IFrame{
		Send: send,
		Recv: recv,
	}
}

func (apci *APCI) parseSFrame() (byte, SFrame) {
	recv := int16(apci.Ctr3)>>1 + int16(apci.Ctr4)<<7
	return sFrame, SFrame{
		Recv: recv,
	}
}

func (apci *APCI) parseUFrame() (byte, UFrame) {
	cmd := [4]byte{apci.Ctr1, apci.Ctr2, apci.Ctr3, apci.Ctr4}
	return uFrame, UFrame{
		cmd: cmd,
	}
}
