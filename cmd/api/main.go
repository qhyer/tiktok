package main

import "tiktok/cmd/api/initialize"

func main() {
	initialize.Jaeger()
	initialize.Rpc()
	initialize.Router()
}
