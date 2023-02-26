package duplicate_request

import (
	"sync"
	"time"

	"github.com/cyiafn/flight_information_system/server/logs"
)

const (
	defaultRequestIDCleanup = time.Minute * 5
)

type Filter struct {
	sync.RWMutex
	KnownRequests map[string]time.Time
	cleanUpTicker *time.Ticker
	cleanupChan   chan struct{}
}

func NewFilter() *Filter {
	filter := &Filter{
		KnownRequests: make(map[string]time.Time),
		cleanUpTicker: time.NewTicker(defaultRequestIDCleanup),
		cleanupChan:   make(chan struct{}),
	}

	go filter.cleanUp()
	return filter
}

func (d *Filter) IsAllowed(requestID string) bool {
	if d.has(requestID) {
		logs.Error("duplicate request ID received: %s, discarding requests", requestID)
		return false
	}
	d.Lock()
	defer d.Unlock()
	d.KnownRequests[requestID] = time.Now()

	return true
}

func (d *Filter) has(requestID string) bool {
	d.RLock()
	defer d.RUnlock()
	_, ok := d.KnownRequests[requestID]
	return ok
}

func (d *Filter) cleanUp() {
	for {
		select {
		case <-d.cleanupChan:
			return
		case <-d.cleanUpTicker.C:
			logs.Info("cleaning up request ids")
			d.Lock()

			for requestID, createTime := range d.KnownRequests {
				if d.isIDExpired(createTime) {
					delete(d.KnownRequests, requestID)
				}
			}

			d.Unlock()
		}
	}
}

func (d *Filter) isIDExpired(createdTime time.Time) bool {
	return createdTime.Add(defaultRequestIDCleanup).Before(time.Now())
}

func (d *Filter) Close() {
	d.cleanupChan <- struct{}{}
	logs.Info("terminating duplicate request filter ticker...")
}
