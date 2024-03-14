package user

import (
	"errors"
	"fmt"
	"log"
	"url-shortener-go/config"
	"url-shortener-go/entity"

	"github.com/google/uuid"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Token struct {
	AccessToken string `json:"access_token"`
}

type UserData struct {
	Name     string
	Email    string
	Password string
}

func FindUserByEmail(email string) (*entity.User, error) {
	userEntity := entity.User{}

	err := config.DBConn.Where(&entity.User{Email: email}).First(&userEntity).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			log.Fatalf("Failed to connect DB: %v", err)
		}

		return nil, err
	}

	return &userEntity, nil
}

func InsertNewUser(name string, email string, password string) (*entity.User, error) {

	user, _ := FindUserByEmail(email)

	if user != nil {
		return nil, errors.New("email is used")
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), 5)

	userEntity := entity.User{
		Id:       uuid.NewString(),
		Name:     name,
		Email:    email,
		Password: string(hashedPassword),
	}

	err := config.DBConn.Create(&userEntity).Error
	if err != nil {
		log.Fatalf("Failed to connect DB: %v", err)

		return nil, err
	}

	return &userEntity, nil
}

func FindUserById(userId string) (*entity.User, error) {
	cacheKey := fmt.Sprintf("user:detail:%v", userId)
	userEntity := entity.User{}

	user, _ := config.CacheRememberV2(cacheKey, 3600, func() interface{} {
		userEntity := entity.User{}

		err := config.DBConn.Where(&entity.User{Id: userId}).First(&userEntity).Error
		if err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				log.Fatalf("Failed to connect DB: %v", err)
			}

			return nil
		}

		return &userEntity
	}, &userEntity)

	dataUser := user.(*entity.User)

	return dataUser, nil
}

func DoLogin(email string, password string) (*Token, error) {

	user, err := FindUserByEmail(email)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, err
	}

	token := config.GenerateJwtToken(user.Id, user.Email, user.Name)
	tokenStruct := Token{AccessToken: token}

	return &tokenStruct, nil
}
