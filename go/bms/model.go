package main

type Book struct {
	ID    int64   `db:"id"` //和数据库相互联系加一个tag
	Title string  `db:"title"`
	Price float64 `db:"price"`
}
type User struct {
	ID       int64  `db:"id"`
	Username string `db:"username"`
	Password string `db:"password"`
}
