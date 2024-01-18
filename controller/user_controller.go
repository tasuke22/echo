package controller

import (
	"github.com/labstack/echo/v4"
	"github.com/tasuke/udemy/model"
	"github.com/tasuke/udemy/usecase"
	"net/http"
	"os"
	"time"
)

type IUserController interface {
	SignUp(c echo.Context) error
	LogIn(c echo.Context) error
	LogOut(c echo.Context) error
}

type userController struct {
	uu usecase.IUserUsecase
}

func NewUserController(uu usecase.IUserUsecase) IUserController {
	return &userController{uu}
}

func (uc userController) SignUp(c echo.Context) error {
	user := model.User{}
	// Bindは、リクエストのボディを構造体にバインドする
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	// userは値渡しなのはなぜ？ => SignUp関数内で変更されないようにするため。これにより、関数外部のuserオブジェクトの状態が保持されます。
	// SignUpしてユーザーの情報が変わるのはおかしいよね
	userRes, err := uc.uu.SignUp(user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, userRes)
}

func (uc userController) LogIn(c echo.Context) error {
	user := model.User{}
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	// Loginメソッドはjwtのトークンを文字列として返す
	tokenString, err := uc.uu.Login(user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	// jwtトークンをサーバーサイドのcookieに保存する
	cookie := new(http.Cookie)
	cookie.Name = "token"
	cookie.Value = tokenString
	cookie.Expires = time.Now().Add(24 * time.Hour)
	cookie.Path = "/"
	cookie.Domain = os.Getenv("API_DOMAIN")
	//cookie.Secure = true // postmanで確認するためtrueにしている
	cookie.HttpOnly = true                  // jsからcookieにアクセスできないようにする
	cookie.SameSite = http.SameSiteNoneMode // フロントエンドとバックエンドが異なるドメインの場合の設定
	c.SetCookie(cookie)                     // httpレスポンスに含める
	return c.NoContent(http.StatusOK)
}

func (uc userController) LogOut(c echo.Context) error {
	cookie := new(http.Cookie)
	cookie.Name = "token"
	cookie.Value = ""
	cookie.Expires = time.Now()
	cookie.Path = "/"
	cookie.Domain = os.Getenv("API_DOMAIN")
	//cookie.Secure = true // postmanで確認するためtrueにしている
	cookie.HttpOnly = true                  // jsからcookieにアクセスできないようにする
	cookie.SameSite = http.SameSiteNoneMode // フロントエンドとバックエンドが異なるドメインの場合の設定
	c.SetCookie(cookie)                     // httpレスポンスに含める
	return c.NoContent(http.StatusOK)
}
