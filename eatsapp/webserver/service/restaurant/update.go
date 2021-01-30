package restaurant

import (
	"errors"
	"fmt"
	"net/http"
)

func (h *RestaurantService) updateOrder(w http.ResponseWriter, r *http.Request) {

	orderID := r.URL.Query().Get("id")
	order, ok := h.state.Orders[orderID]
	if !ok {
		http.Error(w, "Order not found: "+orderID, http.StatusNotFound)
		return
	}

	action := r.URL.Query().Get("action")
	if len(action) == 0 {
		http.Error(w, "No update action specified! "+action, http.StatusUnprocessableEntity)
		return
	}

	h.handleAction(r, order, action)
	fmt.Fprintf(w, "%+v", order)
}

func (h *RestaurantService) handleAction(r *http.Request, order *Order, action string) {
	switch action {
	case "accept":
		// waiter accepted, complete PlaceOrderActivity
		err := h.client.CompleteActivity(order.TaskToken, "ACCEPTED", nil)
		if err != nil {
			order.Status = OSRejected
		} else {
			order.Status = OSPreparing
		}

	case "decline":
		// waiter declined, fail PlaceOrderActivity
		h.client.CompleteActivity(order.TaskToken, "REJECTED", errors.New("Order rejected"))
		order.Status = OSRejected

	case "ready":
		// food is ready, send a signal to the eats.OrderWorkflow
		err := h.client.SignalWorkflow(order.ReadySignal.WorkflowID, order.ReadySignal.RunID, order.ID, "ORDER_READY")
		if err != nil {
			fmt.Printf("%s", err)
		}
		order.Status = OSReady

	case "sent":
		// Courier picked up the food, send a signal to
		// to the courier workflow
		err := h.client.SignalWorkflow(order.PickUpSignal.WorkflowID, order.PickUpSignal.RunID, order.ID, "ORDER_PICKED_UP")
		if err != nil {
			fmt.Printf("%s", err)
		}
		order.Status = OSSent

	case "p_sig":
		// Courier out for pick up, record context
		// for sending signal later
		order.PickUpSignal = getSignalParams(r)
	}
}

func getSignalParams(r *http.Request) *SignalParam {
	return &SignalParam{
		WorkflowID: r.URL.Query().Get("workflow_id"),
		RunID:      r.URL.Query().Get("run_id"),
	}
}
