package main

import(
	"fmt"
	"log"
	"net/http"
	"github.com/graarh/golang-socketio"
	"github.com/graarh/golang-socketio/transport"
)

type Msg struct{
	Title string `json:"title"`
	Content string `json:"content"`
}

func InitSocket(){
	server := gosocketio.NewServer(transport.GetDefaultWebsocketTransport())
	//---------- OnConnection ---------------------------------------
	server.On(gosocketio.OnConnection, func(c *gosocketio.Channel) {
		log.Println("Connected")

		c.Join("author")
	})

	//---------- OnDisconnection ------------------------------------
	server.On(gosocketio.OnDisconnection, func(c *gosocketio.Channel) {
		log.Println("Disconnected")
	})

	//---------- OnRequire ------------------------------------------
	server.On("require", func(c *gosocketio.Channel, msg Msg) string {
			content, err := GetContent(msg.Title)
			if err==nil {
				c.Emit("content", content)
			}else{
				fmt.Println("Failed to Get Article");
				fmt.Println(err);
			}
			return "OK"
	})

	//---------- OnSave ---------------------------------------------
	server.On("save", func(c *gosocketio.Channel, msg Msg) string {
		  fmt.Println(msg.Title)
			err := SaveContent(msg.Title, msg.Content)
			if err==nil {
			}else{
				fmt.Println("Failed to Save Article");
				fmt.Println(err);
			}
			return "OK"
	})

	http.Handle("/socket.io/", server)
}
