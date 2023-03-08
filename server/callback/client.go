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

type Client[T comparable] struct {
	NotifiableClients map[T]*collections.Set[string]
}

func NewClient[T comparable]() *Client[T] {
	return &Client[T]{
		NotifiableClients: make(map[T]*collections.Set[string]),
	}
}

type workerPoolJob struct {
	Payload []byte
	Addr    string
}

func (c *Client[T]) Subscribe(ctx context.Context, item T, expireDuration time.Duration) {
	addr := server.GetIPAddr(ctx)

	if _, ok := c.NotifiableClients[item]; !ok {
		c.NotifiableClients[item] = collections.NewSet[string]()
	}

	c.NotifiableClients[item].MustAdd(addr)
	logs.Info("Client: %s has successfully been subscribed to item: %s", addr, utils.DumpJSON(item))
	go c.cleanup(item, addr, expireDuration)
}

func (c *Client[T]) cleanup(item T, addr string, expireDuration time.Duration) {
	timer := time.NewTimer(expireDuration)

	<-timer.C
	logs.Info("removing address: %s for item: %s from subscription", addr, utils.DumpJSON(item))
	c.NotifiableClients[item].MustRemove(addr)
}

func (c *Client[T]) Notify(item T, respType dto.ResponseType, payload any, err error) error {
	if _, ok := c.NotifiableClients[item]; !ok {
		return nil
	}

	wrappedResp := &dto.Response{StatusCode: status_code.GetStatusCode(err), Data: payload}

	respBody, err := rpc.Marshal(wrappedResp)
	if err != nil {
		logs.Warn("unable to marshal payload for callback, err: %v", err)
		respBody, _ = rpc.Marshal(&dto.Response{
			StatusCode: status_code.Success,
			Data:       nil,
		})
	}

	fullPayload := c.addHeaders(respType, respBody)
	logs.Info("Response Callback: %v", fullPayload)
	load := worker_pools.Load(func(job workerPoolJob) error {
		return net.SendData(job.Payload, job.Addr)
	},
		makeWorkerPoolJobs(c.NotifiableClients[item].ToList(), fullPayload),
		10,
	)

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

func (c *Client[T]) addHeaders(respType dto.ResponseType, body []byte) []byte {
	body = c.addTotalPacketToHeader(body, 1)
	body = c.addPacketNoToHeader(body, 1)
	body = c.addRequestID(body)
	body = c.addResponseTypeHeader(respType, body)
	return body
}

func (c *Client[T]) addPacketNoToHeader(response []byte, packetNo int64) []byte {
	return append(bytes.Int64ToBytes(packetNo), response...)
}

func (c *Client[T]) addTotalPacketToHeader(response []byte, totalPacket int64) []byte {
	return append(bytes.Int64ToBytes(totalPacket), response...)
}

func (c *Client[T]) addResponseTypeHeader(respType dto.ResponseType, body []byte) []byte {
	return append([]byte{uint8(respType)}, body...)
}

func (c *Client[T]) addRequestID(body []byte) []byte {
	reqID := shortid.MustGenerate()
	return append([]byte(reqID), body...)
}
