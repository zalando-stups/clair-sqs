package tokens

import (
	"fmt"
	"log"
	"time"
)

type refreshCallback func(r ManagementRequest)

type scheduler struct {
	req      chan *refreshRequest
	quit     chan struct{}
	callback refreshCallback
}

type refreshRequest struct {
	tokenRequest ManagementRequest
	when         time.Duration
}

func NewScheduler(callback refreshCallback) *scheduler {
	req := make(chan *refreshRequest)
	quit := make(chan struct{})
	s := &scheduler{req, quit, callback}
	go func() {
		m := make(map[string]*time.Timer)

		for {
			select {
			case r := <-req:
				if _, has := m[r.tokenRequest.id]; has {
					log.Printf("Refresh of token %q was already scheduled. Skipping\n", r.tokenRequest.id)
				} else {
					time.AfterFunc(r.when, func() {
						s.callback(r.tokenRequest)
						// delete(m, r.tokenRequest.id)
					})
				}
			}
		}
	}()

	return s
}

func (s *scheduler) scheduleTokenRefresh(tr ManagementRequest, d time.Duration) {
	fmt.Printf("Scheduling refresh of token %q ...\n", tr.id)
	s.req <- &refreshRequest{tokenRequest: tr, when: d}
	fmt.Printf("Refresh of token %q scheduled in %v ...\n", tr.id, d)
}
