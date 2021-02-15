package worker

import (
	"github.com/golang/glog"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/util/workqueue"
	"mlcp/pkg/cache"
	"mlcp/pkg/common"
	"mlcp/pkg/config"
	"mlcp/pkg/database"
	"time"
)

type WorkeQueue struct {
	queue       workqueue.RateLimitingInterface
	workerCount int
	c           cache.Cache
	db          *database.Database
}

func NewWorkQueue(workerCount int, c cache.Cache, db *database.Database) *WorkeQueue {
	return &WorkeQueue{
		queue:       workqueue.NewRateLimitingQueue(workqueue.DefaultItemBasedRateLimiter()),
		workerCount: workerCount,
		c:           c,
		db:          db,
	}
}

func (wq *WorkeQueue) Add(item interface{}) {
	wq.queue.AddRateLimited(item)
}

func (wq *WorkeQueue) Run(stopCh <-chan struct{}) error {
	defer runtime.HandleCrash()
	defer wq.queue.ShutDown()

	for i := 0; i < wq.workerCount; i++ {
		go wait.Until(wq.worker, time.Second, stopCh)
	}
	<-stopCh
	return nil
}

func (wq *WorkeQueue) worker() {
	for wq.processNextItem() {

	}
}

func (wq *WorkeQueue) processNextItem() bool {
	msg, quit := wq.queue.Get()
	if quit {
		return false
	}
	defer wq.queue.Done(msg)
	var m *common.Slot
	var ok bool
	if m, ok = msg.(*common.Slot); !ok {
		glog.Errorf("Unrecognisable message type")
		return false
	}
	err := wq.processMessage(m)
	if err != nil {
		glog.Errorf("Error handling event: %v: %v", msg, err)
		return false
	}
	return true
}

func (wq *WorkeQueue) processMessage(slot *common.Slot) error {
	return wq.assignSlot(slot)
}

func (wq *WorkeQueue) assignSlot(slot *common.Slot) error {
	s := wq.c.AssignSlot(slot)
	if s == nil {
		glog.Errorf("Error assigning slot: No free slot available.")
	}
	wq.db.Write(
		database.SlotData{
			SlotType: slot.Type,
			SlotId:   slot.SlotId,
		},
		database.VehicleData{
			RegnNo:  slot.GetRegNo(),
			Type:    config.VehicleTypeCar,
			InTime:  slot.InTime,
			OutTime: slot.OutTime,
		})
	return nil
}
