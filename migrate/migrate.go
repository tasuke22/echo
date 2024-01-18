package main

import (
	"fmt"
	"github.com/tasuke/udemy/db"
	"github.com/tasuke/udemy/model"
)

func main() {
	dbConn := db.NewDB()
	defer fmt.Println("Successfully Migrated")
	defer db.CloseDB(dbConn)
	dbConn.AutoMigrate(&model.User{}, &model.Task{})
}
