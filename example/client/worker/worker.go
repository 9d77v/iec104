package worker

import (
	"github.com/9d77v/iec104"
)

//Handler 处理接收到的已解析数据
func Handler(c *iec104.Client) {
	c.Logger.Info("数据处理协程启动")
	defer c.Cancel()
	for {
		select {
		case resp := <-c.DataChan:
			c.Logger.Debugf("接收到数据类型:%d,原因:%d,长度:%d", resp.ASDU.TypeID, resp.ASDU.Cause, len(resp.Signals))
			//TODO 数据处理
		case <-c.Ctx.Done():
			c.Logger.Info("数据接收协程停止")
			return
		}
	}
}
