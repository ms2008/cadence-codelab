package courier

import (
	"errors"
	"fmt"
	"net/http"
)

func (h *CourierService) updateJob(w http.ResponseWriter, r *http.Request) {
	jobID := r.URL.Query().Get("id")
	job, ok := h.DeliveryQueue.Jobs[jobID]
	if !ok {
		http.Error(w, "Order not found: "+jobID, http.StatusNotFound)
		return
	}

	action := r.URL.Query().Get("action")
	if len(action) == 0 {
		http.Error(w, "No update action specified! "+action, http.StatusUnprocessableEntity)
		return
	}

	h.handleAction(r, job, action)
	fmt.Fprintf(w, "%s", job)
}

// handleAction takes the action corresponding to the specified action type
func (h *CourierService) handleAction(r *http.Request, job *DeliveryJob, action string) {
	switch action {
	case "accept":
		// driver accepted the trip, complete DispatchCourierActivity
		err := h.client.CompleteActivity(job.AcceptTaskToken, djAccepted, nil)
		if err != nil {
			job.Status = djRejected
		} else {
			job.Status = djAccepted
		}

	case "decline":
		// driver declined the trip, complete DispatchCourierActivity
		h.client.CompleteActivity(job.AcceptTaskToken, djRejected, errors.New("Order rejected"))
		job.Status = djRejected

	case "picked_up":
		// driver picked up from restaurant, complete PickUpOrderActivity
		err := h.client.CompleteActivity(job.PickupTaskToken, djPickedUp, nil)
		if err != nil {
			fmt.Printf("%s", err)
		}
		job.Status = djPickedUp
	case "completed":
		// driver delivered the food, complete the DeliverOrderActivity
		err := h.client.CompleteActivity(job.CompletTaskToken, djCompleted, nil)
		if err != nil {
			fmt.Printf("%s", err)
		}
		job.Status = djCompleted

	case "p_token":
		// record the task token for PickUpOrderActivity
		job.PickupTaskToken = []byte(r.URL.Query().Get("task_token"))

	case "c_token":
		// record the task token for DeliverOrderActivity
		job.CompletTaskToken = []byte(r.URL.Query().Get("task_token"))
	}
}
