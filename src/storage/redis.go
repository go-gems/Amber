package storage

import (
	"fmt"
	"log"
	"time"
	"amber/src/encoder"

	redisClient "github.com/gomodule/redigo/redis"
)

const LIST_REDIS = "Shorteners"

type redis struct{ pool *redisClient.Pool }

func (r *redis) saveToList(code string) {
	conn := r.pool.Get()
	defer conn.Close()
	_, err := conn.Do("LPUSH", LIST_REDIS, code)
	log.Println(err)
}

func (r *redis) GetAll() []Item {
	conn := r.pool.Get()
	defer conn.Close()
	var links = []Item{}
	values, e := redisClient.Values(conn.Do("SORT", LIST_REDIS,
		"BY", "*",
		"GET", "*->id",
		"GET", "*->url",
		"GET", "*->visits"))
	log.Println(values, e)
	if err := redisClient.ScanSlice(values, &links); err != nil {
		fmt.Println(err)
	}
	return links
}

func (r *redis) Load(code string) (string, error) {
	conn := r.pool.Get()
	defer conn.Close()

	urlString, err := redisClient.String(conn.Do("HGET", "Shortener:"+code, "url"))

	if err != nil {
		return "", err
	} else if len(urlString) == 0 {
		return "", ErrNoLink
	}

	_, err = conn.Do("HINCRBY", "Shortener:"+code, "visits", 1)

	return urlString, nil
}

func (r *redis) isAvailable(id string) bool {
	conn := r.pool.Get()
	defer conn.Close()

	exists, err := redisClient.Bool(conn.Do("EXISTS", "Shortener:"+id))
	if err != nil {
		return false
	}
	return !exists
}

func (r *redis) LoadInfo(code string) (*Item, error) {
	conn := r.pool.Get()
	defer conn.Close()

	values, err := redisClient.Values(conn.Do("HGETALL", "Shortener:"+code))
	if err != nil {
		return nil, err
	} else if len(values) == 0 {
		return nil, ErrNoLink
	}
	var shortLink Item
	err = redisClient.ScanStruct(values, &shortLink)
	if err != nil {
		return nil, err
	}

	return &shortLink, nil
}

func (r *redis) Close() error {
	return r.pool.Close()
}

func New(host, port string) (Service, error) {
	pool := &redisClient.Pool{
		MaxIdle:     10,
		IdleTimeout: 10 * time.Second,
		Dial: func() (redisClient.Conn, error) {
			return redisClient.Dial("tcp", fmt.Sprintf("%s:%s", host, port))
		},
	}

	return &redis{pool}, nil
}

func (r *redis) isUsed(id string) bool {
	conn := r.pool.Get()
	defer conn.Close()

	exists, err := redisClient.Bool(conn.Do("EXISTS", "Shortener:"+id))
	if err != nil {

		log.Println(err)
		return false
	}
	return exists
}

func (r *redis) Save(url string) (string, error) {
	conn := r.pool.Get()
	defer conn.Close()

	var id string
	l := 3

	for used := true; used; used = r.isUsed(id) {
		id = encoder.GetRandomId(l)
		l++
	}

	shortLink := Item{id, url, 0}
	log.Println(shortLink)
	key := "Shortener:" + id

	_, err := conn.Do("HMSET", redisClient.Args{key}.AddFlat(shortLink)...)
	if err != nil {
		return "", err
	}
	r.saveToList(key)
	return id, nil
}
