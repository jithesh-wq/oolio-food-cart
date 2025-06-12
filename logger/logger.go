package logger

import (
	"log"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

var Log *zap.SugaredLogger

func Init() error {
	if err := godotenv.Load("/home/hp-jv/interview/oolio-food-cart/cmd/.env"); err != nil {
		log.Println("No .env file found, relying on system environment variables")
	}
	logger, err := zap.NewProduction()
	if err != nil {
		return err
	}
	Log = logger.Sugar()
	return nil
}

func Close() {
	if Log != nil {
		_ = Log.Sync()
	}
}
