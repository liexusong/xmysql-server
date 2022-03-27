package valueImpl

import (
	"github.com/piex/transcode"
	"github.com/zhukovaskychina/xmysql-server/server/innodb"
)

type BoolValue struct {
	value bool
}

func NewBoolValue(value bool) innodb.Value {
	var boolValue = new(BoolValue)
	boolValue.value = value
	return boolValue
}

func (b *BoolValue) Raw() interface{} {
	return b.value
}

func (b BoolValue) ToByte() []byte {
	panic("implement me")
}

func (b BoolValue) DataType() innodb.ValType {
	panic("implement me")
}

func (b BoolValue) Compare(x innodb.Value) (innodb.CompareType, error) {
	panic("implement me")
}

func (b BoolValue) UnaryPlus() (innodb.Value, error) {
	panic("implement me")
}

func (b BoolValue) UnaryMinus() (innodb.Value, error) {
	panic("implement me")
}

func (b BoolValue) Add(value innodb.Value) (innodb.Value, error) {
	panic("implement me")
}

func (b BoolValue) Sub(value innodb.Value) (innodb.Value, error) {
	panic("implement me")
}

func (b BoolValue) Mul(value innodb.Value) (innodb.Value, error) {
	panic("implement me")
}

func (b BoolValue) Div(value innodb.Value) (innodb.Value, error) {
	panic("implement me")
}

func (b BoolValue) Pow(value innodb.Value) (innodb.Value, error) {
	panic("implement me")
}

func (b BoolValue) Mod(value innodb.Value) (innodb.Value, error) {
	panic("implement me")
}

func (b BoolValue) Equal(value innodb.Value) (innodb.Value, error) {
	panic("implement me")
}

func (b BoolValue) NotEqual(value innodb.Value) (innodb.Value, error) {
	panic("implement me")
}

func (b BoolValue) GreaterThan(value innodb.Value) (innodb.Value, error) {
	panic("implement me")
}

func (b BoolValue) LessThan(value innodb.Value) (innodb.Value, error) {
	panic("implement me")
}

func (b BoolValue) GreaterOrEqual(value innodb.Value) (innodb.Value, error) {
	panic("implement me")
}

func (b BoolValue) LessOrEqual(value innodb.Value) (innodb.Value, error) {
	panic("implement me")
}

func (b BoolValue) And(value innodb.Value) (innodb.Value, error) {
	panic("implement me")
}

func (b BoolValue) Or(value innodb.Value) (innodb.Value, error) {
	panic("implement me")
}

type VarcharVal struct {
	innodb.Value
	value []byte
}

func NewVarcharVal(content []byte) innodb.Value {
	var varcharVal = new(VarcharVal)

	varcharVal.value = content
	return varcharVal
}

func (v *VarcharVal) Raw() interface{} {
	return transcode.FromByteArray(v.value).Decode("GBK").ToString()
}

func (v VarcharVal) ToByte() []byte {
	return v.value
}

func (v VarcharVal) DataType() innodb.ValType {
	return innodb.StrVal
}

func (v VarcharVal) Compare(x innodb.Value) (innodb.CompareType, error) {
	panic("implement me")
}

func (v VarcharVal) UnaryPlus() (innodb.Value, error) {
	panic("implement me")
}

func (v VarcharVal) UnaryMinus() (innodb.Value, error) {
	panic("implement me")
}

func (v VarcharVal) Add(value innodb.Value) (innodb.Value, error) {

	v.value = append(v.value, value.ToByte()...)
	return NewVarcharVal(v.value), nil
}

func (v VarcharVal) Sub(value innodb.Value) (innodb.Value, error) {
	panic("implement me")
}

func (v VarcharVal) Mul(value innodb.Value) (innodb.Value, error) {
	panic("implement me")
}

func (v VarcharVal) Div(value innodb.Value) (innodb.Value, error) {
	panic("implement me")
}

func (v VarcharVal) Pow(value innodb.Value) (innodb.Value, error) {
	panic("implement me")
}

func (v VarcharVal) Mod(value innodb.Value) (innodb.Value, error) {
	panic("implement me")
}

func (v VarcharVal) Equal(value innodb.Value) (innodb.Value, error) {
	panic("implement me")
}

func (v VarcharVal) NotEqual(value innodb.Value) (innodb.Value, error) {
	panic("implement me")
}

func (v VarcharVal) GreaterThan(value innodb.Value) (innodb.Value, error) {
	panic("implement me")
}

func (v VarcharVal) LessThan(value innodb.Value) (innodb.Value, error) {
	panic("implement me")
}

func (v VarcharVal) GreaterOrEqual(value innodb.Value) (innodb.Value, error) {
	panic("implement me")
}

func (v VarcharVal) LessOrEqual(value innodb.Value) (innodb.Value, error) {
	panic("implement me")
}

func (v VarcharVal) And(value innodb.Value) (innodb.Value, error) {
	panic("implement me")
}

func (v VarcharVal) Or(value innodb.Value) (innodb.Value, error) {
	panic("implement me")
}
