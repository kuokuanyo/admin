// Copyright 2019 GoAdmin Core Team. All rights reserved.
// Use of this source code is governed by a Apache-2.0 style
// license that can be found in the LICENSE file.

package plugins

import (
	"bytes"
	"errors"
	template2 "html/template"
	"net/http"
	"plugin"

	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/auth"
	"github.com/GoAdminGroup/go-admin/modules/config"
	"github.com/GoAdminGroup/go-admin/modules/db"
	"github.com/GoAdminGroup/go-admin/modules/logger"
	"github.com/GoAdminGroup/go-admin/modules/menu"
	"github.com/GoAdminGroup/go-admin/modules/service"
	"github.com/GoAdminGroup/go-admin/modules/ui"
	"github.com/GoAdminGroup/go-admin/plugins/admin/models"
	"github.com/GoAdminGroup/go-admin/template"
	"github.com/GoAdminGroup/go-admin/template/types"
)

// Plugin as one of the key components of goAdmin has three
// methods. GetRequest return all the path registered in the
// plugin. GetHandler according the url and method return the
// corresponding handler. InitPlugin init the plugin which do
// something like init the database and set the config and register
// the routes. The Plugin must implement the three methods.
// GetRequest�^�Ǵ��󤤪��Ҧ����|
// InitPlugin��l�ƴ���A�������l�Ƹ�Ʈw�ó]�m�ΰt�m���|
type Plugin interface {
	GetHandler() context.HandlerMap
	InitPlugin(services service.List)
	Name() string
	Prefix() string
}

// Base(struct)�]�OPlugin(interface)
type Base struct {
	// context.App�bcontext\context.go��
	App       *context.App
	Services  service.List
	Conn      db.Connection
	UI        *ui.Service
	PlugName  string
	URLPrefix string
}

// �^��Base.App.Handlers
func (b *Base) GetHandler() context.HandlerMap {
	return b.App.Handlers
}

// �^��Base.PlugName
func (b *Base) Name() string {
	return b.PlugName
}

// �^��Base.URLPrefix
func (b *Base) Prefix() string {
	return b.URLPrefix
}

// �N�Ѽ�srv(map[string]Service)�]�m��Base(struct)
func (b *Base) InitBase(srv service.List) {
	b.Services = srv
	// �N�Ѽ�b.Services�ഫ��Connect(interface)�^�Ǩæ^��
	b.Conn = db.GetConnection(b.Services)
	// �N�Ѽ�b.Services�ഫ��Service(struct)��^��
	b.UI = ui.GetService(b.Services)
}

// �N�ѼƳ]�m��ExecuteParam(struct)�A���۱N���w���ƾ�(types.Page(struct))�g�Jbuf(struct)�æ^��
func (b *Base) ExecuteTmpl(ctx *context.Context, panel types.Panel, animation ...bool) *bytes.Buffer {
	return Execute(ctx, b.Conn, b.UI.NavButtons, auth.Auth(ctx), panel, animation...)
}

// �N�ѼƳ]�m��ExecuteParam(struct)�A���۱N���w���ƾ�(types.Page(struct))�g�Jbuf(struct)�A�ÿ�XHTML��Context.response.Body
func (b *Base) HTML(ctx *context.Context, panel types.Panel, animation ...bool) {
	// ExecuteTmpl�N�ѼƳ]�m��ExecuteParam(struct)�A���۱N���w���ƾ�(types.NewPageParam...)�g�Jbuf(struct)�æ^��
	buf := b.ExecuteTmpl(ctx, panel, animation...)
	// �bcontext/context.go
	// ��XHTML�A�Nbody�ѼƳ]�m��Context.response.Body
	ctx.HTMLByte(http.StatusOK, buf.Bytes())
}

// ���ѪR�ɮ׫�N�ѼƳ]�m��ExecuteParam(struct)�A���۱N���w���ƾ�(types.Page(struct))�g�Jbuf(struct)�A�ÿ�XHTML��Context.response.Body
func (b *Base) HTMLFile(ctx *context.Context, path string, data map[string]interface{}, animation ...bool) {

	buf := new(bytes.Buffer)
	var panel types.Panel
	// html\template�M��
	// ParseFiles�ѪR�ɮרóЫؤ@��Template(struct)
	t, err := template2.ParseFiles(path)
	if err != nil {
		// IsProductionEnvironment�P�_globalCfg(Config).Env�O�_�O"prod"
		panel = template.WarningPanel(err.Error()).GetContent(config.IsProductionEnvironment())
	} else {
		// Execute�N�Ѽ�data�g�J�Ѽ�buf��
		if err := t.Execute(buf, data); err != nil {
			panel = template.WarningPanel(err.Error()).GetContent(config.IsProductionEnvironment())
		} else {
			panel = types.Panel{
				// HTML�ഫ���A��string
				Content: template.HTML(buf.String()),
			}
		}
	}
	// �N�ѼƳ]�m��ExecuteParam(struct)�A���۱N���w���ƾ�(types.Page(struct))�g�Jbuf(struct)�A�ÿ�XHTML��Context.response.Body
	b.HTML(ctx, panel, animation...)
}

// ���ѪR�h���ɮ׫�N�ѼƳ]�m��ExecuteParam(struct)�A���۱N���w���ƾ�(types.Page(struct))�g�Jbuf(struct)�A�ÿ�XHTML��Context.response.Body
func (b *Base) HTMLFiles(ctx *context.Context, data map[string]interface{}, files []string, animation ...bool) {
	buf := new(bytes.Buffer)
	var panel types.Panel

	t, err := template2.ParseFiles(files...)
	if err != nil {
		panel = template.WarningPanel(err.Error()).GetContent(config.IsProductionEnvironment())
	} else {
		// Execute�N�Ѽ�data�g�J�Ѽ�buf��
		if err := t.Execute(buf, data); err != nil {
			panel = template.WarningPanel(err.Error()).GetContent(config.IsProductionEnvironment())
		} else {
			panel = types.Panel{
				Content: template.HT
				// HTML�ഫ���A��stringML(buf.String()),
			}
		}
	}
	// �N�ѼƳ]�m��ExecuteParam(struct)�A���۱N���w���ƾ�(types.Page(struct))�g�Jbuf(struct)�A�ÿ�XHTML��Context.response.Body
	b.HTML(ctx, panel, animation...)
}

// �N�Ѽ�mod���}���oPlugin(struct)��M��"Plugin"���Ÿ��A�̫��ഫ��Plugin(interface)���O�^��
func LoadFromPlugin(mod string) Plugin {

	// plugin�M��
	// Open���}go plugin�A�p�G�Ѽ�mod�w�g�s�b�h�^��Plugin(struct)
	plug, err := plugin.Open(mod)
	if err != nil {
		logger.Error("LoadFromPlugin err", err)
		panic(err)
	}

	// �M��plug(struct)���Ѽ�"Plugin"���Ÿ��A�^��symPlugin(interface)
	symPlugin, err := plug.Lookup("Plugin")
	if err != nil {
		logger.Error("LoadFromPlugin err", err)
		panic(err)
	}

	var p Plugin
	//�NsymPlugin�ഫ��Plugin(interface)
	p, ok := symPlugin.(Plugin)
	if !ok {
		logger.Error("LoadFromPlugin err: unexpected type from module symbol")
		panic(errors.New("LoadFromPlugin err: unexpected type from module symbol"))
	}

	return p
}

// GetHandler is a help method for Plugin GetHandler.
// �^�ǰѼ�app(App.Handlers)
func GetHandler(app *context.App) context.HandlerMap { return app.Handlers }

// �N�ѼƳ]�m��ExecuteParam(struct)�A���۱N���w���ƾ�(types.Page(struct))�g�Jbuf(struct)�æ^��
func Execute(ctx *context.Context, conn db.Connection, navButtons types.Buttons, user models.UserModel,
	panel types.Panel, animation ...bool) *bytes.Buffer {
	// GetTheme�^��globalCfg.Theme
	// IsPjax�P�_�O�_header X-PJAX:true
	// Get�P�_templateMap(map[string]Template)��key��O�_�Ѽ�config.GetTheme()�A���h�^��Template(interface)
	// GetTemplate��Template(interface)����k
	tmpl, tmplName := template.Get(config.GetTheme()).GetTemplate(ctx.IsPjax())

	// template\template.go��
	return template.Execute(template.ExecuteParam{
		User:       user,
		TmplName:   tmplName,
		Tmpl:       tmpl,
		Panel:      panel,
		// �ƻsglobalCfg(Config struct)��NConfig.Databases[key].Driver�]�m��Config.Databases[key]��^��
		Config:     *config.Get(),
		// GetGlobalMenu�^�ǰѼ�user(struct)��Menu(�]�mmenuList�BmenuOption�BMaxOrder)
		// SetActiveClass�]�wmenu��active
		// URLRemovePrefixglobalCfg(Config struct).prefix�NURL���e��h��
		Menu:       menu.GetGlobalMenu(user, conn).SetActiveClass(config.URLRemovePrefix(ctx.Path())),
		Animation:  len(animation) > 0 && animation[0] || len(animation) == 0,
		Buttons:    navButtons.CheckPermission(user),
		NoCompress: len(animation) > 1 && animation[1],
	})
}
