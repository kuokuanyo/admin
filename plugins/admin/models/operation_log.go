package models

import (
	"github.com/GoAdminGroup/go-admin/modules/db"
	"github.com/GoAdminGroup/go-admin/modules/db/dialect"
)

// OperationLogModel is operation log model structure.
// OperationLogModel�����ϥΪ̾ާ@�欰
type OperationLogModel struct {
	//Base(struct)�bplugins\admin\models\base.go
	Base

	Id        int64
	UserId    int64
	Path      string
	Method    string
	Ip        string
	Input     string
	CreatedAt string
	UpdatedAt string
}

// OperationLog return a default operation log model.
// �^�ǹw�]��OperationLogModel(struct)�A��ƪ�W��goadmin_operation_log
// goadmin_operation_log��ƪ������ϥΪ̾ާ@�欰
func OperationLog() OperationLogModel {
	return OperationLogModel{Base: Base{TableName: "goadmin_operation_log"}}
}

// Find return a default operation log model of given id.
// �z�L�Ѽ�(id)�M��ŦX��ơA�N��T�]�m��OperationLogModel(struct)
func (t OperationLogModel) Find(id interface{}) OperationLogModel {
	// �bplugins\admin\models\base.go
	// Table�ǥѵ��w���Ѽ�(t.TableName)�^��sql(struct)
	// Find�bmodules\db\statement.go��
	// Find�ǥѰѼ�id���o�ŦX���
	item, _ := t.Table(t.TableName).Find(id)
	// �z�L�Ѽ�(m map[string]interface{})�N��T�]�m��OperationLogModel(struct)
	return t.MapToModel(item)
}

// �N�Ѽ�conn(Connection(interface))�]�m��OperationLogModel.Base.Conn(struct)
func (t OperationLogModel) SetConn(con db.Connection) OperationLogModel {
	t.Conn = con
	return t
}

// New create a new operation log model.
// �s�W�@���ϥΪ̦欰��Ʀܸ�ƪ�A�^��OperationLogModel(struct)
func (t OperationLogModel) New(userId int64, path, method, ip, input string) OperationLogModel {
	// ���J�ϥΪ̦欰���
	// OperationLogModel.Base.TableName
	// Table�ǥѵ��w���Ѽ�(t.TableName)�^��sql(struct)
	id, _ := t.Table(t.TableName).Insert(dialect.H{
		"user_id": userId,
		"path":    path,
		"method":  method,
		"ip":      ip,
		"input":   input,
	})

	t.Id = id
	t.UserId = userId
	t.Path = path
	t.Method = method
	t.Ip = ip
	t.Input = input

	return t
}

// MapToModel get the operation log model from given map.
// �z�L�Ѽ�(m map[string]interface{})�N��T�]�m��OperationLogModel(struct)
func (t OperationLogModel) MapToModel(m map[string]interface{}) OperationLogModel {
	t.Id = m["id"].(int64)
	t.UserId = m["user_id"].(int64)
	t.Path, _ = m["path"].(string)
	t.Method, _ = m["method"].(string)
	t.Ip, _ = m["ip"].(string)
	t.Input, _ = m["input"].(string)
	t.CreatedAt, _ = m["created_at"].(string)
	t.UpdatedAt, _ = m["updated_at"].(string)
	return t
}
