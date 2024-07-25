package goredispubsub

import (
	"sync"

	"github.com/kvizyx/wera"
	"github.com/redis/go-redis/v9"
)

type RedisSubscriber struct {
	sub      *redis.PubSub
	messages chan wera.Message
	msgOnce  *sync.Once
}

var _ wera.Subscriber = &RedisSubscriber{}

func newSub(sub *redis.PubSub) *RedisSubscriber {
	return &RedisSubscriber{
		sub:      sub,
		messages: make(chan wera.Message),
		msgOnce:  &sync.Once{},
	}
}

func (r *RedisSubscriber) Messages() <-chan wera.Message {
	r.msgOnce.Do(func() {
		r.messages = make(chan wera.Message)

		// this go-routine will exit with call to close subscriber
		go func() {
			for msg := range r.sub.Channel() {
				r.messages <- wera.Message{
					Data: []byte(msg.Payload),
				}
			}

			close(r.messages)
		}()
	})

	return r.messages
}

func (r *RedisSubscriber) Close() error {
	return r.sub.Close()
}