package main

import (
	"fmt"
	"github.com/jmoiron/sqlx"
)

type MySqlCon struct {
	db *sqlx.DB
}

var err error

func (con *MySqlCon) InitCon(username, password, urlAndPort, database string) {
	if 0 == len(username) {
		fmt.Println("username is empty")
	}
	if 0 == len(password) {
		fmt.Println("password is empty")
	}
	if 0 == len(urlAndPort) {
		fmt.Println("urlAndPort is empty")
	}
	if 0 == len(database) {
		fmt.Println("database is empty")
	}
	con.db, err = sqlx.Open("mysql", username+":"+password+"@tcp("+urlAndPort+")/"+database)
	if err != nil {
		fmt.Println("Mysql Connetecd Error")
	} else {
		fmt.Println("Mysql Connection success")
	}
}

func (con MySqlCon) Insert(item *XingZhengChuFaItem) {
	_, err := con.db.Exec("insert into administrative_punishment_decision (index_num,name,sort,issuer,issue_date,"+
		"issue_num,keywords,content,create_time)"+" values(?,?,?,?,?,?,?,?,?)", item.indexNum, item.name,
		item.sort, item.issuer, item.issueDate, item.issueNum, item.keywords, item.content, item.createdTime)

	if err != nil {
		panic(err)
	}
}
