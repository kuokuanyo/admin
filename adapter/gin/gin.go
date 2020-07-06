// Copyright 2019 GoAdmin Core Team. All rights reserved.
// Use of this source code is governed by a Apache-2.0 style
// license that can be found in the LICENSE file.

package gin

import (
	"bytes"
	"errors"
	"net/http"
	"net/url"
	"strings"

	"github.com/GoAdminGroup/go-admin/adapter"
	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/engine"
	"github.com/GoAdminGroup/go-admin/modules/config"
	"github.com/GoAdminGroup/go-admin/plugins"
	"github.com/GoAdminGroup/go-admin/plugins/admin/models"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/constant"
	"github.com/GoAdminGroup/go-admin/template/types"
	"github.com/gin-gonic/gin"
)

// Gin structure value is a Gin GoAdmin adapter.
// Gin�P�ɤ]�ŦXadapter.WebFrameWork(interface)
type Gin struct {
	// adapter.BaseAdapter�badapter/adapter.go��
	// adapter.BaseAdapter(struct)�̭���db.Connection(interface)
	adapter.BaseAdapter
	// gin-gonic�M��
	// gin.Context(struct)��gin�̭��n�������A���\�bmiddleware�ǻ��ܼ�(�Ҧp���ҽШD�B�޲z�y�{)
	ctx *gin.Context
	// gin-gonic�M��
	// app���ج[������ҡA�]�tmuxer,middleware ,configuration�A�ǥ�New() or Default()�إ�Engine
	app *gin.Engine
}

// ��l��
func init() {
	// �bengine\engine.go
	// �إߤ����w�]���t�A��
	engine.Register(new(Gin))
}

//-------------------------------------
// �U�C��adapter.WebFrameWork(interface)����k
// Gin(struct)�]�Oadapter.WebFrameWork(interface)
//------------------------------------

// User implements the method Adapter.User.
// �qctx�����o�Τ�ҫ�(UserModel)�A��UserModel.Base.Conn = nil(�]ReleaseConn��k)
func (gins *Gin) User(ctx interface{}) (models.UserModel, bool) {
	// GetUser�qctx�����o�Τ�ҫ�(��adapter.BaseAdapter����k)
	// ���o�Τᨤ��B�v���H�Υi��menu
	return gins.GetUser(ctx, gins)
}

// Use implements the method Adapter.Use.
// plugins.Plugin�bplugins/plugins.go���A�Ointerface
// �W�[�B�z�{��(Handler)
func (gins *Gin) Use(app interface{}, plugs []plugins.Plugin) error {
	// GetUse�W�[����ܮج[(��adapter.BaseAdapter����k)
	// �W�[�B�z�{��(Handler)
	return gins.GetUse(app, plugs, gins)
}

// Content implements the method Adapter.Content.
// �Q��cookie���ҨϥΪ̡A���orole�Bpermission�Bmenu�A�����ˬd�v���A����ҪO�þɤJHTML
func (gins *Gin) Content(ctx interface{}, getPanelFn types.GetPanelFn, btns ...types.Button) {
	// GetContent�badapter\adapter.go
	// GetContent(Gin.adapter.BaseAdapter(struct)����k)
	// �Q��cookie���ҨϥΪ̡A���orole�Bpermission�Bmenu�A�����ˬd�v���A����ҪO�þɤJHTML
	gins.GetContent(ctx, getPanelFn, gins, btns)
}

type HandlerFunc func(ctx *gin.Context) (types.Panel, error)

// �K�[html��ج[��
// �Q��cookie���ҨϥΪ̡A�����ˬd�v���A�إ߼ҪO�ð���ҪO
func Content(handler HandlerFunc) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Content�bengine/engine.go��
		// Engine.Adapter���ରnil�A���۲K�[html��ج[��
		// �̫�@�˰���W����func (gins *Gin) Content(ctx inter...)�禡
		engine.Content(ctx, func(ctx interface{}) (types.Panel, error) {
			return handler(ctx.(*gin.Context))
		})
	}
}

func (gins *Gin) DisableLog()                {}
func (gins *Gin) Static(prefix, path string) {}

// SetApp implements the method Adapter.SetApp.
// �]�mGin.app(gin.Engine)�Agin.Engine(gin-gonic�M��)
func (gins *Gin) SetApp(app interface{}) error {
	var (
		eng *gin.Engine
		ok  bool
	)
	// app.(*gin.Engine)�Ninterface{}�ഫ��gin.Engine���A
	if eng, ok = app.(*gin.Engine); !ok {
		return errors.New("gin adapter SetApp: wrong parameter")
	}
	gins.app = eng
	return nil
}

// AddHandler implements the method Adapter.AddHandler.
// �K�[�B�z�{��
// context.Handlers(context.context.go��)�AHandlers������[]context.Handler�AHandler������function(*context.Context)
// AddHandler�ǥ�method�Bpath�W�[�B�z�{��(Handler)
// �]�mcontext.Context�W�[handlers�B�B�zurl�μg�Jheader�A�̫���o�s��request handle�Pmiddleware
func (gins *Gin) AddHandler(method, path string, handlers context.Handlers) {

	// gins.app������gin.Engine
	// Handle��k���ǥ�path��method���orequest handle�Pmiddleware�A���\�ର�j�qloading
	// Handle�ĤT�ӰѼ�(�D�n�B�z�{��)��funcion(*gin.Context)�Agin.Context��struct(gin-gonic�M��)
	gins.app.Handle(strings.ToUpper(method), path, func(c *gin.Context) {

		// �]�m�sContext(struct)�A�]�mRequest(�ШD)�H��UserValue�BResponse(�w�]��slice)
		// NewContext�bcontext\context.go
		ctx := context.NewContext(c.Request)

		// Context.Params������[]Context.Param�AParam�̦�key�H��value(�L�Ourl�Ѽƪ���P��)
		// �N�ѼƳ]�m�burl��
		for _, param := range c.Params {
			if c.Request.URL.RawQuery == "" {
				c.Request.URL.RawQuery += strings.Replace(param.Key, ":", "", -1) + "=" + param.Value
			} else {
				c.Request.URL.RawQuery += "&" + strings.Replace(param.Key, ":", "", -1) + "=" + param.Value
			}
		}

		// SetHandlers�bcontext\context.go
		// �]�mHandlers�A�Nhandlers�]�m��Context.handlers
		// Next�u�bmiddleware���ϥ�
		ctx.SetHandlers(handlers).Next()

		// ctx.Response.Header �i�঳�h����P��(map[string][]string)
		// Header()�bgin-gonic�M��
		// Header()�g�Jheader
		for key, head := range ctx.Response.Header {
			c.Header(key, head[0])
		}

		// �Ȥ�P�ǿ��Body���ରnil
		if ctx.Response.Body != nil {
			buf := new(bytes.Buffer)
			_, _ = buf.ReadFrom(ctx.Response.Body)
			c.String(ctx.Response.StatusCode, buf.String())
		} else {
			c.Status(ctx.Response.StatusCode)
		}
	})
}

// Name implements the method Adapter.Name.
// �^�Ǯج[�W��
func (gins *Gin) Name() string {
	return "gin"
}

// SetContext implements the method Adapter.SetContext.
// �]�mGin.ctx(struct)
func (gins *Gin) SetContext(contextInterface interface{}) adapter.WebFrameWork {
	var (
		// gin.Context(struct)�Ogin-gonic�M��
		ctx *gin.Context
		ok  bool
	)

	// �NcontextInterface���O�ܦ�gin.Context(struct)
	if ctx, ok = contextInterface.(*gin.Context); !ok {
		panic("gin adapter SetContext: wrong parameter")
	}

	return &Gin{ctx: ctx}
}

// Redirect implements the method Adapter.Redirect.
// ���s�ɦV�ܵn�J����(�X�{���~)
func (gins *Gin) Redirect() {
	// Redirect()��gin-gonic�M��̪���k
	// http.StatusFound = 302
	// config.GetLoginUrl()�n�J������url
	gins.ctx.Redirect(http.StatusFound, config.Url(config.GetLoginUrl()))
	gins.ctx.Abort()
}

// SetContentType implements the method Adapter.SetContentType.
func (gins *Gin) SetContentType() {
	return
}

// Write implements the method Adapter.Write.
func (gins *Gin) Write(body []byte) {
	// Data��k�bgin-gonic�M��
	// Data�N��Ƽg�Jbody�ç�shttp�N�X
	// gins.HTMLContentType() return "text/html; charset=utf-8"
	gins.ctx.Data(http.StatusOK, gins.HTMLContentType(), body)
}

// GetCookie implements the method Adapter.GetCookie.
// ���ocookie value�ǥ�cookie�R�W�M��
func (gins *Gin) GetCookie() (string, error) {
	// Cookie()�bgin-gonic�M���Context(struct)����k
	// Cookie()�^��cookie(�ǥѰѼƸ̪��R�W�^�Ǫ�)
	// gins.CookieKey()�O�Q��Gin.adapter.BaseAdapter�̪�CookieKey��k���ocookie���R�W
	// gins.CookieKey() = go_admin_session
	return gins.ctx.Cookie(gins.CookieKey())
}

// Path implements the method Adapter.Path.
// �^�Ǹ��|
func (gins *Gin) Path() string {
	return gins.ctx.Request.URL.Path
}

// Method implements the method Adapter.Method.
// �^�Ǥ�k
func (gins *Gin) Method() string {
	return gins.ctx.Request.Method
}

// FormParam implements the method Adapter.FormParam.
// �ѪR�Ѽ�(multipart/form-data�̪�)
func (gins *Gin) FormParam() url.Values {
	// http�M��
	// �ѪRmultipart/form-data�̪��Ѽ�
	_ = gins.ctx.Request.ParseMultipartForm(32 << 20)
	return gins.ctx.Request.PostForm
}

// IsPjax implements the method Adapter.IsPjax.
// �]�m���Y X-PJAX = true
func (gins *Gin) IsPjax() bool {
	// http�M��
	// constant.PjaxHeader = X-PJAX
	return gins.ctx.Request.Header.Get(constant.PjaxHeader) == "true"
}
