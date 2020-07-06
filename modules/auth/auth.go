// Copyright 2019 GoAdmin Core Team. All rights reserved.
// Use of this source code is governed by a Apache-2.0 style
// license that can be found in the LICENSE file.

package auth

import (
	"sync"

	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/db"
	"github.com/GoAdminGroup/go-admin/modules/service"
	"github.com/GoAdminGroup/go-admin/plugins/admin/models"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules"
	"golang.org/x/crypto/bcrypt"
)

// Auth get the user model from Context.
// 硓筁把计ctx肚ヘ玡祅ノめ(Context.UserValue["user"])锣传ΘUserModel
func Auth(ctx *context.Context) models.UserModel {
	// User肚ヘ玡祅ノめ(Context.UserValue["user"])锣传ΘUserModel
	return ctx.User().(models.UserModel)
}

// Check check the password and username and return the user model.
// 浪琩user盞絏琌タ絋ぇ眔userrolepermissionのノmenu程穝戈(goadmin_users)盞絏(盞)
func Check(password string, username string, conn db.Connection) (user models.UserModel, ok bool) {
	// plugins\admin\models\user.go
	// User砞竚UserModel.Base.TableName(struct)肚砞竚UserModel(struct)
	// SetConn盢把计conn(db.Connection)砞竚UserModel.conn(UserModel.Base.Conn)
	user = models.User().SetConn(conn).FindByUserName(username)

	// 耞user琌
	if user.IsEmpty() {
		ok = false
	} else {
		// 浪琩盞絏
		if comparePassword(password, user.Password) {
			ok = true
			//眔userrolepermissionのノmenu
			user = user.WithRoles().WithPermissions().WithMenus()
			// EncodePassword盢把计pwd盞
			// UpdatePwd盢把计砞竚UserModel.UserModel穝dialect.H{"password": password,}
			user.UpdatePwd(EncodePassword([]byte(password)))
		} else {
			ok = false
		}
	}
	return
}

// 浪琩盞絏琌才
func comparePassword(comPwd, pwdHash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(pwdHash), []byte(comPwd))
	return err == nil
}

// EncodePassword encode the password.
// 盢把计pwd盞
func EncodePassword(pwd []byte) string {
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.DefaultCost)
	if err != nil {
		return ""
	}
	return string(hash[:])
}

// SetCookie set the cookie.
// 砞竚cookie(struct)纗response header Set-Cookieい
func SetCookie(ctx *context.Context, user models.UserModel, conn db.Connection) error {
	// 砞竚Session(struct)戈癟眔cookieの砞竚cookie
	ses, err := InitSession(ctx, conn)

	if err != nil {
		return err
	}
	// Add盢把计"user_id"user.IdSession.Values浪琩琌Τ才Session.Sid戈耞础┪琌穝戈
	// 程砞竚cookie(struct)纗response header Set-Cookieい
	return ses.Add("user_id", user.Id)
}

// DelCookie delete the cookie from Context.
// 睲埃cookie(session)戈
func DelCookie(ctx *context.Context, conn db.Connection) error {
	// 砞竚Session(struct)戈癟眔cookieの砞竚cookie
	ses, err := InitSession(ctx, conn)

	if err != nil {
		return err
	}
    // 睲埃cookie(session)
	return ses.Clear()
}

type TokenService struct {
	tokens CSRFToken //[]string
	lock   sync.Mutex
}

// 肚"token_csrf_helper"
func (s *TokenService) Name() string {
	// TokenServiceKey = token_csrf_helper
	return TokenServiceKey
}


func init() {
	// TokenServiceKey = token_csrf_helper
	// Register盢把计TokenServiceKeygen(ㄧΑ)盢services(map[string]Generator)い
	service.Register(TokenServiceKey, func() (service.Service, error) {
		// 砞竚肚TokenService.tokens(struct)
		return &TokenService{
			tokens: make(CSRFToken, 0),
		}, nil
	})
}

const (
	TokenServiceKey = "token_csrf_helper"
	ServiceKey      = "auth"
)

// 盢把计s锣传ΘTokenService(struct)摸肚
func GetTokenService(s interface{}) *TokenService {
	if srv, ok := s.(*TokenService); ok {
		return srv
	}
	panic("wrong service")
}

// AddToken add the token to the CSRFToken.
// ミuuid砞竚TokenService.tokens肚uuid(string)
func (s *TokenService) AddToken() string {
	s.lock.Lock()
	defer s.lock.Unlock()
	// 砞竚uuid
	tokenStr := modules.Uuid()
	s.tokens = append(s.tokens, tokenStr)
	return tokenStr
}

// CheckToken check the given token with tokens in the CSRFToken, if exist
// return true.
// 浪琩TokenService.tokens([]string)柑琌Τ才把计toCheckToken
// 狦才盢TokenService.tokens([]string)柑盢才toCheckToken眖[]string
func (s *TokenService) CheckToken(toCheckToken string) bool {
	for i := 0; i < len(s.tokens); i++ {
		if (s.tokens)[i] == toCheckToken {
			s.tokens = append((s.tokens)[:i], (s.tokens)[i+1:]...)
			return true
		}
	}
	return false
}

// CSRFToken is type of a csrf token list.
type CSRFToken []string

type Processor func(ctx *context.Context) (model models.UserModel, exist bool, msg string)

type Service struct {
	P Processor
}

// 肚auth(string)
func (s *Service) Name() string {
	return "auth"
}

// 盢把计s肚传ΘService(struct)肚
func GetService(s interface{}) *Service {
	if srv, ok := s.(*Service); ok {
		return srv
	}
	panic("wrong service")
}

// 盢把计processor砞竚Service.P(struct)肚
func NewService(processor Processor) *Service {
	return &Service{
		P: processor,
	}
}
