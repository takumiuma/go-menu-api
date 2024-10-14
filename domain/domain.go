package domain

type Menu struct{
	MenuId			uint   `json:"menu_id"`
	MenuName	string `json:"menu_name"`
	EatingGenreId uint `json:"eating_genre_id"`
	EatingCategoryId uint `json:"eating_category_id"`
}