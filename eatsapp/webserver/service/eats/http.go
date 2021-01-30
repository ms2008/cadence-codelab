package eats

import (
	"net/http"

	"go.uber.org/cadence"
	s "go.uber.org/cadence/.gen/go/shared"

	common "github.com/venkat1109/cadence-codelab/eatsapp/webserver/service"
)

type (
	// EatsService implements the handler for requests sent
	// to the Eats http service
	EatsService struct {
		menu   *common.Menu
		client cadence.Client
	}

	// EatsOrderListPage models the data to be displayed in response to
	// GET requests to the Eats service.
	EatsOrderListPage struct {
		ShowOrderExistError bool
		Orders              *s.ListOpenWorkflowExecutionsResponse
	}
)

const (
	cadenceTaskList = "cadence-bistro"
)

// NewService returns a new EatsService instance
func NewService(c cadence.Client, menu *common.Menu) *EatsService {
	return &EatsService{
		client: c,
		menu:   menu,
	}
}

func (h *EatsService) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		h.show(w, r)
	case "POST":
		h.create(w, r)
	default:
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
}
