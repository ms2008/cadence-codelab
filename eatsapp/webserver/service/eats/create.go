package eats

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"go.uber.org/cadence"

	"github.com/venkat1109/cadence-codelab/eatsapp/worker/workflow/eats"
)

// create creates a new eats order
func (h *EatsService) create(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	items := r.Form["item-id"]

	if len(items) == 0 {
		http.Error(w, "Order constains no items!", http.StatusUnprocessableEntity)
		return
	}

	execution, err := h.startOrderWorkflow(items)
	if err != nil {
		if strings.HasPrefix(err.Error(), "WorkflowExecutionAlreadyStartedError") {
			http.Redirect(w, r, "/eats-orders?error=order_exist", http.StatusFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	url := fmt.Sprintf("/eats-orders?id=%s&run_id=%s&page=eats-order-status", execution.ID, execution.RunID)
	http.Redirect(w, r, url, http.StatusFound)
}

// startOrderWorkflow starts the eats order workflow
func (h *EatsService) startOrderWorkflow(items []string) (*cadence.WorkflowExecution, error) {

	orderID := fmt.Sprintf("EO-USR-JOE-%v", time.Now().Unix())

	workflowOptions := cadence.StartWorkflowOptions{
		ID:                              orderID,
		TaskList:                        cadenceTaskList,
		ExecutionStartToCloseTimeout:    10 * time.Minute,
		DecisionTaskStartToCloseTimeout: 10 * time.Minute,
	}

	return h.client.StartWorkflow(workflowOptions, eats.OrderWorkflow, orderID, items)
}
