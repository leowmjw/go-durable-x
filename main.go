package main

import (
	"fmt"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/celestix/gotgproto"
	"github.com/celestix/gotgproto/dispatcher"
	gphandlers "github.com/celestix/gotgproto/dispatcher/handlers"
	"github.com/celestix/gotgproto/dispatcher/handlers/filters"
	gpext "github.com/celestix/gotgproto/ext"
	"github.com/celestix/gotgproto/sessionMaker"
	"github.com/davecgh/go-spew/spew"
	"github.com/gotd/td/tg"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/message"
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
	RunTGBot(botToken)
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
