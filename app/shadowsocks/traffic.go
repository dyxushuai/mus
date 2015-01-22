//just add some methods to process server events
package shadowsocks

import (
	"github.com/JohnSmithX/mus/app/shadowsocks/lib"
	"sync"
	"log"
	"github.com/dropbox/godropbox/errors"
	"io"
)


//pipe between client and remote
type traffic struct {
	mu sync.Mutex
	//just local to remote => request
	//local to client
	flow int64

	//this field is for database method to record traffic flow
	recordFunc func(i *int)
}

func (t *traffic) doWithLock(fn func()) {
	t.mu.Lock()
	defer t.mu.Unlock()
	fn()
}

func (t *traffic) NewClient(c lib.SSClienter) {}

func (t *traffic) ClientConnClosed(c lib.SSClienter, err error){
	//err => EOF | i/o timeout
	log.Println(errors.New(err.Error()))
}
func (t *traffic) NewRemote(c lib.SSClienter){}

func (t *traffic) RemoteConnClosed(c lib.SSClienter, err error){
	c.Close()
}


func (t *traffic) ClientNewData(c lib.SSClienter, data []byte) {
	//do anything with data
	_, err := c.Remote().Write(data)
	if err != nil {
		log.Println(errors.New(err.Error()))
		c.Remote().Close()
	}
}

func (t *traffic) RemoteNewData(c lib.SSClienter, data []byte) {

	_, err := c.Write(data)
	if err != nil {
		if err != io.ErrClosedPipe {
			log.Println(errors.New(err.Error()))
		}
		c.Close()
	}
}

func (t *traffic) Record(i *int) {
	t.doWithLock(func () {
		t.flow += int64(*i)
	})
	t.recordFunc(i)
}
