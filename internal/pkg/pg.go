package pkg

import (
	"context"
	"fmt"
	"net/url"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	DbName   string
	TimeOut  int
}

func NewPoolConfig(cfg *Config) (*pgxpool.Config, error) {
	connStr := fmt.Sprintf("%s://%s:%s@%s:%s/%s?sslmode=disable&connect_timeout=%d", "postgres", url.QueryEscape(cfg.Username), url.QueryEscape(cfg.Password), cfg.Host, cfg.Port, cfg.DbName, cfg.TimeOut)
	poolConfig, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v", err)
		return nil, err
	}
	return poolConfig, nil
}
func NewConnection(poolconfig *pgxpool.Config) (*pgxpool.Pool, error) {
	conn, err := pgxpool.ConnectConfig(context.Background(), poolconfig)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v", err)
		return nil, err
	}
	return conn, nil
}
