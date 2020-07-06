package guard

import (
	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/db"
	"github.com/GoAdminGroup/go-admin/modules/errors"
	"github.com/GoAdminGroup/go-admin/modules/service"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/constant"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/response"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/table"
)

type Guard struct {
	services  service.List
	conn      db.Connection
	tableList table.GeneratorList
}

// �N�Ѽ�s�Bc�Bt�]�m��Guard(struct)��^��
func New(s service.List, c db.Connection, t table.GeneratorList) *Guard {
	return &Guard{
		services:  s,
		conn:      c,
		tableList: t,
	}
}

func (g *Guard) table(ctx *context.Context) (table.Table, string) {
	prefix := ctx.Query(constant.PrefixKey)
	return g.tableList[prefix](ctx), prefix
}

// �d��url�̪��Ѽ�(__prefix)�A�p�GGuard.tableList�s�b��prefix(key)�h����j��
func (g *Guard) CheckPrefix(ctx *context.Context) {
	// PrefixKey = __prefix
	// ���oRequest url�̪��Ѽ�(__prefix)
	prefix := ctx.Query(constant.PrefixKey)

	if _, ok := g.tableList[prefix]; !ok {
		if ctx.Headers(constant.PjaxHeader) == "" && ctx.Method() != "GET" {
			response.BadRequest(ctx, errors.Msg)
		} else {
			response.Alert(ctx, errors.Msg, errors.Msg, "table model not found", g.conn)
		}
		ctx.Abort()
		return
	}
	// ����j��Context.handlers[ctx.index](ctx)
	ctx.Next()
}

const (
	editFormParamKey   = "edit_form_param"
	deleteParamKey     = "delete_param"
	exportParamKey     = "export_param"
	deleteMenuParamKey = "delete_menu_param"
	editMenuParamKey   = "edit_menu_param"
	newMenuParamKey    = "new_menu_param"
	newFormParamKey    = "new_form_param"
	updateParamKey     = "update_param"
	showFormParamKey   = "show_form_param"
	showNewFormParam   = "show_new_form_param"
)
