package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"audiman/projects/recharger-x/entity"
	"audiman/projects/recharger-x/errors"
	"audiman/projects/recharger-x/service"
)

type controller struct{}

var rechargeService service.RechargeService

// RechargeController creates type for recharges controller
type RechargeController interface {
	GetRecharges(res http.ResponseWriter, req *http.Request)
	AddRecharge(res http.ResponseWriter, req *http.Request)
}

// NewRechargeController creates a controller for recharges
func NewRechargeController(s service.RechargeService) RechargeController {
	rechargeService = s
	return &controller{}
}

func (*controller) GetRecharges(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	recharges, err := rechargeService.List()
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(res).Encode(errors.ServiceError{Message: "Error getting the recharges"})
	}
	res.WriteHeader(http.StatusOK)
	json.NewEncoder(res).Encode(recharges)
}

func (*controller) AddRecharge(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	var recharge entity.Recharge
	err := json.NewDecoder(req.Body).Decode(&recharge)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(res).Encode(errors.ServiceError{Message: "Error unmarshalling data"})
		return
	}
	err = rechargeService.Validate(&recharge)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(res).Encode(errors.ServiceError{Message: err.Error()})
		return
	}
	err = rechargeService.Create(&recharge)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err)
		json.NewEncoder(res).Encode(errors.ServiceError{Message: "Error saving the recharge"})
		return
	}
	res.WriteHeader(http.StatusNoContent)
}
