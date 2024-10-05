package main

import (
	"fmt"
	"goproxy/internal/application/controllers"
	"goproxy/internal/infrastructure/config"
	"goproxy/internal/infrastructure/database"
	"goproxy/internal/infrastructure/logger"
	"goproxy/internal/infrastructure/repositories"
	"goproxy/internal/infrastructure/server"
	"log"
	"net/http"
	"os"
	"runtime"

	"go.uber.org/fx"
)

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
		config.Module,
		controllers.Module,
		repositories.Module,
		database.Module,
		logger.Module,
		server.Module,
		fx.Invoke(func(*http.Server) {}),
		fx.NopLogger,
	)
	if app.Err() != nil {
		fmt.Println(app.Err())
	}
	app.Run()
}
