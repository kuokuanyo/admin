package guard

import (
	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/errors"
)

type MenuDeleteParam struct {
	Id string
}

// �d��url���Ѽ�id���ȫ�Nid�]�m��MenuDeleteParam(struct)�A���۱N�ȳ]�m��Context.UserValue[delete_menu_param]���A�̫����j��Context.handlers[ctx.index](ctx)
func (g *Guard) MenuDelete(ctx *context.Context) {
	// �d��url���Ѽ�id����
	id := ctx.Query("id")

	if id == "" {
		// alertWithTitleAndDesc�bplugins\admin\modules\guard\edit.go
		// WrongID = wrong id
		// Alert�z�L�Ѽ�ctx�^�ǥثe�n�J���Τ�(Context.UserValue["user"])���ഫ��UserModel�A���۱N���w���ƾ�(types.Page(struct))�g�Jbuf(struct)�æ^�ǡA�̫��XHTML
		// �N�Ѽ�Menu�Bmenu�Berrors.WrongID�g�JPanel
		alertWithTitleAndDesc(ctx, "Menu", "menu", errors.WrongID, g.conn)
		ctx.Abort()
		return
	}

	// TODO: check the user permission

	// deleteMenuParamKey = delete_menu_param
	// �Nid�]�m��MenuDeleteParam(struct)
	// SetUserValue�ǥѰѼ�delete_menu_param�B&MenuDeleteParam{...}(struct)�]�wContext.UserValue
	// �N�ѼƳ]�m��Context.UserValue[delete_menu_param]
	ctx.SetUserValue(deleteMenuParamKey, &MenuDeleteParam{
		Id: id,
	})

	// ����j��Context.handlers[ctx.index](ctx)
	ctx.Next()
}

// �NContext.UserValue(map[string]interface{})[delete_menu_param]�����ഫ��MenuDeleteParam(struct)���O
func GetMenuDeleteParam(ctx *context.Context) *MenuDeleteParam {
	// deleteMenuParamKey = delete_menu_param
	// �NContext.UserValue(map[string]interface{})[delete_menu_param]�����ഫ��MenuDeleteParam(struct)���O
	return ctx.UserValue[deleteMenuParamKey].(*MenuDeleteParam)
}
