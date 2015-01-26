package db

import (
	"github.com/garyburd/redigo/redis"
	"time"
	"fmt"
)

const (
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



func (self * Storage) Test() (err error) {
	_, err = self.doWithConn("PING")
	return
}



func (self *Storage) doWithConn(keyName string, arg... interface {}) (reply interface{}, err error) {
	conn := self.pool.Get()
	defer conn.Close()
	reply, err = conn.Do(keyName, arg...)
	return
}

func (self *Storage) getKeyBy(k int, id string) (key string) {

	var year, month, day, hour int
	now := time.Now()
	year = now.Year()
	month = int(now.Month())
	day = now.Day()
	hour = now.Hour()

	switch k {
	case Hour:
		key = fmt.Sprintf("%s:%d:%d:%d:%d", id, year, month, day, hour)
	case Day:
		key = fmt.Sprintf("%s:%d:%d:%d", id, year, month, day)
	case Month:
		key = fmt.Sprintf("%s:%d:%d", id, year, month)
	case Year:
		key = fmt.Sprintf("%s:%d", id, year)
	default:

	}
	return
}



func (self *Storage) Keys(pat string) (keys []string, err error) {
	keys, err = redis.Strings(self.doWithConn("KEYS", pat))
	return
}

func (self *Storage) GetByt(key string) (data []byte, err error) {
	data, err = redis.Bytes(self.doWithConn("GET", key))
	return
}
func (self *Storage) GetStr(key string) (str string, err error) {
	str, err = redis.String(self.doWithConn("GET", key))
	return
}

func (self *Storage) GetNum(key string) (score int64, err error) {
	score, err = redis.Int64(self.doWithConn("GET", key))
	return
}

func (self *Storage) Set(key string, data []byte) (err error) {
	_, err = self.doWithConn("SET", key, data)
	return
}

func (self *Storage) Del(key string) (err error) {
	_, err = self.doWithConn("DEL", key)
	return
}

func (self *Storage) Incr(key string, incr int) (score int64, err error) {
	score, err = redis.Int64(self.doWithConn("INCRBY", key, incr))
	return
}



func (self *Storage) IncreaseByHour(port string, incr int) (score int64, err error) {
	key := self.getKeyBy(Hour, port)
	if key != "" {
		return
	}
	score, err = self.Incr(key, incr)
	return
}

func (self *Storage) IncreaseByDay(port string, incr int) (score int64, err error) {
	key := self.getKeyBy(Day, port)
	if key != "" {
		return
	}
	score, err = self.Incr(key, incr)
	return
}

func (self *Storage) IncreaseByMonth(port string, incr int) (score int64, err error) {
	key := self.getKeyBy(Month, port)
	if key != "" {
		return
	}
	score, err = self.Incr(key, incr)
	return
}

func (self *Storage) IncreaseByYear(port string, incr int) (score int64, err error) {
	key := self.getKeyBy(Year, port)
	if key != "" {
		return
	}
	score, err = self.Incr(key, incr)
	return
}
