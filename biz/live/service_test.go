package live

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/qbox/livekit/biz/model"
	"github.com/qbox/livekit/common/mysql"
)

func TestService_listOnlineRooms(t *testing.T) {
	roomsSetup()
	defer roomsTearDown()

	s := &Service{}
	liveIds, err := s.listOnlineRooms(context.Background())
	assert.Nil(t, err)
	assert.Equal(t, 5, len(liveIds))
}

func roomsSetup() {
	mysql.Init(&mysql.ConfigStructure{
		Host:     "localhost",
		Port:     3306,
		Username: "root",
		Password: "123456",
		Database: "live_test",
		Default:  "live",
	}, &mysql.ConfigStructure{
		Host:     "localhost",
		Port:     3306,
		Username: "root",
		Password: "123456",
		Database: "live_test",
		Default:  "live",
		ReadOnly: true,
	})

	mysql.GetLive().AutoMigrate(model.LiveEntity{})
	for i := 0; i < 10; i++ {
		liveEntity := model.LiveEntity{
			LiveId: fmt.Sprintf("test_live_%d", i),
			Status: i % 2,
		}
		mysql.GetLive().Save(&liveEntity)
	}
}

func roomsTearDown() {
	db := mysql.GetLive()
	db.DropTableIfExists(model.LiveEntity{})
}
