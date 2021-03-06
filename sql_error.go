package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
)

const (
	dbDriverName = "sqlite3"
	dbName       = "./chogori.db"
)

var db *sql.DB

type User struct {
	Name  string
	Age   int
	Job   string
	Hobby string
}

// 处理思路:
//dao:
//
// return errors.Wrapf(errors.NotFound, fmt.Sprintf("sql: %s error: %v", sql, err))
//
//
//biz:
//
//if errors.Is(err, errors.NotFound} {
//
//}
//
//if errors.Reason(err, xxxx) == xxxx {
//
//}

func DBInit() {
	var e error
	db, e = sql.Open(dbDriverName, dbName)
	if e != nil {
		log.Fatalln("connect db failed.")
	}
	log.Println("db connnect success")
}

func DBClose() {
	db.Close()
	log.Println("db connect closed.")
}

func createTable() error {
	q := `create table if not exists "user" (
		"id" integer primary key autoincrement,
		"name"  text not null,
		"age" integer not null,
		"job" text,
		"hobby" text
	)`
	_, e := db.Exec(q)
	if e != nil {
		return errors.Wrap(e, "create user table error.")
	}
	fmt.Println("db table create success.")
	return nil
}

func InsertUser(user User) error {
	q := "insert into user(name, age, job, hobby) values(?,?,?,?)"
	stmt, e := db.Prepare(q)
	if e != nil {
		return e
	}
	defer stmt.Close()
	_, e = stmt.Exec(user.Name, user.Age, user.Job, user.Hobby)
	return e
}

func QueryUserNameById(id int) (string, error) {
	var name string
	q := `select name from user where id = ?`
	e := db.QueryRow(q, id).Scan(&name)
	if errors.Is(e, sql.ErrNoRows) {
		return "", errors.Wrapf(sql.ErrNoRows, "user id=%v not found", id)
	} else if e != nil {
		return "", errors.Wrap(e, "db error.")
	}
	return name, nil
}

func main() {
	fmt.Println("vim-go")
	os.Remove(dbName)

	DBInit()
	defer DBClose()

	createTable()

	_, e := QueryUserNameById(1)
	if errors.Cause(e) == sql.ErrNoRows {
		log.Printf("user not found, %v\n", e)
		log.Printf("%+v\n", e)
	} else if e != nil {
		log.Println("db error.", e)
	}

	for i := 0; i < 5; i++ {
		user := User{
			Name:  fmt.Sprintf("user-%v", i),
			Age:   20 + 2*i,
			Job:   "Gopher",
			Hobby: "Play Game",
		}
		if e := InsertUser(user); e != nil {
			log.Fatalln("insert user to db failed.")
		}
	}

	log.Println("insert user success.")

	for i := 0; i < 6; i++ {
		name, e := QueryUserNameById(i)
		if errors.Cause(e) == sql.ErrNoRows {
			log.Printf("user not found, %v\n", e)
			log.Printf("%+v\n", e)
		} else if e != nil {
			log.Println("db error.", e)
		} else {
			log.Printf("find user id=%v, name is %v\n", i, name)
		}
	}
}
