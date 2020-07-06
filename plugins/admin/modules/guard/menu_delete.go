package guard

import (
	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/errors"
)

type MenuDeleteParam struct {
	Id string
}

// 琩高urlい把计id盢id砞竚MenuDeleteParam(struct)钡帝盢砞竚Context.UserValue[delete_menu_param]い程磅︽癹伴Context.handlers[ctx.index](ctx)
func (g *Guard) MenuDelete(ctx *context.Context) {
	// 琩高urlい把计id
	id := ctx.Query("id")

	if id == "" {
		// alertWithTitleAndDescplugins\admin\modules\guard\edit.go
		// WrongID = wrong id
		// Alert硓筁把计ctx肚ヘ玡祅ノめ(Context.UserValue["user"])锣传ΘUserModel钡帝盢倒﹚计沮(types.Page(struct))糶buf(struct)肚程块HTML
		// 盢把计Menumenuerrors.WrongID糶Panel
		alertWithTitleAndDesc(ctx, "Menu", "menu", errors.WrongID, g.conn)
		ctx.Abort()
		return
	}

	// TODO: check the user permission

	// deleteMenuParamKey = delete_menu_param
	// 盢id砞竚MenuDeleteParam(struct)
	// SetUserValue虑パ把计delete_menu_param&MenuDeleteParam{...}(struct)砞﹚Context.UserValue
	// 盢把计砞竚Context.UserValue[delete_menu_param]
	ctx.SetUserValue(deleteMenuParamKey, &MenuDeleteParam{
		Id: id,
	})

	// 磅︽癹伴Context.handlers[ctx.index](ctx)
	ctx.Next()
}

// 盢Context.UserValue(map[string]interface{})[delete_menu_param]锣传ΘMenuDeleteParam(struct)摸
func GetMenuDeleteParam(ctx *context.Context) *MenuDeleteParam {
	// deleteMenuParamKey = delete_menu_param
	// 盢Context.UserValue(map[string]interface{})[delete_menu_param]锣传ΘMenuDeleteParam(struct)摸
	return ctx.UserValue[deleteMenuParamKey].(*MenuDeleteParam)
}
