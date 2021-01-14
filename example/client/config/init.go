package config

import (
	"io"
	"os"
	"strconv"

	"github.com/sirupsen/logrus"

	"github.com/9d77v/iec104/example/client/utils"
)

var (
	//ServerHost 服务器ip
	ServerHost = os.Getenv("SERVER_HOST")
	//ServerPort 服务器端口
	ServerPort, _ = strconv.Atoi(os.Getenv("SERVER_PORT"))
	//SubServerHost 备用服务器ip
	SubServerHost = os.Getenv("SUB_SERVER_HOST")
	//SubServerPort 备用服务器端口
	SubServerPort, _ = strconv.Atoi(os.Getenv("SUB_SERVER_PORT"))
	//Debug 是否debug模式
	Debug, _ = strconv.ParseBool(os.Getenv("DEBUG"))
	//Logger 日志
	Logger *logrus.Logger
)

func init() {
	if ServerPort == 0 {
		ServerPort = 2404
	}
	initLogger()
}

func initLogger() {
	logger := logrus.New()
	//debug模式下打印行号
	if Debug {
		logger.SetLevel(logrus.DebugLevel)
		logger.Hooks.Add(utils.NewContextHook())
	} else {
		// 设置为json格式的日志
		logger.Formatter = &logrus.JSONFormatter{}
	}

	f, err := os.OpenFile("./logs/iec104.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644) // 创建一个log日志文件
	if err != nil {
		logger.Fatalln("打开日志文件失败")
	}
	writers := []io.Writer{
		f}
	//debug模式下输出到控制台
	if Debug {
		writers = append(writers, os.Stdout)
	}
	fileAndStdoutWriter := io.MultiWriter(writers...)

	logger.Out = fileAndStdoutWriter
	Logger = logger
}
