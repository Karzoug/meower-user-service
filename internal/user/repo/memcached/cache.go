package memcached

import (
	"bytes"
	"encoding/gob"
	"errors"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/rs/xid"

	"github.com/Karzoug/meower-common-go/memcached"

	"github.com/Karzoug/meower-user-service/internal/user/entity"
	"github.com/Karzoug/meower-user-service/internal/user/repo"
)

type cache struct {
	client memcached.Client
}

func NewUserCache(client memcached.Client) cache {
	return cache{
		client: client,
	}
}

func (c cache) GetOne(id xid.ID) (entity.UserShortProjection, error) {
	item, err := c.client.Get(id.String())
	if err != nil {
		if errors.Is(err, memcache.ErrCacheMiss) {
			return entity.UserShortProjection{}, repo.ErrRecordNotFound
		}
		return entity.UserShortProjection{}, err
	}

	b := bytes.NewReader(item.Value)
	var ui entity.UserShortProjection
	if err := gob.NewDecoder(b).Decode(&ui); err != nil {
		return entity.UserShortProjection{}, err
	}

	return ui, nil
}

func (c cache) GetMany(ids []xid.ID) (users []entity.UserShortProjection, missed []xid.ID, err error) {
	keys := make([]string, len(ids))
	for i, id := range ids {
		keys[i] = id.String()
	}

	items, err := c.client.GetMulti(keys)
	if err != nil {
		return nil, nil, err
	}

	res := make([]entity.UserShortProjection, 0, len(items))

	for _, item := range items {
		b := bytes.NewReader(item.Value)
		var ui entity.UserShortProjection
		if err := gob.NewDecoder(b).Decode(&res); err != nil {
			return nil, nil, err
		}

		res = append(res, ui)
	}

	missed = make([]xid.ID, 0)
	for _, key := range keys {
		if _, ok := items[key]; !ok {
			id, _ := xid.FromString(key)
			missed = append(missed, id)
		}
	}

	return res, missed, nil
}

func (c cache) Set(id xid.ID, u entity.UserShortProjection, ttl int32) error {
	var b bytes.Buffer
	if err := gob.NewEncoder(&b).Encode(u); err != nil {
		return err
	}

	return c.client.Set(&memcache.Item{
		Key:        id.String(),
		Value:      b.Bytes(),
		Expiration: ttl,
	})
}
