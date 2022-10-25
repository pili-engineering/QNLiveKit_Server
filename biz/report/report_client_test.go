package report

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/qbox/livekit/biz/model"
	"github.com/qbox/livekit/common/mysql"
	"github.com/qbox/livekit/utils/timestamp"
)

func TestRClient_SaveStatsSingleLive(t *testing.T) {
	saveStatsSingleLiveSetup()
	defer saveStatsSingleLiveTearDown()

	entity1 := model.StatsSingleLiveEntity{
		Type:      1,
		LiveId:    "live-1",
		UserId:    "user-1",
		BizId:     "",
		Count:     1,
		UpdatedAt: timestamp.Now(),
	}

	entity2 := model.StatsSingleLiveEntity{
		Type:      1,
		LiveId:    "live-1",
		UserId:    "user-2",
		BizId:     "",
		Count:     2,
		UpdatedAt: timestamp.Now(),
	}

	entity3 := model.StatsSingleLiveEntity{
		Type:      1,
		LiveId:    "live-1",
		UserId:    "user-3",
		BizId:     "",
		Count:     3,
		UpdatedAt: timestamp.Now(),
	}

	s := &RClient{}
	err := s.SaveStatsSingleLive(context.Background(), []*model.StatsSingleLiveEntity{&entity1, &entity2})
	assert.Nil(t, err)

	entity1.Count = 5
	entity2.Count = 10

	err = s.SaveStatsSingleLive(context.Background(), []*model.StatsSingleLiveEntity{&entity1, &entity2, &entity3})
	assert.Nil(t, err)
}

func saveStatsSingleLiveSetup() {
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

}

func saveStatsSingleLiveTearDown() {
}
