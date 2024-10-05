package main

import (
	"context"
	"fmt"
	"goproxy/internal/infrastructure/config"
	"goproxy/internal/infrastructure/database"
	"goproxy/internal/infrastructure/logger"
	"goproxy/internal/infrastructure/repositories"
	"log"
	"os"
	"runtime"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

var Module = fx.Options(
	fx.Invoke(RegisterHandler),
)

func RegisterHandler(lf fx.Lifecycle, log *zap.Logger, repo *repositories.BlockedIpsRepository, shutdowner fx.Shutdowner) {
	lf.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				url := os.Args
				if len(url) == 1 {
					log.Fatal("missing command")
				}
				switch url[1] {
				case "block":
					if len(url) <= 2 {
						log.Fatal("missing ip to block")
					}
					log.Info(fmt.Sprintf("blocking ip %s", url[1]))
					err := repo.InsertBlockedIP(url[2])
					if err != nil {
						log.Error(fmt.Sprintf("failed to block ip %s", url[2]))
						shutdowner.Shutdown()
						return nil
					}
					log.Info(fmt.Sprintf("ip %s blocked", url[2]))
					shutdowner.Shutdown()
					return nil
				case "remove":
					if len(url) <= 2 {
						log.Fatal("missing ip to remove")
					}
					log.Info(fmt.Sprintf("removing ip %s", url[1]))
					err := repo.RemoveBlockedIP(url[2])
					if err != nil {
						log.Error(fmt.Sprintf("failed to remove ip %s", url[2]))
						shutdowner.Shutdown()
						return nil
					}
					log.Info(fmt.Sprintf("ip %s removed", url[2]))
					shutdowner.Shutdown()
					return nil
				}
				log.Info(fmt.Sprintf("command %s not supported", url[1]))
				shutdowner.Shutdown()
				return nil
			},
		},
	)
}

func main() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Fprintf(os.Stderr, "exception: %v\n", err)
			buf := make([]byte, 1024)
			n := runtime.Stack(buf, false)
			log.Printf("Stack trace:\n%s", buf[:n])
			os.Exit(1)
		}
	}()
	app := fx.New(
		Module,
		config.Module,
		repositories.Module,
		database.Module,
		logger.Module,
		fx.NopLogger,
	)
	if app.Err() != nil {
		fmt.Println(app.Err())
	}
	app.Run()
}
