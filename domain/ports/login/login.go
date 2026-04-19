package login

import "github.com/mlbautomation/Ecommmerce_MLB/model"

type Service interface {
	Login(email, password, jwtSecretKey string) (model.User, string, error)
}
