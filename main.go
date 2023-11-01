package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/message"
	"github.com/celestix/gotgproto"
	"github.com/celestix/gotgproto/dispatcher"
	gphandlers "github.com/celestix/gotgproto/dispatcher/handlers"
	"github.com/celestix/gotgproto/dispatcher/handlers/filters"
	gpext "github.com/celestix/gotgproto/ext"
	//gpfunctions "github.com/celestix/gotgproto/functions"
	"github.com/celestix/gotgproto/sessionMaker"
	"github.com/davecgh/go-spew/spew"
	"github.com/gotd/td/tg"

	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"

	ggtelegram "github.com/amarnathcjd/gogram/telegram"
)

func main() {
	fmt.Println("Test out Telegram App + Bot APIs")
	appID, err := strconv.Atoi(os.Getenv("TELEGRAM_APPID"))
	if err != nil {
		panic(err)
	}
	appHash := os.Getenv("TELEGRAM_APPHASH")
	botToken := os.Getenv("TELEGRAM_BOT_KEY")
	fmt.Println("DUMP:", appID, appHash, botToken)
	//Run(appID, appHash, botToken)

	// Test as Bot ..
	//RunTGBot(botToken)
	// tets new .. as App
	//RunTGProto(appID, appHash)

	RunGoGram(appID, appHash)
}

func RunTGBot(botToken string) {
	if botToken == "" {
		panic("TELEGRAM_BOT_KEY cannot be empty!!!!")
	}

	// Create bot from environment value.
	b, err := gotgbot.NewBot(botToken, &gotgbot.BotOpts{
		BotClient: &gotgbot.BaseBotClient{
			Client: http.Client{},
			DefaultRequestOpts: &gotgbot.RequestOpts{
				Timeout: time.Second,           // Customise the default request timeout here
				APIURL:  gotgbot.DefaultAPIURL, // As well as the Default API URL here (in case of using local bot API servers)
			},
		},
	})

	if err != nil {
		panic("failed to create new bot: " + err.Error())
	}

	// Create updater and dispatcher.
	updater := ext.NewUpdater(&ext.UpdaterOpts{
		Dispatcher: ext.NewDispatcher(&ext.DispatcherOpts{
			// If an error is returned by a handler, log it and continue going.
			Error: func(b *gotgbot.Bot, ctx *ext.Context, err error) ext.DispatcherAction {
				log.Println("an error occurred while handling update:", err.Error())
				return ext.DispatcherActionNoop
			},
			MaxRoutines: ext.DefaultMaxRoutines,
		}),
	})
	dispatcher := updater.Dispatcher

	sumf := func(b *gotgbot.Bot, ctx *ext.Context) error {
		fmt.Println("INSIDE COMMAND /sum")
		//spew.Dump(ctx.EffectiveChat)
		spew.Dump(ctx)

		fmt.Println("************************************")
		spew.Dump(ctx.EffectiveSender)
		spew.Dump(ctx.EffectiveChat)
		spew.Dump(ctx.EffectiveUser)
		fmt.Println("======================================")
		return nil
	}
	dispatcher.AddHandler(handlers.NewCommand("sum", sumf))
	// NOTE: Order seems to matter ..
	// Add echo handler to reply to all text messages.
	dispatcher.AddHandler(handlers.NewMessage(message.Text, echo))

	// Start receiving updates.
	err = updater.StartPolling(b, &ext.PollingOpts{
		DropPendingUpdates: true,
		GetUpdatesOpts: &gotgbot.GetUpdatesOpts{
			Timeout: 9,
			RequestOpts: &gotgbot.RequestOpts{
				Timeout: time.Second * 10,
			},
		},
	})

	if err != nil {
		panic("failed to start polling: " + err.Error())
	}
	log.Printf("%s has been started...\n", b.User.Username)

	// Idle, to keep updates coming in, and avoid bot stopping.
	updater.Idle()
}

// echo replies to a messages with its own contents.
func echo(b *gotgbot.Bot, ctx *ext.Context) error {
	_, err := ctx.EffectiveMessage.Reply(b, ctx.EffectiveMessage.Text, nil)
	if err != nil {
		return fmt.Errorf("failed to echo message: %w", err)
	}
	return nil
}

// RunGoGram uses gogram to authenticate as an app ..
func RunGoGram(appID int, appHash string) {
	// Test runing basic ..

	// Create a new client
	client, ncerr := ggtelegram.NewClient(ggtelegram.ClientConfig{
		AppID:    int32(appID),
		AppHash:  appHash,
		LogLevel: ggtelegram.LogDebug,
		//StringSession: "", // Uncomment this line to use string session
		Session: "/tmp/gogram",
	})

	if ncerr != nil {
		panic(ncerr)
	}

	// Connect to the server
	if err := client.Connect(); err != nil {
		panic(err)
	}

	// Authenticate the client using the bot token
	// This will send a code to the phone number if it is not already authenticated
	if _, err := client.Login("+60162332450"); err != nil {
		panic(err)
	}

	uo, gmerr := client.GetMe()
	if gmerr != nil {
		panic(gmerr)
	}

	// DEBUG
	//spew.Dump(uo)
	//client.Cache.UpdateUser(uo)

	dialogs, gderr := client.GetDialogs(&ggtelegram.DialogOptions{
		//OffsetID:      0,
		//OffsetDate: 1698220267,
		OffsetPeer: &ggtelegram.InputPeerUser{
			UserID:     uo.ID,
			AccessHash: uo.AccessHash,
		},
		Limit: 10,
		//ExcludePinned: true,
		//FolderID:      0,
	})
	if gderr != nil {
		panic(gderr)
	}
	// DEBUG
	//spew.Dump(dialogs)
	//if err != nil {
	//	panic(err)
	//}
	for _, dialog := range dialogs {
		switch d := dialog.(type) {
		case *ggtelegram.DialogObj:
			//fmt.Println(d.TopMessage)
			switch pt := d.Peer.(type) {

			case *ggtelegram.PeerChat:
				id := pt.ChatID
				fmt.Println(">>> CHAT: ", id)
				co, gcerr := client.GetChat(id)
				if gcerr != nil {
					fmt.Println("ERR: ", gcerr)
					continue
				}
				spew.Dump(co.Title)
			case *ggtelegram.PeerUser:
				id := pt.UserID
				fmt.Println(">>> USER: ", id)
				uo, guerr := client.GetUser(id)
				if guerr != nil {
					fmt.Println("ERR: ", guerr)
					continue
				}
				// For Saved Message is chat with self?
				if uo.Username == "leowmjw" {
					uuf, ugferr := client.UsersGetFullUser(&ggtelegram.InputUserSelf{})
					if ugferr != nil {
						fmt.Println("ERR: ", ugferr)
					}
					spew.Dump(uuf)

					continue
				}
				fmt.Println("USER: ", uo.Username)
			case *ggtelegram.PeerChannel:
				id := pt.ChannelID
				fmt.Println(">>> CHANNEL: ", id)
				c, gcerr := client.GetChannel(id)
				if gcerr != nil {
					//panic(gcerr)
					fmt.Println("ERR: ", gcerr)
					//mf, cgferr := client.ChannelsGetFullChannel(&ggtelegram.InputChannelObj{
					//	ChannelID: id,
					//})
					//if cgferr != nil {
					//	fmt.Println("ERR: ", cgferr)
					//	continue
					//}
					//spew.Dump(mf)
					continue
				}
				spew.Dump(c.Title)
			default:
				fmt.Println("UNKNOWN TYPE: ", pt)
			}
		default:
			fmt.Println("UNKNOWN ...")
			spew.Dump(d)
		}
	}
	//d, gderr := client.GetDialogs()
	//if gderr != nil {
	//	panic(gderr)
	//}
	//spew.Dump(d)

	// Below is the answer we are llking for ..
	mm, msgerr := client.MessagesGetHistory(&ggtelegram.MessagesGetHistoryParams{
		Peer: &ggtelegram.InputPeerSelf{},
		//OffsetDate: 1698824064, // Returns items before this timestamp ..
		Limit: 2,
	})
	if msgerr != nil {
		//fmt.Println("ERR: ", msgerr)
		panic(msgerr)
	}
	spew.Dump(mm)

	//client.Idle()

}

// RunTGProto uses gotgproto to authenticate as an app
func RunTGProto(appID int, appHash string) {

	// Identify as app using the phone number
	clientType := gotgproto.ClientType{
		Phone: "+60162332450",
	}

	// Now newclient to get the context for filtering out subscribed channels/chats
	client, err := gotgproto.NewClient(appID, appHash, clientType,
		&gotgproto.ClientOpts{
			Session:          sessionMaker.NewSession("/tmp/tgproto", sessionMaker.Session),
			AutoFetchReply:   true,
			DisableCopyright: true,
		})

	if err != nil {
		panic(err)
	}
	me := client.Self

	// this works ..
	//client.API().MessagesGetHistory(context.Background(), &tg.MessagesGetHistoryRequest{
	//	Peer:       nil,
	//})

	//inputChannel := &tg.InputChannel{
	//	ChannelID:  123445,
	//	AccessHash: 0,
	//}
	//client.API().ChannelsGetChannels(context.Background(), []tg.InputUser{})
	state, gserr := client.API().UpdatesGetState(context.Background())
	if gserr != nil {
		panic(gserr)
	}
	fmt.Println("STATE:", state.String())
	// DEBUG State ..
	//spew.Dump(state)

	fmt.Printf("client (@%s) has been started...\n", me.Username)

	//gpfunctions.GetMessages(context.Background(), nil, 0, nil)
	//msrc, mgscerr := client.API().MessagesGetSearchResultsCalendar(context.Background(), &tg.MessagesGetSearchResultsCalendarRequest{
	//	Peer:   me.AsInputPeer(),
	//	Filter: filters.MessageFilter(),
	//})
	//if mgscerr != nil {
	//	panic(mgscerr)
	//}
	// FInally can get all the chats!!
	// Is the first one Saved Messages?
	//fmt.Println("DATA:", msrc.GetChats()[0].String())
	// now dump out raw data ...
	//fmt.Println("RAW DATA next ..")
	//spew.Dump(msrc.GetChats())
	//c := gotgproto.Client.API(client)
	//MessagesGetSearchResultsCalendar(context.Background(), &tg.MessagesGetSearchResultsCalendarRequest{
	//	Peer:       nil,
	//	Filter:     nil,
	//	OffsetID:   0,
	//	OffsetDate: 0,
	//})

	fmt.Println("========================= MessagesGetDialogUnreadMarks ====================================")
	mgdu, mgduerr := client.API().MessagesGetDialogUnreadMarks(context.Background())

	if mgduerr != nil {
		panic(mgduerr)
	}
	spew.Dump(mgdu)

	fmt.Println("========================= MessagesGetDialogs ====================================")
	mgd, mgderr := client.API().MessagesGetDialogs(context.Background(), &tg.MessagesGetDialogsRequest{
		//Flags:         0,
		ExcludePinned: true,
		//FolderID:      0,
		//OffsetDate:    0,
		//OffsetID:      0,
		OffsetPeer: me.AsInputPeer(),
		Limit:      10,
		//Hash:          0,
	})
	if mgderr != nil {
		panic(mgderr)
	}
	spew.Dump(mgd)

	fmt.Println("========================= MessagesGetDialogFilters ====================================")
	mgdf, mgdferr := client.API().MessagesGetDialogFilters(context.Background())
	if mgdferr != nil {
		panic(mgdferr)
	}
	spew.Dump(mgdf)

	ierr := client.Idle()
	if ierr != nil {
		panic(ierr)
	}

}

// Run uses gotgproto .. not working  :(
func Run(appID int, botToken, appHash string) {
	// Identify as bot ..
	clientType := gotgproto.ClientType{
		BotToken: botToken,
	}
	// Setuo pmemory session only; what;s the diff?
	c, cerr := gotgproto.NewClient(appID, appHash, clientType,
		&gotgproto.ClientOpts{
			// we can use file location for session; the sqlite will be stored there ..
			Session: sessionMaker.NewSession("/tmp/summarizer", sessionMaker.Session),
			//Session: sessionMaker.NewInMemorySession(
			//	"summarizer",
			//	sessionMaker.Session,
			//),
			AutoFetchReply:   true,
			DisableCopyright: true,
		})
	if cerr != nil {
		fmt.Println("ERR cerr!!")
		panic(cerr)
	}

	fmt.Println("OK!!!")
	// Identify the user
	fmt.Println("USER: ", c.Self.Username)
	// secret stuff
	xss, xerr := c.ExportStringSession()
	if xerr != nil {
		fmt.Println("ERR xerr!!")
		panic(xerr)
	}
	fmt.Println("SECRET stuff ..", xss)

	dispatcher := c.Dispatcher
	//c.Self
	//c.CreateContext().GetChat(300).
	// Command Handler for /start
	dispatcher.AddHandler(gphandlers.NewCommand("start", start))
	// Callback Query Handler with prefix filter for recieving specific query
	dispatcher.AddHandler(gphandlers.NewCallbackQuery(filters.CallbackQuery.Prefix("cb_"), buttonCallback))

	//dispatcher.AddHandlerToGroup(handlers.NewAnyUpdate(handlers.CallbackResponse), 1)
	dispatcher.AddHandler(gphandlers.NewMessage(filters.Message.All, summarize))
	// Too low level ..
	//c.Run(context.TODO(), func(ctx context.Context) error {
	//	fmt.Println("INSIDE ...")
	//	return nil
	//})

	fmt.Printf("client (@%s) has been started...\n", c.Self.Username)

	ierr := c.Idle()
	if ierr != nil {
		panic(ierr)
	}
}

func summarize(ctx *gpext.Context, update *gpext.Update) error {
	msg := update.EffectiveMessage
	spew.Dump(msg)
	_, err := ctx.Reply(update, msg.Text, nil)
	return err
}

// callback function for /start command
func start(ctx *gpext.Context, update *gpext.Update) error {
	user := update.EffectiveUser()
	_, _ = ctx.Reply(update, fmt.Sprintf("Hello %s, I am @%s and will repeat all your messages.\nI was made using gotd and gotgproto.", user.FirstName, ctx.Self.Username), &gpext.ReplyOpts{
		Markup: &tg.ReplyInlineMarkup{
			Rows: []tg.KeyboardButtonRow{
				{
					Buttons: []tg.KeyboardButtonClass{
						&tg.KeyboardButtonURL{
							Text: "gotd/td",
							URL:  "https://github.com/gotd/td",
						},
						&tg.KeyboardButtonURL{
							Text: "gotgproto",
							URL:  "https://github.com/celestix/gotgproto",
						},
					},
				},
				{
					Buttons: []tg.KeyboardButtonClass{
						&tg.KeyboardButtonCallback{
							Text: "Click Here",
							Data: []byte("cb_pressed"),
						},
					},
				},
			},
		},
	})
	// End dispatcher groups so that bot doesn't echo /start command usage
	return dispatcher.EndGroups
}

func buttonCallback(ctx *gpext.Context, update *gpext.Update) error {
	query := update.CallbackQuery
	_, _ = ctx.AnswerCallback(&tg.MessagesSetBotCallbackAnswerRequest{
		Alert:   true,
		QueryID: query.QueryID,
		Message: "This is an example bot!",
	})
	return nil
}
