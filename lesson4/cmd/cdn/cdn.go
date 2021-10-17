package main

import (
	"context"
	"gb-backend/lesson4/api/handler"
	"gb-backend/lesson4/api/server"
	"gb-backend/lesson4/app/repos/files"
	"gb-backend/lesson4/app/starter"
	"os"
	"os/signal"
	"sync"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)

	f := files.NewFiles("upload")
	a := starter.NewApp(f)
	h := handler.NewRouter(f)
	srv := server.NewServer(":8080", h)

	wg := &sync.WaitGroup{}
	wg.Add(1)

	go a.Serve(ctx, wg, srv)

	<-ctx.Done()
	cancel()
	wg.Wait()
}
