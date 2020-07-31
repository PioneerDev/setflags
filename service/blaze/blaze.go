package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"log"
	"set-flags/global"
	"set-flags/models"
	"set-flags/schemas"
	"time"

	sdk "github.com/fox-one/mixin-sdk"
)

// Handler Handler
type Handler struct{}

// OnMessage OnMessage call back
func (h Handler) OnMessage(ctx context.Context, msgView *sdk.MessageView, userID string) error {
	// log.Println("I received a msg", msgView)

	if msgView.Category == sdk.MessageCategorySystemAccountSnapshot {
		data, err := base64.StdEncoding.DecodeString(msgView.Data)
		if err != nil {
			return err
		}
		var snapshot schemas.AccountSnapshot
		if err := json.Unmarshal(data, &snapshot); err != nil {
			return err
		}

		log.Println(snapshot.TraceID, snapshot.Amount)

		// transfer in
		// now ignore transfer out
		if snapshot.Amount > 0 {
			models.UpdatePaymentAndFlag(global.Db, snapshot)
		}
	}

	return nil
}

// OnAckReceipt OnAckReceipt
func (h Handler) OnAckReceipt(ctx context.Context, msg *sdk.MessageView, userID string) error {
	return nil
}

// Run Run
func (h Handler) Run(ctx context.Context, user *sdk.User) {
	for {
		select {
		case <-ctx.Done():
			break

		default:
		}
		if err := sdk.NewBlazeClient(user).Loop(ctx, h); err != nil {
			log.Println("something is wrong", err)
			time.Sleep(1 * time.Second)
		}
	}
}

func main() {
	global.BotInit()
	global.InitDB()
	ctx := context.Background()
	log.Println("start bot")
	handler := Handler{}

	handler.Run(ctx, global.Bot)
}
