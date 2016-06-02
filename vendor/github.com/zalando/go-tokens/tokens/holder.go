package tokens

import "fmt"

type holder struct {
	operationChannel chan *operation
	quitChannel      chan struct{}
}

type operationEnum int

const (
	get operationEnum = iota
	set
)

type operation struct {
	op          operationEnum
	tokenID     string
	accessToken *AccessToken
	response    chan *AccessToken
}

func newHolder() *holder {
	req := make(chan *operation)
	quit := make(chan struct{})

	go func() {
		m := make(map[string]*AccessToken)

		for {
			select {
			case r := <-req:
				r.response <- doOp(m, r)
			case <-quit:
				return
			}

		}
	}()

	return &holder{req, quit}
}

func doOp(m map[string]*AccessToken, o *operation) *AccessToken {
	switch o.op {
	case set:
		old := m[o.tokenID]
		m[o.tokenID] = o.accessToken
		return old
	case get:
		return m[o.tokenID]
	}
	panic(fmt.Errorf("Unknown operation: %v", o.op))
}

func (h *holder) get(tokenID string) *AccessToken {
	response := make(chan *AccessToken)
	h.operationChannel <- &operation{tokenID: tokenID, response: response, op: get}
	return <-response
}

func (h *holder) set(tokenID string, token *AccessToken) *AccessToken {
	response := make(chan *AccessToken)
	h.operationChannel <- &operation{tokenID: tokenID, accessToken: token, response: response, op: set}
	return <-response
}

func (h *holder) shutdown() {
	close(h.quitChannel)
}
