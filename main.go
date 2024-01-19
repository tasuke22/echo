package main

import (
	"github.com/tasuke/udemy/controller"
	"github.com/tasuke/udemy/db"
	"github.com/tasuke/udemy/repository"
	"github.com/tasuke/udemy/router"
	"github.com/tasuke/udemy/usecase"
)

func main() {
	// ここから連続で依存性の注入。外側でインスタンス化したものを注入
	// IFが返ってきてそれを注入しているので疎結合
	db := db.NewDB()
	userRepository := repository.NewUserRepository(db)
	userUsecase := usecase.NewUserUsecase(userRepository)
	userController := controller.NewUserController(userUsecase)
	e := router.NewRouter(userController)
	// エラーが発生した場合はログ情報を出力してプログラムを強制終了
	e.Logger.Fatal(e.Start(":8080"))
}
