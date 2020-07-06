package response

import (
	"net/http"

	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/auth"
	"github.com/GoAdminGroup/go-admin/modules/config"
	"github.com/GoAdminGroup/go-admin/modules/db"
	"github.com/GoAdminGroup/go-admin/modules/errors"
	"github.com/GoAdminGroup/go-admin/modules/language"
	"github.com/GoAdminGroup/go-admin/modules/menu"
	"github.com/GoAdminGroup/go-admin/template"
	"github.com/GoAdminGroup/go-admin/template/types"
)

// ���\�A�^��code:200 and msg:ok
func Ok(ctx *context.Context) {
	ctx.JSON(http.StatusOK, map[string]interface{}{
		"code": http.StatusOK,
		"msg":  "ok",
	})
}

// ���\�A�^��code:200 and msg
func OkWithMsg(ctx *context.Context, msg string) {
	ctx.JSON(http.StatusOK, map[string]interface{}{
		"code": http.StatusOK,
		"msg":  msg,
	})
}

// ���\�A�^��code:200 and msg:ok and data
func OkWithData(ctx *context.Context, data map[string]interface{}) {
	ctx.JSON(http.StatusOK, map[string]interface{}{
		"code": http.StatusOK,
		"msg":  "ok",
		"data": data,
	})
}

// ���~�ШD�A�^��code:400 and msg
func BadRequest(ctx *context.Context, msg string) {
	ctx.JSON(http.StatusBadRequest, map[string]interface{}{
		"code": http.StatusBadRequest,
		// Get�̷ӳ]�w���y�������T��
		"msg": language.Get(msg),
	})
}

// �z�L�Ѽ�ctx�^�ǥثe�n�J���Τ�(Context.UserValue["user"])���ഫ��UserModel�A���۱N���w���ƾ�(types.Page(struct))�g�Jbuf(struct)�æ^�ǡA�̫��XHTML
// �N�Ѽ�desc�Btitle�Bmsg�g�JPanel
func Alert(ctx *context.Context, desc, title, msg string, conn db.Connection) {
	// �z�L�Ѽ�ctx�^�ǥثe�n�J���Τ�(Context.UserValue["user"])���ഫ��UserModel
	user := auth.Auth(ctx)

	// Get�P�_templateMap(map[string]Template)��key��O�_�Ѽ�config.GetTheme()�A���h�^��Template(interface)
	// GetTemplate��Template(interface)����k
	tmpl, tmplName := template.Get(config.GetTheme()).GetTemplate(ctx.IsPjax())

	// �N���w���ƾ�(types.Page(struct))�g�Jbuf(struct)�æ^��
	buf := template.Execute(template.ExecuteParam{
		User:     user,
		TmplName: tmplName,
		Tmpl:     tmpl,
		Panel: types.Panel{
			Content:     template.Get(config.GetTheme()).Alert().Warning(msg),
			Description: template.HTML(desc),
			Title:       template.HTML(title),
		},
		Config:    *config.Get(),
		// GetGlobalMenu�^�ǰѼ�user(struct)��Menu(�]�mmenuList�BmenuOption�BMaxOrder)
		// �]�wmenu��active
		// URLRemovePrefix globalCfg(Config struct).prefix�NURL���e��h��
		Menu:      menu.GetGlobalMenu(user, conn).SetActiveClass(config.URLRemovePrefix(ctx.Path())),
		Animation: true,
	})

	// �Nbuf��X��HTML
	ctx.HTML(http.StatusOK, buf.String())
}

// ���~�A�^��code:500 and msg
func Error(ctx *context.Context, msg string) {
	ctx.JSON(http.StatusInternalServerError, map[string]interface{}{
		"code": http.StatusInternalServerError,
		"msg":  language.Get(msg),
	})
}

// ���~�A�^��code:403 and msg
func Denied(ctx *context.Context, msg string) {
	ctx.JSON(http.StatusInternalServerError, map[string]interface{}{
		"code": http.StatusForbidden,
		"msg":  language.Get(msg),
	})
}

// OffLineHandler(function)
// �P�_���I�O�_�n�����A�p�n�����A�P�_method�O�_��get�H��header�̥]�taccept:html���XHTML
var OffLineHandler = func(ctx *context.Context) {
	// GetSiteOff�^��globalCfg.SiteOff(���I����)
	if config.GetSiteOff() {
		// �P�_method�O�_��get�H��header�̥]�taccept:html
		if ctx.WantHTML() {
			//��XHTML
			ctx.HTML(http.StatusOK, `<html><body><h1>The website is offline</h1></body></html>`)
		} else {
			ctx.JSON(http.StatusForbidden, map[string]interface{}{
				"code": http.StatusForbidden,
				"msg":  language.Get(errors.SiteOff),
			})
		}
		// Context.index = 63
		ctx.Abort()
	}
}
