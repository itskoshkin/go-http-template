package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"time"
	"unicode/utf8"

	"github.com/gin-gonic/gin"

	"go-http-template/internal/utils/ginutils"
	"go-http-template/internal/utils/useragent"
)

var (
	logFile     *os.File
	logFileOnce sync.Once
)

func GetLogFile() *os.File {
	logFileOnce.Do(func() {
		var err error
		logFile, err = os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0640)
		if err != nil {
			log.Printf("Error: failed to open log file: %v", err)
		}
	})
	return logFile
}

func SetupLogging() {
	file := GetLogFile()
	log.SetFlags(log.Ldate | log.Ltime)
	if file != nil {
		log.SetOutput(io.MultiWriter(os.Stdout, file))
		_, _ = file.WriteString("==== New run at " + time.Now().Format("2006-01-02 15:04:05") + " ====\n")
	}
}

func CustomGinLogger(out io.Writer) gin.HandlerFunc {
	return gin.LoggerWithConfig(gin.LoggerConfig{
		Formatter: func(param gin.LogFormatterParams) string {
			ipAddr, referer := ginutils.ExtractIPAndReferer(param)

			if utf8.RuneCountInString(param.Path) > ginutils.PathWidth {
				return ginutils.TwoLinedAccessLog(param)
			}

			return fmt.Sprintf("[GIN] %s | %7s %-42s | %3d | %10v | %-15s | %s%s\n",
				param.TimeStamp.Format("2006/01/02 - 15:04:05"),
				param.Method,
				param.Path,
				param.StatusCode,
				param.Latency,
				ipAddr,
				useragent.ShortenUserAgent(param.Request.UserAgent()),
				referer,
			)
		},
		Output: out,
	})
}
