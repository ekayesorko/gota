package gota

import (
	"context"
	"errors"
	"net/http"
	"time"

	"cloud.google.com/go/auth/credentials/idtoken"
	"github.com/ekayesorko/gota/resterr"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type UserModel struct {
	Model

	Name       string
	Email      string
	Role       string
	Picture    string
	Provider   string
	ProviderID string
}

type UserRepository struct {
	CommonRepository[UserModel]
}

type UserService struct {
	CommonService[UserModel]
}

type UserController struct{}

var UserCtrl UserController
var UserRepo UserRepository
var UserSvc UserService

var GoogleClientId string

func (ctrl *UserController) Login(c echo.Context) error {
	var req LoginReq
	bindErr := c.Bind(&req)
	if bindErr != nil {
		//todo
		panic(bindErr)
	}
	ctx, err := GetContext(c)
	if err != nil {
		return err.Respond(c)
	}

	resp, err := UserSvc.Login(ctx, req)
	if err != nil {
		return err.Respond(c)
	}
	return c.JSON(http.StatusOK, resp)
}

func (s *UserService) Login(c Context, req LoginReq) (*LoginResp, *resterr.RestError) {
	var user UserModel
	var err error
	user, err = googleGetUserInfo(c.C, req.Token)
	if err != nil {
		return nil, resterr.NewNotFoundError("user")
	}

	existingUser, err := UserRepo.GetByParam(c, UserModel{Email: user.Email})
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, resterr.InternalServerError
	}
	if existingUser != nil {
		return s.createLoginResp(existingUser)
	}

	createUserErr := s.Create(c, &user)
	if createUserErr != nil {
		return nil, createUserErr
	}

	existingUser = &user
	return s.createLoginResp(existingUser)
}

func (s *UserService) createLoginResp(user *UserModel) (*LoginResp, *resterr.RestError) {
	accessToken, err := s.createJWT(user, "todoaccesskey", time.Hour*24)
	if err != nil {
		return nil, err
	}
	refreshToken, err := s.createJWT(user, "todorefreshkey", time.Hour*24*7)
	if err != nil {
		return nil, err
	}
	return &LoginResp{
		AccessToken:      accessToken,
		RefreshToken:     refreshToken,
		AccessExpiresIn:  24 * 60 * 60,
		RefreshExpiresIn: 7 * 24 * 60 * 60,
		User:             *user,
	}, nil
}

func (s *UserService) createJWT(user *UserModel, key string, expiresIn time.Duration) (string, *resterr.RestError) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  user.Id,
		"role": user.Role,
		"exp":  time.Now().Add(expiresIn).Unix(),
	})
	tokenStr, err := token.SignedString([]byte(key))
	if err != nil {
		return "", resterr.InternalServerError
	}
	return tokenStr, nil
}

func googleGetUserInfo(ctx context.Context, token string) (UserModel, error) {
	payload, err := idtoken.Validate(ctx, token, gInstance.config.Auth.GoogleClientId)
	if err != nil {
		return UserModel{}, err
	}

	providerEmail, found := payload.Claims["email"].(string)
	if !found {
		return UserModel{}, errors.New("email not found in token")
	}

	name, _ := payload.Claims["name"].(string)
	picture, _ := payload.Claims["picture"].(string)
	providerID, _ := payload.Claims["sub"].(string)

	user := UserModel{
		Email:      providerEmail,
		Name:       name,
		Picture:    picture,
		Provider:   "google",
		ProviderID: providerID,
	}

	return user, nil
}

func SetupAuth(router *echo.Group) {
	router.POST("/login", UserCtrl.Login)
}
