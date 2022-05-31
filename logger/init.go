package logger

import (
	"path"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	log "github.com/sirupsen/logrus"

	"github.com/teleport/stress-raw/utils"
)

// 全局日志配置初始化
// 1.日志同时输出到到文件和标准输出；
// 2.保留最近 7 天日志，一天一个日志文件；
//
func init() {
	// dir, _ := os.Getwd()
	dir := "./"
	logDir := dir + "/logs"
	utils.PathNEAC(logDir)

	logfile := path.Join(logDir, "app")
	fsWriter, err := rotatelogs.New(
		logfile+"_%Y-%m-%d.log",
		rotatelogs.WithMaxAge(time.Duration(24)*time.Hour),
		rotatelogs.WithRotationTime(time.Duration(24)*time.Hour),
	)
	if err != nil {
		panic(err)
	}
	defer fsWriter.Close()

	// multiWriter := io.MultiWriter(fsWriter, os.Stdout)
	// log.SetOutput(multiWriter)
	log.SetReportCaller(true)
	log.SetFormatter(&log.JSONFormatter{})
	// log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(fsWriter)

	log.SetLevel(log.InfoLevel)
}
