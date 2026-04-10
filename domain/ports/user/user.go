package user

import "github.com/mlbautomation/ProyectoEMLB/model"

type Repository interface {
	Create(m *model.User) error
	GetByEmail(email string) (model.User, error)
	GetAll() (model.Users, error)
}

type Service interface {
	Create(m *model.User) error
	GetByEmail(email string) (model.User, error)
	GetAll() (model.Users, error)
}

type ServiceLogin interface {
	Login(email, password string) (model.User, error)
}
