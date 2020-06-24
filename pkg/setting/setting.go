package setting

// import (
// 	"log"
// 	"time"

// 	"github.com/go-ini/ini"
// 	"github.com/gofrs/uuid"
// )

// var (
// 	// Cfg config
// 	Cfg *ini.File
// 	// RunMode debug or release
// 	RunMode string
// 	// HTTPPort port
// 	HTTPPort int
// 	// ReadTimeout read timeout
// 	ReadTimeout time.Duration
// 	// WriteTimeout write timeout
// 	WriteTimeout time.Duration
// 	// Name app name
// 	Name string
// 	// PageSize default page size
// 	PageSize int
// 	// JwtSecret jwt secret
// 	JwtSecret string
// 	// SessionAssetPIN session asset pin
// 	SessionAssetPIN string
// 	// ClientID robot id
// 	ClientID uuid.UUID
// 	// ClientSecret robot client secret
// 	ClientSecret string
// 	// SessionID robot session id
// 	SessionID string
// 	// PINToken robot pin token
// 	PINToken string
// 	// SessionKey robot session key
// 	SessionKey string
// 	// CodeVerifier auth code verifier
// 	CodeVerifier string
// 	// S3AccessKey    string
// 	// S3SecretKey    string
// 	// S3EndPoint     string
// 	// S3Region       string
// 	// S3Bucket       string

// 	// MixinAPIDomain mixin api domain
// 	MixinAPIDomain string
// )

// func init() {
// 	var err error
// 	Cfg, err = ini.Load("/home/ubuntu/setflags/secrets/app.ini")

// 	if err != nil {
// 		log.Fatalf("Fail to parse 'secrets/app.ini': %v", err)
// 	}

// 	LoadBase()
// 	LoadBot()
// 	LoadServer()
// 	LoadApp()
// 	// LoadAWSS3()
// }

// // func LoadAWSS3() {
// // 	sec, err := Cfg.GetSection("s3")
// // 	if err != nil {
// // 		log.Fatalf("Fail to get section 's3': %v", err)
// // 	}

// // 	S3AccessKey = sec.Key("access_key").MustString("debug")
// // 	S3SecretKey = sec.Key("secret_key").MustString("debug")
// // 	S3EndPoint = sec.Key("end_point").MustString("debug")
// // 	S3Region = sec.Key("region").MustString("debug")
// // 	S3Bucket = sec.Key("bucket").MustString("debug")
// // }

// // LoadBot load bot config
// func LoadBot() {
// 	sec, err := Cfg.GetSection("bot")
// 	if err != nil {
// 		log.Fatalf("Fail to get section 'bot': %v", err)
// 	}
// 	ClientID, _ = uuid.FromString(sec.Key("client_id").MustString("debug"))
// 	SessionAssetPIN = sec.Key("session_asset_pin").MustString("debug")
// 	ClientSecret = sec.Key("client_secret").MustString("debug")
// 	SessionID = sec.Key("session_id").MustString("debug")
// 	PINToken = sec.Key("pin_token").MustString("debug")
// 	SessionKey = sec.Key("session_key").MustString("debug")
// 	CodeVerifier = sec.Key("code_verifier").MustString("debug")
// }

// // LoadBase load base config
// func LoadBase() {
// 	RunMode = Cfg.Section("").Key("RUN_MODE").MustString("debug")
// }

// // LoadServer load service config
// func LoadServer() {
// 	sec, err := Cfg.GetSection("server")

// 	if err != nil {
// 		log.Fatalf("Fail to get section 'server': %v", err)
// 	}

// 	HTTPPort = sec.Key("HTTP_PORT").MustInt(8000)
// 	ReadTimeout = time.Duration(sec.Key("READ_TIMEOUT").MustInt(60)) * time.Second
// 	WriteTimeout = time.Duration(sec.Key("WRITE_TIMEOUT").MustInt(60)) * time.Second
// }

// // LoadApp load app config
// func LoadApp() {
// 	sec, err := Cfg.GetSection("app")
// 	if err != nil {
// 		log.Fatalf("Fail to get section 'app': %v", err)
// 	}
// 	Name = sec.Key("NAME").MustString("debug")
// 	JwtSecret = sec.Key("JWT_SECRET").MustString("!@)*#)!@U#@*!@!)")
// 	PageSize = sec.Key("PAGE_SIZE").MustInt(10)
// }
