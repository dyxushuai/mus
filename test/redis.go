package models

import (
	"encoding/json"
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

type Server struct {
	Port     string `json:"port"`
	Password string `json:"password"`
	Method   string `json:"method"`
	Current  int64	`json:"current"`
	Limit    int64  `json:"limit"`
	Timeout  int64  `json:"timeout"`
	store 	 *Storage
}

type Storage struct {
	*redis.Pool
}

func NewStorage(host string) (s *Storage) {
	pool := redis.NewPool(func() (conn redis.Conn, err error) {
			conn, err = redis.Dial("tcp", host)
			return
		}, 3)
	s = &Storage{pool}
	return
}

func (self *Server) Init(s *Storage) {
	self.store = s
	if self.Timeout == 0 {
		self.Timeout = 60
	}
}

func GetServer(port string, s *Storage) (server *Server, err error) {
	server, err = s.GetBy(port)
	if err != nil {
		return
	}
	server.Init(s)
	return
}

func NewServer(server *Server, s *Storage) (err error) {
	err = s.SetBy(server)
	if err != nil {
		return
	}
	server.Init(s)
	return
}

func (self *Storage) GetBy(key string) (server *Server, err error) {
	var data []byte

	var conn = self.Get()
	defer conn.Close()

	data, err = redis.Bytes(conn.Do("GET", PREFIX + key))
	if err != nil {
		return
	}

	err = json.Unmarshal(data, server)
	return
}

func (self *Storage) SetBy(server *Server) (err error) {
	data, err := json.Marshal(server)
	if err != nil {
		return err
	}

	conn := self.Get()
	defer conn.Close()

	_, err = conn.Do("SET", PREFIX + server.Port, data)
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

func (self *Server) OverFlow() bool {
	return self.Current > self.Limit
}

func (self *Server) increase(key string, incr int) (score int64, err error) {
	score, err = self.store.IncrSize(key, incr)
	return
}

func (self *Server) getKeyBy(k int) (key string) {

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

func (self *Server) Config() (port, method, password string, timeout int64) {
	port = self.Port
	method = self.Method
	password = self.Password
	timeout = self.Timeout
	return
}

func (self *Server) IncreaseByHour(incr int) (score int64, err error) {
	key := self.getKeyBy(Hour)
	if key != "" {
		return
	}
	score, err = self.increase(key, incr)
	return
}

func (self *Server) IncreaseByDay(incr int) (score int64, err error) {
	key := self.getKeyBy(Day)
	if key != "" {
		return
	}
	score, err = self.increase(key, incr)
	return
}

func (self *Server) IncreaseByMonth(incr int) (score int64, err error) {
	key := self.getKeyBy(Month)
	if key != "" {
		return
	}
	score, err = self.increase(key, incr)
	return
}

func (self *Server) IncreaseByYear(incr int) (score int64, err error) {
	key := self.getKeyBy(Year)
	if key != "" {
		return
	}
	score, err = self.increase(key, incr)
	return
}

