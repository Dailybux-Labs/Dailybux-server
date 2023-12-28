package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"dailybux/model"
)

var (
	Conf *model.Config
	Db0  *gorm.DB
	Rds  *RedisConn
)

type RedisConn struct {
	redisClient *redis.Client
	projectName string
}

func Init() {
	logInit()
	var err error
	if Conf, err = loadConfig("config/config.toml"); err != nil {
		zap.S().Errorf("[%v] error %v", "config/config.toml", err.Error())
	}

	Rds, err = rdsInit(Conf.Redis)
	if err != nil {

	}
	Db0, err = dbInit(Conf.MySQL.Host, Conf.MySQL.User, Conf.MySQL.Password, Conf.MySQL.Port, Conf.MySQL.DBName)
	if err != nil {

	}
}

func loadConfig(path string) (*model.Config, error) {
	var c model.Config
	_, err := toml.DecodeFile(path, &c)
	if err != nil {
		zap.S().Error(err.Error())
		return nil, err
	}
	return &c, nil
}

func logInit() {
	atomicLevel := zap.NewAtomicLevel()
	atomicLevel.SetLevel(zapcore.InfoLevel)
	syncer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "./logs/dailybux-server.log",
		MaxSize:    6,
		MaxAge:     7,
		MaxBackups: 10,
		LocalTime:  true,
		Compress:   false,
	})

	zapLog := zap.New(zapcore.NewCore(zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "name",
		CallerKey:      "line",
		MessageKey:     "msg",
		FunctionKey:    "F",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05"),
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
		EncodeName:     zapcore.FullNameEncoder,
	}),
		zapcore.NewMultiWriteSyncer(syncer, zapcore.AddSync(os.Stdout)), atomicLevel), zap.AddCaller())

	zap.ReplaceGlobals(zapLog)
}

func dbInit(host, user, pass string, port int64, dbName string) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, pass, host, port, dbName)), &gorm.Config{Logger: logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			Colorful:                  true,        // Disable color
		},
	)})
	if err != nil {
		zap.S().Errorf("[%v] error %v", host, err.Error())
		return nil, err
	}
	return db, nil
}

func rdsInit(redis2 model.Redis) (*RedisConn, error) {
	redisClient := redis.NewClient(&redis.Options{
		Addr: redis2.Addr, Password: redis2.Password,
	})

	_, err := redisClient.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}
	redisConn := &RedisConn{}
	redisConn.redisClient = redisClient
	redisConn.projectName = "unify-service"

	return redisConn, nil
}
