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
// �z�L�Ѽ�ctx�^�ǥثe�n�J���Τ�(Context.UserValue["user"])���ഫ��UserModel
func Auth(ctx *context.Context) models.UserModel {
	// User�^�ǥثe�n�J���Τ�(Context.UserValue["user"])���ഫ��UserModel
	return ctx.User().(models.UserModel)
}

// Check check the password and username and return the user model.
// �ˬduser�K�X�O�_���T������ouser��role�Bpermission�Υi��menu�A�̫��s��ƪ�(goadmin_users)���K�X��(�[�K)
func Check(password string, username string, conn db.Connection) (user models.UserModel, ok bool) {
	// plugins\admin\models\user.go
	// User�]�mUserModel.Base.TableName(struct)�æ^�ǳ]�mUserModel(struct)
	// SetConn�N�Ѽ�conn(db.Connection)�]�m��UserModel.conn(UserModel.Base.Conn)
	user = models.User().SetConn(conn).FindByUserName(username)

	// �P�_user�O�_����
	if user.IsEmpty() {
		ok = false
	} else {
		// �ˬd�K�X
		if comparePassword(password, user.Password) {
			ok = true
			//���ouser��role�Bpermission�Υi��menu
			user = user.WithRoles().WithPermissions().WithMenus()
			// EncodePassword�N�Ѽ�pwd�[�K
			// UpdatePwd�N�ѼƳ]�m��UserModel.UserModel�åB��sdialect.H{"password": password,}
			user.UpdatePwd(EncodePassword([]byte(password)))
		} else {
			ok = false
		}
	}
	return
}

// �ˬd�K�X�O�_�۲�
func comparePassword(comPwd, pwdHash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(pwdHash), []byte(comPwd))
	return err == nil
}

// EncodePassword encode the password.
// �N�Ѽ�pwd�[�K
func EncodePassword(pwd []byte) string {
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.DefaultCost)
	if err != nil {
		return ""
	}
	return string(hash[:])
}

// SetCookie set the cookie.
// �]�mcookie(struct)���x�s�bresponse header Set-Cookie��
func SetCookie(ctx *context.Context, user models.UserModel, conn db.Connection) error {
	// �]�mSession(struct)��T�è��ocookie�γ]�mcookie��
	ses, err := InitSession(ctx, conn)

	if err != nil {
		return err
	}
	// Add�N�Ѽ�"user_id"�Buser.Id�[�JSession.Values���ˬd�O�_���ŦXSession.Sid����ơA�P�_���J�άO��s���
	// �̫�]�mcookie(struct)���x�s�bresponse header Set-Cookie��
	return ses.Add("user_id", user.Id)
}

// DelCookie delete the cookie from Context.
// �M��cookie(session)���
func DelCookie(ctx *context.Context, conn db.Connection) error {
	// �]�mSession(struct)��T�è��ocookie�γ]�mcookie��
	ses, err := InitSession(ctx, conn)

	if err != nil {
		return err
	}
    // �M��cookie(session)
	return ses.Clear()
}

type TokenService struct {
	tokens CSRFToken //[]string
	lock   sync.Mutex
}

// �^��"token_csrf_helper"
func (s *TokenService) Name() string {
	// TokenServiceKey = token_csrf_helper
	return TokenServiceKey
}


func init() {
	// TokenServiceKey = token_csrf_helper
	// Register�N�Ѽ�TokenServiceKey�Bgen(�禡)�N�Jservices(map[string]Generator)��
	service.Register(TokenServiceKey, func() (service.Service, error) {
		// �]�m�æ^��TokenService.tokens(struct)
		return &TokenService{
			tokens: make(CSRFToken, 0),
		}, nil
	})
}

const (
	TokenServiceKey = "token_csrf_helper"
	ServiceKey      = "auth"
)

// �N�Ѽ�s�ഫ��TokenService(struct)���O��^��
func GetTokenService(s interface{}) *TokenService {
	if srv, ok := s.(*TokenService); ok {
		return srv
	}
	panic("wrong service")
}

// AddToken add the token to the CSRFToken.
// �إ�uuid�ó]�m��TokenService.tokens�A�^��uuid(string)
func (s *TokenService) AddToken() string {
	s.lock.Lock()
	defer s.lock.Unlock()
	// �]�muuid
	tokenStr := modules.Uuid()
	s.tokens = append(s.tokens, tokenStr)
	return tokenStr
}

// CheckToken check the given token with tokens in the CSRFToken, if exist
// return true.
// �ˬdTokenService.tokens([]string)�̬O�_���ŦX�Ѽ�toCheckToken����
// �p�G�ŦX�A�N�bTokenService.tokens([]string)�̱N�ŦX��toCheckToken�q[]string���X
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

// �^��auth(string)
func (s *Service) Name() string {
	return "auth"
}

// �N�Ѽ�s�Ǵ���Service(struct)��^��
func GetService(s interface{}) *Service {
	if srv, ok := s.(*Service); ok {
		return srv
	}
	panic("wrong service")
}

// �N�Ѽ�processor�]�m��Service.P(struct)�æ^��
func NewService(processor Processor) *Service {
	return &Service{
		P: processor,
	}
}
