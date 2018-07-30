package api

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func GetDomainsInfo() httprouter.Handle {
	is_update, read_data := updateAsyncDB()
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		w.Header().Set("Content-Type", "text/plain")
		is_update <- false
		fmt.Fprintf(w, <-read_data)
	}
}
