package main

import "tiktok/cmd/api/initialize"

func main() {
	initialize.Viper()
	initialize.Router()
	initialize.Rpc()
}
