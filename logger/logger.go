package main

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger
var err error
var Undo func()

func init() {
	Logger, err = newProductionLogger()
	if err != nil {
		panic(err)
	}
	defer Logger.Sync()

	Undo = zap.ReplaceGlobals(Logger)
}

func newProductionLogger() (*zap.Logger, error) {
	productionLoggerConfig := zap.NewProductionConfig()
	productionLoggerConfig.OutputPaths = append(productionLoggerConfig.OutputPaths, "./tmp/logs")
	productionLoggerConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

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
