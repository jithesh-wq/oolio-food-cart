package handler

import (
	"io"


	"net/http"
	"strings"

	"github.com/jithesh-wq/oolio-food-cart/logger"
	"github.com/jithesh-wq/oolio-food-cart/service"
	"github.com/jithesh-wq/oolio-food-cart/store"
)

const (
	adminAuthKey string = "92563ae6-b76c-4c87-ac12-544995c03523"
	customer     string = "CUSTOMER"
	admin        string = "ADMIN"
)

type APIHandler struct {
	service service.IService
	store   *store.MemoryStore
}

func CreateApiHandler(s service.IService, store *store.MemoryStore) *APIHandler {
	return &APIHandler{service: s, store: store}
}
func (h *APIHandler) HandleRequest(w http.ResponseWriter, r *http.Request) {
	var body []byte
	var err error
	logger.Log.Infoln("Received request:", r.Method, r.RequestURI)
	status, user := ValidateRequest(r, h.store)
	if !status {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	

	if strings.Contains(r.RequestURI, "/generate-table-session") {
		tableId := r.URL.Query().Get("tableId")
		body = []byte(tableId)
	} else if strings.Contains(r.RequestURI, "/products") {
		productId := r.URL.Query().Get("id")
		body = []byte(productId)
	} else {
		body, err = io.ReadAll(r.Body)
		if err != nil {
			logger.Log.Errorln(err.Error())
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	// Decode and validate the request
	req, err := h.service.DecodeAndValidate(body, user, r.Header.Get("oolio-auth-key"))
	if err != nil {
		logger.Log.Errorln(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Process the request
	response, err := h.service.ProcessRequest(req)
	if err != nil {
		if err.Error() == "notfound" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		logger.Log.Errorln(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(response)
	if err != nil {
		logger.Log.Infoln("Failed to write response:", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	logger.Log.Infoln("Response sent successfully")
}

func ValidateRequest(r *http.Request, store *store.MemoryStore) (bool, string) {
	if strings.Contains(r.RequestURI, "/generate-table-session") {
		return true, customer
	} else {
		key := r.Header.Get("oolio-auth-key")
		logger.Log.Infoln("Authorization header:", key)
		if key == "" {
			logger.Log.Infoln("Unauthorized request")
			return false, ""
		} else {
			//validate api key with the one in memory store also validate the case of admin access the api to do admin operations
			if key == adminAuthKey {
				if strings.Contains(r.RequestURI, "add-products") || strings.Contains(r.RequestURI, "checkout-table-session") ||
					strings.Contains(r.RequestURI, "orders") {
					logger.Log.Infoln("admin API key is valid")
					return true, admin

				}
				return false, ""
			} else if store.ValidateSession(key) {
				logger.Log.Infoln("Authorization header is valid")
				return true, customer
			}
		}
		logger.Log.Infoln("Unauthorized request")
		return false, ""
	}
}
