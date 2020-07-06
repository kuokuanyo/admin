package guard

import (
	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/auth"
	"github.com/GoAdminGroup/go-admin/modules/errors"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/form"
	"html/template"
	"strconv"
)

type MenuNewParam struct {
	Title    string
	Header   string
	ParentId int64
	Icon     string
	Uri      string
	Roles    []string
	Alert    template.HTML
}

// �P�_MenuNewParam.Alert�O�_�X�{ĵ�i(���O�ŭ�)
func (e MenuNewParam) HasAlert() bool {
	return e.Alert != template.HTML("")
}

// �ǥѰѼƨ��omultipart/form-data���]�m���ȡA��������token�ñNmultipart/form-data��key�Bvalue�ȳ]�m��Context.UserValue[new_menu_param]�A�̫����j��Context.handlers[ctx.index](ctx)
func (g *Guard) MenuNew(ctx *context.Context) {

	// �ǥѰѼƨ��omultipart/form-data����parent_id��
	parentId := ctx.FormValue("parent_id")
	if parentId == "" {
		parentId = "0"
	}

	var (
		alert template.HTML
		// TokenKey = __go_admin_t_
		// �ǥѰѼƨ��omultipart/form-data����__go_admin_t_��
		token = ctx.FormValue(form.TokenKey)
	)

	// TokenServiceKey = token_csrf_helper
	// Get�z�L�Ѽ�(token_csrf_helper)���o�ǰt��Service(interface)
	// GetTokenService�N�Ѽ�s�ഫ��TokenService(struct)���O��^��
	// CheckToken�ˬdTokenService.tokens([]string)�̬O�_���ŦX�Ѽ�token����
    // �p�G�ŦX�A�N�bTokenService.tokens([]string)�̱N�ŦX��token�q[]string���X
	if !auth.GetTokenService(g.services.Get(auth.TokenServiceKey)).CheckToken(token) {
		alert = getAlert(errors.EditFailWrongToken)
	}

	// title�Picon�Ȥ@�w�n�]�m(multipart/form-data)
	// checkEmpty�ˬd�Ѽ�(�h��key)���bmultipart/form-data�̳]�m(�p�G�Ȭ��ūh�X�{���~)
	if alert == "" {
		alert = checkEmpty(ctx, "title", "icon")
	}

	parentIdInt, _ := strconv.Atoi(parentId)

	// newMenuParamKey = new_menu_param
	// SetUserValue�ǥѰѼ�new_menu_param�B&MenuNewParam{...}(struct)�]�wContext.UserValue
	// �Nmultipart/form-data��key�Bvalue�ȳ]�m��Context.UserValue[new_menu_param]
	ctx.SetUserValue(newMenuParamKey, &MenuNewParam{
		Title:    ctx.FormValue("title"),
		Header:   ctx.FormValue("header"),
		ParentId: int64(parentIdInt),
		Icon:     ctx.FormValue("icon"),
		Uri:      ctx.FormValue("uri"),
		Roles:    ctx.Request.Form["roles[]"],
		Alert:    alert,
	})

	// ����j��Context.handlers[ctx.index](ctx)
	ctx.Next()
}

// �NContext.UserValue(map[string]interface{})[new_menu_param]�����ഫ��MenuNewParam(struct)���O
func GetMenuNewParam(ctx *context.Context) *MenuNewParam {
	// newMenuParamKey = new_menu_param
	// �NContext.UserValue(map[string]interface{})[new_menu_param]�����ഫ��MenuNewParam(struct)���O
	return ctx.UserValue[newMenuParamKey].(*MenuNewParam)
}
