package main

import (
	"context"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/MixinNetwork/bot-api-go-client"
	"github.com/MixinNetwork/bot-api-go-client/config"
	"github.com/MixinNetwork/go-number"
	_ "github.com/go-sql-driver/mysql"
	"github.com/satori/go.uuid"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

// CNBAssetID is the CNB's ID in Mixin Network
const CNBAssetID = "965e5c6e-434c-3fa9-b780-c50f43cd955c"

const defaultResponse = `命令错误
命令一：领糖果
命令二：创建社群#社群名称#总量#份数
～～～～～～～～～～～～～～～～～～
注：总量和份数为数字，暂时只支持CNB`

var client *bot.BlazeClient

var Db *sql.DB

//  链接mysql数据
func init() {
	// 数据库名称 group_helper
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/group_helper?charset=utf8")
	dbCheckError(err)
	Db = db
}

func dbCheckError(err error) {
	if err != nil {
		//log.Panicf("Error: %s\n", err)
	}
}

// Handler is an implementation for interface bot.BlazeListener
// check out the url for more details: https://github.com/MixinNetwork/bot-api-go-client/blob/master/blaze.go#L89.
type Handler struct{}

// OnMessage is a general method of bot.BlazeListener
func (r Handler) OnMessage(ctx context.Context, msgView bot.MessageView, botID string) error {
	// I handle PLAIN_TEXT message only and make sure respond to current conversation.
	// 转账类型 && 当前会话 //支付完成后提示成功，
	if msgView.Category == bot.MessageCategorySystemAccountSnapshot &&
		msgView.ConversationId == bot.UniqueConversationId(config.GetConfig().ClientID, msgView.UserId) {
		// candySave(msgView.UserId)

		data, _ := base64.StdEncoding.DecodeString(msgView.Data)
		var resp struct {
			TraceId string `json:"trace_id"`
		}
		json.Unmarshal(data, &resp)

		println("正在更新昨天", resp.TraceId)
		//change app status
		changeAppStatus(ctx, msgView, resp.TraceId)
		Respond(ctx, msgView, fmt.Sprintf("创建社群成功，用户加入社群回复'领糖果'即可领取糖果！"))
	}

	if msgView.Category == bot.MessageCategoryPlainText &&
		msgView.ConversationId == bot.UniqueConversationId(config.GetConfig().ClientID, msgView.UserId) {
		var data []byte
		var err error
		if data, err = base64.StdEncoding.DecodeString(msgView.Data); err != nil {
			log.Panicf("Error: %s\n", err)
			return err
		}

		inst := string(data)
		log.Printf("I got a message from %s, it said: `%s`\n", msgView.UserId, inst)
		if "test" == inst {
			Respond(ctx, msgView, "test")
		} else if "糖果" == inst {
			// 转账给用户
			Transfer(ctx, msgView)
		} else if 0 < strings.Index(inst, "#") {
			app := strings.Split(inst, "#")
			amount := app[2]
			println(app[0], app[1], app[2], app[3])
			//num[0] num[1] 转账后存放至数据库
			trace := bot.UniqueConversationId(strconv.Itoa(rand.Intn(1111)), strconv.Itoa(rand.Intn(1111)))
			createApp(amount, app[3], msgView.UserId, app[1], CNBAssetID, trace)
			user_id := config.GetConfig().ClientID //转账给机器人
			// 拼接支付链接
			msg := "https://mixin.one/pay?recipient=" + user_id + "&asset=" + CNBAssetID + "&amount=" + amount + "&trace=" + trace + "&memo=TEXT"
			RespondButton(ctx, msgView, fmt.Sprintf(msg))
		} else {
			Respond(ctx, msgView, defaultResponse)
		}
	}
	return nil
}

// 根据用户抵押资产 发送糖果
func Transfer(ctx context.Context, msgView bot.MessageView) {

	candyNum := candyNum(msgView.UserId)
	payload := bot.TransferInput{
		AssetId:     CNBAssetID,
		RecipientId: msgView.UserId,
		Amount:      number.FromString(string(candyNum)),
		TraceId:     uuid.Must(uuid.NewV4()).String(),
		Memo:        "Enjoy!",
	}
	err := bot.CreateTransfer(ctx, &payload,
		config.GetConfig().ClientID,
		config.GetConfig().SessionID,
		config.GetConfig().PrivateKey,
		config.GetConfig().Pin,
		config.GetConfig().PinToken,
	)
	if err != nil {
		Respond(ctx, msgView, fmt.Sprintf("Oops, %s\n", err))
	}
}

//获取用户糖果数量
func candyNum(user_id string) int {
	var num int
	Db.QueryRow("select SUN(num) from candy where user_id = ? ", user_id).Scan(&num)
	//num = "1"
	if num <= 0 {

	}
	return num
}

//创建app
func createApp(num, lot, user_id, name, asset_id, trace string) {
	tx, _ := Db.Begin()
	stmt, _ := tx.Prepare("INSERT INTO apps (num,lot, user_id,name,asset_id,status,trace) VALUES (?,?,?,?,?,?,?)")
	stmt.Exec(num, lot, user_id, name, asset_id, 0, trace)
	tx.Commit()
}

// 更新状态
func changeAppStatus(ctx context.Context, msgView bot.MessageView, trace string) {

	stmt, _ := Db.Prepare("UPDATE  apps  set status=1 where trace =?")
	rs, _ := stmt.Exec(trace)
	app_id, _ := rs.LastInsertId()
	createGroup(ctx, msgView, app_id)
}

func createGroup(ctx context.Context, msgView bot.MessageView, app_id int64) {
	conversationId := CreateConversation(ctx, msgView)
	tx, _ := Db.Begin()
	stmt, _ := tx.Prepare("INSERT INTO groups (app_id,name,conversation_id) VALUES (?,?,?)")
	stmt.Exec(app_id, "群组", conversationId)
	tx.Commit()
}

// 创建一个群组
func CreateConversation(ctx context.Context, msgView bot.MessageView) string {

	conversationId := bot.UniqueConversationId(strconv.Itoa(rand.Intn(1111)), msgView.UserId)
	participant := bot.Participant{
		UserId: msgView.UserId,
		Role:   "ADMIN",
	}
	participants := []bot.Participant{
		participant,
	}
	_, err := bot.CreateConversation(ctx, "GROUP", conversationId, participants, config.GetConfig().ClientID, config.GetConfig().SessionID, config.GetConfig().PrivateKey)

	if err != nil {
		Respond(ctx, msgView, fmt.Sprintf("error, %s\n", err))
		return ""
	}
	return conversationId
}

// Respond to user.
func Respond(ctx context.Context, msgView bot.MessageView, msg string) {
	if err := client.SendPlainText(ctx, msgView, msg); err != nil {
		log.Panicf("Error: %s\n", err)
	}
}

func RespondButton(ctx context.Context, msgView bot.MessageView, msg string) {
	if err := client.SendAppButton(ctx, msgView.ConversationId, msgView.UserId, "支付并确认建立!", msg, "#FF0000"); err != nil {
		log.Panicf("Error: %s\n", err)
	}
}

func main() {

	ctx := context.Background()
	log.Println("start bot")
	handler := Handler{}

	// Create a bot client
	client = bot.NewBlazeClient(config.GetConfig().ClientID, config.GetConfig().SessionID, config.GetConfig().PrivateKey)

	// Start the loop
	for {
		if err := client.Loop(ctx, handler); err != nil {
			log.Printf("Error: %v\n", err)
		}
		log.Println("connection loop end")
		time.Sleep(time.Second)
	}
}
