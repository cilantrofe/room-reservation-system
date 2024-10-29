package main

import (
	"BookingSvc/internal/app"
	"context"
	"fmt"
)

func main() {
	service := app.NewApp()
	ctxWithCancel, cancel := context.WithCancel(context.Background())
	defer cancel()
	err := service.Init(ctxWithCancel)
	if err != nil {
		return
	}

	err = service.Start(ctxWithCancel)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

}
