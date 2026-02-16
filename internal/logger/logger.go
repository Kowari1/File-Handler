package logger

import "go.uber.org/zap"

var L *zap.Logger

func Init() error {
	log, err := zap.NewProduction()
	if err != nil {
		return err
	}

	L = log
	return nil
}
