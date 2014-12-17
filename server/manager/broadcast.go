
/*
*accept all of message eg. error, feedback from gorountine
*
*/
package manager


type message interface {
}

var msgChan chan message

func addMsg(msg message) {
	go func() {
		msgChan <- msg
	}()

}


type Broadcast struct {
	errMessage chan
}
