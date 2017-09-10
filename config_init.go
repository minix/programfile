package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/BurntSushi/toml"
	_ "github.com/lib/pq"
	"reflect"
)

type Config struct {
	Database dbConfig
	Server   ServerConfig
}

type ServerConfig struct {
	Host string       `toml:"host"`
	Port string       `toml:"port"`
	Init []InitServer `toml:"init"`
}

type InitServer struct {
	Command string `toml:"command"`
	User    string `toml:"user"`
	Port    string `toml:"port"`
}

type dbConfig struct {
	Adapter  string `toml:"adapter"`
	Database string `toml:"database"`
	User     string `toml:"user"`
	Password string `toml:"password"`
	Sslmode  string `toml:"sslmode"`
}

type ipData struct {
	Ip string
	InitServer
	Data []InitServer
}

func main() {
	var conf Config
	var targetBit bool = true
	var s, queryIp ipData

	str := `{"Ip":"10.1.1.99","data":[{"Command":"grpc","User":"root","Port":"8999"},{"Command":"agentd","User":"root","Port":"10090"},{"Command":"nginx","User":"root","Port":"80"},{"Command":"redis","User":"root","Port":"6379"}]}`

	_, err := toml.DecodeFile("./init.toml", &conf)
	checkErr(err)

	var adapter string = conf.Database.Adapter

	connString := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=%s", conf.Database.User, conf.Database.Password, conf.Database.Database, conf.Database.Sslmode)
	db, err := sql.Open(adapter, connString)
	checkErr(err)

	if err := json.Unmarshal([]byte(str), &s); err != nil {
		panic(err)
	}
	rows, err := db.Query("select * from regs where ip=$1", s.Ip)
	checkErr(err)

	for _, comeinData := range s.Data {
		//如果数据相同，就不用跳出循环，并将标志位设置为false.当不一样时，标志位的结果是true, 这样循环后会被打印
		for _, initData := range conf.Server.Init {
			if reflect.DeepEqual(comeinData, initData) {
				targetBit = false
				break
			}
		}
		//如果标志位为true,刚打开不一致的数据
		if targetBit {
			for rows.Next() {
				var id int
				var ip, process, user, port, created_at, updated_at string

				err = rows.Scan(&id, &ip, &process, &user, &port, &created_at, &updated_at)
				checkErr(err)
				queryIp.Ip = ip
				queryIp.Port = port
				queryIp.Command = process
				queryIp.User = user
			}
			if reflect.DeepEqual(comeinData, queryIp.InitServer) {
				continue
			} else {
				fmt.Println(comeinData)
			}
		}
		//设置标志位为默认值
		targetBit = true
	}

	defer db.Close()
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
