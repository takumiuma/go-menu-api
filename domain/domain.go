package domain

// レスポンス用のメニュー情報
type Menu struct{
	MenuId			uint   `json:"menu_id"`
	MenuName	string `json:"menu_name"`
	GenreIds []uint `json:"genre_ids"`
	CategoryIds []uint `json:"category_ids"`
}