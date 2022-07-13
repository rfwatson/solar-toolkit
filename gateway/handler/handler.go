package handler

import (
	"net/http"

	"git.netflux.io/rob/solar-toolkit/inverter"
)

type Store interface {
	InsertETRuntimeData(*inverter.ETRuntimeData) error
}

type Handler struct {
	store Store
}

func NewHandler(store Store) *Handler { return &Handler{store: store} }

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

}
