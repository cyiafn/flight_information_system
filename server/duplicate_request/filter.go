package duplicate_request

import (
	"sync"
	"time"

	"github.com/cyiafn/flight_information_system/server/logs"
)

const (
	// defaultRequestIDCleanup is the time it will run a cleanup job on all the requestIDs to reduce memory use
	// we assume at the end of 5 minutes, all duplicate requests would have been received.
	defaultRequestIDCleanup = time.Minute * 5
)

// Filter is a guard against duplicate requests in the event that at-least-once mode is set.
// This is CONCURRENT-SAFE
type Filter struct {
	sync.RWMutex
	KnownRequests       map[string]time.Time
	KnownRequestReplies map[string][][]byte
	cleanUpTicker       *time.Ticker
	cleanupChan         chan struct{}
}

// NewFilter instantiates a new filter object
func NewFilter() *Filter {
	filter := &Filter{
		KnownRequests: make(map[string]time.Time),
		cleanUpTicker: time.NewTicker(defaultRequestIDCleanup),
		cleanupChan:   make(chan struct{}),
	}

	go filter.cleanUp()
	return filter
}

func (d *Filter) RegisterResponse(requestID string, response [][]byte) {
	d.Lock()
	defer d.Unlock()
	d.KnownRequestReplies[requestID] = response
}

func (d *Filter) GetKnownResponse(requestID string) [][]byte {
	d.RLock()
	defer d.RUnlock()
	return d.KnownRequestReplies[requestID]
}

// IsAllowed simply checks if there is a duplicate requestID previously or not, if there is, it will return false, else it will add the new request and return true
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

// has checks if the requestID exists in the filter
func (d *Filter) has(requestID string) bool {
	d.RLock()
	defer d.RUnlock()
	_, ok := d.KnownRequests[requestID]
	return ok
}

// cleanup is ran as a goroutine where based on the ticker created on instantiation, it will cleanup the all requestIDs
// that have been around for more than the specified timing
func (d *Filter) cleanUp() {
	for {
		select {
		// stops the cleanup
		case <-d.cleanupChan:
			return
			// every 5 minutes this will trigger
		case <-d.cleanUpTicker.C:
			logs.Info("cleaning up request ids")
			d.Lock()

			// deletes all requestIDs that have been around for the defaultExpiryTiming
			for requestID, createTime := range d.KnownRequests {
				if d.isIDExpired(createTime) {
					delete(d.KnownRequests, requestID)
					delete(d.KnownRequestReplies, requestID)
				}
			}

			d.Unlock()
		}
	}
}

// isIDExpired checks if the requestID has existed for more than or equal to defaultRequestIDCleanup
func (d *Filter) isIDExpired(createdTime time.Time) bool {
	return createdTime.Add(defaultRequestIDCleanup).Before(time.Now())
}

// Close gracefully closes this filter
func (d *Filter) Close() {
	d.cleanupChan <- struct{}{}
	logs.Info("terminating duplicate request filter ticker...")
}
