package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

type exampleData struct {
	Process string
	User    string
	Port    string
	Ip      string
}

type ipData struct {
	Ip string
	exampleData
}

func main() {
	db, err := sql.Open("postgres", "user=gtuiw password=NG4HYQP4bCg@ dbname=db_test sslmode=disable")
	checkErr(err)

	rows, err := db.Query("select * from regs")
	checkErr(err)
	for rows.Next() {
		var id int
		var ip, process, user, port, created_at, updated_at string

		var queryIp ipData
		err = rows.Scan(&id, &ip, &process, &user, &port, &created_at, &updated_at)
		checkErr(err)
		queryIp.Ip = ip
		queryIp.Port = port
		queryIp.Process = process
		queryIp.User = user

		fmt.Println(queryIp.exampleData)
	}
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
