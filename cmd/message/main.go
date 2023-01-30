package main

import (
	"log"

	message "tiktok/kitex_gen/message/messagesrv"
)

func main() {
	svr := message.NewServer(new(MessageSrvImpl))

	err := svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
