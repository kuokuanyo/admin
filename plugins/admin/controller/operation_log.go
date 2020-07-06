package controller

import (
	"encoding/json"
	"github.com/GoAdminGroup/go-admin/plugins/admin/models"

	"github.com/GoAdminGroup/go-admin/context"
)

// RecordOperationLog record all operation logs, store into database.
// 癘魁┮Τ巨︽戈(goadmin_operation_log)い
func (h *Handler) RecordOperationLog(ctx *context.Context) {
	// 琩高Context.UserValue["user"]ぇ锣传ΘUserModel摸
	if user, ok := ctx.UserValue["user"].(models.UserModel); ok {
		var input []byte
		// 秆猂虫(form)把计
		form := ctx.Request.MultipartForm
		if form != nil {
			// 絪絏竚input
			input, _ = json.Marshal((*form).Value)
		}

		// OperationLogplugins\admin\models\operation_log.go
		// OperationLog肚箇砞OperationLogModel(struct)戈goadmin_operation_log
		// goadmin_operation_log戈魁ㄏノ幫巨︽
		// SetConn盢把计h.conn(Connection(interface))砞竚OperationLogModel.Base.Conn(struct)
		// New穝糤掸ㄏノ幫︽戈戈肚OperationLogModel(struct)
		// 戈input逆纗ㄏノ把计(ㄒ穝ㄏノ幫把计(form-data把计){"__go_admin_previous_":["/admin/info/manager?__page=1\u0026__pageSize=10\u0026__sort=id\u0026__sort_type=desc"],"__go_admin_t_":["972c6941-35fc-4401-9e95-e07a53c5370e"],"avatar":[""],"avatar__delete_flag":["0"],"name":["iiiii"],"password":["admin"],"password_again":["admin"],"username":["iiiii"]})
		models.OperationLog().SetConn(h.conn).New(user.Id, ctx.Path(), ctx.Method(), ctx.LocalIP(), string(input))
	}
}
