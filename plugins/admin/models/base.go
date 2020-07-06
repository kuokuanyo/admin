package models

import (
	"database/sql"
	"github.com/GoAdminGroup/go-admin/modules/db"
)

// Base is base model structure.
type Base struct {
	TableName string

	Conn db.Connection
	Tx   *sql.Tx
}

// �N�Ѽ�con(Connection(interface))�]�m��Base.Conn
func (b Base) SetConn(con db.Connection) Base {
	b.Conn = con
	return b
}

// �ǥѵ��w��table�^��sql(struct)
func (b Base) Table(table string) *db.SQL {
	// Table�bmodules/db/statement.go��
	// Table�ǥѵ��w��table�^��sql(struct)
	// WithDriver�ǥѵ��w��conn�^��sql(struct)
	return db.Table(table).WithDriver(b.Conn)
}
