package domain

const AppStorageKey = "app"

type App struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}
