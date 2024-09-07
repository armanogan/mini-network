package postgresql

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
)

func Run(host, port, user, password, dbname string) (*pgx.Conn, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	ctx := context.Background()
	connect, err := pgx.Connect(ctx, psqlInfo)
	if err != nil {
		return nil, err
	}
	if err = connect.Ping(ctx); err != nil {
		connect.Close(ctx)
		return nil, err
	}
	return connect, nil
}
