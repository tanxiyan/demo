package main

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var db *sqlx.DB

func initDB() (err error) {
	dsn := "root:123456@tcp(127.0.0.1:3306)/book"
	db, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		fmt.Printf("connent DB failed , err:%v\n", err)
		return
	}
	db.SetMaxOpenConns(20) //设置最大连接数
	db.SetMaxIdleConns(10)
	return
}

// 查询所以有数据
func queryAllBook() (bookList []*Book, err error) {
	sqlStr := "select id , title , price from book"
	err = db.Select(&bookList, sqlStr)
	if err != nil {
		fmt.Printf("查询信息失败err = %v\n", err)
		return
	}
	return
}

// 查询单条书籍
func querySingalBook(id int64) (book Book, err error) {
	sqlstr := "select id , title , price from book where id = ?"
	err = db.Get(&book, sqlstr, id)
	if err != nil {
		fmt.Printf("查询信息失败err=%v\n", err)
		return
	}
	return
}

//插入数据

func insertAllBook(title string, price float64) (err error) {
	sqlstr := "insert into book (title ,price) values (? , ?)"
	_, err = db.Exec(sqlstr, title, price)
	if err != nil {
		fmt.Printf("添加信息失败err=%v\n", err)
		return
	}
	return
}

// 删除数据
func deleteBook(id int64) (err error) {
	sqlstr := "delete from book where id = ?"
	_, err = db.Exec(sqlstr, id)
	if err != nil {
		fmt.Printf("删除信息失败err=%v\n", err)
		return
	}
	return
}

// 更新数据
func updateBook(title string, price float64, id int64) (err error) {
	sqlstr := "update book set title = ? , price = ? where id = ?"
	_, err = db.Exec(sqlstr, title, price, id)
	if err != nil {
		fmt.Printf("更新信息失败err=%v\n", err)
		return
	}
	return
}

// 查询用户信息
func GetUserInfo(username string) (user User, err error) {
	sqlstr := "select id , username , password from user where username = ?"
	err = db.Get(&user, sqlstr, username)
	if err != nil {
		fmt.Printf("查询信息失败err=%v\n", err)
		return
	}
	return
}
