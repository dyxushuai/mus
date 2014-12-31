package config

import (
	"github.com/garyburd/redigo/redis"
	"time"
)



type Storage struct {
	*redis.Pool
}



func NewStorage() (s *Storage) {
	server, password := REDIS_SERVER, REDIS_PASSWORD
	pool := &redis.Pool{
		MaxIdle: 3,
		IdleTimeout: 240 * time.Second,
		Dial: func () (redis.Conn, error) {
			c, err := redis.Dial("tcp", server)
			if err != nil {
				return nil, err
			}
			if _, err := c.Do("AUTH", password); err != nil {
				c.Close()
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}

	s = &Storage{pool}
	return
}

func (self *Storage) Keys(pat string) (data []byte, err error) {
	var data []byte
	var conn = self.Get()
	defer conn.Close()
	data, err = redis.Bytes(conn.Do("KEYS", PREFIX + "*"))
	return
}

func (self *Storage) GetBy(key string) (data []byte, err error) {
	var data []byte

	var conn = self.Get()
	defer conn.Close()

	data, err = redis.Bytes(conn.Do("GET", PREFIX + key))
	return
}

func (self *Storage) SetBy(key string, data []byte) (err error) {
	conn := self.Get()
	defer conn.Close()

	_, err = conn.Do("SET", PREFIX + key, data)
	return
}

func (self *Storage) IncrSize(key string, incr int) (score int64, err error) {
	var conn = self.Get()
	defer conn.Close()

	score, err = redis.Int64(conn.Do("INCRBY", PREFIX + key, incr))
	return
}

func (self *Storage) GetSize(key string) (score int64, err error) {
	var conn = self.Get()
	defer conn.Close()

	score, err = redis.Int64(conn.Do("GET", PREFIX + key))
	return
}

