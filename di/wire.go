package di

import (
	"go-menu/gateway"
	"go-menu/handler"
	"go-menu/resource"
	"go-menu/resource/menu"
	"go-menu/usecase"
)

func InitSystemHandler() *handler.SystemHandler {
	systemHandler := handler.NewSystemHandler()
	return systemHandler
}

func InitTodoHandler() *handler.MenuHandler {
	db := resource.ConnectToDatabase()
	menuDriver := menu.ProvideMenuDriver(db)
	menuPort := gateway.ProvideMenuPort(menuDriver)
	menuUsecase := usecase.ProvideMenuUsecase(menuPort)
	menuHandler := handler.ProvideMenuHandler(menuUsecase)
	return menuHandler
}