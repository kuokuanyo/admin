package controller

import (
	"encoding/json"
	"github.com/GoAdminGroup/go-admin/plugins/admin/models"

	"github.com/GoAdminGroup/go-admin/context"
)

// RecordOperationLog record all operation logs, store into database.
// �O���Ҧ��ާ@�欰�ܸ�ƪ�(goadmin_operation_log)��
func (h *Handler) RecordOperationLog(ctx *context.Context) {
	// �d��Context.UserValue["user��"]�����ഫ��UserModel���O
	if user, ok := ctx.UserValue["user"].(models.UserModel); ok {
		var input []byte
		// �ѪR�����(form)�Ѽ�
		form := ctx.Request.MultipartForm
		if form != nil {
			// �s�X��minput
			input, _ = json.Marshal((*form).Value)
		}

		// OperationLog�bplugins\admin\models\operation_log.go
		// OperationLog�^�ǹw�]��OperationLogModel(struct)�A��ƪ�W��goadmin_operation_log
		// goadmin_operation_log��ƪ������ϥΎ;ާ@�欰
		// SetConn�N�Ѽ�h.conn(Connection(interface))�]�m��OperationLogModel.Base.Conn(struct)
		// New�s�W�@���ϥΎͦ欰��Ʀܸ�ƪ�A�^��OperationLogModel(struct)
		// ��ƪ�input��쬰�x�s�ϥΪ��Ѽ�(�Ҧp�s�بϥΎͪ��Ѽ�(form-data�Ѽ�)�A{"__go_admin_previous_":["/admin/info/manager?__page=1\u0026__pageSize=10\u0026__sort=id\u0026__sort_type=desc"],"__go_admin_t_":["972c6941-35fc-4401-9e95-e07a53c5370e"],"avatar":[""],"avatar__delete_flag":["0"],"name":["iiiii"],"password":["admin"],"password_again":["admin"],"username":["iiiii"]})
		models.OperationLog().SetConn(h.conn).New(user.Id, ctx.Path(), ctx.Method(), ctx.LocalIP(), string(input))
	}
}
