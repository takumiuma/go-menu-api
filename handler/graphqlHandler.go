package handler

import (
	"go-menu/domain"
	"go-menu/usecase"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
)

type GraphQLHandler struct {
	menuUsecase usecase.MenuUsecase
	schema      graphql.Schema
}

func ProvideGraphQLHandler(u usecase.MenuUsecase) *GraphQLHandler {
	h := &GraphQLHandler{menuUsecase: u}
	h.initSchema()
	return h
}

// GraphQL型定義
func (h *GraphQLHandler) initSchema() {
	// Menu型の定義
	menuType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Menu",
		Fields: graphql.Fields{
			"menuId": &graphql.Field{
				Type: graphql.Int,
			},
			"menuName": &graphql.Field{
				Type: graphql.String,
			},
			"genreIds": &graphql.Field{
				Type: graphql.NewList(graphql.Int),
			},
			"categoryIds": &graphql.Field{
				Type: graphql.NewList(graphql.Int),
			},
		},
	})

	// CreateMenuInput型の定義
	createMenuInputType := graphql.NewInputObject(graphql.InputObjectConfig{
		Name: "CreateMenuInput",
		Fields: graphql.InputObjectConfigFieldMap{
			"menuName": &graphql.InputObjectFieldConfig{
				Type: graphql.NewNonNull(graphql.String),
			},
			"genreIds": &graphql.InputObjectFieldConfig{
				Type: graphql.NewList(graphql.Int),
			},
			"categoryIds": &graphql.InputObjectFieldConfig{
				Type: graphql.NewList(graphql.Int),
			},
		},
	})

	// UpdateMenuInput型の定義
	updateMenuInputType := graphql.NewInputObject(graphql.InputObjectConfig{
		Name: "UpdateMenuInput",
		Fields: graphql.InputObjectConfigFieldMap{
			"menuId": &graphql.InputObjectFieldConfig{
				Type: graphql.NewNonNull(graphql.Int),
			},
			"menuName": &graphql.InputObjectFieldConfig{
				Type: graphql.NewNonNull(graphql.String),
			},
			"genreIds": &graphql.InputObjectFieldConfig{
				Type: graphql.NewList(graphql.Int),
			},
			"categoryIds": &graphql.InputObjectFieldConfig{
				Type: graphql.NewList(graphql.Int),
			},
		},
	})

	// UpdateGenreRelationsInput型の定義
	updateGenreRelationsInputType := graphql.NewInputObject(graphql.InputObjectConfig{
		Name: "UpdateGenreRelationsInput",
		Fields: graphql.InputObjectConfigFieldMap{
			"menuId": &graphql.InputObjectFieldConfig{
				Type: graphql.NewNonNull(graphql.Int),
			},
			"genreIds": &graphql.InputObjectFieldConfig{
				Type: graphql.NewList(graphql.Int),
			},
		},
	})

	// UpdateCategoryRelationsInput型の定義
	updateCategoryRelationsInputType := graphql.NewInputObject(graphql.InputObjectConfig{
		Name: "UpdateCategoryRelationsInput",
		Fields: graphql.InputObjectConfigFieldMap{
			"menuId": &graphql.InputObjectFieldConfig{
				Type: graphql.NewNonNull(graphql.Int),
			},
			"categoryIds": &graphql.InputObjectFieldConfig{
				Type: graphql.NewList(graphql.Int),
			},
		},
	})

	// Queryの定義
	queryType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"menus": &graphql.Field{
				Type: graphql.NewList(menuType),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return h.menuUsecase.GetAll()
				},
			},
		},
	})

	// Mutationの定義
	mutationType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Mutation",
		Fields: graphql.Fields{
			"createMenu": &graphql.Field{
				Type: menuType,
				Args: graphql.FieldConfigArgument{
					"input": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(createMenuInputType),
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					input := p.Args["input"].(map[string]interface{})
					menuName := input["menuName"].(string)

					var genreIds []uint
					if gIds, ok := input["genreIds"].([]interface{}); ok {
						for _, id := range gIds {
							genreIds = append(genreIds, uint(id.(int)))
						}
					}

					var categoryIds []uint
					if cIds, ok := input["categoryIds"].([]interface{}); ok {
						for _, id := range cIds {
							categoryIds = append(categoryIds, uint(id.(int)))
						}
					}

					menu := domain.Menu{
						MenuName:    menuName,
						GenreIds:    genreIds,
						CategoryIds: categoryIds,
					}

					return h.menuUsecase.CreateMenu(menu)
				},
			},
			"updateMenu": &graphql.Field{
				Type: menuType,
				Args: graphql.FieldConfigArgument{
					"input": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(updateMenuInputType),
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					input := p.Args["input"].(map[string]interface{})
					menuId := uint(input["menuId"].(int))
					menuName := input["menuName"].(string)

					var genreIds []uint
					if gIds, ok := input["genreIds"].([]interface{}); ok {
						for _, id := range gIds {
							genreIds = append(genreIds, uint(id.(int)))
						}
					}

					var categoryIds []uint
					if cIds, ok := input["categoryIds"].([]interface{}); ok {
						for _, id := range cIds {
							categoryIds = append(categoryIds, uint(id.(int)))
						}
					}

					menu := domain.Menu{
						MenuId:      menuId,
						MenuName:    menuName,
						GenreIds:    genreIds,
						CategoryIds: categoryIds,
					}

					return h.menuUsecase.UpdateMenu(menu)
				},
			},
			"deleteMenu": &graphql.Field{
				Type: graphql.Boolean,
				Args: graphql.FieldConfigArgument{
					"menuId": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Int),
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					menuId := uint(p.Args["menuId"].(int))
					err := h.menuUsecase.DeleteMenu(menuId)
					if err != nil {
						return false, err
					}
					return true, nil
				},
			},
			"updateGenreRelations": &graphql.Field{
				Type: menuType,
				Args: graphql.FieldConfigArgument{
					"input": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(updateGenreRelationsInputType),
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					input := p.Args["input"].(map[string]interface{})
					menuId := uint(input["menuId"].(int))

					var genreIds []uint
					if gIds, ok := input["genreIds"].([]interface{}); ok {
						for _, id := range gIds {
							genreIds = append(genreIds, uint(id.(int)))
						}
					}

					return h.menuUsecase.UpdateGenreRelations(menuId, genreIds)
				},
			},
			"updateCategoryRelations": &graphql.Field{
				Type: menuType,
				Args: graphql.FieldConfigArgument{
					"input": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(updateCategoryRelationsInputType),
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					input := p.Args["input"].(map[string]interface{})
					menuId := uint(input["menuId"].(int))

					var categoryIds []uint
					if cIds, ok := input["categoryIds"].([]interface{}); ok {
						for _, id := range cIds {
							categoryIds = append(categoryIds, uint(id.(int)))
						}
					}

					return h.menuUsecase.UpdateCategoryRelations(menuId, categoryIds)
				},
			},
		},
	})

	// スキーマの作成
	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query:    queryType,
		Mutation: mutationType,
	})
	if err != nil {
		panic(err)
	}

	h.schema = schema
}

// GraphQLHandlerFunc はGraphQLリクエストを処理する
func (h *GraphQLHandler) GraphQLHandlerFunc() gin.HandlerFunc {
	// 環境変数でPlaygroundの有効/無効を制御
	// GIN_MODE=release または GRAPHQL_PLAYGROUND=false で無効化
	enablePlayground := true
	if os.Getenv("GIN_MODE") == "release" {
		enablePlayground = false
	}
	if os.Getenv("GRAPHQL_PLAYGROUND") == "false" {
		enablePlayground = false
	}

	graphqlHandler := handler.New(&handler.Config{
		Schema:     &h.schema,
		Pretty:     enablePlayground, // 本番ではPretty出力も無効化
		GraphiQL:   false,
		Playground: enablePlayground,
	})

	return func(c *gin.Context) {
		graphqlHandler.ServeHTTP(c.Writer, c.Request)
	}
}
