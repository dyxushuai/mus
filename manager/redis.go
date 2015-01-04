package manager

import (
	"github.com/garyburd/redigo/redis"
	"time"
	"fmt"
)

const (
	PREFIX = "mus:"
	Hour = iota
	Day
	Month
	Year
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

func (self *Storage) getKeyBy(k int) (key string) {

	var year, month, day, hour int
	now := time.Now()
	year = now.Year()
	month = int(now.Month())
	day = now.Day()
	hour = now.Hour()

	switch k {
	case Hour:
		key = fmt.Sprintf("%s:%d:%d:%d:%d", self.Port, year, month, day, hour)
	case Day:
		key = fmt.Sprintf("%s:%d:%d:%d", self.Port, year, month, day)
	case Month:
		key = fmt.Sprintf("%s:%d:%d", self.Port, year, month)
	case Year:
		key = fmt.Sprintf("%s:%d", self.Port, year)
	default:

	}
	return
}

func (self *Storage) Keys(pat string) (data []byte, err error) {
	var data []byte
	data, err = redis.Bytes(self.withConnDo("KEYS", PREFIX + pat))
	return
}

func (self *Storage) GetServer(key string) (data []byte, err error) {
	var data []byte
	data, err = redis.Bytes(self.withConnDo("GET", PREFIX + key))
	return
}

func (self *Storage) SetServer(key string, data []byte) (err error) {
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

func (self *Storage) IncreaseByHour(incr int) (score int64, err error) {
	key := self.getKeyBy(Hour)
	if key != "" {
		return
	}
	score, err = self.IncrSize(key, incr)
	return
}

func (self *Storage) IncreaseByDay(incr int) (score int64, err error) {
	key := self.getKeyBy(Day)
	if key != "" {
		return
	}
	score, err = self.IncrSize(key, incr)
	return
}

func (self *Storage) IncreaseByMonth(incr int) (score int64, err error) {
	key := self.getKeyBy(Month)
	if key != "" {
		return
	}
	score, err = self.IncrSize(key, incr)
	return
}

func (self *Storage) IncreaseByYear(incr int) (score int64, err error) {
	key := self.getKeyBy(Year)
	if key != "" {
		return
	}
	score, err = self.IncrSize(key, incr)
	return
}
