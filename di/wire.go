package di

import (
	"go-menu/gateway"
	"go-menu/handler"
	"go-menu/resource"
	"go-menu/resource/menu"
	"go-menu/resource/user"
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

func InitFavoriteHandler() *handler.FavoriteHandler {
	db := resource.ConnectToDatabase()
	userDriver := user.ProvideUserDriver(db)
	favoriteHandler := handler.ProvideFavoriteHandler(userDriver)
	return favoriteHandler
}

func InitUserDriver() user.UserDriver {
	db := resource.ConnectToDatabase()
	userDriver := user.ProvideUserDriver(db)
	return userDriver
}

func InitUserHandler() *handler.UserHandler {
	db := resource.ConnectToDatabase()
	userDriver := user.ProvideUserDriver(db)
	userHandler := handler.ProvideUserHandler(userDriver)
	return userHandler
}
