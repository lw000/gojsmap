package global

import (
	tyip2region "github.com/lw000/gocommon/ip2region"
	"gojsmap/config"
	SourceMap "gojsmap/sourcemap"
	"path"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/pkg/errors"
	"github.com/rifflock/lfshook"
	log "github.com/sirupsen/logrus"
)

var (
	// ProjectConfig 工程配置
	ProjectConfig *config.IniConfig
	// IpServer IP地址转换
	IpServer *tyip2region.IpRegionServer
	// SourceMapServer SourceMap文件解析对象
	SourceMapManager *SourceMap.Manager
)

func init() {
	IpServer = tyip2region.NewIpRegionServer()
	SourceMapManager = SourceMap.NewSourceMapManager()
}

// config logrus log to local filesystem, with file rotation
func configLocalFilesystemLogger(logPath string, logFileName string, maxAge time.Duration, rotationTime time.Duration) {
	baseLogPath := path.Join(logPath, logFileName)
	writer, err := rotatelogs.New(
		baseLogPath+".%Y%m%d_%H%M",
		// rotatelogs.WithLinkName(baseLogPath), // 生成软链，指向最新日志文件
		rotatelogs.WithMaxAge(maxAge), // 文件最大保存时间
		// rotatelogs.WithRotationCount(365),  // 最多存365个文件
		rotatelogs.WithRotationTime(rotationTime), // 日志切割时间间隔
	)

	if err != nil {
		log.Errorf("config local file system logger error. %+v", errors.WithStack(err))
	}

	lfHook := lfshook.NewHook(lfshook.WriterMap{
		log.DebugLevel: writer, // 为不同级别设置不同的输出目的
		log.InfoLevel:  writer,
		log.WarnLevel:  writer,
		log.ErrorLevel: writer,
		log.FatalLevel: writer,
		log.PanicLevel: writer,
	}, &log.TextFormatter{TimestampFormat: "2006-01-02 15:04:05"})

	// 打印函数，和函数所在行号
	// log.SetReportCaller(true)

	log.AddHook(lfHook)
}

// LoadGlobalConfig 加载全局配置
func LoadGlobalConfig() error {
	configLocalFilesystemLogger("log", "gojsmap", time.Hour*24*365, time.Hour*24)

	var err error
	ProjectConfig, err = config.LoadIniConfig("conf/conf.ini")
	if err != nil {
		log.Error(err)
		return err
	}

	// 日志分割 1按天分割，2按周分割, 3 按月分割，4按年分割
	var logname = "gojsmap"
	switch ProjectConfig.SplitLog {
	case 1:
		configLocalFilesystemLogger("log", logname, time.Hour*24*365, time.Hour*24)
	case 2:
		configLocalFilesystemLogger("log", logname, time.Hour*24*365, time.Hour*24*7)
	case 3:
		configLocalFilesystemLogger("log", logname, time.Hour*24*365, time.Hour*24*30)
	case 4:
		configLocalFilesystemLogger("log", logname, time.Hour*24*365, time.Hour*24*365)
	default:
		configLocalFilesystemLogger("log", logname, time.Hour*24*365, time.Hour*24)
	}

	return nil
}
