package starter

import (
	"context"
	"gb-backend/lesson4/app/repos/files"
	"sync"
)

type App struct {
	f *files.Files
}

func NewApp(f *files.Files) *App {
	a := &App{
		f: f,
	}
	return a
}

type APIServer interface {
	Start(f *files.Files)
	Stop()
}

func (a *App) Serve(ctx context.Context, wg *sync.WaitGroup, hs APIServer) {
	defer wg.Done()
	hs.Start(a.f)
	<-ctx.Done()
	hs.Stop()
}
