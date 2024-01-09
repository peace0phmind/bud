package infra

import (
	"errors"
	"fmt"
	"github.com/mcuadros/go-defaults"
	"github.com/peace0phmind/bud/factory"
	_struct "github.com/peace0phmind/bud/struct"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"log"
	"os"
	"path/filepath"
)

var _config = factory.Singleton[Config]().MustBuilder()

type Config struct {
	EnvType   string `default:"prod"`
	MySQLHost string `default:"127.0.0.1"`
	MySQLPort string `default:"3306"`
	MySQLUser string `default:"vm"`
	MySQLPass string `default:"vm123456"`
	MySQLDb   string `default:"vm"`

	LogFileDir    string `default:"."`
	LogFileName   string `default:"mix.log"`
	LogLevel      string `default:"info"`
	DBLogLevel    string `default:"info"`
	JwtSigningKey string `default:"2C9+UxzL7yasmuirhYeZ0WMm093QNwK9gjjrEIyjnSvJ2AwkrDze3Wp0qjzbQJJNJHhc6DppvSMNYd+t6svTcAk9vlHyVh3ccB6z7mYGq+k4yEOxzrz1xBUCX2mCrg3UzhzoFIPXLoA7CaajOIY0gs7k+GcPAABaXe2K5BeFSXw="` // command: openssl rand -base64 128

	_dsn string
}

func (cfg *Config) InitOnce() error {

	defaults.SetDefaults(cfg)

	// 下面默认值可以写入配置文件
	fh := _struct.FieldHelper(cfg)
	for _, fieldName := range []string{"MySQLHost", "MySQLPort", "UseGPU"} {
		fh.GetValue2Do(fieldName, func(fn string, value any) bool {
			viper.SetDefault(fn, value)
			return true
		})
	}

	// 访问.env文件
	viper.SetConfigType("env")
	viper.SetConfigName(".env")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()
	envErr := viper.ReadInConfig()                                             // 这里不会报错，如果 .env 文件不存在，viper 会忽略它
	if envErr != nil && !errors.As(envErr, &viper.ConfigFileNotFoundError{}) { // 如果有其他错误除了文件不存在，则报错
		log.Fatalf("Error reading .env file, %s", envErr)
	}

	// 检查 mix.yml 文件是否存在，如果即不存在.env，也不存在mix.yaml,则创建mix.yaml文件
	mixConfigPath := filepath.Join(".", "mix.yml")
	if _, mixErr := os.Stat(mixConfigPath); os.IsNotExist(mixErr) && envErr != nil && errors.As(envErr, &viper.ConfigFileNotFoundError{}) {
		// mix.yml 文件不存在，创建文件并写入默认值
		err := viper.SafeWriteConfigAs(mixConfigPath) // 使用 SafeWriteConfigAs 防止覆盖已存在的文件
		if err != nil {
			log.Fatalf("Failed to create mix.yaml file with default values, %s", err)
		}
	}

	// 读取 mix.yml 文件
	viper.SetConfigType("yml")
	viper.SetConfigName("mix")
	viper.AddConfigPath(".")
	err := viper.MergeInConfig() // 合并配置文件，而不是覆盖
	if err != nil && errors.As(envErr, &viper.ConfigFileNotFoundError{}) && envErr != nil && !errors.As(envErr, &viper.ConfigFileNotFoundError{}) {
		log.Fatalf("Error reading mix.yaml file, %s. '.env' or 'mix.yaml' file must be present", err)
	}

	// 将读取的配置与结构体中的默认值合并
	err = viper.Unmarshal(cfg)
	if err != nil {
		log.Fatalf("Unable to decode into struct, %v", err)
	}

	return cfg.initLog()
}

func (cfg *Config) initLog() error {
	customFormatter := &logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		FullTimestamp:   true,
	}
	logrus.SetFormatter(customFormatter)

	if err := os.MkdirAll(cfg.LogFileDir, os.ModePerm); err != nil {
		return err
	}

	if cfg.EnvType != "dev" {
		fullLogFile := filepath.Join(cfg.LogFileDir, cfg.LogFileName)
		// set log file
		logFile, err := os.OpenFile(fullLogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			logFile, err = os.OpenFile("./video_mixing.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
			if err != nil {
				log.Panic(err)
			}
		}
		logrus.SetOutput(logFile)
	}

	// set log level
	lvl, err := logrus.ParseLevel(cfg.LogLevel)
	if err != nil {
		lvl = logrus.InfoLevel
	}
	logrus.SetLevel(lvl)
	return nil
}

func (cfg *Config) DSN() string {
	if len(cfg._dsn) > 0 {
		return cfg._dsn
	} else {
		cfg._dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			cfg.MySQLUser, cfg.MySQLPass, cfg.MySQLHost, cfg.MySQLPort, cfg.MySQLDb)
		return cfg._dsn
	}
}
