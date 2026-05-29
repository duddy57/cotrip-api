package config

import (
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewLogger(env string) (*zap.Logger, error) {
	var l *zap.Logger
	var err error

	switch strings.ToLower(env) {
	case "prd":
		l, err = zap.NewProduction()
		if err != nil {
			return nil, err
		}
	default:
		cfg := zap.NewDevelopmentConfig()
		cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

		// Usamos '=' em vez de ':=' aqui também
		l, err = cfg.Build()
		if err != nil {
			return nil, err
		}
	}
	l = l.Named("journey_logger")
	// e não quando a aplicação fechar. Geralmente o Sync() é chamado no main.
	// 

	return l, nil
}
