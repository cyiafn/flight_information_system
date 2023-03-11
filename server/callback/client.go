package callback

import (
	"context"
	"time"

	"github.com/cyiafn/flight_information_system/server/dto"
	"github.com/cyiafn/flight_information_system/server/dto/status_code"
	"github.com/cyiafn/flight_information_system/server/logs"
	"github.com/cyiafn/flight_information_system/server/net"
	"github.com/cyiafn/flight_information_system/server/server"
	"github.com/cyiafn/flight_information_system/server/utils"
	"github.com/cyiafn/flight_information_system/server/utils/bytes"
	"github.com/cyiafn/flight_information_system/server/utils/collections"
	"github.com/cyiafn/flight_information_system/server/utils/predicates"
	"github.com/cyiafn/flight_information_system/server/utils/rpc"
	"github.com/cyiafn/flight_information_system/server/utils/worker_pools"
	"github.com/teris-io/shortid"
)

// Client is a callback client designed to handle generic subscribers and notifying of those subscribers.
type Client[T comparable] struct {
	// NotifiableClients are the set of IP:Port addresses for each item subscribed.
	NotifiableClients map[T]*collections.Set[string]
}

// NewClient is an instantiate for the Client.
func NewClient[T comparable]() *Client[T] {
	return &Client[T]{
		NotifiableClients: make(map[T]*collections.Set[string]),
	}
}

// workerPoolJob is a request object designed to store the necessary details for the job.
type workerPoolJob struct {
	// Payload to deliver to subscriber
	Payload []byte
	// IP:Port of subscriber
	Addr string
}

// Subscribe subscribes a client to be notified on change of a particular item with an expiry duration defined.
// note that IP addresses are propagated through the program in the context object.
func (c *Client[T]) Subscribe(ctx context.Context, item T, expireDuration time.Duration) {
	// gets the IP address from the ctx
	addr := server.GetIPAddr(ctx)

	// If the item to subscribe to doesn't exist yet, we need to allocate memory for a new set at that key.
	if _, ok := c.NotifiableClients[item]; !ok {
		c.NotifiableClients[item] = collections.NewSet[string]()
	}

	// Adds the client to that set to be subscribed. We don't care if it replaces.
	c.NotifiableClients[item].MustAdd(addr)
	logs.Info("Client: %s has successfully been subscribed to item: %s", addr, utils.DumpJSON(item))
	// Runs a go routine that removes that user from the subscription list once they have expired.
	// Note that we can do this as goroutines are extremely cheap and only take up a minimal amount of memory.
	// A good optimization will be to continue to store this data, and only clean it up after 5 minutes or so, with a condition
	// that notify will not deliver changes to that particular user if already expired.
	go c.cleanup(item, addr, expireDuration)
}

// cleanup simply blocks and removes the users from the subscription once their subscription duration ends.
func (c *Client[T]) cleanup(item T, addr string, expireDuration time.Duration) {
	// timer
	timer := time.NewTimer(expireDuration)

	// on timer finish
	<-timer.C
	logs.Info("removing address: %s for item: %s from subscription", addr, utils.DumpJSON(item))
	// removes person from subscription
	c.NotifiableClients[item].MustRemove(addr)
}

// Notify notifies all subscribers for that particular item
func (c *Client[T]) Notify(item T, respType dto.ResponseType, payload any, err error) error {
	// if that item does not exist in the map, we don't do anything
	if _, ok := c.NotifiableClients[item]; !ok {
		return nil
	}

	// if there is no clients for that item, we don't do anything
	if clients := c.NotifiableClients[item]; clients.Len() == 0 {
		logs.Info("no client to notify")
		return nil
	}

	// wrap it in the default response wrapper
	wrappedResp := &dto.Response{StatusCode: status_code.GetStatusCode(err), Data: payload}

	// marshal response into bytes
	respBody, err := rpc.Marshal(wrappedResp)
	if err != nil {
		logs.Warn("unable to marshal payload for callback, err: %v", err)
		respBody, _ = rpc.Marshal(&dto.Response{
			StatusCode: status_code.Success,
			Data:       nil,
		})
	}

	// add the headers to the payload, over here, as per our use case, we assume that it will fit in a 512 buffer.
	fullPayload := c.addHeaders(respType, respBody)
	logs.Info("Response Callback: %v, sending to: %v", fullPayload, c.NotifiableClients[item].ToList())

	// We spawn max of 10 workers (limit resource usage) for a worker pool pattern to concurrently send the callback to users
	load := worker_pools.Load(func(job workerPoolJob) error {
		return net.SendData(job.Payload, job.Addr)
	},
		makeWorkerPoolJobs(c.NotifiableClients[item].ToList(), fullPayload),
		10,
	)

	// checks if there is error if any of them
	if ok := predicates.One(load, func(a error) bool { return a != nil }); ok {
		logs.Warn("Not all callbacks completed successfully, errs: %s", utils.DumpJSON(load))
		for _, err := range load {
			err := err
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// makeWorkerPoolJobs simply creates the worker pool jobs
func makeWorkerPoolJobs(addrs []string, payload []byte) []workerPoolJob {
	jobs := make([]workerPoolJob, len(addrs))
	for i, v := range addrs {
		jobs[i] = workerPoolJob{
			Payload: payload,
			Addr:    v,
		}
	}
	return jobs
}

// addHeaders adds the responseType, requestID, byteBufferArrayNo and totalByteBufferArray to header. For the sake of simplicity, we assumed
// that all callbacks will only use max of 1 byte array buffer (512 bytes - headers)
func (c *Client[T]) addHeaders(respType dto.ResponseType, body []byte) []byte {
	body = c.addTotalByteBufferArrayToHeader(body, 1)
	body = c.addByteBufferArrayNoToHeader(body, 1)
	body = c.addRequestID(body)
	body = c.addResponseTypeHeader(respType, body)
	return body
}

// addByteBufferArrayNoToHeader appends the byteBufferArray number to the header
func (c *Client[T]) addByteBufferArrayNoToHeader(response []byte, byteBufferArrayNo int64) []byte {
	return append(bytes.Int64ToBytes(byteBufferArrayNo), response...)
}

// addTotalByteBufferArrayToHeader appends the total byteBufferArray number to the header
func (c *Client[T]) addTotalByteBufferArrayToHeader(response []byte, totalByteBufferArrays int64) []byte {
	return append(bytes.Int64ToBytes(totalByteBufferArrays), response...)
}

// addResponseTypeHeader appends the responseType to the header
func (c *Client[T]) addResponseTypeHeader(respType dto.ResponseType, body []byte) []byte {
	return append([]byte{uint8(respType)}, body...)
}

// addRequestID appends the requestID to the header
func (c *Client[T]) addRequestID(body []byte) []byte {
	// this generates a random 9 character string guaranteed for uniqueness until 2050.
	reqID := shortid.MustGenerate()
	return append([]byte(reqID), body...)
}
