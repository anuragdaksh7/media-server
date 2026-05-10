package logger

import (
	"fileserver/internal/config"

	"go.uber.org/zap"
)

var Logger *zap.Logger

func InitLogger(config config.Config) {

	//if config.Environment == "dev" {
	devLogger, _ := zap.NewDevelopment()
	Logger = devLogger
	return
	//}

	//
	//core, err := adapter.New(
	//	adapter.SetDataset(config.AxiomDataset),
	//	adapter.SetClientOptions(
	//		axiom.SetToken(config.AxiomToken),
	//	),
	//)
	//if err != nil {
	//	log.Fatalf("Could not create axiom core %v", err)
	//}
	//
	//consoleCore := zapcore.NewCore(
	//	zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig()),
	//	zapcore.AddSync(os.Stdout),
	//	zapcore.DebugLevel,
	//)
	//
	//Logger = zap.New(zapcore.NewTee(core, consoleCore))
}
