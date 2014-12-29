/*
*accept all of message eg. error, feedback from gorountine
*
*/
package manager



type message interface {
}

type Broadcast struct {
	errChan chan error
	msgChan chan message
}

func NewBroadcast() (bd *Broadcast) {
	bd = &Broadcast{}
	bd.errChan = make(chan error, numOfErr)
	bd.msgChan = make(chan message, numOfMsg)
	return
}


func (self *Broadcast) addError(err error) {
	select {
	case self.errChan <- err:
	default://avoid block with default
	}
}

func (self *Broadcast) addMsg(msg message) {
	select {
	case self.msgChan <- msg:
	default://avoid block with default
	}
}
