package router

import (
	"go-menu/di"

	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func NewServer() *gin.Engine {
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		// アクセスを許可したいアクセス元
		AllowOrigins: []string{"*"},
		// アクセスを許可したいHTTPメソッド
		AllowMethods: []string{
			"POST",
			"GET",
			"OPTIONS",
			"PUT",
			"DELETE",
			"PATCH",
		},
		// 許可したいHTTPリクエストヘッダ
		AllowHeaders: []string{
			"Accept",
			"Accept-Encoding",
			"Accept-Language",
			"Access-Control-Allow-Credentials",
			"Access-Control-Allow-Headers",
			"Content-Type",
			"Content-Length",
			"Authorization",
			"Connection",
			"Host",
			"Origin",
			"Referer",
			"Sec-Ch-Ua",
			"Sec-Ch-Ua-Mobile",
			"Sec-Ch-Ua-Platform",
			"Sec-Fetch-Dest",
			"Sec-Fetch-Mode",
			"Sec-Fetch-Site",
			"User-Agent",
		},
		// cookieなどの情報を必要とするかどうか
		AllowCredentials: false,
		// preflightリクエストの結果をキャッシュする時間
		MaxAge: 24 * time.Hour,
	}))

	v1 := r.Group("/v1")

	{
		menuHandler := di.InitTodoHandler()
		v1.GET("/menus", menuHandler.GetAll)
		v1.POST("/menus", menuHandler.CreateMenu)
		v1.PUT("/menus/:menu_id", menuHandler.UpdateMenu)
		v1.DELETE("/menus/:menu_id", menuHandler.DeleteMenu)
		v1.PATCH("/menus/:menu_id/genres", menuHandler.UpdateGenreRelations)
		v1.PATCH("/menus/:menu_id/categories", menuHandler.UpdateCategoryRelations)
	}

	return r
}
