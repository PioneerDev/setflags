package global

import (
	"crypto/x509"
	"encoding/pem"
	"set-flags/src/pkg/setting"

	"github.com/fox-one/mixin-sdk"
)

// Bot global bot
var Bot *mixin.User

// BotInit bot init
func BotInit() {
	Bot = &mixin.User{
		UserID:    setting.GetConfig().Bot.ClientID.String(),
		SessionID: setting.GetConfig().Bot.SessionID,
		PINToken:  setting.GetConfig().Bot.PinToken,
	}
	block, _ := pem.Decode([]byte(setting.GetConfig().Bot.PrivateKey))
	if block != nil {
		privateKey, _ := x509.ParsePKCS1PrivateKey(block.Bytes)
		Bot.SetPrivateKey(privateKey)
	}
}
