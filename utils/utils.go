package utils

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/JAYENDRA06/apiproject/types"
	"github.com/go-playground/validator/v10"
)

var Validate = validator.New()

type allowedTypes interface {
	types.RegisterUserPayload | types.LoginUserPayload | types.ProductItem | types.CartCheckOutPayload
}

func ParseJSON[T allowedTypes](r *http.Request, payload *T) error {
	if r.Body == nil {
		return fmt.Errorf("payload is empty")
	}
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		return err
	}
	return nil
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

func WriteError(w http.ResponseWriter, status int, err error) {
	WriteJSON(w, status, map[string]string{"error": err.Error()})
}
