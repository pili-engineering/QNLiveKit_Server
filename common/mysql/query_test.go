package mysql

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQueryString(t *testing.T) {
	querys := NewQuerys()
	querys.Set("uid", 123)
	querys.Like("sign", "七牛%")
	querys.Between("created", "22", "33")
	querys.LessThan("count", 23)
	querys.GreaterThan("age", 11)
	querys.LessThanEqual("top", 99)
	querys.GreaterThanEqual("job", 123)
	querys.NotEqual("name", "bob")
	querys.Null("job")
	querys.NotNull("salary")
	str := querys.String()

	assert.Equal(t, "uid = 123 AND sign LIKE '七牛%' AND created BETWEEN '22' AND '33' AND count < 23 AND age > 11 AND top <= 99 AND job >= 123 AND name != 'bob' AND job IS NULL AND salary IS NOT NULL", str)
}
