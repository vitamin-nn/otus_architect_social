package feed

import (
	"context"
	"encoding/json"

	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"github.com/vitamin-nn/otus_architect_social/server/internal/cache"
	"github.com/vitamin-nn/otus_architect_social/server/internal/repository"
)

const limitPerRequest = 100

func GetConsumerFunc(profileRepo repository.ProfileRepo, c cache.Feed) func(<-chan amqp.Delivery) {
	return func(msgCh <-chan amqp.Delivery) {
		for msg := range msgCh {
			f := new(repository.Feed)
			err := json.Unmarshal(msg.Body, f)
			if err != nil {
				log.Errorf("unmarshal error: %v", err)
				err = msg.Nack(false, false)
				if err != nil {
					log.Errorf("nack error: %v", err)
				}

				continue
			}

			err = msg.Ack(false)
			if err != nil {
				log.Errorf("ack error: %v", err)
			}

			cacheFriendsFeed(f, profileRepo, c)
		}
	}
}

func cacheFriendsFeed(f *repository.Feed, profileRepo repository.ProfileRepo, c cache.Feed) {
	err := c.Add(f.UserID, f)
	if err != nil {
		log.Errorf("add to cache error: %v", err)
	}

	ctx := context.Background()
	offset := 0

	for {
		friendList, err := profileRepo.GetFriendsProfileList(ctx, f.UserID, limitPerRequest, offset)
		if err != nil {
			log.Errorf("get friends list error: %v", err)

			break
		}

		for _, friend := range friendList {
			err := c.Add(friend.ID, f)
			if err != nil {
				log.Errorf("add to cache error: %v", err)
			}
		}

		if len(friendList) < limitPerRequest {
			break
		} else {
			offset += limitPerRequest
		}
	}
}
