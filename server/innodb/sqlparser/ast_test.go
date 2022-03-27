/*
Copyright 2017 Google Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package sqlparser

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/zhukovaskychina/xmysql-server/server/innodb"
	"github.com/zhukovaskychina/xmysql-server/server/innodb/planning/plan"
	"github.com/zhukovaskychina/xmysql-server/server/innodb/sqlparser/dependency/sqltypes"
	"reflect"
	"strings"
	"testing"
	"unsafe"
)

func TestAppend(t *testing.T) {
	query := "select * from t left join m on t.x=m.y where a = 1  and b <2 and c=3"
	tree, err := Parse(query)
	if err != nil {
		t.Error(err)
	}

	var b bytes.Buffer
	Append(&b, tree)
	got := b.String()
	want := query
	if got != want {
		t.Errorf("Append: %s, want %s", got, want)
	}
	Append(&b, tree)
	got = b.String()
	want = query + query
	if got != want {
		t.Errorf("Append: %s, want %s", got, want)
	}
}

func TestSelect(t *testing.T) {
	tree, err := Parse("select * from t where a = 1")
	if err != nil {
		t.Error(err)
	}
	expr := tree.(*Select).Where.Expr

	sel := &Select{}
	sel.AddWhere(expr)
	buf := NewTrackedBuffer(nil)
	sel.Where.Format(buf)
	want := " where a = 1"
	if buf.String() != want {
		t.Errorf("where: %q, want %s", buf.String(), want)
	}
	sel.AddWhere(expr)
	buf = NewTrackedBuffer(nil)
	sel.Where.Format(buf)
	want = " where a = 1 and a = 1"
	if buf.String() != want {
		t.Errorf("where: %q, want %s", buf.String(), want)
	}
	sel = &Select{}
	sel.AddHaving(expr)
	buf = NewTrackedBuffer(nil)
	sel.Having.Format(buf)
	want = " having a = 1"
	if buf.String() != want {
		t.Errorf("having: %q, want %s", buf.String(), want)
	}
	sel.AddHaving(expr)
	buf = NewTrackedBuffer(nil)
	sel.Having.Format(buf)
	want = " having a = 1 and a = 1"
	if buf.String() != want {
		t.Errorf("having: %q, want %s", buf.String(), want)
	}

	// OR clauses must be parenthesized.
	tree, err = Parse("select * from t where a = 1 or b = 1")
	if err != nil {
		t.Error(err)
	}
	expr = tree.(*Select).Where.Expr
	sel = &Select{}
	sel.AddWhere(expr)
	buf = NewTrackedBuffer(nil)
	sel.Where.Format(buf)
	want = " where (a = 1 or b = 1)"
	if buf.String() != want {
		t.Errorf("where: %q, want %s", buf.String(), want)
	}
	sel = &Select{}
	sel.AddHaving(expr)
	buf = NewTrackedBuffer(nil)
	sel.Having.Format(buf)
	want = " having (a = 1 or b = 1)"
	if buf.String() != want {
		t.Errorf("having: %q, want %s", buf.String(), want)
	}
}

func TestRemoveHints(t *testing.T) {
	for _, query := range []string{
		"select * from t use index (i)",
		"select * from t force index (i)",
	} {
		tree, err := Parse(query)
		if err != nil {
			t.Fatal(err)
		}
		sel := tree.(*Select)
		sel.From = TableExprs{
			sel.From[0].(*AliasedTableExpr).RemoveHints(),
		}
		buf := NewTrackedBuffer(nil)
		sel.Format(buf)
		if got, want := buf.String(), "select * from t"; got != want {
			t.Errorf("stripped plan: %s, want %s", got, want)
		}
	}
}

func TestAddOrder(t *testing.T) {
	src, err := Parse("select foo, bar from baz order by foo")
	if err != nil {
		t.Error(err)
	}
	order := src.(*Select).OrderBy[0]
	dst, err := Parse("select * from t")
	if err != nil {
		t.Error(err)
	}
	dst.(*Select).AddOrder(order)
	buf := NewTrackedBuffer(nil)
	dst.Format(buf)
	want := "select * from t order by foo asc"
	if buf.String() != want {
		t.Errorf("order: %q, want %s", buf.String(), want)
	}
	dst, err = Parse("select * from t union select * from s")
	if err != nil {
		t.Error(err)
	}
	dst.(*Union).AddOrder(order)
	buf = NewTrackedBuffer(nil)
	dst.Format(buf)
	want = "select * from t union select * from s order by foo asc"
	if buf.String() != want {
		t.Errorf("order: %q, want %s", buf.String(), want)
	}
}

func TestSetLimit(t *testing.T) {
	src, err := Parse("select foo, bar from baz limit 4")
	if err != nil {
		t.Error(err)
	}
	limit := src.(*Select).Limit
	dst, err := Parse("select * from t")
	if err != nil {
		t.Error(err)
	}
	dst.(*Select).SetLimit(limit)
	buf := NewTrackedBuffer(nil)
	dst.Format(buf)
	want := "select * from t limit 4"
	if buf.String() != want {
		t.Errorf("limit: %q, want %s", buf.String(), want)
	}
	dst, err = Parse("select * from t union select * from s")
	if err != nil {
		t.Error(err)
	}
	dst.(*Union).SetLimit(limit)
	buf = NewTrackedBuffer(nil)
	dst.Format(buf)
	want = "select * from t union select * from s limit 4"
	if buf.String() != want {
		t.Errorf("order: %q, want %s", buf.String(), want)
	}
}

func TestWhere(t *testing.T) {
	var w *Where
	buf := NewTrackedBuffer(nil)
	w.Format(buf)
	if buf.String() != "" {
		t.Errorf("w.Format(nil): %q, want \"\"", buf.String())
	}
	w = NewWhere(WhereStr, nil)
	buf = NewTrackedBuffer(nil)
	w.Format(buf)
	if buf.String() != "" {
		t.Errorf("w.Format(&Where{nil}: %q, want \"\"", buf.String())
	}
}

func TestIsAggregate(t *testing.T) {
	f := FuncExpr{Name: NewColIdent("avg")}
	if !f.IsAggregate() {
		t.Error("IsAggregate: false, want true")
	}

	f = FuncExpr{Name: NewColIdent("Avg")}
	if !f.IsAggregate() {
		t.Error("IsAggregate: false, want true")
	}

	f = FuncExpr{Name: NewColIdent("foo")}
	if f.IsAggregate() {
		t.Error("IsAggregate: true, want false")
	}
}

func TestReplaceExpr(t *testing.T) {
	tcases := []struct {
		in, out string
	}{{
		in:  "select * from t where (select a from b)",
		out: ":a",
	}, {
		in:  "select * from t where (select a from b) and b",
		out: ":a and b",
	}, {
		in:  "select * from t where a and (select a from b)",
		out: "a and :a",
	}, {
		in:  "select * from t where (select a from b) or b",
		out: ":a or b",
	}, {
		in:  "select * from t where a or (select a from b)",
		out: "a or :a",
	}, {
		in:  "select * from t where not (select a from b)",
		out: "not :a",
	}, {
		in:  "select * from t where ((select a from b))",
		out: "(:a)",
	}, {
		in:  "select * from t where (select a from b) = 1",
		out: ":a = 1",
	}, {
		in:  "select * from t where a = (select a from b)",
		out: "a = :a",
	}, {
		in:  "select * from t where a like b escape (select a from b)",
		out: "a like b escape :a",
	}, {
		in:  "select * from t where (select a from b) between a and b",
		out: ":a between a and b",
	}, {
		in:  "select * from t where a between (select a from b) and b",
		out: "a between :a and b",
	}, {
		in:  "select * from t where a between b and (select a from b)",
		out: "a between b and :a",
	}, {
		in:  "select * from t where (select a from b) is null",
		out: ":a is null",
	}, {
		// exists should not replace.
		in:  "select * from t where exists (select a from b)",
		out: "exists (select a from b)",
	}, {
		in:  "select * from t where a in ((select a from b), 1)",
		out: "a in (:a, 1)",
	}, {
		in:  "select * from t where a in (0, (select a from b), 1)",
		out: "a in (0, :a, 1)",
	}, {
		in:  "select * from t where (select a from b) + 1",
		out: ":a + 1",
	}, {
		in:  "select * from t where 1+(select a from b)",
		out: "1 + :a",
	}, {
		in:  "select * from t where -(select a from b)",
		out: "-:a",
	}, {
		in:  "select * from t where interval (select a from b) aa",
		out: "interval :a aa",
	}, {
		in:  "select * from t where (select a from b) collate utf8",
		out: ":a collate utf8",
	}, {
		in:  "select * from t where func((select a from b), 1)",
		out: "func(:a, 1)",
	}, {
		in:  "select * from t where func(1, (select a from b), 1)",
		out: "func(1, :a, 1)",
	}, {
		in:  "select * from t where group_concat((select a from b), 1 order by a)",
		out: "group_concat(:a, 1 order by a asc)",
	}, {
		in:  "select * from t where group_concat(1 order by (select a from b), a)",
		out: "group_concat(1 order by :a asc, a asc)",
	}, {
		in:  "select * from t where group_concat(1 order by a, (select a from b))",
		out: "group_concat(1 order by a asc, :a asc)",
	}, {
		in:  "select * from t where substr(a, (select a from b), b)",
		out: "substr(a, :a, b)",
	}, {
		in:  "select * from t where substr(a, b, (select a from b))",
		out: "substr(a, b, :a)",
	}, {
		in:  "select * from t where convert((select a from b), json)",
		out: "convert(:a, json)",
	}, {
		in:  "select * from t where convert((select a from b) using utf8)",
		out: "convert(:a using utf8)",
	}, {
		in:  "select * from t where match((select a from b), 1) against (a)",
		out: "match(:a, 1) against (a)",
	}, {
		in:  "select * from t where match(1, (select a from b), 1) against (a)",
		out: "match(1, :a, 1) against (a)",
	}, {
		in:  "select * from t where match(1, a, 1) against ((select a from b))",
		out: "match(1, a, 1) against (:a)",
	}, {
		in:  "select * from t where case (select a from b) when a then b when b then c else d end",
		out: "case :a when a then b when b then c else d end",
	}, {
		in:  "select * from t where case a when (select a from b) then b when b then c else d end",
		out: "case a when :a then b when b then c else d end",
	}, {
		in:  "select * from t where case a when b then (select a from b) when b then c else d end",
		out: "case a when b then :a when b then c else d end",
	}, {
		in:  "select * from t where case a when b then c when (select a from b) then c else d end",
		out: "case a when b then c when :a then c else d end",
	}, {
		in:  "select * from t where case a when b then c when d then c else (select a from b) end",
		out: "case a when b then c when d then c else :a end",
	}}
	to := NewValArg([]byte(":a"))
	for _, tcase := range tcases {
		tree, err := Parse(tcase.in)
		if err != nil {
			t.Fatal(err)
		}
		var from *Subquery
		_ = Walk(func(node SQLNode) (kontinue bool, err error) {
			if sq, ok := node.(*Subquery); ok {
				from = sq
				return false, nil
			}
			return true, nil
		}, tree)
		if from == nil {
			t.Fatalf("from is nil for %s", tcase.in)
		}
		expr := ReplaceExpr(tree.(*Select).Where.Expr, from, to)
		got := String(expr)
		if tcase.out != got {
			t.Errorf("ReplaceExpr(%s): %s, want %s", tcase.in, got, tcase.out)
		}
	}
}

func TestExprFromValue(t *testing.T) {
	tcases := []struct {
		in  sqltypes.Value
		out SQLNode
		err string
	}{{
		in:  sqltypes.NULL,
		out: &NullVal{},
	}, {
		in:  sqltypes.NewInt64(1),
		out: NewIntVal([]byte("1")),
	}, {
		in:  sqltypes.NewFloat64(1.1),
		out: NewFloatVal([]byte("1.1")),
	}, {
		in:  sqltypes.MakeTrusted(sqltypes.Decimal, []byte("1.1")),
		out: NewFloatVal([]byte("1.1")),
	}, {
		in:  sqltypes.NewVarChar("aa"),
		out: NewStrVal([]byte("aa")),
	}, {
		in:  sqltypes.MakeTrusted(sqltypes.Expression, []byte("rand()")),
		err: "cannot convert valueImpl EXPRESSION(rand()) to AST",
	}}
	for _, tcase := range tcases {
		got, err := ExprFromValue(tcase.in)
		if tcase.err != "" {
			if err == nil || err.Error() != tcase.err {
				t.Errorf("ExprFromValue(%v) err: %v, want %s", tcase.in, err, tcase.err)
			}
			continue
		}
		if err != nil {
			t.Error(err)
		}
		if got, want := got, tcase.out; !reflect.DeepEqual(got, want) {
			t.Errorf("ExprFromValue(%v): %v, want %s", tcase.in, got, want)
		}
	}
}

func TestColNameEqual(t *testing.T) {
	var c1, c2 *ColName
	if c1.Equal(c2) {
		t.Error("nil columns equal, want unequal")
	}
	c1 = &ColName{
		Name: NewColIdent("aa"),
	}
	c2 = &ColName{
		Name: NewColIdent("bb"),
	}
	if c1.Equal(c2) {
		t.Error("columns equal, want unequal")
	}
	c2.Name = NewColIdent("aa")
	if !c1.Equal(c2) {
		t.Error("columns unequal, want equal")
	}
}

func TestColIdent(t *testing.T) {
	str := NewColIdent("Ab")
	if str.String() != "Ab" {
		t.Errorf("String=%s, want Ab", str.String())
	}
	if str.String() != "Ab" {
		t.Errorf("Val=%s, want Ab", str.String())
	}
	if str.Lowered() != "ab" {
		t.Errorf("Val=%s, want ab", str.Lowered())
	}
	if !str.Equal(NewColIdent("aB")) {
		t.Error("str.Equal(NewColIdent(aB))=false, want true")
	}
	if !str.EqualString("ab") {
		t.Error("str.EqualString(ab)=false, want true")
	}
	str = NewColIdent("")
	if str.Lowered() != "" {
		t.Errorf("Val=%s, want \"\"", str.Lowered())
	}
}

func TestColIdentMarshal(t *testing.T) {
	str := NewColIdent("Ab")
	b, err := json.Marshal(str)
	if err != nil {
		t.Fatal(err)
	}
	got := string(b)
	want := `"Ab"`
	if got != want {
		t.Errorf("json.Marshal()= %s, want %s", got, want)
	}
	var out ColIdent
	if err := json.Unmarshal(b, &out); err != nil {
		t.Errorf("Unmarshal err: %v, want nil", err)
	}
	if !reflect.DeepEqual(out, str) {
		t.Errorf("Unmarshal: %v, want %v", out, str)
	}
}

func TestColIdentSize(t *testing.T) {
	size := unsafe.Sizeof(NewColIdent(""))
	want := 2 * unsafe.Sizeof("")
	if size != want {
		t.Errorf("Size of ColIdent: %d, want 32", want)
	}
}

func TestTableIdentMarshal(t *testing.T) {
	str := NewTableIdent("Ab")
	b, err := json.Marshal(str)
	if err != nil {
		t.Fatal(err)
	}
	got := string(b)
	want := `"Ab"`
	if got != want {
		t.Errorf("json.Marshal()= %s, want %s", got, want)
	}
	var out TableIdent
	if err := json.Unmarshal(b, &out); err != nil {
		t.Errorf("Unmarshal err: %v, want nil", err)
	}
	if !reflect.DeepEqual(out, str) {
		t.Errorf("Unmarshal: %v, want %v", out, str)
	}
}

func TestHexDecode(t *testing.T) {
	testcase := []struct {
		in, out string
	}{{
		in:  "313233",
		out: "123",
	}, {
		in:  "ag",
		out: "encoding/hex: invalid byte: U+0067 'g'",
	}, {
		in:  "777",
		out: "encoding/hex: odd length hex string",
	}}
	for _, tc := range testcase {
		out, err := newHexVal(tc.in).HexDecode()
		if err != nil {
			if err.Error() != tc.out {
				t.Errorf("Decode(%q): %v, want %s", tc.in, err, tc.out)
			}
			continue
		}
		if !bytes.Equal(out, []byte(tc.out)) {
			t.Errorf("Decode(%q): %s, want %s", tc.in, out, tc.out)
		}
	}
}

func TestCompliantName(t *testing.T) {
	testcases := []struct {
		in, out string
	}{{
		in:  "aa",
		out: "aa",
	}, {
		in:  "1a",
		out: "_a",
	}, {
		in:  "a1",
		out: "a1",
	}, {
		in:  "a.b",
		out: "a_b",
	}, {
		in:  ".ab",
		out: "_ab",
	}}
	for _, tc := range testcases {
		out := NewColIdent(tc.in).CompliantName()
		if out != tc.out {
			t.Errorf("ColIdent(%s).CompliantNamt: %s, want %s", tc.in, out, tc.out)
		}
		out = NewTableIdent(tc.in).CompliantName()
		if out != tc.out {
			t.Errorf("TableIdent(%s).CompliantNamt: %s, want %s", tc.in, out, tc.out)
		}
	}
}

func TestColumns_FindColumn(t *testing.T) {
	cols := Columns{NewColIdent("a"), NewColIdent("c"), NewColIdent("b"), NewColIdent("0")}

	testcases := []struct {
		in  string
		out int
	}{{
		in:  "a",
		out: 0,
	}, {
		in:  "b",
		out: 2,
	},
		{
			in:  "0",
			out: 3,
		},
		{
			in:  "f",
			out: -1,
		}}

	for _, tc := range testcases {
		val := cols.FindColumn(NewColIdent(tc.in))
		if val != tc.out {
			t.Errorf("FindColumn(%s): %d, want %d", tc.in, val, tc.out)
		}
	}
}

func TestSplitStatementToPieces(t *testing.T) {
	testcases := []struct {
		input  string
		output string
	}{{
		input: "select * from table",
	}, {
		input:  "select * from table1; select * from table2;",
		output: "select * from table1; select * from table2",
	}, {
		input:  "select * from /* comment ; */ table;",
		output: "select * from /* comment ; */ table",
	}, {
		input:  "select * from table where semi = ';';",
		output: "select * from table where semi = ';'",
	}, {
		input:  "select * from table1;--comment;\nselect * from table2;",
		output: "select * from table1;--comment;\nselect * from table2",
	}, {
		input: "CREATE TABLE `total_data` (`id` int(11) NOT NULL AUTO_INCREMENT COMMENT 'id', " +
			"`region` varchar(32) NOT NULL COMMENT 'region name, like zh; th; kepler'," +
			"`data_size` bigint NOT NULL DEFAULT '0' COMMENT 'data size;'," +
			"`createtime` datetime NOT NULL DEFAULT NOW() COMMENT 'create time;'," +
			"`comment` varchar(100) NOT NULL DEFAULT '' COMMENT 'comment'," +
			"PRIMARY KEY (`id`))",
	}}

	for _, tcase := range testcases {
		if tcase.output == "" {
			tcase.output = tcase.input
		}

		stmtPieces, err := SplitStatementToPieces(tcase.input)
		if err != nil {
			t.Errorf("input: %s, err: %v", tcase.input, err)
			continue
		}

		out := strings.Join(stmtPieces, ";")
		if out != tcase.output {
			t.Errorf("out: %s, want %s", out, tcase.output)
		}
	}
}

func TestDAGTree(t *testing.T) {

	//query := "select * from t left join m on t.x=m.y where a like 1 and cc=32  and (dd=212 and ad=890)"

	query := "select a from t left join m on t.x=m.y where t.ddd=3+(5*(2+3)+7)"
	stmt, _ := Parse(query)

	switch stmt := stmt.(type) {

	case SelectStatement:
		{
			fmt.Println(stmt)

			sSelect := stmt.(*Select)

			fmt.Println(sSelect)
			dagWhere := NewDagWhere()

			RecursiveBuildDAGTree(dagWhere, sSelect.Where.Expr, 0)
			fmt.Println(dagWhere)
			mapList := make(map[int]([]*DAGWhere))
			RecursiveDAGTree(dagWhere, mapList, 0)

			for k, _ := range mapList {
				fmt.Println(k)

			}
		}

	}

}

func buildDAGWhere(Operator string, leftExpr Expr, rightExpr Expr, dagWhere *DAGWhere) *DAGWhere {
	dagWhere.Left = &DAGWhere{Expr: leftExpr}
	dagWhere.Right = &DAGWhere{Expr: rightExpr}
	return dagWhere
}

type DAGWhere struct {
	Expr      Expr
	Cost      float64 //单个表达式的代价   后面用于估计扫描page的代价，预估值
	Operator  string
	Left      *DAGWhere
	Right     *DAGWhere
	TreeDepth int

	LeftColName *ColName

	RightValue innodb.Value

	TableFilter *plan.TableFilter
}

func NewDagWhere() *DAGWhere {
	var dagRoot = new(DAGWhere)
	dagRoot.Cost = 0
	dagRoot.TreeDepth = 0
	return dagRoot
}

func RecursiveBuildDAGTree(root *DAGWhere, expr Expr, depth int) {

	root.Expr = expr
	root.TreeDepth = depth
	switch whereExpr := expr.(type) {

	case *ComparisonExpr:
		{

			root.Expr = whereExpr
			root.Operator = whereExpr.Operator
			root.Left = &DAGWhere{Expr: whereExpr.Left}
			root.Right = &DAGWhere{Expr: whereExpr.Right}
			//RecursiveBuildDAGTree(root.Left, whereExpr.Left, depth+1)
			rightValue, _ := whereExpr.Right.Eval()
			root.Left = &DAGWhere{Expr: whereExpr.Left}
			root.RightValue = rightValue
			switch leftExpr := whereExpr.Left.(type) {
			case *ColName:
				{
					root.TableFilter = &plan.TableFilter{
						TableName: leftExpr.Name.String(),
						Cost:      0,
						LeftValue: &plan.ColNameVariable{
							Name:      leftExpr.Name.String(),
							Qualifier: nil,
						},
					}
				}
			}

			root.TableFilter.RightValue = rightValue
			fmt.Println(rightValue.Raw())
			//RecursiveBuildDAGTree(root.Right, whereExpr.Right, depth+1)

		}

	case *AndExpr:
		{
			fmt.Println(whereExpr)
			root.Expr = whereExpr
			root.Operator = "&"
			root.Left = &DAGWhere{Expr: whereExpr.Left}
			root.Right = &DAGWhere{Expr: whereExpr.Right}
			RecursiveBuildDAGTree(root.Left, whereExpr.Left, depth+1)
			RecursiveBuildDAGTree(root.Right, whereExpr.Right, depth+1)
		}
	case *OrExpr:
		{
			fmt.Println(whereExpr)
			root.Expr = whereExpr
			root.Operator = "|"
			root.Left = &DAGWhere{Expr: whereExpr.Left}
			root.Right = &DAGWhere{Expr: whereExpr.Right}
			RecursiveBuildDAGTree(root.Left, whereExpr.Left, depth+1)
			RecursiveBuildDAGTree(root.Right, whereExpr.Right, depth+1)
		}
	case *ParenExpr:
		{
			fmt.Println(whereExpr)
			root.Expr = whereExpr
			root.Operator = "()"
			root.Left = &DAGWhere{Expr: whereExpr.Expr}
			RecursiveBuildDAGTree(root.Left, whereExpr.Expr, depth+1)
		}

	}

}

func RecursiveDAGTree(root *DAGWhere, mapList map[int]([]*DAGWhere), depth int) {
	dagList := make([]*DAGWhere, 0)
	if depth == 0 {
		if root != nil {
			dagList = append(dagList, root)
			RecursiveDAGTree(root.Left, mapList, depth+1)
			RecursiveDAGTree(root.Right, mapList, depth+1)
			mapList[depth] = dagList
		}
		return
	}
	if root != nil {
		if mapList[depth] == nil {
			mapList[depth] = dagList
		}

		mapList[depth] = append(mapList[depth], root)
		RecursiveDAGTree(root.Left, mapList, depth+1)
		RecursiveDAGTree(root.Right, mapList, depth+1)
	}
}

func RecursiveDagWalk(root *DAGWhere) {

}

func TestBig(t *testing.T) {
	//	var buff=make([]byte,0)
	//buff= append(buff,49)
	//buff = append(buff, 0,0,0)
	//
	//fmt.Println(buff)
	//
	//fmt.Println(util.ReadUB4Byte2Int(buff))
	//fmt.Println(util.ConvertInt4Bytes(1))
	var m uint8 = 49
	var buff = make([]byte, 0)
	buff = append(buff, m)
	fmt.Println(BytesToInt(buff))

	var resBuff = IntToByte(1)
	fmt.Println(resBuff)

}
func IntToByte(n int) []byte {
	x := int32(n)
	bytesBuffer := bytes.NewBuffer([]byte{})

	binary.Write(bytesBuffer, binary.BigEndian, x)
	return bytesBuffer.Bytes()
}

func BytesToInt(b []byte) int {
	bytesBuffer := bytes.NewBuffer(b)
	var x int32
	binary.Read(bytesBuffer, binary.BigEndian, x)
	return int(x)
}

func TestRune(t *testing.T) {
	var c uint8 = 49
	var d uint8 = 48
	fmt.Println(string(rune(c)))
	fmt.Println(string(rune(d)))

	stringBuffer := bytes.NewBufferString("")
	stringBuffer.WriteRune(rune(c))
	stringBuffer.WriteRune(rune(d))
	fmt.Println(stringBuffer.String())

}

type Element interface{} //可存入任何类型

type Stack struct {
	list []Element
}

//初始化栈
func NewStack() *Stack {
	return &Stack{list: make([]Element, 0)}
}

func (s *Stack) Len() int {
	return len(s.list)
}

//判断栈是否空
func (s *Stack) IsEmpty() bool {
	return len(s.list) == 0
}

//入栈
func (s *Stack) Push(x interface{}) {
	s.list = append(s.list, x)
}

//连续传入
func (s *Stack) PushList(x []Element) {
	s.list = append(s.list, x...)
}

//出栈
func (s *Stack) Pop() Element {
	if len(s.list) <= 0 {
		//fmt.Println("Stack is Empty")
		return nil
	} else {
		ret := s.list[len(s.list)-1]
		s.list = s.list[:len(s.list)-1]
		return ret
	}
}

//返回栈顶元素，空栈返nil
func (s *Stack) Top() Element {
	if s.IsEmpty() == true {
		//fmt.Println("Stack is Empty")
		return nil
	} else {
		return s.list[len(s.list)-1]
	}
}

//清空栈
func (s *Stack) Clear() {
	if len(s.list) == 0 {
		return
	}
	for i := 0; i < s.Len(); i++ {
		s.list[i] = nil
	}
	s.list = make([]Element, 0)
}

//打印测试
func (s *Stack) Show() {
	_len := len(s.list)
	for i := 0; i != _len; i++ {
		fmt.Println(s.Pop()) //这个注意:show会清空栈
	}
}
