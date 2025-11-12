package repositories_users

import (
	"errors"
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	users "users/dao_users"
	domain "users/domain_users"
	errores "users/extras"
)

type MySQLConfig struct {
	Host     string
	Port     string
	Database string
	Username string
	Password string
}

type MySQL struct {
	db *gorm.DB
}

func NewMySQL(config MySQLConfig) MySQL {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", config.Username, config.Password, config.Host, config.Port, config.Database)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to MySQL: %s", err.Error())
	}

	// AutoMigrate para mantener el esquema sincronizado
	if err := db.AutoMigrate(&users.User{}); err != nil {
		log.Fatalf("error running Automigrate: %s", err.Error())
	}

	return MySQL{db: db}
}

func (repository MySQL) GetUserById(id int64) (users.User, errores.ApiError) {
	var user users.User
	if err := repository.db.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return user, errores.NewNotFoundApiError("user not found")
		}
		return user, errores.NewInternalServerApiError("error fetching user by id", err)
	}
	return user, nil
}

func (repository MySQL) GetUserByEmail(email string) (users.User, errores.ApiError) {
	var user users.User
	if err := repository.db.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return user, errores.NewNotFoundApiError("user not found")
		}
		return user, errores.NewInternalServerApiError("error fetching user by email", err)
	}
	return user, nil
}

func (repository MySQL) CreateUser(user users.User) (int64, errores.ApiError) {
	if err := repository.db.Create(&user).Error; err != nil {
		return 0, errores.NewInternalServerApiError("error creating user", err)
	}
	return user.User_id, nil
}

func (repository MySQL) Login(login domain.Login) (users.User, errores.ApiError) {
	var user users.User
	if err := repository.db.Where("email = ?", login.Email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return user, errores.NewNotFoundApiError("invalid credentials")
		}
		return user, errores.NewInternalServerApiError("error fetching user for login", err)
	}

	// Comparar contraseña (en texto plano por ahora, o podés agregar hash)
	if user.Password != login.Password {
		return user, errores.NewUnauthorizedApiError("invalid password")
	}

	return user, nil
}
