package main

import (
	"fmt"
	"os"
	"github.com/wehw93/task-list/internal/server"
	"github.com/wehw93/task-list/internal/pkg"
	
)

func main() {
	cfg := &pkg.Config{}
	cfg.Host = "localhost"
	cfg.Password = "pwd123"
	cfg.Username = "db_user"
	cfg.Port = "54320"
	cfg.DbName = "db_task_list"
	cfg.TimeOut = 5
	poolConfig, err := pkg.NewPoolConfig(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "erorr %v", err)
		os.Exit(1)
	}
	conn, err := pkg.NewConnection(poolConfig)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v", err)
		os.Exit(1)
	}
	var users []server.User
	countUsers := 0
	server.Hello(&users, &countUsers, conn)
}
