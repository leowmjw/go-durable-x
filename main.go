package main

import (
	"fmt"
	"github.com/celestix/gotgproto"
	handlers "github.com/celestix/gotgproto/dispatcher/handlers"
	"github.com/celestix/gotgproto/dispatcher/handlers/filters"
	"github.com/celestix/gotgproto/ext"
	"github.com/celestix/gotgproto/sessionMaker"
	"github.com/davecgh/go-spew/spew"
	"os"
	"strconv"
)

func main() {
	fmt.Println("Test out Telegram App + Bot APIs")
	appID, err := strconv.Atoi(os.Getenv("TELEGRAM_APPID"))
	if err != nil {
		panic(err)
	}
	appHash := os.Getenv("TELEGRAM_APPHASH")
	botToken := os.Getenv("TELEGRAM_BOT_KEY")
	spew.Dump("DUMP:", appID, appHash, botToken)
	Run(appID, appHash, botToken)
}

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
	//dispatcher.AddHandlerToGroup(handlers.NewAnyUpdate(handlers.CallbackResponse), 1)
	dispatcher.AddHandler(handlers.NewMessage(filters.Message.All, summarize))
	// Too low level ..
	//c.Run(context.TODO(), func(ctx context.Context) error {
	//	fmt.Println("INSIDE ...")
	//	return nil
	//})
	ierr := c.Idle()
	if ierr != nil {
		panic(ierr)
	}
}

func summarize(ctx *ext.Context, update *ext.Update) error {
	msg := update.EffectiveMessage
	spew.Dump(msg)
	_, err := ctx.Reply(update, msg.Text, nil)
	return err
}
