package setup

import "go.uber.org/zap"

const (
	devEnvType  = "dev"
	prodEnvType = "prod"
)

func MustSetupLogger(configEnvType string) *zap.Logger {
	switch configEnvType {
	case devEnvType:
		return zap.Must(zap.NewDevelopment())
	case prodEnvType:
		return zap.Must(zap.NewProduction())
	default:
		zap.L().Fatal("unknown env type", zap.String("env", configEnvType))
		return nil
	}
}
