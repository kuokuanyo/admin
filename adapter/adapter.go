// Copyright 2019 GoAdmin Core Team. All rights reserved.
// Use of this source code is governed by a Apache-2.0 style
// license that can be found in the LICENSE file.

package adapter

import (
	"bytes"
	"fmt"
	"net/url"

	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/auth"
	"github.com/GoAdminGroup/go-admin/modules/config"
	"github.com/GoAdminGroup/go-admin/modules/db"
	"github.com/GoAdminGroup/go-admin/modules/errors"
	"github.com/GoAdminGroup/go-admin/modules/logger"
	"github.com/GoAdminGroup/go-admin/modules/menu"
	"github.com/GoAdminGroup/go-admin/plugins"
	"github.com/GoAdminGroup/go-admin/plugins/admin/models"
	"github.com/GoAdminGroup/go-admin/template"
	"github.com/GoAdminGroup/go-admin/template/types"
)

// WebFrameWork is an interface which is used as an adapter of
// framework and goAdmin. It must implement two methods. Use registers
// the routes and the corresponding handlers. Content writes the
// response to the corresponding context of framework.
// WebFrameWork�\�ೣ�]�w�b�ج[��(�ϥ�/adapter/gin/gin.go�ج[)
type WebFrameWork interface {
	// Name return the web framework name.
	// �^�ǨϥΪ�web�ج[�W��
	Name() string

	// Use method inject the plugins to the web framework engine which is the
	// first parameter.
	// �N���󴡤J�ج[������
	Use(app interface{}, plugins []plugins.Plugin) error

	// Content add the panel html response of the given callback function to
	// the web framework context which is the first parameter.
	// �K�[html��ج[��
	Content(ctx interface{}, fn types.GetPanelFn, navButtons ...types.Button)

	// User get the auth user model from the given web framework context.
	// �q���w���W�U�夤���o�Τ�ҫ�
	User(ctx interface{}) (models.UserModel, bool)

	// AddHandler inject the route and handlers of GoAdmin to the web framework.
	// �N����(���|)�γB�z�{���[�J�ج[
	AddHandler(method, path string, handlers context.Handlers)

	DisableLog()

	Static(prefix, path string)

	// Helper functions
	// ================================

	SetApp(app interface{}) error
	SetConnection(db.Connection)
	GetConnection() db.Connection
	SetContext(ctx interface{}) WebFrameWork
	GetCookie() (string, error)
	Path() string
	Method() string
	FormParam() url.Values
	IsPjax() bool
	Redirect()
	SetContentType()
	Write(body []byte)
	CookieKey() string
	HTMLContentType() string
}

// BaseAdapter is a base adapter contains some helper functions.
// �򥻰t�����]�t���U�\��
// db.Connection(interface)
type BaseAdapter struct {
	db db.Connection
}

// SetConnection set the db connection.
// �]�w��Ʈw�s��
func (base *BaseAdapter) SetConnection(conn db.Connection) {
	base.db = conn
}

// GetConnection get the db connection.
// ���o�s�u
func (base *BaseAdapter) GetConnection() db.Connection {
	return base.db
}

// HTMLContentType return the default content type header.
// �^�ǹw�]��content type
func (base *BaseAdapter) HTMLContentType() string {
	return "text/html; charset=utf-8"
}

// CookieKey return the cookie key.
func (base *BaseAdapter) CookieKey() string {
	//go_admin_session
	return auth.DefaultCookieKey
}

// GetUser is a helper function get the auth user model from the context.
// �q�W�U��(�Ѽ�ctx)�����o�Τ�ҫ�(UserModel)�A��UserModel.Base.Conn = nil(�]ReleaseConn��k)
// ���o�Τᨤ��B�v���H�Υi��menu
func (base *BaseAdapter) GetUser(ctx interface{}, wf WebFrameWork) (models.UserModel, bool) {
	// ���ocookie
	cookie, err := wf.SetContext(ctx).GetCookie()
	// models.UserModel�bplugins/admin/modules/user.go��
	if err != nil {
		return models.UserModel{}, false
	}
	// auth.GetCurUser�bmodules/auth/middleware.go���A�^�ǨϥΪ̸�T
	// WebFrameWork.GetConnection()�^��BaseAdapter.db
	// �ǥ�cookie�Bconn�i�H�o�쨤��B�v���H�Υi�ϥε��
	user, exist := auth.GetCurUser(cookie, wf.GetConnection())
	// ReleaseConn�bplugins/admin/modules/user.go��
	// �NUserModel.Conn(UserModel.Base.Conn) = nil
	return user.ReleaseConn(), exist
}

// GetUse is a helper function adds the plugins to the framework.
// �W�[����ܮج[
// plugins.Plugin�bplugins/plugins.go���A�Ointerface
// WebFrameWork interface
// �ǥ�method�Burl�W�[�B�z�{��(Handler)
func (base *BaseAdapter) GetUse(app interface{}, plugin []plugins.Plugin, wf WebFrameWork) error {
	// adapter\gin\gin.go��
	// �]�mGin.app(gin.Engine)
	if err := wf.SetApp(app); err != nil {
		return err
	}

	// �bplugins/plugins.go
	// plug interface
	for _, plug := range plugin {
		// ��^���ѩM�����k
		// GetHandler()�b context\context.go�A����map[Path]Handlers
		// path struct�A�]�turl�Bmethod
		// handlers����[]Handler�A Handler���� func(ctx *Context)
		for path, handlers := range plug.GetHandler() {
			// �����k(WebFrameWork interface)�A adapter\gin\gin.go���]�m�Ӥ�k
			// �ǥ�method�Burl�W�[�B�z�{��(Handler)
			// �]�mcontext.Context�P�]�murl�P�g�Jheader�A���o�s��request�Pmiddleware
			wf.AddHandler(path.Method, path.URL, handlers)
		}
	}

	return nil
}

// GetContent is a helper function of adapter.Content
// �Q��cookie���ҨϥΪ̡A���orole�Bpermission�Bmenu�A�����ˬd�v���A����ҪO�þɤJHTML
func (base *BaseAdapter) GetContent(ctx interface{}, getPanelFn types.GetPanelFn, wf WebFrameWork, navButtons types.Buttons) {

	// SetContext�]�mGin.ctx(struct)
	newBase := wf.SetContext(ctx)

	// ����cookie value
	cookie, hasError := newBase.GetCookie()

	// �p�X�{���~���s�ɦV�ܵn�J����
	if hasError != nil || cookie == "" {
		newBase.Redirect()
		return
	}

	// GetCurUser�^�ǨϥΪ̼ҫ��A���orole�Bpermission�Bmenu
	// wf.GetConnection()�^��BaseAdapter.db(interface)
	user, authSuccess := auth.GetCurUser(cookie, wf.GetConnection())

	// �p�X�{���~���s�ɦV�ܵn�J����
	if !authSuccess {
		newBase.Redirect()
		return
	}

	var (
		panel types.Panel
		err   error
	)

	// CheckPermissions�ˬd�Τ��v��(�bmodules\auth\middleware.go)
	if !auth.CheckPermissions(user, newBase.Path(), newBase.Method(), newBase.FormParam()) {
		// �S���v��
		// errors.NoPermission = no permission
		panel = template.WarningPanel(errors.NoPermission)
	} else {
		panel, err = getPanelFn(ctx)
		if err != nil {
			panel = template.WarningPanel(err.Error())
		}
	}

	// Default()���o�w�]��template(�D�D�W�٤w�g�q�L�����t�m)
	// tmpl���O��template.Template(interface)�A�btemplate/template.go��
	// template.Template��ui�ե󪺤�k�A�N�bplugins���۩w�qui
	// IsPjax()�bgin/gin.go���A�]�m���Y X-PJAX = true
	// GetTemplate(bool)��template.Template(interface)����k
	tmpl, tmplName := template.Default().GetTemplate(newBase.IsPjax())

	buf := new(bytes.Buffer)

	// ExecuteTemplate����ҪO(html\template\template.go��Template����k)
	// �ǥѵ���tmplName���μҪO����w����H(�ĤT�ӰѼ�)
	hasError = tmpl.ExecuteTemplate(buf, tmplName, types.NewPage(types.NewPageParam{
		User: user,
		// GetGlobalMenu ��^user��menu(modules\menu\menu.go��)
		// Menu(struct�]�t)List�BOptions�BMaxOrder
		Menu: menu.GetGlobalMenu(user, wf.GetConnection()).SetActiveClass(config.URLRemovePrefix(newBase.Path())),
		// IsProductionEnvironment�ˬd�Ͳ�����
		// GetContent�btemplate\types\page.go
		// Panel(struct)�D�n���e�ϥ�pjax���ҪO
		// GetContent������e(�]�m�e��HTML)�A�]�mPanel�æ^��
		Panel: panel.GetContent(config.IsProductionEnvironment()),
		// Assets���O��template.HTML(string)
		// �B�zasset��æ^��HTML�y�k
		Assets: template.GetComponentAssetImportHTML(),
		// �ˬd�v���A�^��Buttons([]Button(interface))
		// �btemplate\types\button.go
		Buttons: navButtons.CheckPermission(user),
	}))

	if hasError != nil {
		logger.Error(fmt.Sprintf("error: %s adapter content, ", newBase.Name()), hasError)
	}

	// �]�mContentType
	newBase.SetContentType()
	// �g�J
	newBase.Write(buf.Bytes())
}
