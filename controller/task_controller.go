package controller

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/tasuke/udemy/model"
	"github.com/tasuke/udemy/usecase"
	"net/http"
	"strconv"
)

type ITaskController interface {
	GetAllTasks(c echo.Context) error
	GetTaskById(c echo.Context) error
	CreateTask(c echo.Context) error
	UpdateTask(c echo.Context) error
	DeleteTask(c echo.Context) error
}

type taskController struct {
	tu usecase.ITaskUsecase
}

func NewTaskController(tu usecase.ITaskUsecase) ITaskController {
	return &taskController{tu}
}

func (tc taskController) GetAllTasks(c echo.Context) error {
	// ユーザーから送られてくるjwtトークンに組み込まれているユーザーIDの値を取り出す middlewareでデコードされたものが入ってる
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userId := claims["user_id"]

	tasksRes, err := tc.tu.GetAllTasks(uint(userId.(float64)))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, tasksRes)
}

func (tc taskController) GetTaskById(c echo.Context) error {
	// ユーザーから送られてくるjwtトークンに組み込まれているユーザーIDの値を取り出す middlewareでデコードされたものが入ってる
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userId := claims["user_id"]

	// requestのパラメーターを取得してintに変換
	id := c.Param("taskId")
	taskId, _ := strconv.Atoi(id)

	tasksRes, err := tc.tu.GetTaskById(uint(userId.(float64)), uint(taskId))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, tasksRes)
}

func (tc taskController) CreateTask(c echo.Context) error {
	// ユーザーから送られてくるjwtトークンに組み込まれているユーザーIDの値を取り出す middlewareでデコードされたものが入ってる
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userId := claims["user_id"]

	// requestのボディを取得して構造体にバインド
	task := model.Task{}
	if err := c.Bind(&task); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	// コンテキストから取得としたユーザーIDを構造体ユーザーIDにセット
	task.UserId = uint(userId.(float64))
	taskRes, err := tc.tu.CreateTask(task)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusCreated, taskRes)
}

func (tc taskController) UpdateTask(c echo.Context) error {
	// ユーザーから送られてくるjwtトークンに組み込まれているユーザーIDの値を取り出す middlewareでデコードされたものが入ってる
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userId := claims["user_id"]

	// requestのパラメーターを取得してintに変換
	id := c.Param("taskId")
	taskId, _ := strconv.Atoi(id)

	// requestのボディを取得して構造体にバインド
	task := model.Task{}
	if err := c.Bind(&task); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	// コンテキストから取得としたユーザーIDを構造体ユーザーIDにセット
	task.UserId = uint(userId.(float64))
	taskRes, err := tc.tu.UpdateTask(task, uint(userId.(float64)), uint(taskId))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusCreated, taskRes)
}

func (tc taskController) DeleteTask(c echo.Context) error {
	// ユーザーから送られてくるjwtトークンに組み込まれているユーザーIDの値を取り出す middlewareでデコードされたものが入ってる
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userId := claims["user_id"]

	// requestのパラメーターを取得してintに変換
	id := c.Param("taskId")
	taskId, _ := strconv.Atoi(id)

	err := tc.tu.DeleteTask(uint(userId.(float64)), uint(taskId))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.NoContent(http.StatusNoContent)
}
