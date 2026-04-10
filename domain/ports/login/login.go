package login

import "github.com/mlbautomation/ProyectoEMLB/model"

type Service interface {
	Login(email, password, jwtSecretKey string) (model.User, string, error)
}
