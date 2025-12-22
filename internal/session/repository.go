package session

type Repository interface {
	func(accessToken string) (*int64, error)
}
