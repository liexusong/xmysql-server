package innodb

import "github.com/zhukovaskychina/xmysql-server/server/common"

// ValType specifies the type for SQLVal.
type ValType int

// These are the possible Valtype values.
// HexNum represents a 0x... valueImpl. It cannot
// be treated as a simple valueImpl because it can
// be interpreted differently depending on the
// context.
const (
	UNKVAL   = 0
	IntVal   = common.COLUMN_TYPE_INT24
	StrVal   = common.COLUMN_TYPE_STRING
	FloatVal = common.COLUMN_TYPE_FLOAT
	HexNum   = common.COLUMN_TYPE_VARCHAR
	HexVal   = common.COLUMN_TYPE_SHORT
	ValArg   = common.COLUMN_TYPE_BLOB
	BitVal   = common.COLUMN_TYPE_BIT
	RowIdVal = common.COLUMN_TYPE_LONG
)

type CompareType string

// ComparisonExpr.Operator
const (
	EqualStr             CompareType = "="
	LessThanStr          CompareType = "<"
	GreaterThanStr       CompareType = ">"
	LessEqualStr         CompareType = "<="
	GreaterEqualStr      CompareType = ">="
	NotEqualStr          CompareType = "!="
	NullSafeEqualStr     CompareType = "<=>"
	InStr                CompareType = "in"
	NotInStr             CompareType = "not in"
	LikeStr              CompareType = "like"
	NotLikeStr           CompareType = "not like"
	RegexpStr            CompareType = "regexp"
	NotRegexpStr         CompareType = "not regexp"
	JSONExtractOp        CompareType = "->"
	JSONUnquoteExtractOp CompareType = "->>"
)

// UnaryExpr.Operator
const (
	UPlusStr   = "+"
	UMinusStr  = "-"
	TildaStr   = "~"
	BangStr    = "!"
	BinaryStr  = "binary "
	UBinaryStr = "_binary "
)

// BinaryExpr.Operator
const (
	BitAndStr     = "&"
	BitOrStr      = "|"
	BitXorStr     = "^"
	PlusStr       = "+"
	MinusStr      = "-"
	MultStr       = "*"
	DivStr        = "/"
	IntDivStr     = "div"
	ModStr        = "%"
	ShiftLeftStr  = "<<"
	ShiftRightStr = ">>"
)

// this string is "character set" and this comment is required
const (
	CharacterSetStr = " character set"
)

func (s CompareType) String() string {
	return string(s)
}

type Value interface {
	Raw() interface{}
	ToByte() []byte
	DataType() ValType
	Compare(x Value) (CompareType, error)
	UnaryPlus() (Value, error)
	UnaryMinus() (Value, error)
	Add(Value) (Value, error)
	Sub(Value) (Value, error)
	Mul(Value) (Value, error)
	Div(Value) (Value, error)
	Pow(Value) (Value, error)
	Mod(Value) (Value, error)
	Equal(Value) (Value, error)
	NotEqual(Value) (Value, error)
	GreaterThan(Value) (Value, error)
	LessThan(Value) (Value, error)
	GreaterOrEqual(Value) (Value, error)
	LessOrEqual(Value) (Value, error)
	And(Value) (Value, error)
	Or(Value) (Value, error)
}
