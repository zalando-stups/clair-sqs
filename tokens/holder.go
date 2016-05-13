package tokens

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
	tokenId     string
	accessToken *AccessToken
	response    chan *AccessToken
}

func NewHolder() *holder {
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
		old := m[o.tokenId]
		m[o.tokenId] = o.accessToken
		return old
	case get:
		return m[o.tokenId]
	}
	return nil
}

func (h *holder) get(tokenId string) *AccessToken {
	response := make(chan *AccessToken)
	h.operationChannel <- &operation{tokenId: tokenId, response: response, op: get}
	return <-response
}

func (h *holder) set(tokenId string, token *AccessToken) *AccessToken {
	response := make(chan *AccessToken)
	h.operationChannel <- &operation{tokenId: tokenId, accessToken: token, response: response, op: set}
	return <-response
}
