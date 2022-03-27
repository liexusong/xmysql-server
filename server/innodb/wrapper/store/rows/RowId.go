package rows

import (
	"github.com/zhukovaskychina/xmysql-server/server/common"
	"github.com/zhukovaskychina/xmysql-server/server/innodb"
	"github.com/zhukovaskychina/xmysql-server/server/innodb/wrapper/store/valueImpl"
	"github.com/zhukovaskychina/xmysql-server/util"
)

type DBRowId struct {
	innodb.Value
	value []byte //6 byte
}

func NewDBRowId(rowId uint64) innodb.Value {
	var buff = make([]byte, 0)
	util.WriteUB6(buff, rowId)
	var dbRowId = new(DBRowId)
	dbRowId.value = buff
	return dbRowId
}

func NewDBRowIdWithBytes(rowIdByte []byte) innodb.Value {

	var dbRowId = new(DBRowId)
	dbRowId.value = rowIdByte
	return dbRowId
}

func (D *DBRowId) Raw() interface{} {
	return D.value
}

func (D DBRowId) ToByte() []byte {

	return D.value
}

func (D DBRowId) DataType() innodb.ValType {
	return common.COLUMN_TYPE_LONG
}

func (D DBRowId) Compare(x innodb.Value) (innodb.CompareType, error) {
	panic("implement me")
}

func (D DBRowId) UnaryPlus() (innodb.Value, error) {
	panic("implement me")
}

func (D DBRowId) UnaryMinus() (innodb.Value, error) {
	panic("implement me")
}

func (D DBRowId) Add(value innodb.Value) (innodb.Value, error) {
	panic("implement me")
}

func (D DBRowId) Sub(value innodb.Value) (innodb.Value, error) {
	panic("implement me")
}

func (D DBRowId) Mul(value innodb.Value) (innodb.Value, error) {
	panic("implement me")
}

func (D DBRowId) Div(value innodb.Value) (innodb.Value, error) {
	panic("implement me")
}

func (D DBRowId) Pow(value innodb.Value) (innodb.Value, error) {
	panic("implement me")
}

func (D DBRowId) Mod(value innodb.Value) (innodb.Value, error) {
	panic("implement me")
}

func (D DBRowId) Equal(value innodb.Value) (innodb.Value, error) {
	panic("implement me")
}

func (D DBRowId) NotEqual(value innodb.Value) (innodb.Value, error) {
	panic("implement me")
}

func (D DBRowId) GreaterThan(value innodb.Value) (innodb.Value, error) {
	panic("implement me")
}

func (D *DBRowId) LessThan(value innodb.Value) (innodb.Value, error) {

	thanPK := value.Raw().([]byte)
	thisPK := D.value

	thanVal := util.ReadUB8Byte2Long(thanPK)
	thisVal := util.ReadUB8Byte2Long(thisPK)
	if thanVal > thisVal {
		return valueImpl.NewBoolValue(true), nil
	}
	return valueImpl.NewBoolValue(false), nil
}

func (D DBRowId) GreaterOrEqual(value innodb.Value) (innodb.Value, error) {
	panic("implement me")
}

func (D DBRowId) LessOrEqual(value innodb.Value) (innodb.Value, error) {
	panic("implement me")
}

func (D DBRowId) And(value innodb.Value) (innodb.Value, error) {
	panic("implement me")
}

func (D DBRowId) Or(value innodb.Value) (innodb.Value, error) {
	panic("implement me")
}
