package store

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type User struct {
	Name     string
	Password string
	id       int
}
type Instance struct {
	Db *pgxpool.Pool
}

func (i *Instance) CreateUser(user User) {
	commandTag, err := i.Db.Exec(context.Background(), "INSERT INTO users (name,password) VALUES($1,$2)", user.Name, user.Password)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v", err)
		return
	}
	fmt.Println(commandTag.String())
	fmt.Println(commandTag.RowsAffected())
}
func (i *Instance) GetAllUsers(ctx context.Context) {
	var users []User
	rows, err := i.Db.Query(ctx, "SELECT name,password,id FROM users;")
	if err == pgx.ErrNoRows {
		fmt.Println("No rows")
		return
	} else if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		user := User{}
		rows.Scan(&user.Name, &user.Password, &user.id)
		users = append(users, user)
	}
	fmt.Println(users)

}
