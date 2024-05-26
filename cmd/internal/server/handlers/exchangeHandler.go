package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"

	"gitlab.com/llcmediatel/recruiting/golang-junior-dev/cmd/internal/server/exchange"
)

// ExchangeHandler(w http.ResponseWriter, r *http.Request)  handles the exchange function from processor via exchange package
func ExchangeHandler(w http.ResponseWriter, r *http.Request) {
	var input exchange.Input

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	output, err := exchange.PresentResults(input)
	if err != nil {
		logrus.Error("processing error by timeout, your amount is too big for processing via 10 sec: ", err)
		http.Error(w, "processing took too long", http.StatusRequestTimeout)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(output); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	logrus.Info("request processed successfuly")
}
