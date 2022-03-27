package valueImpl

import (
	"bytes"
	"github.com/zhukovaskychina/xmysql-server/server/innodb"
	"github.com/zhukovaskychina/xmysql-server/util"
)

type BigIntValue struct {
	innodb.Value
	value []byte
}

func NewBigIntValue(value []byte) innodb.Value {

	var bigIntValue = new(BigIntValue)
	bigIntValue.value = value
	return bigIntValue
}

func (b BigIntValue) Raw() interface{} {
	return util.ReadUB8Bytes2Long(b.value)
}

func (b BigIntValue) ToByte() []byte {
	return b.value
}

func (b BigIntValue) DataType() innodb.ValType {
	return innodb.RowIdVal
}

func (b BigIntValue) Compare(x innodb.Value) (innodb.CompareType, error) {
	panic("implement me")
}

func (b BigIntValue) UnaryPlus() (innodb.Value, error) {
	panic("implement me")
}

func (b BigIntValue) UnaryMinus() (innodb.Value, error) {
	panic("implement me")
}

func (b BigIntValue) Add(value innodb.Value) (innodb.Value, error) {
	panic("implement me")
}

func (b BigIntValue) Sub(value innodb.Value) (innodb.Value, error) {
	panic("implement me")
}

func (b BigIntValue) Mul(value innodb.Value) (innodb.Value, error) {
	panic("implement me")
}

func (b BigIntValue) Div(value innodb.Value) (innodb.Value, error) {
	panic("implement me")
}

func (b BigIntValue) Pow(value innodb.Value) (innodb.Value, error) {
	panic("implement me")
}

func (b BigIntValue) Mod(value innodb.Value) (innodb.Value, error) {
	panic("implement me")
}

func (b BigIntValue) Equal(value innodb.Value) (innodb.Value, error) {

	return NewBoolValue(bytes.Compare(b.value, value.ToByte()) == 0), nil
}

func (b BigIntValue) NotEqual(value innodb.Value) (innodb.Value, error) {
	panic("implement me")
}

func (b BigIntValue) GreaterThan(value innodb.Value) (innodb.Value, error) {
	panic("implement me")
}

func (b BigIntValue) LessThan(value innodb.Value) (innodb.Value, error) {
	second := value.Raw().(int64)
	first := b.Raw().(int64)
	return NewBoolValue(first < second), nil
}

func (b BigIntValue) GreaterOrEqual(value innodb.Value) (innodb.Value, error) {
	panic("implement me")
}

func (b BigIntValue) LessOrEqual(value innodb.Value) (innodb.Value, error) {
	second := value.Raw().(int64)
	first := b.Raw().(int64)
	return NewBoolValue(first <= second), nil
}

func (b BigIntValue) And(value innodb.Value) (innodb.Value, error) {
	panic("implement me")
}

func (b BigIntValue) Or(value innodb.Value) (innodb.Value, error) {
	panic("implement me")
}
