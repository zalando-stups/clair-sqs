package tokens

import (
	"fmt"
	"time"
)

type (
	scheduler struct {
		req      chan *refreshRequest
		quit     chan struct{}
		callback refreshCallback
	}

	refreshRequest struct {
		mgmtReq ManagementRequest
		when    time.Duration
		err     chan error
	}

	refreshCallback func(r ManagementRequest)
	scheduleFunc    func(d time.Duration, f func()) *time.Timer
)

var runner = time.AfterFunc

func NewScheduler(callback refreshCallback) *scheduler {
	req := make(chan *refreshRequest)
	quit := make(chan struct{})
	s := &scheduler{req, quit, callback}
	go func() {
		m := make(map[string]*time.Timer)

		for {
			select {
			case r := <-req:
				if _, has := m[r.mgmtReq.id]; has {
					r.err <- fmt.Errorf("Refresh of token %q was already scheduled. Skipping\n", r.mgmtReq.id)
				} else {
					m[r.mgmtReq.id] = runner(r.when, func() {
						s.callback(r.mgmtReq)
						delete(m, r.mgmtReq.id)
					})
					r.err <- nil
				}
			case <-quit:
				return
			}
		}
	}()

	return s
}

func (s *scheduler) scheduleTokenRefresh(mr ManagementRequest, d time.Duration) error {
	e := make(chan error)
	s.req <- &refreshRequest{mgmtReq: mr, when: d, err: e}
	return <-e
}

func (s *scheduler) Stop() {
	close(s.quit)
}
