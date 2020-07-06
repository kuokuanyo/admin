package guard

import (
	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/auth"
	"github.com/GoAdminGroup/go-admin/modules/errors"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/form"
	"html/template"
	"strconv"
)

type MenuEditParam struct {
	Id       string
	Title    string
	Header   string
	ParentId int64
	Icon     string
	Uri      string
	Roles    []string
	Alert    template.HTML
}

// �ˬdMenuEditParam.Alert�O�_�X�{ĵ�i(���O�ŭ�)
func (e MenuEditParam) HasAlert() bool {
	return e.Alert != template.HTML("")
}

// �ǥѰѼƨ��omultipart/form-data���]�m���ȡA��������token�ñNmultipart/form-data��key�Bvalue�ȳ]�m��Context.UserValue[edit_menu_param]�A�̫����j��Context.handlers[ctx.index](ctx)
func (g *Guard) MenuEdit(ctx *context.Context) {

	// �ǥѰѼƨ��omultipart/form-data����parent_id��
	parentId := ctx.FormValue("parent_id")
	if parentId == "" {
		parentId = "0"
	}

	var (
		parentIdInt, _ = strconv.Atoi(parentId)
		// TokenKey = __go_admin_t_
		// �ǥѰѼƨ��omultipart/form-data����__go_admin_t_��
		token          = ctx.FormValue(form.TokenKey)
		alert          template.HTML
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
		alert = checkEmpty(ctx, "id", "title", "icon")
	}

	// TODO: check the user permission
	// editMenuParamKey = edit_menu_param
	// SetUserValue�ǥѰѼ�edit_menu_param�B&MenuEditParam{...}(struct)�]�wContext.UserValue
	// �Nmultipart/form-data��key�Bvalue�ȳ]�m��Context.UserValue[edit_menu_param]
	ctx.SetUserValue(editMenuParamKey, &MenuEditParam{
		Id:       ctx.FormValue("id"),
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

// �NContext.UserValue(map[string]interface{})[edit_menu_param]�����ഫ��MenuEditParam(struct)���O
func GetMenuEditParam(ctx *context.Context) *MenuEditParam {
	// editMenuParamKey = edit_menu_param
	// �NContext.UserValue(map[string]interface{})[edit_menu_param]�����ഫ��MenuEditParam(struct)���O
	return ctx.UserValue[editMenuParamKey].(*MenuEditParam)
}

// �ˬd�Ѽ�(�h��key)���bmultipart/form-data�̳]�m(�p�G�Ȭ��ūh�X�{���~)
func checkEmpty(ctx *context.Context, key ...string) template.HTML {
	for _, k := range key {
		if ctx.FormValue(k) == "" {
			return getAlert("wrong " + k)
		}
	}
	return template.HTML("")
}
