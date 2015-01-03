package manager

import (
	"github.com/garyburd/redigo/redis"
	"time"
)

type Storage struct {
	pool *redis.Pool
}



func NewStorage(server, password string) (s *Storage) {
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

	s = &Storage{pool: pool}
	return
}

func (self * Storage) withConnDo(arg... interface {}) (reply interface{}, err error) {
	conn := self.pool.Get()
	defer conn.Close()
	reply, err = conn.DO(arg)
	return
}

func (self *Storage) Keys(pat string) (data []byte, err error) {
	var data []byte
	data, err = redis.Bytes(self.withConnDo("KEYS", PREFIX + pat))
	return
}

func (self *Storage) GetBy(key string) (data []byte, err error) {
	var data []byte
	data, err = redis.Bytes(self.withConnDo("GET", PREFIX + key))
	return
}

func (self *Storage) SetBy(key string, data []byte) (err error) {
	_, err = self.withConnDo("SET", key, data)
	return
}

func (self *Storage) IncrSize(key string, incr int) (score int64, err error) {
	score, err = redis.Bytes(self.withConnDo("INCRBY", PREFIX + key, incr))
	return
}

func (self *Storage) GetSize(key string) (score int64, err error) {
	score, err = redis.Bytes(self.withConnDo("GET", PREFIX + key))
	return
}

