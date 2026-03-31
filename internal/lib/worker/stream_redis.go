package worker

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"iter"
	"strings"

	"github.com/redis/go-redis/v9"
)

var ErrRedisMissingData = errors.New("missing data field in redis stream message")

var _ Stream = (*RedisStream)(nil)

type RedisStream struct {
	streamName    string
	consumerGroup string
	consumerName  string
	read          int64
	client        *redis.Client
	ctx           context.Context
}

func NewRedisStream(client *redis.Client) (*RedisStream, error) {
	s := &RedisStream{
		client:        client,
		streamName:    "jobs:stream",
		consumerGroup: "jobs:group",
		consumerName:  "jobs:consumer:1",
		read:          10,
		ctx:           context.Background(),
	}
	err := s.prepare()
	return s, err
}

func (s *RedisStream) prepare() error {
	err := s.client.XGroupCreateMkStream(context.Background(), s.streamName, s.consumerGroup, "0").Err()
	if err != nil && !strings.Contains(err.Error(), "BUSYGROUP") {
		return fmt.Errorf("failed to create consumer group: %w", err)
	}
	return nil
}

func (s *RedisStream) Publish(ctx context.Context, job Job) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	data, err := json.Marshal(job)
	if err != nil {
		return fmt.Errorf("failed to marshal job: %w", err)
	}

	return s.client.XAdd(ctx, &redis.XAddArgs{
		Stream: s.streamName,
		ID:     fmt.Sprintf("%d", job.ID),
		Values: map[string]any{"job": string(data)},
	}).Err()
}

func (s *RedisStream) Read(ctx context.Context) iter.Seq2[*Message, error] {
	return func(yield func(*Message, error) bool) {
		for {
			if err := ctx.Err(); err != nil {
				yield(nil, err)
				return
			}

			streams, err := s.client.XReadGroup(s.ctx, &redis.XReadGroupArgs{
				Group:    s.consumerGroup,
				Consumer: s.consumerName,
				Streams:  []string{s.streamName, ">"},
				Count:    s.read,
			}).Result()

			if err != nil {
				if !yield(nil, err) {
					return
				}
				continue
			}

			stream := streams[0]
			for _, message := range stream.Messages {
				data, ok := message.Values["data"].(string)
				if !ok {
					if !yield(nil, ErrRedisMissingData) {
						return
					}
					continue
				}

				var job Job
				if err := json.Unmarshal([]byte(data), &job); err != nil {
					if !yield(nil, fmt.Errorf("failed json unmarshall: %w", err)) {
						return
					}
					continue
				}

				m := &Message{
					Job: job,
					Ack: func(ctx context.Context) error {
						return s.client.XAck(ctx, s.streamName, s.consumerGroup, message.ID).Err()
					},
				}
				if !yield(m, nil) {
					return
				}
			}
		}
	}
}
