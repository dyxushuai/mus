package models

import (
	"github.com/JohnSmithX/mus/app/utils"
	"github.com/JohnSmithX/mus/app/db"
	ss "github.com/JohnSmithX/mus/app/shadowsocks"
	"encoding/json"
	uuid "github.com/satori/go.uuid"
	"time"
	"strings"

	"github.com/dropbox/godropbox/errors"
	"sync"

)

//for redis key string
const (
	serverPrefix = "mus:server:"
	flowPrefix = "mus:flow:"
)

var (
	rdPool *db.Storage
)

func InitDb(store *db.Storage) {
	rdPool = store
}

type ServerI interface {
	Update() (err error)
	Delete() (err error)
	JSON() (result []byte, err error)
	IsStopped() bool
	Stop()
	Start()
	Restart()
	Key() string
}


type Server struct {
	mu   				sync.Mutex
	proxy				ss.Proxyer

	//current flow
	current       		int64

	Id 					uuid.UUID		`json:"id"`
	CreateTime			utils.Time		`json:"create_at"`
	UpdateTime			utils.Time		`json:"update_at"`
	Port 				string			`json:"port"`
	Method       		string       	`json:"method"`
	Password      		string       	`json:"password"`
	Limit         		int64        	`json:"limit"`
	Timeout       		int64        	`json:"timeout"`
}


func New(port, method, password string, limit, timeout int64) (server *Server,err error) {
	
	server = &Server{}

	server.proxy, err = ss.New(":" + port, method, password, time.Duration(timeout), func(int){})
	if err != nil {
		err = errors.New(err.Error())
		return
	}
	server.Port = port
	server.Method = method
	server.Password = password
	server.Limit = limit
	server.Timeout = timeout
	err = server.initialize()
	
	return
}

func addPrefix(key, prefix string) string {
	if strings.HasPrefix(key, prefix) {
		return key
	}
	return prefix + key
}

func (self *Server) serverKey() string {
	return addPrefix(self.Port, serverPrefix)
}

func (self *Server) flowKey() string {
	return addPrefix(self.Port, flowPrefix)
}

func (self *Server) initialize() (err error) {
	defer func() {
		if err != nil {
			err = errors.New(err.Error())
		}
	}()
	self.Id = uuid.NewV4()
	self.upTime()
	self.crTime()
	err = self.save()
	if err != nil {
		return
	}
	_, err = self.getCurrent()
	return
}

func (self *Server) doWithLock(fn func()) {
	self.mu.Lock()
	defer self.mu.Unlock()
	fn()
}

func (self *Server) getCurrent() (n int64, err error){
	n, err = rdPool.GetNum(self.flowKey())
	if err != nil {
		err = errors.New(err.Error())
		self.doWithLock(func(){
			self.current = 0
		})
		return
	}
	self.doWithLock(func(){
		self.current = n
	})

	return
}

//update time at Now
func (self *Server) upTime() {
	self.UpdateTime = utils.Time(time.Now())
}

//create time at Now
func (self *Server) crTime() {
	self.CreateTime = utils.Time(time.Now())
}


func (self *Server) save() (err error) {

	defer func() {
		if err != nil {
			err = errors.New(err.Error())
		}
	}()
	data, err := json.Marshal(self)
	if err != nil {
		return
	}
	err = rdPool.Set(self.serverKey(), data)
	return
}


func (self *Server) Update() (err error) {
	self.upTime()
	err = self.save()
	return
}

func (self *Server) Delete() (err error) {
	err = rdPool.Del(self.serverKey())
	err = rdPool.Del(self.flowKey())
	self.Stop()
	return
}

func (self *Server) JSON() (result []byte, err error) {
	result, err = json.Marshal(self)
	return
}

func (self *Server) IsStopped() bool {
	return self.proxy.IsStopped()
}

func (self *Server) Stop() {
	if self.IsStopped() {
		return
	}
	self.proxy.Stop()
}

func (self *Server) Restart() {
	self.Stop()
	self.Start()
}

func (self *Server) Start() {
	if self.IsStopped() {
		self.proxy.Listen()
	}
}

func (self *Server) Key() string {
	return self.Port
}

//operate servers from redis
func GetServerFromRedis(port string) (server *Server, err error) {
	data, err :=  rdPool.GetByt(addPrefix(port, serverPrefix))

	if err != nil {

		return
	}

	server = &Server{}
	err = json.Unmarshal(data, server)

	if err != nil {

		return
	}

	size, _ := rdPool.GetNum(addPrefix(port, flowPrefix))
	if err != nil {
		server.current = 0
	} else {
		server.current = size
	}

	return
}

func GetServersFromRedis( ports ...string) (servers []*Server, err error) {
	if len(ports) == 0 {
		err = errors.New("Need port but port is nil")
		return
	}

	for _, port := range ports {
		if server, er := GetServerFromRedis(string(port)); er == nil {
			servers = append(servers, server)
		} else {
			err = er
			return
		}
	}
	return
}

func GetAllServersFromRedis() (servers []*Server, err error) {

	defer func() {
		if err != nil {
			err = errors.New(err.Error())
		}
	}()
	keys, err := rdPool.Keys(serverPrefix + "**")

	if err != nil {
		return
	}
	for _, key := range keys {

		if server, er := GetServerFromRedis(key); er == nil {
			servers = append(servers, server)
		} else {
			err = er

			return
		}
	}
	return
}


