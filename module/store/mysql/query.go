package mysql

import (
	"bytes"
	"fmt"
	"reflect"
	"sort"
)

type querier interface {
	String() string
	Index() int
}

type operator struct {
	Column string
	Value  interface{}
	Op     string
	idx    int
}

// String build operator string
func (e operator) String() string {
	ret := ""
	val := reflect.ValueOf(e.Value)
	if val.Kind() == reflect.String {
		ret = fmt.Sprintf("%s %s '%v'", e.Column, e.Op, e.Value)
	} else {
		ret = fmt.Sprint(e.Column, " ", e.Op, " ", e.Value)
	}

	return ret
}

// Index implement querier
func (e operator) Index() int {
	return e.idx
}

type between struct {
	Column string
	Start  interface{}
	End    interface{}
	idx    int
}

// String build between string
func (b between) String() string {
	ret := ""
	if _, ok := b.Start.(string); ok {
		ret = fmt.Sprintf("%s BETWEEN '%v' AND '%v'", b.Column, b.Start, b.End)
	} else {
		ret = fmt.Sprintf("%s BETWEEN %v AND %v", b.Column, b.Start, b.End)
	}

	return ret
}

// Index implement querier
func (b between) Index() int {
	return b.idx
}

type null struct {
	Column string
	IsNull bool
	idx    int
}

// String build null string
func (n null) String() string {
	if n.IsNull {
		return fmt.Sprintf("%s IS NULL", n.Column)
	}

	return fmt.Sprintf("%s IS NOT NULL", n.Column)
}

// Index implement querier
func (n null) Index() int {
	return n.idx
}

// Querys mysql query, donot support concurrency call
type Querys struct {
	queries map[int]querier
	count   int
}

// NewQuerys new mysql querys
func NewQuerys() *Querys {
	return &Querys{
		queries: make(map[int]querier),
	}
}

// Set 设置 query 参数
func (q *Querys) Set(column string, value interface{}) {
	q.Equal(column, value)
}

func (q *Querys) setOperator(column, op string, value interface{}) {
	q.count++
	q.queries[q.count] = operator{
		Column: column,
		Value:  value,
		Op:     op,
		idx:    q.count,
	}
}

// String 生成 mysql query where
func (q Querys) String() string {
	if len(q.queries) == 0 {
		return ""
	}

	queries := []querier{}
	for _, v := range q.queries {
		queries = append(queries, v)
	}

	sort.Slice(queries, func(i, j int) bool {
		return queries[i].Index() < queries[j].Index()
	})

	buf := bytes.NewBuffer(nil)
	if len(queries) == 1 {
		buf.WriteString(queries[0].String())
		return buf.String()
	}

	buf.WriteString(queries[0].String())
	for _, v := range queries[1:] {
		buf.WriteString(" AND ")
		buf.WriteString(v.String())
	}

	return buf.String()
}

// Like  column like 'value'
func (q *Querys) Like(column, value string) *Querys {
	q.setOperator(column, "LIKE", value)
	return q
}

// LessThan column < value
func (q *Querys) LessThan(column string, value interface{}) *Querys {
	q.setOperator(column, "<", value)
	return q
}

// GreaterThan column > value
func (q *Querys) GreaterThan(column string, value interface{}) *Querys {
	q.setOperator(column, ">", value)
	return q
}

// LessThanEqual column <= value
func (q *Querys) LessThanEqual(column string, value interface{}) *Querys {
	q.setOperator(column, "<=", value)
	return q
}

// GreaterThanEqual column >= value
func (q *Querys) GreaterThanEqual(column string, value interface{}) *Querys {
	q.setOperator(column, ">=", value)
	return q
}

// NotEqual column != value
func (q *Querys) NotEqual(column string, value interface{}) *Querys {
	q.setOperator(column, "!=", value)
	return q
}

// Equal column = value
func (q *Querys) Equal(column string, value interface{}) *Querys {
	q.setOperator(column, "=", value)
	return q
}

// Between column Between start and end
func (q *Querys) Between(column string, start, end interface{}) *Querys {
	q.count++
	q.queries[q.count] = between{
		Column: column,
		Start:  start,
		End:    end,
		idx:    q.count,
	}

	return q
}

// Null column IS NULL
func (q *Querys) Null(column string) *Querys {
	q.count++
	q.queries[q.count] = null{column, true, q.count}
	return q
}

// NotNull column IS NOT NULL
func (q *Querys) NotNull(column string) *Querys {
	q.count++
	q.queries[q.count] = null{column, false, q.count}
	return q
}
