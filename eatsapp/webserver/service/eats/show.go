package eats

import (
	"net/http"
	"reflect"
	"runtime"
	"time"

	s "go.uber.org/cadence/.gen/go/shared"

	"github.com/venkat1109/cadence-codelab/eatsapp/webserver/service"
	"github.com/venkat1109/cadence-codelab/eatsapp/worker/workflow/eats"
)

func (h *EatsService) show(w http.ResponseWriter, r *http.Request) {
	orderID := r.URL.Query().Get("id")
	runID := r.URL.Query().Get("run_id")

	if len(orderID) == 0 || len(runID) == 0 {
		err := h.listOrders(w, r)
		if err != nil {
			return
		}
	} else {
		err := h.showOrder(w, r, orderID, runID)
		if err != nil {
			return
		}
	}
}

func (h *EatsService) showOrder(
	w http.ResponseWriter, r *http.Request, orderID string, runID string) error {
	data, err := h.processExecution(orderID, runID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	return service.ViewHandler(w, r, *data)
}

func (h *EatsService) processExecution(workflowID string, runID string) (*TaskGroup, error) {
	tf := NewTaskGroupExecution(h.client)
	return tf.Transform(workflowID, runID)
}

func (h *EatsService) listOrders(w http.ResponseWriter, r *http.Request) error {
	page := EatsOrderListPage{
		ShowOrderExistError: (r.URL.Query().Get("error") == "order_exist"),
	}

	var err error
	page.Orders, err = h.listOpenWorkflows()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	return service.ViewHandler(w, r, page)
}

func getWorkflowName(workflowFunc interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(workflowFunc).Pointer()).Name()
}

// listOpenWorkflows returns all the open eats order workflows
// created over the past ten hours
func (h *EatsService) listOpenWorkflows() (*s.ListOpenWorkflowExecutionsResponse, error) {

	// list all the open workflows in the past 10 hours
	startTime := time.Now().Add(-10 * time.Hour).UnixNano()
	latestTime := time.Now().Add(time.Minute).UnixNano()

	// convert the workflow function into a fully qualified name
	workflowName := getWorkflowName(eats.OrderWorkflow)

	// list all the eats order workflows during the past 10 hours
	req := s.ListOpenWorkflowExecutionsRequest{
		StartTimeFilter: &s.StartTimeFilter{
			EarliestTime: &startTime,
			LatestTime:   &latestTime,
		},
		TypeFilter: &s.WorkflowTypeFilter{
			Name: &workflowName,
		},
	}
	return h.client.ListOpenWorkflow(&req)
}
