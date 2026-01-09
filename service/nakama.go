package service

import (
	"context"
	"log"
	"time"

	"github.com/ascii8/nakama-go"
)

const (
	timeout = 1 * time.Second
)

type Service struct {
	cl     *nakama.Client
	authed bool
}

func New(serverUrl string, serverKey string) *Service {
	svc := &Service{}
	svc.cl = nakama.New(
		nakama.WithURL(serverUrl),
		nakama.WithServerKey(serverKey),
	)
	return svc
}

func (svc *Service) Auth(deviceId string) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	var (
		create   = true
		username = ""
	)
	if err := svc.cl.AuthenticateDevice(ctx, deviceId, create, username); err != nil {
		log.Printf("error: %s\n", err)
	} else {
		svc.authed = true
	}
}

func (svc *Service) RegisterScore(id string, score int64) {
	if !svc.authed {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	svc.cl.WriteLeaderboardRecord(ctx, &nakama.WriteLeaderboardRecordRequest{
		LeaderboardId: id,
		Record: &nakama.LeaderboardRecordWrite{
			Score:    score,
			Subscore: 0,
			Metadata: `{}`,
			Operator: nakama.OpType_BEST,
		},
	})
}

func (svc *Service) GetHighScore(id string) (int, bool) {
	if !svc.authed {
		return 0, false
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	resp, err := svc.cl.LeaderboardRecords(ctx, &nakama.LeaderboardRecordsRequest{
		LeaderboardId: id,
	})
	if err != nil {
		log.Printf("GetHighScore: %v\n", err)
		return 0, false
	}
	if resp == nil {
		return 0, false
	}

	if len(resp.Records) == 0 {
		return 0, true
	}
	return int(resp.Records[0].Score), true
}
