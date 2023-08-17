package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/kirillApanasiuk/toll-calculator/types"
	"net/http"
	"strconv"
)

func main() {
	listenAddr := flag.String("listenaddr", ":3001", "the listenaddress of the HTTP server")
	store := NewMemoryStore()
	var (
		svc = NewInvoiceAggregatore(store)
	)
	svc = NewLogMiddleware(svc)
	makeHTTPTransport(*listenAddr, svc)
}

func handleAggregate(svc Aggregator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var distance types.Distance
		if err := json.NewDecoder(r.Body).Decode(&distance); err != nil {
			writeJson(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		if err := svc.AggregateDistance(distance); err != nil {
			writeJson(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		writeJson(w, http.StatusOK, map[string]string{"Success": "the result was processed"})
		return
	}
}

func makeHTTPTransport(listenAddr string, svc Aggregator) {
	fmt.Println("HTTP transport running on port", listenAddr)
	http.HandleFunc("/aggregate", handleAggregate(svc))
	http.HandleFunc("/invoice", handleGetInvoice(svc))
	http.ListenAndServe(listenAddr, nil)
}

func handleGetInvoice(svc Aggregator) func(rw http.ResponseWriter, req *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		values, ok := req.URL.Query()["obuId"]
		if !ok {
			writeJson(rw, http.StatusBadRequest, map[string]string{"error": "missing OBU ID"})
		}
		obuId, err := strconv.Atoi(values[0])
		if err != nil {
			writeJson(rw, http.StatusBadRequest, map[string]string{"error": "invalid obuId"})
		}

		invoice, err := svc.CalculateInvoice(obuId)
		if err != nil {
			writeJson(rw, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		writeJson(rw, http.StatusOK, map[string]any{"invoice": invoice})
	}
}

func writeJson(rw http.ResponseWriter, status int, v any) error {
	rw.WriteHeader(status)
	rw.Header().Add("Content-Type", "application/json")
	return json.NewEncoder(rw).Encode(v)
}
