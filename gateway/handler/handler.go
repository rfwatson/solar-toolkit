package handler

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"git.netflux.io/rob/solar-toolkit/inverter"
)

const timestampMinimumYear = 2022

type Store interface {
	InsertDataFrame(*inverter.ETDataFrame) error
}

type Handler struct {
	store Store
}

func New(store Store) *Handler { return &Handler{store: store} }

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet && r.URL.Path == "/healthz" {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK\n"))
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if r.URL.Path != "/gateway/et_runtime_data" {
		http.Error(w, "endpoint not found", http.StatusNotFound)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("could not read body: %v", err)
		http.Error(w, "unexpected error", http.StatusInternalServerError)
		return
	}

	var runtimeData inverter.ETRuntimeData
	err = json.Unmarshal(body, &runtimeData)
	if err != nil {
		log.Printf("could not unmarshal runtime data body: %v", err)
		http.Error(w, "unexpected error", http.StatusInternalServerError)
		return
	}

	if runtimeData.Timestamp.Year() < timestampMinimumYear {
		log.Printf("invalid timestamp: %v", runtimeData.Timestamp)
		http.Error(w, "invalid data", http.StatusBadRequest)
		return
	}

	var meterData inverter.ETMeterData
	err = json.Unmarshal(body, &meterData)
	if err != nil {
		log.Printf("could not unmarshal meterData body: %v", err)
		http.Error(w, "unexpected error", http.StatusInternalServerError)
		return
	}

	dataFrame := inverter.ETDataFrame{ETRuntimeData: &runtimeData, ETMeterData: &meterData}
	if err = h.store.InsertDataFrame(&dataFrame); err != nil {
		log.Printf("error storing data: %v", err)
		http.Error(w, "unexpected error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK\n"))
}
