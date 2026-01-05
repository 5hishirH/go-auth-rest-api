package jwt

type Jwt interface {
	Generate(data any, expiry int) (*string, error)
	Compare(t *string) (*bool, error)
}

// type JwtService struct {
// 	secret string
// }

// func (j *JwtService) Generate(data *map[string]any)
