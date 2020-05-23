package setting

import (
	"github.com/go-ini/ini"
	uuid "github.com/gofrs/uuid"
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
	PageSize        int
	JwtSecret       string
	SessionAssetPIN string
	ClientId        uuid.UUID
	ClientSecret    string
	SessionID       string
	PINToken        string
	SessionKey      string
)

func init() {
	var err error
	Cfg, err = ini.Load("secrets/app.ini")

	if err != nil {
		log.Fatalf("Fail to parse 'conf/app.ini': %v", err)
	}

	LoadBase()
	LoadBot()
	LoadServer()
	LoadApp()
}

func LoadBot() {
	sec, err := Cfg.GetSection("bot")
	if err != nil {
		log.Fatalf("Fail to get section 'bot': %v", err)
	}
	ClientId, _ = uuid.FromString(sec.Key("client_id").MustString("debug"))
	SessionAssetPIN = sec.Key("session_asset_pin").MustString("debug")
	ClientSecret = sec.Key("client_secret").MustString("debug")
	SessionID = sec.Key("session_id").MustString("debug")
	PINToken = sec.Key("pin_token").MustString("debug")
	SessionKey = sec.Key("session_key").MustString("debug")
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
	PageSize = sec.Key("PAGE_SIZE").MustInt(10)
}
