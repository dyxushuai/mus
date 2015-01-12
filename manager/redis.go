package manager

import (
	"github.com/garyburd/redigo/redis"
	"time"
	"fmt"
	"encoding/json"
	"strings"
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

func (self * Storage) Test() (err error) {
	_, err = self.doWithConn("PING")
	return
}
//remove PREFIX form key if it has
func removePrefix(key string) string {
	if strings.HasPrefix(key, PREFIX) {
		return strings.Replace(key, PREFIX, "", 1)
	}
	return key
}

func (self *Storage) doWithConn(keyName string, arg... interface {}) (reply interface{}, err error) {
	conn := self.pool.Get()
	defer conn.Close()
	reply, err = conn.Do(keyName, arg...)
	if err != nil {
		err = newError(err.Error())
	}
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
	keys, err = redis.Strings(self.doWithConn("KEYS", PREFIX + pat))
	return
}

func (self *Storage) GetServer(key string) (server *Server, err error) {
	data, err := redis.Bytes(self.doWithConn("GET", PREFIX + key))

	err = json.Unmarshal(data, &server)
	return
}

//pat -> "server:**" will get all exsited servers
func (self *Storage) GetServers(pat string) (servers []*Server, err error) {
	keys, err := self.Keys(pat)
	if err != nil {
		return
	}
	for _, key := range keys {
		key = removePrefix(key)
		if server, err := self.GetServer(key); err == nil {
			servers = append(servers, server)
		}
	}
	return
}

func (self *Storage) SetServer(key string, server *Server) (err error) {
	data, err := json.Marshal(server)
	if err != nil {
		return
	}
	_, err = self.doWithConn("SET", key, data)
	return
}

func (self *Storage) IncrSize(key string, incr int) (score int64, err error) {
	score, err = redis.Int64(self.doWithConn("INCRBY", PREFIX + key, incr))
	return
}

func (self *Storage) GetSize(key string) (score int64, err error) {
	score, err = redis.Int64(self.doWithConn("GET", PREFIX + key))
	return
}

func (self *Storage) IncreaseByHour(port string, incr int) (score int64, err error) {
	key := self.getKeyBy(Hour, port)
	if key != "" {
		return
	}
	score, err = self.IncrSize(key, incr)
	return
}

func (self *Storage) IncreaseByDay(port string, incr int) (score int64, err error) {
	key := self.getKeyBy(Day, port)
	if key != "" {
		return
	}
	score, err = self.IncrSize(key, incr)
	return
}

func (self *Storage) IncreaseByMonth(port string, incr int) (score int64, err error) {
	key := self.getKeyBy(Month, port)
	if key != "" {
		return
	}
	score, err = self.IncrSize(key, incr)
	return
}

func (self *Storage) IncreaseByYear(port string, incr int) (score int64, err error) {
	key := self.getKeyBy(Year, port)
	if key != "" {
		return
	}
	score, err = self.IncrSize(key, incr)
	return
}
