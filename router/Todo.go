package router

import (
	"goapi_base/service"

	"github.com/labstack/echo/v4"
)

func TodoRouter(e *echo.Group) {
	e.POST("/Todo", service.TodoAdd)
	e.GET("/Todo", service.TodoList)
	e.GET("/Todo/:id", service.TodoGet)
	e.DELETE("/Todo/:id", service.TodoDelete)
	e.PUT("/Todo/:id", service.TodoSet)
}
