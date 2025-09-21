package router

import (
	"go-menu/di"
	"go-menu/middleware"

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

	// システム関連エンドポイント
	{
		systemHandler := di.InitSystemHandler()
		v1.GET("/ping", systemHandler.Ping)
	}

	// メニュー関連エンドポイント（認証不要）
	{
		menuHandler := di.InitTodoHandler()
		v1.GET("/menus", menuHandler.GetAll)
		v1.POST("/menus", menuHandler.CreateMenu)
		v1.PUT("/menus/:menu_id", menuHandler.UpdateMenu)
		v1.DELETE("/menus/:menu_id", menuHandler.DeleteMenu)
		v1.PATCH("/menus/:menu_id/genres", menuHandler.UpdateGenreRelations)
		v1.PATCH("/menus/:menu_id/categories", menuHandler.UpdateCategoryRelations)
	}

	// お気に入り関連エンドポイント（認証必要）
	{
		// Auth0設定とミドルウェアの初期化
		auth0Config := middleware.NewAuth0Config()
		userDriver := di.InitUserDriver()
		authMiddleware := middleware.AuthMiddleware(userDriver, auth0Config)

		favoriteHandler := di.InitFavoriteHandler()

		// 認証が必要なエンドポイントグループ
		authGroup := v1.Group("/favorites")
		authGroup.Use(authMiddleware)
		{
			authGroup.GET("", favoriteHandler.GetFavorites)
			authGroup.POST("", favoriteHandler.AddFavorite)
			authGroup.DELETE("/:menu_id", favoriteHandler.RemoveFavorite)
		}
	}

	return r
}
