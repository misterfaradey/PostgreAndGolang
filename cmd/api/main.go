package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/misterfaradey/PostgreAndGolang/internal/default_logger"
	"github.com/misterfaradey/PostgreAndGolang/internal/server"
	"github.com/misterfaradey/PostgreAndGolang/internal/server/controllers"
	"github.com/misterfaradey/PostgreAndGolang/internal/storage"
	"log"

	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const mode = gin.DebugMode //gin.DebugMode gin.ReleaseMode

func main() {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger := default_logger.NewDefaultLogger()

	db := storage.NewCarStorage()

	err := db.Connect(
		&storage.StorageConf{ //todo read from file
			DbHost:     "localhost",
			DbPort:     "5432",
			DbName:     "test",
			DbScheme:   "test",
			DbUser:     "postgres",
			DbPassword: "postgres",
		})
	if err != nil {
		logger.Println(err)
		cancel()
		return
	}

	err = db.InitDB()
	if err != nil {
		logger.Println(err)
		cancel()
		return
	}

	go db.HealthChecker(logger)

	srv := server.NewServer(
		controllers.NewMethodController(db),
		&server.ServerConf{ //todo read from file
			GinMode:        mode,
			Address:        ":8080",
			ReadTimeout:    time.Second * 3,
			WriteTimeout:   time.Second * 3,
			MaxHeaderBytes: 10240,
		},
	)

	go runServer(cancel, srv, logger)

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGKILL)

	logger.Println("Running ...")

	select {
	case <-signalChan:
	case <-ctx.Done():
	}

	ctx, cancel = context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	logger.Println("Closing application ...")

	closeServer(ctx, srv, logger)
	closeStorage(ctx, db, logger)

}

func runServer(cancel context.CancelFunc, server server.Server, logger *log.Logger) {

	err := server.Run()
	if err == http.ErrServerClosed {
		logger.Println("server closed")
		return
	}

	if err != nil {
		logger.Println(err)
		cancel()
	}
}

func closeServer(ctx context.Context, server server.Server, logger *log.Logger) {

	err := server.Shutdown(ctx)
	if err != nil {
		logger.Println(err)
	}
}

func closeStorage(ctx context.Context, store storage.DB, logger *log.Logger) {

	select {
	case <-ctx.Done():
		logger.Println(ctx.Err())
	default:
		err := store.Close()
		if err != nil {
			logger.Println(err)
		}
		logger.Println("storage closed")
	}

}
