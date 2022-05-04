package logger

import (
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger
var err error
var Undo func()

func Init() {
	Logger, err = newProductionLogger()
	if err != nil {
		panic(err)
	}
	defer Logger.Sync()

	Undo = zap.ReplaceGlobals(Logger)
}

const logPath = "logs"

func newProductionLogger() (*zap.Logger, error) {
	productionLoggerConfig := zap.NewProductionConfig()
	productionLoggerConfig.OutputPaths = append(productionLoggerConfig.OutputPaths, logPath)
	productionLoggerConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	if _, err := os.Stat(logPath); err != nil {
		if os.IsNotExist(err) {
			logfile, iErr := os.Create(logPath)
			if iErr != nil {
				panic(err)
			}
			iErr = logfile.Close()
			if iErr != nil {
				panic(err)
			}
		}
	}

	return productionLoggerConfig.Build()

}

func main() {
	zap.L().Error("init logger", zap.String("1", "a"))
	zap.S().Errorw("hahaha", "1", "2", "3", "4")
	Undo()
	fmt.Println(zap.L())
	zap.L().Error("init logger", zap.String("1", "a"))
	zap.S().Errorw("hahaha", "1", "2", "3", "4")
}
