package iec104

//Signal 104信号
type Signal struct {
	TypeID  uint    `json:"type_id"` //类型id，1:单点遥信，9:单点遥测
	Address uint32  `json:"address"` //地址
	Value   float64 `json:"value"`   //值
	Quality byte    `json:"quality"` //品质描述
	Ts      float64 `json:"ts"`      //毫秒时间戳
}
