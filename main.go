package main

import (
	"github.com/flevin58/apiserver/api"
)

func main() {
	srv := api.NewServer().WithAddress(":2000")
	srv.Run()
}
