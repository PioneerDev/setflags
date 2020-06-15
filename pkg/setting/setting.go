package setting

import (
	"github.com/go-ini/ini"
	"log"
	"time"
)

var (
	Cfg             *ini.File
	RunMode         string
	HTTPPort        int
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	Name            string
	PageSize        string
	JwtSecret       string
	SessionAssetPIN string
	ClientId        string
	ClientSecret    string
	SessionID       string
	PINToken        string
	SessionKey      string
	CodeVerifier    string
	S3AccessKey     string
	S3SecretKey     string
	S3EndPoint      string
	S3Region        string
	S3Bucket        string
	MixinAPIDomain  string
)

func init() {
	var err error
	Cfg, err = ini.Load("/home/ubuntu/setflags/secrets/app.ini")

	if err != nil {
		log.Fatalf("Fail to parse 'secrets/app.ini': %v", err)
	}

	LoadBase()
	LoadBot()
	LoadServer()
	LoadApp()
	LoadAWSS3()
	LoadMixin()
}

func LoadMixin() {
	sec, err := Cfg.GetSection("mixin")
	if err != nil {
		log.Fatalf("Fail to get section 'mixin': %v", err)
	}

	MixinAPIDomain = sec.Key("api_domain").MustString("https://api.mixin.one")
}

func LoadAWSS3() {
	sec, err := Cfg.GetSection("s3")
	if err != nil {
		log.Fatalf("Fail to get section 's3': %v", err)
	}

	S3AccessKey = sec.Key("access_key").MustString("debug")
	S3SecretKey = sec.Key("secret_key").MustString("debug")
	S3EndPoint = sec.Key("end_point").MustString("debug")
	S3Region = sec.Key("region").MustString("debug")
	S3Bucket = sec.Key("bucket").MustString("debug")
}

func LoadBot() {
	sec, err := Cfg.GetSection("bot")
	if err != nil {
		log.Fatalf("Fail to get section 'bot': %v", err)
	}
	ClientId = sec.Key("client_id").MustString("debug")
	SessionAssetPIN = sec.Key("session_asset_pin").MustString("debug")
	ClientSecret = sec.Key("client_secret").MustString("debug")
	SessionID = sec.Key("session_id").MustString("debug")
	PINToken = sec.Key("pin_token").MustString("debug")
	SessionKey = sec.Key("session_key").MustString("debug")
	CodeVerifier = sec.Key("code_verifier").MustString("debug")
}

func LoadBase() {
	RunMode = Cfg.Section("").Key("RUN_MODE").MustString("debug")
}

func LoadServer() {
	sec, err := Cfg.GetSection("server")

	if err != nil {
		log.Fatalf("Fail to get section 'server': %v", err)
	}

	HTTPPort = sec.Key("HTTP_PORT").MustInt(8000)
	ReadTimeout = time.Duration(sec.Key("READ_TIMEOUT").MustInt(60)) * time.Second
	WriteTimeout = time.Duration(sec.Key("WRITE_TIMEOUT").MustInt(60)) * time.Second
}

func LoadApp() {
	sec, err := Cfg.GetSection("app")
	if err != nil {
		log.Fatalf("Fail to get section 'app': %v", err)
	}
	Name = sec.Key("NAME").MustString("debug")
	JwtSecret = sec.Key("JWT_SECRET").MustString("!@)*#)!@U#@*!@!)")
	PageSize = sec.Key("PAGE_SIZE").MustString("10")
}
