package zaplog

import (
	"fmt"
	"go-spider/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"path/filepath"
)

var Logger *zap.Logger

func InitLogger() {
	logFile, err := getLogFile()
	if err != nil {
		fmt.Printf("Failed to create path：%s\n", err.Error())
	}
	writeSyncer := zapcore.AddSync(logFile)
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.MessageKey = "MESSAGE"
	encoderConfig.TimeKey = "datetime"
	encoderConfig.LevelKey = "level_name"
	//Specify time format
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	//Display different colors according to levels. If not needed, just use zapcore.CapitalLevelEncoder.
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	//encoderConfig.EncodeCaller = zapcore.FullCallerEncoder
	encoder := zapcore.NewJSONEncoder(encoderConfig)
	core := zapcore.NewCore(encoder, writeSyncer, zapcore.DebugLevel)

	Logger = zap.New(core)
}

func getLogFile() (*os.File, error) {
	filePath := config.ConfData.Logger.Path
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		dir := filepath.Dir(filePath)
		err = os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return nil, err
		}
		file, err := os.Create(filePath)
		if err != nil {
			return nil, err
		}
		return file, nil
	} else if err == nil {
		file, err := os.OpenFile(filePath, os.O_APPEND|os.O_RDWR, os.ModePerm)
		if err != nil {
			fmt.Printf("open file %s failed：%s\n", filePath, err.Error())
			return nil, err
		}
		return file, nil
	} else {
		fmt.Printf("check path %s failed：%s\n", filePath, err.Error())
		return nil, err
	}

}
