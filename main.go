package main

import (
	"AuthService/db"
	"AuthService/server"
	"log"
)

func main() {
	err := db.InitDB("root", "vdySqAwCIwMHUfdUyqaQlBOBlCrZovdD", "centerbeam.proxy.rlwy.net:36885", "railway")
	if err != nil {
		log.Fatalf("DB error: %v", err)
	}
	//log.Println("БД успешно подключена")

	server.Start()
}
