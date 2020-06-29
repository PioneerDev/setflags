package main

import (
	"context"
	"crypto/md5"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"log"
	"math"
	"set-flags/models"
	"set-flags/pkg/logging"
	"set-flags/pkg/setting"
	"set-flags/routers"
	"strings"
	"syscall"
	"time"

	number "github.com/MixinNetwork/go-number"
	sdk "github.com/fox-one/mixin-sdk"
	"github.com/fvbock/endless"
	uuid "github.com/gofrs/uuid"
	cron "github.com/robfig/cron/v3"
)

func newWithSeconds() *cron.Cron {
	var secondParser = cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.DowOptional | cron.Descriptor)
	return cron.New(cron.WithParser(secondParser), cron.WithChain())
}

func contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}

func UniqueConversationId(userId, recipientId uuid.UUID) uuid.UUID {
	minId, maxId := userId.String(), recipientId.String()
	if strings.Compare(minId, maxId) > 0 {
		maxId, minId = userId.String(), recipientId.String()
	}
	h := md5.New()
	io.WriteString(h, minId)
	io.WriteString(h, maxId)
	sum := h.Sum(nil)
	sum[6] = (sum[6] & 0x0f) | 0x30
	sum[8] = (sum[8] & 0x3f) | 0x80
	return uuid.FromBytesOrNil(sum)
}

func paySFCPrize(ctx context.Context, bot *sdk.User, userID uuid.UUID, flag *models.Flag) error {
	asset := models.FindAsset(flag.AssetID)
	price := asset.PriceUSD
	if asset.Symbol == "SFC" {
		price = 1
	}
	_, err := bot.Transfer(ctx, &sdk.TransferInput{
		TraceID:    uuid.Must(uuid.NewV4()).String(),
		AssetID:    "f80b5f5f-8e4d-3072-b618-bd6c0d8ccaa5", // SFC
		OpponentID: userID.String(),
		Amount:     number.FromFloat(flag.Amount).Mul(number.FromFloat(price)).Div(number.FromString("10")).Persist(),
		Memo:       setting.GetConfig().App.Name,
	}, setting.GetConfig().Bot.Pin)
	return err
}

func payFee(ctx context.Context, bot *sdk.User, userID uuid.UUID, flag *models.Flag, amount number.Decimal) error {
	_, err := bot.Transfer(ctx, &sdk.TransferInput{
		TraceID:    uuid.Must(uuid.NewV4()).String(),
		AssetID:    flag.AssetID.String(),
		OpponentID: userID.String(),
		Amount:     amount.Persist(),
		Memo:       setting.GetConfig().App.Name,
	}, setting.GetConfig().Bot.Pin)
	return err
}

func remainingDays(flag *models.Flag) int {
	days := int(flag.CreatedAt.Add(24*time.Hour*time.Duration(flag.Days)).Sub(time.Now()).Hours()/24) + 1
	return days
}

func getTask(flag *models.Flag) string {
	return flag.Task
}

func countVotes(flag *models.Flag) (int, int, int) {
	nWitnesses := len(flag.Witnesses())
	var yesVotes int
	var noVotes int
	yesVotes = 0
	noVotes = 0
	for _, p := range flag.Witnesses() {
		if p.PayeeID == flag.PayerID {
			continue
		}
		if p.Verified == 1 {
			yesVotes = yesVotes + 1
		}
		if p.Verified == -1 {
			noVotes = noVotes + 1
		}
	}
	return nWitnesses, yesVotes, noVotes
}

func payWitnesses(ctx context.Context, bot *sdk.User, flag *models.Flag, nCorrect, yesVotes, noVotes, remainingDays int) {
	amount := number.FromFloat(flag.Amount).Div(number.FromFloat(float64(10) * float64(flag.Days) * float64(nCorrect)))
	for _, p := range flag.Witnesses() {
		if p.PayeeID != flag.PayerID {
			if yesVotes >= noVotes && p.Verified == 1 || yesVotes <= noVotes && p.Verified == -1 {
				payFee(ctx, bot, p.PayeeID, flag, amount)
			}
			p.Verified = 0
		}
	}
}

func payWitnessesUnconditionally(ctx context.Context, bot *sdk.User, flag *models.Flag, nWitnesses, remainingDays int, task string) {
	amount := number.FromFloat(flag.Amount)
	if nWitnesses <= 0 {
		return
	}
	amount = number.FromString(amount.Div(number.FromFloat(float64(flag.Days) * float64(10) * float64(nWitnesses))).PresentFloor())
	for _, p := range flag.Witnesses() {
		if p.PayeeID != flag.PayerID {
			payFee(ctx, bot, p.PayeeID, flag, amount)
		}
	}
	flag.Status = "PAID"
	flag.RemainingAmount = number.FromFloat(flag.Amount).Div(number.FromFloat(float64(flag.Days))).Mul(number.FromFloat(float64(remainingDays) - 1)).Float64()
}

func rewardPayer(ctx context.Context, bot *sdk.User, flag *models.Flag, nCorrect, yesVotes, noVotes, remainingDays int, task string) {
	if yesVotes > noVotes && 0 <= int64(nCorrect) {
		payFee(ctx, bot, flag.PayerID, flag, number.FromFloat(flag.Amount).Div(number.FromFloat(float64(flag.Days))).Mul(number.FromString("0.9")))
		if flag.RemainingDays == flag.Days {
			paySFCPrize(ctx, bot, flag.PayerID, flag)
		}
		flag.TimesAchieved = flag.TimesAchieved + 1
	}
	flag.RemainingAmount = number.FromFloat(flag.Amount).Div(number.FromFloat(float64(flag.Days))).Mul(number.FromFloat(float64(remainingDays))).Float64()
	flag.RemainingDays = remainingDays
}

func sendTextMessage(ctx context.Context, bot *sdk.User, conversationId uuid.UUID, message string) error {
	err := bot.SendMessage(ctx, &sdk.MessageRequest{
		ConversationID: conversationId.String(),
		MessageID:      uuid.Must(uuid.NewV4()).String(),
		Category:       "PLAIN_TEXT",
		Data:           base64.StdEncoding.EncodeToString([]byte(message)),
	})
	if err != nil {
		log.Println(err)
	}
	return err
}

func sendUserAppCard(ctx context.Context, bot *sdk.User, userId uuid.UUID, flag *models.Flag) error {
	payer := models.FindUser(flag.PayerID)
	card, _ := json.Marshal(map[string]string{
		"app_id":      setting.GetConfig().Bot.ClientID.String(),
		"icon_url":    "https://images.mixin.one/X44V48LK9oEBT3izRGKqdVSPfiH5DtYTzzF0ch5nP-f7tO4v0BTTqVhFEHqd52qUeuVas-BSkLH1ckxEI51-jXmF=s256",
		"title":       "励志定投群红包",
		"description": fmt.Sprintf("来自@%s 的红包", payer.IdentityNumber),
		"action":      "https://group-redirect.droneidentity.eu" + "/flags/" + flag.ID.String(),
	})
	cID := UniqueConversationId(setting.GetConfig().Bot.ClientID, userId)
	err := bot.SendMessage(ctx, &sdk.MessageRequest{
		ConversationID: cID.String(),
		MessageID:      uuid.Must(uuid.NewV4()).String(),
		Category:       "APP_CARD",
		Data:           base64.StdEncoding.EncodeToString(card),
	})
	if err != nil {
		log.Println(err)
	}
	return nil
}

func remindWitnesses(ctx context.Context, bot *sdk.User, flag *models.Flag, remainingDays int, task string) {
	for _, p := range flag.Witnesses() {
		if p.Verified == 0 && p.PayeeID != flag.PayerID {
			appMsg := "请您验证:@%d第%d天完成%s了吗？"
			cID := UniqueConversationId(setting.GetConfig().Bot.ClientID, p.PayeeID)
			payer := models.FindUser(flag.PayerID)
			sendTextMessage(ctx, bot, cID, fmt.Sprintf(appMsg, payer.IdentityNumber, int(flag.Days)-remainingDays+1, task))
			sendUserAppCard(ctx, bot, p.PayeeID, flag)
		}
	}
}

func encouragePayer(ctx context.Context, bot *sdk.User, flag *models.Flag, remainingDays int, task string) {
	payMsg := "谢谢@%d, 收到你第%d天的红包，希望你能够坚持每天完成'%s'，记得分享证据。确定你做到了！"
	cID := UniqueConversationId(setting.GetConfig().Bot.ClientID, flag.PayerID)
	payer := models.FindUser(flag.PayerID)
	sendTextMessage(ctx, bot, cID, fmt.Sprintf(payMsg, payer.IdentityNumber, int(flag.Days)-remainingDays+1, task))
	sendUserAppCard(ctx, bot, flag.PayerID, flag)
}

func remindPayerForEvidence(ctx context.Context, bot *sdk.User, flag *models.Flag, task string) {
	done := false
	for _, p := range flag.Witnesses() {
		if p.PayeeID == flag.PayerID {
			done = (p.Verified == 2)
			break
		}
	}
	if !done {
		cID := UniqueConversationId(setting.GetConfig().Bot.ClientID, flag.PayerID)
		payMsg := "今天@%s, 你完成'%s'了吗？请先上传证据，然后点击确认"
		payer := models.FindUserByID(flag.PayerID)
		sendTextMessage(ctx, bot, cID, fmt.Sprintf(payMsg, payer.IdentityNumber, task))
		sendUserAppCard(ctx, bot, flag.PayerID, flag)
	}
}

func upsertAsset(ctx context.Context, bot *sdk.User) {
	assets, _ := bot.ReadAssets(ctx)

	for _, asset := range assets {
		models.UpsertAsset(asset)
	}
}

// Reminder Reminder
func Reminder(ctx context.Context, bot *sdk.User, newDay bool) {
	flags := models.ListActiveFlags(true)
	for _, flag := range flags {
		task := flag.Task
		remainingDays := flag.RemainingDays
		nWitnesses, yesVotes, noVotes := countVotes(flag)
		if remainingDays <= 0 {
			continue
		}
		isVerified := false
		for _, pp := range flag.Witnesses() {
			if pp.Verified == 2 && pp.PayeeID == flag.PayerID {
				isVerified = true
				break
			}
		}
		if isVerified {
			if newDay {
				nCorrect := int(math.Max(float64(yesVotes), float64(noVotes)))
				if nCorrect > 0 {
					payWitnesses(ctx, bot, flag, nCorrect, yesVotes, noVotes, remainingDays)
					rewardPayer(ctx, bot, flag, nCorrect, yesVotes, noVotes, remainingDays, task)
				}
				remindWitnesses(ctx, bot, flag, remainingDays, task)
			}
		} else {
			if newDay {
				payWitnessesUnconditionally(ctx, bot, flag, nWitnesses, remainingDays, task)
				encouragePayer(ctx, bot, flag, remainingDays, task)
			} else {
				remindPayerForEvidence(ctx, bot, flag, task)
			}
		}
	}
}

func addTimers(ctx context.Context, cron *cron.Cron, bot *sdk.User) {
	/*
		cron.AddFunc("0 * * * * ?", func() {
			log.Println("Sceduled every minute to test")
			Reminder(ctx, bot, false)
		})
	*/
	cron.AddFunc("0 0 8 * * ?", func() {
		Reminder(ctx, bot, false)
	})
	cron.AddFunc("0 0 20 * * ?", func() {
		Reminder(ctx, bot, false)
	})
	cron.AddFunc("0 0 23 * * ?", func() {
		Reminder(ctx, bot, true)
	})

	cron.AddFunc("0 * * * * ?", func() {
		upsertAsset(ctx, bot)
	})
}

func main() {
	logging.Setup()
	models.InitDB()
	bot := &sdk.User{
		UserID:    setting.GetConfig().Bot.ClientID.String(),
		SessionID: setting.GetConfig().Bot.SessionID,
		PINToken:  setting.GetConfig().Bot.PinToken,
	}
	block, _ := pem.Decode([]byte(setting.GetConfig().Bot.PrivateKey))
	if block != nil {
		privateKey, _ := x509.ParsePKCS1PrivateKey(block.Bytes)
		bot.SetPrivateKey(privateKey)
	}

	cron := newWithSeconds()
	cron.Start()
	defer cron.Stop()
	ctx := context.Background()
	addTimers(ctx, cron, bot)

	endless.DefaultReadTimeOut = time.Duration(setting.GetConfig().App.ReadTimeOut) * time.Second
	endless.DefaultWriteTimeOut = time.Duration(setting.GetConfig().App.WriteTimeOut) * time.Second
	endless.DefaultMaxHeaderBytes = 1 << 20
	endPoint := fmt.Sprintf(":%d", setting.GetConfig().App.HTTPPort)

	server := endless.NewServer(endPoint, routers.InitRouter())
	server.BeforeBegin = func(add string) {
		logging.Info(fmt.Sprintf("Actual pid is %d", syscall.Getpid()))
		log.Printf("Actual pid is %d", syscall.Getpid())
	}

	err := server.ListenAndServe()
	if err != nil {
		logging.Info(fmt.Sprintf("Server err: %v", err))
		log.Printf("Server err: %v", err)
	}
}
