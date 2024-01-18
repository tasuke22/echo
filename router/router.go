package router

import (
	"github.com/labstack/echo/v4"
	"github.com/tasuke/udemy/controller"
)

// なんでここには、今までにあった構造体やインターフェースがいらないの？ => NewRouter関数はルーティングの設定を行うための「接続点」として機能し、自身で状態を持つ構造体や追加のインターフェースを必要としません。
func NewRouter(uc controller.IUserController) *echo.Echo {
	e := echo.New()
	e.POST("/signup", uc.SignUp)
	e.POST("/login", uc.LogIn)
	e.POST("/logout", uc.LogOut)
	return e
}
