package h_log

import (
	"fmt"
	"io"
	"load-config/h_config"
	"load-config/h_file"
	"os"
	"path"
	"strings"
	"sync"
	"time"

	"github.com/rs/zerolog"
)

type HLogger struct {
	zerolog.Logger
	mux          sync.Mutex
	interval     int       //日志切割时间间隔，单位：小时
	lastFileTime time.Time //上次日志文件的创建时间
	path         string    //日志文件存放路径
	serviceName  string    //日志服务器名
	env          string    //环境
}

var logger *HLogger

// 初始化全局HLogger对象logger
// 输入日志等级限制，日志路径，切割周期
func Init(level, pathName string, interval int, serviceName, env string) {
	//1.根据输入设置输出的最低日志级别阈值，低于该级别的日志将会被自动过滤
	switch strings.ToLower(level) {
	case LOG_LEVEL_DEBUG:
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case LOG_LEVEL_INFO:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case LOG_LEVEL_WARN:
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case LOG_LEVEL_ERROR:
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case LOG_LEVEL_FATAL:
		zerolog.SetGlobalLevel(zerolog.FatalLevel)
	case LOG_LEVEL_PANIC:
		zerolog.SetGlobalLevel(zerolog.PanicLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	//2.设置日志格式
	zerolog.TimeFieldFormat = "2006-01-02 15:04:05.000" //定义日志中时间字段的显示格式
	zerolog.TimestampFieldName = "timestamp"            //修改日志输出中的时间戳字段名
	zerolog.LevelFieldName = "Level"                    //修改日志中的日志级别字段名
	zerolog.MessageFieldName = "msg"                    //修改日志中的消息内容字段名

	//3.调用newOutPut()方法获取io流
	io, err := newOutPut(pathName)
	if err != nil {
		panic(err)
	}

	//4.实例化全局HLogger对象logger
	logger = &HLogger{
		Logger:       zerolog.New(io).With().Logger(),
		mux:          sync.Mutex{},
		interval:     interval,
		lastFileTime: time.Now(),
		path:         pathName,
		serviceName:  serviceName,
		env:          env,
	}
}

// 创建一个新的输出的io流
func newOutPut(pathName string) (io.Writer, error) {
	if pathName != "" {
		now := time.Now().Format("2006-01-02 15:04;05")
		fileName := fmt.Sprintf("%s.log", now)

		//文件不存在，则创建
		if !h_file.IsExist(pathName) {
			err := os.MkdirAll(pathName, os.ModePerm)
			//创建失败，保存
			if err != nil {
				return nil, err
			}
		}
		//在指定目录下创建文件
		file, err := os.Create(path.Join(pathName, fileName))
		if err != nil {
			return nil, err
		}
		return file, nil
	}
	fmt.Println("直接输出到控制台")
	//默认情况下直接输出到控制台
	return os.Stdout, nil
}

func Debug() *zerolog.Event {
	return newEvent(zerolog.DebugLevel)
}
func Info() *zerolog.Event {
	return newEvent(zerolog.InfoLevel)
}
func Error() *zerolog.Event {
	return newEvent(zerolog.ErrorLevel)
}
func Warn() *zerolog.Event {
	return newEvent(zerolog.WarnLevel)
}
func Fatal() *zerolog.Event {
	return newEvent(zerolog.FatalLevel)
}
func Panic() *zerolog.Event {
	return newEvent(zerolog.PanicLevel)
}
func Log() *zerolog.Event {
	return newEvent(zerolog.NoLevel)
}

// 创建一个zerolog.Event对象，用于构建和输出单条日志记录
// 通过链式调用的设计实现高性能结构化日志记录
func newEvent(level zerolog.Level) *zerolog.Event {
	//1.检测logger引擎是否初始化，是否要切割
	logger.check()

	//2.根据level返回Event
	switch level {
	case zerolog.DebugLevel:
		return logger.Logger.Debug().Str("env", logger.env).Str("service", logger.serviceName).Timestamp()
	case zerolog.InfoLevel:
		fmt.Println("打印Info")
		return logger.Logger.Info().Str("env", logger.env).Str("service", logger.serviceName).Timestamp()
	case zerolog.WarnLevel:
		return logger.Logger.Warn().Str("env", logger.env).Str("service", logger.serviceName).Timestamp()
	case zerolog.ErrorLevel:
		return logger.Logger.Error().Str("env", logger.env).Str("service", logger.serviceName).Timestamp()
	case zerolog.FatalLevel:
		return logger.Logger.Fatal().Str("env", logger.env).Str("service", logger.serviceName).Timestamp()
	case zerolog.PanicLevel:
		return logger.Logger.Panic().Str("env", logger.env).Str("service", logger.serviceName).Timestamp()
	case zerolog.NoLevel:
		return logger.Logger.Log().Str("env", logger.env).Str("service", logger.serviceName).Timestamp()
	default:
		return logger.Logger.Debug().Str("env", logger.env).Str("service", logger.serviceName).Timestamp()
	}
}

// 检测当前logger是否存在且没进入分页周期
func (log *HLogger) check() {
	if log == nil {
		Init(h_config.Cfg.Log.Level, h_config.Cfg.Log.Path, h_config.Cfg.Log.Interval, h_config.Cfg.ServiceName, h_config.Cfg.Env)
	} else {
		//日志文件切割
		if log.interval > 0 && time.Now().Add(-time.Hour*time.Duration(log.interval)).After(log.lastFileTime) {
			log.mux.Lock()
			defer log.mux.Unlock()
			if log.interval > 0 && time.Now().Add(-time.Hour*time.Duration(log.interval)).After(log.lastFileTime) {
				Init("debug", log.path, log.interval, log.serviceName, log.env)
			}
		}
	}
}
