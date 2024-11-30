package server

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"github.com/wehw93/task-list/internal/store"

	"github.com/jackc/pgx/v4/pgxpool"
)

type User struct {
	Name      string
	Password  string
	Task      []Task
	Summtasks uint16
	id        int
}

func Hello(users *[]User, countUSers *int, conn *pgxpool.Pool) {
	for {
		fmt.Println("Hello, what do you want?\n1 - Signup; 2 - Login; 3 - Exit")
		var choice uint8
		fmt.Scan(&choice)
		switch choice {
		case 1:
			var newUser User
			*countUSers++
			newUser.signup(conn)
			*users = append(*users, newUser)
			fmt.Printf("Ok! Your name: %s, password: %s\n", newUser.Name, newUser.Password)
		case 2:
			var thisUser User
			if loginin(&thisUser, conn) {
				After_login(&thisUser)
			}
		case 3:
			fmt.Println("Exiting...")
			return
		default:
			fmt.Println("Incorrect input")
		}
	}
}
func After_See(user *User) {
	fmt.Println("1 - check task; 2 - remove task; 3 - exit")
	var choice int
	fmt.Scan(&choice)
	switch choice {
	case 1:
		fmt.Println("Choose task to check: ")
		var choose int
		fmt.Scan(&choose)
		check_task(user, choose)
	case 2:
		user.RemoveTask()
	case 3:
		return
	default:
		fmt.Println("Error")
	}
}

type Task struct {
	description  string
	task_name    string
	date_of_task string
}

func See_All_Tasks(user *User) {
	for i, k := range user.Task {
		fmt.Println(i+1, k.date_of_task, ":", k.task_name)
	}
	fmt.Println("That's all")
	After_See(user)
}
func check_task(user *User, k int) {
	fmt.Println(user.Task[k-1].description)
}

func (user *User) Create_new_task() {
	var thisTask Task
	reader := bufio.NewReader(os.Stdin)
	shit, _ := reader.ReadString('\n')
	shit = strings.TrimSpace(shit)
	fmt.Println("Enter the task name:")
	taskName, _ := reader.ReadString('\n')
	taskName = strings.TrimSpace(taskName)
	thisTask.task_name = taskName

	fmt.Println("Enter the task description:")
	taskDescription, _ := reader.ReadString('\n')
	taskDescription = strings.TrimSpace(taskDescription)
	thisTask.description = taskDescription

	fmt.Println("Enter the due date:")
	taskDueDate, _ := reader.ReadString('\n')
	taskDueDate = strings.TrimSpace(taskDueDate)
	thisTask.date_of_task = taskDueDate

	user.Task = append(user.Task, thisTask)
	userfile := "Users/" + user.Name + ".txt"
	file, err := os.OpenFile(userfile, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("Error of createnewtask")
		os.Exit(1)
	}
	defer file.Close()
	file.WriteString(thisTask.task_name + "\n")
	file.WriteString(thisTask.date_of_task + "\n")
	file.WriteString(thisTask.description + "\n")
	fmt.Println("Task Created Successfully!")
	fmt.Printf("Task Name: %s\n", taskName)
	fmt.Printf("Task Description: %s\n", taskDescription)
	fmt.Printf("Task Due Date: %s\n", taskDueDate)
}

func (user *User) RemoveTask() {
	fmt.Println("Number of task which you wnt to remove")
	var choice uint8
	fmt.Scan(&choice)
	removeIndex(&user.Task, choice-1)

	fmt.Println("Succes")
}
func removeIndex(tasks *[]Task, n uint8) {
	copy((*tasks)[n:], (*tasks)[n+1:])
	*tasks = (*tasks)[:len(*tasks)-1]

}
func After_login(user *User) {
	for {
		var choice int
		fmt.Println("What do you want?\n1 - Create new task; 2 - See all tasks; 3 - Remove tasks; 4 - Exit account;")
		fmt.Scan(&choice)
		switch choice {
		case 1:
			user.Create_new_task()
		case 2:
			See_All_Tasks(user)
		case 3:
			user.RemoveTask()
		case 4:
			fmt.Println("Exiting account...")
			return
		default:
			fmt.Println("Incorrect input")
		}
	}
}

func loginin(thisUser *User, conn *pgxpool.Pool) bool {
	var name_ string
	var pass_ string
	db_user := store.User{}
	var id int
	fmt.Println("Your name: ")
	fmt.Scan(&name_)
	ins := store.Instance{conn}
	row := ins.Db.QueryRow(context.Background(), "SELECT password, id FROM users WHERE $1 = name", name_)
	err := row.Scan(&db_user.Password, &id)
	if err != nil {
		fmt.Println("This name hasn't been found")
		return false
	}
	db_user.Name = name_
	fmt.Println("Password: ")
	fmt.Scan(&pass_)
	if db_user.Password != pass_ {
		fmt.Println("Password incorrect")
		return false
	}
	fmt.Printf("Succes: id :%d ", id)
	thisUser.id = id
	thisUser.Name = name_
	thisUser.Password = db_user.Password
	return true
}

func (u *User) signup(conn *pgxpool.Pool) {
	fmt.Println("Your name: ")
	fmt.Scan(&u.Name)
	fmt.Println("Password: ")
	fmt.Scan(&u.Password)
	userBd := store.User{}
	userBd.Name = u.Name
	userBd.Password = u.Password
	ins := store.Instance{conn}
	ins.CreateUser(userBd)
	fmt.Println("User information saved successfully.")
}
