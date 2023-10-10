package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"github.com/mdp/qrterminal"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
	"google.golang.org/protobuf/proto"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	conn := Connect()
	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"contact": "nandarusfikri@gmail.com",
			"message": "Whatsapp Gateway",
			"Author":  "NandaRusfikri",
		})
	})
	router.POST("/send_message", func(c *gin.Context) {
		var input Message

		if err := c.ShouldBindJSON(&input); err != nil {
			c.Abort()
			return
		}
		if err := SendWAMessage(conn, input.DestinationNumber, input.Message); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "failed",
				"message": err.Error(),
			})
		}

		c.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": "Pesan berhasil dikirim",
			"nomor":   input.DestinationNumber,
			"pesan":   input.Message,
		})
	})

	err := router.Run(":33133")
	if err != nil {
		log.Fatal(err)
	}

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	conn.Disconnect()
}

type Message struct {
	DestinationNumber string `json:"destination_number" binding:"required"`
	Message           string `json:"message" binding:"required"`
}

func eventHandler(evt interface{}) {
	switch v := evt.(type) {
	case *events.Message:
		fmt.Println("Received a message!", v.Message.GetConversation())
	}
}

func SendWAMessage(client *whatsmeow.Client, phone string, message string) error {
	// Buat konteks
	ctx := context.Background()

	targetJID := types.NewJID(phone, types.DefaultUserServer)

	msg := &waProto.Message{
		Conversation: proto.String(message),
	}

	// Kirim pesan menggunakan client
	_, err := client.SendMessage(ctx, targetJID, msg)
	if err != nil {
		return err
	}

	return nil
}

func Connect() *whatsmeow.Client {
	dbLog := waLog.Stdout("Database", "DEBUG", true)
	// Make sure you add appropriate DB connector imports, e.g. github.com/mattn/go-sqlite3 for SQLite
	container, err := sqlstore.New("sqlite3", "file:store_whatsapp.db?_foreign_keys=on", dbLog)
	if err != nil {
		log.Fatalln(err)
	}
	// If you want multiple sessions, remember their JIDs and use .GetDevice(jid) or .GetAllDevices() instead.
	deviceStore, err := container.GetFirstDevice()
	if err != nil {
		panic(err)
	}
	clientLog := waLog.Stdout("Client", "DEBUG", true)
	client := whatsmeow.NewClient(deviceStore, clientLog)
	client.AddEventHandler(eventHandler)

	if client.Store.ID == nil {
		// No ID stored, new login
		qrChan, _ := client.GetQRChannel(context.Background())
		err = client.Connect()
		if err != nil {
			panic(err)
		}
		for evt := range qrChan {
			if evt.Event == "code" {
				// Render the QR code here
				// e.g. qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stdout)
				qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stdout)
				// or just manually `echo 2@... | qrencode -t ansiutf8` in a terminal
				fmt.Println("QR code:", evt.Code)
			} else {
				fmt.Println("Login event:", evt.Event)
			}
		}
	} else {
		// Already logged in, just connect
		err = client.Connect()
		if err != nil {
			panic(err)
		}
	}

	return client
}
