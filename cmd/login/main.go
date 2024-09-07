package main

import (
	"context"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	"log"
	postgresql2 "mini-network/internal/v1/repository/postgresql"
	"mini-network/internal/v1/token"
	"mini-network/internal/v1/transport/rest"
	"mini-network/internal/v1/usecase"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func dbMigration(connect *pgx.Conn) error {
	config := connect.Config()
	m, err := migrate.New("file://migrations", fmt.Sprintf("pgx5://%s:%s@%s:%d/%s", config.User, config.Password, config.Host, config.Port, config.Database))
	if err != nil {
		return err
	}
	if err = m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}
	return nil
}

func main() {
	err := godotenv.Load("./config/production.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}
	connect, err := postgresql2.Run(os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_NAME"))
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM, syscall.SIGKILL)
	defer connect.Close(ctx)
	err = dbMigration(connect)
	if err != nil {
		log.Fatalf("Migration creation failed: %v", err)
	}
	userCase := usecase.NewUserUseCase(postgresql2.NewUserRepo(connect), token.NewJwt(os.Getenv("JWT_TOKEN"), time.Minute*30))
	httpServer := rest.NewHttpTransport(userCase)
	go func() {
		errorCh := httpServer.Run(os.Getenv("PORT"))
		if errorCh != nil {
			signal.Stop(ch)
			close(ch)
			log.Fatal(errorCh)
		}
	}()
	<-ch
	contextTimeout, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	httpServer.ShutDown(contextTimeout)
	connect.Close(contextTimeout)
	time.Sleep(10 * time.Second)
}
