package mysql_test

import (
	"fmt"
	"testing"

	mysql2 "github.com/qbox/livekit/module/store/mysql"

	"github.com/stretchr/testify/assert"

	"github.com/qbox/livekit/utils/timestamp"
)

type BatchTest struct {
	ID        uint `gorm:"primary_key"`
	Content   string
	MessageID string
	CreatedAt timestamp.Timestamp
	DeletedAt *timestamp.Timestamp
}

func (BatchTest) TableName() string {
	return "batch_tests"
}

func TestBatchInsert(t *testing.T) {
	mysql2.Init(&mysql2.ConfigStructure{
		Host:     "localhost",
		Port:     3306,
		Username: "root",
		Password: "123456",
		Database: "live_test",
	})

	mysql2.Get().DropTableIfExists(&BatchTest{})
	mysql2.Get().AutoMigrate(&BatchTest{})

	bts := make([]interface{}, 0)
	for i := 0; i < 100; i++ {
		bts = append(bts, BatchTest{
			Content:   fmt.Sprintf("您的验证码是 %d", i),
			MessageID: fmt.Sprintf("1%d2", i),
			CreatedAt: timestamp.Now(),
		})
	}

	rows, err := mysql2.BatchInsert(mysql2.Get(), bts)
	assert.Nil(t, err)
	assert.Equal(t, rows, int64(100))
}
