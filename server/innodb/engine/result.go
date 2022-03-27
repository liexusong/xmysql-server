package engine

import "github.com/zhukovaskychina/xmysql-server/server/innodb"

type Result struct {
	StatementID int64
	Rows        innodb.Rows
	Err         error
	ResultType  string
}

func NewResult() *Result {
	var rows = make([]innodb.Row, 0)
	return &Result{
		StatementID: 0,
		Rows:        rows,
		Err:         nil,
	}
}

func (result *Result) AddRows(row innodb.Row) {
	result.Rows.AddRow(row)
}
