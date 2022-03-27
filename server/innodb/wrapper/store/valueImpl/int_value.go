package valueImpl

import "github.com/zhukovaskychina/xmysql-server/server/innodb"

type IntValue struct {
	innodb.Value
	value []byte
}

func NewIntValue(value []byte) innodb.Value {

	var bigIntValue = new(IntValue)
	bigIntValue.value = value
	return bigIntValue
}

func (i IntValue) Raw() interface{} {
	panic("implement me")
}

func (i IntValue) ToByte() []byte {
	panic("implement me")
}

func (i IntValue) DataType() innodb.ValType {
	panic("implement me")
}

func (i IntValue) Compare(x innodb.Value) (innodb.CompareType, error) {
	panic("implement me")
}

func (i IntValue) UnaryPlus() (innodb.Value, error) {
	panic("implement me")
}

func (i IntValue) UnaryMinus() (innodb.Value, error) {
	panic("implement me")
}

func (i IntValue) Add(value innodb.Value) (innodb.Value, error) {
	panic("implement me")
}

func (i IntValue) Sub(value innodb.Value) (innodb.Value, error) {
	panic("implement me")
}

func (i IntValue) Mul(value innodb.Value) (innodb.Value, error) {
	panic("implement me")
}

func (i IntValue) Div(value innodb.Value) (innodb.Value, error) {
	panic("implement me")
}

func (i IntValue) Pow(value innodb.Value) (innodb.Value, error) {
	panic("implement me")
}

func (i IntValue) Mod(value innodb.Value) (innodb.Value, error) {
	panic("implement me")
}

func (i IntValue) Equal(value innodb.Value) (innodb.Value, error) {
	panic("implement me")
}

func (i IntValue) NotEqual(value innodb.Value) (innodb.Value, error) {
	panic("implement me")
}

func (i IntValue) GreaterThan(value innodb.Value) (innodb.Value, error) {
	panic("implement me")
}

func (i IntValue) LessThan(value innodb.Value) (innodb.Value, error) {
	panic("implement me")
}

func (i IntValue) GreaterOrEqual(value innodb.Value) (innodb.Value, error) {
	panic("implement me")
}

func (i IntValue) LessOrEqual(value innodb.Value) (innodb.Value, error) {
	panic("implement me")
}

func (i IntValue) And(value innodb.Value) (innodb.Value, error) {
	panic("implement me")
}

func (i IntValue) Or(value innodb.Value) (innodb.Value, error) {
	panic("implement me")
}
