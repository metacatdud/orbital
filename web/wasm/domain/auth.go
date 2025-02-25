package domain

const (
	AuthStorageKey = "auth"
)

type Auth struct {
	SecretKey string `json:"secretKey"`
}
