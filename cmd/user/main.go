package main

import (
	"log"

	user "tiktok/kitex_gen/user/usersrv"
)

func main() {
	svr := user.NewServer(new(UserSrvImpl))

	err := svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
