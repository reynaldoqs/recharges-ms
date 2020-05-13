package service

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"

	"audiman/projects/recharger-x/entity"
	msgbroker "audiman/projects/recharger-x/messageBroker"
	"audiman/projects/recharger-x/repository"
)

// RechargeService interface for recharge service
type RechargeService interface {
	Validate(recharge *entity.Recharge) error
	Create(recharge *entity.Recharge) error
	List() ([]*entity.Recharge, error)
	//methods for message broker
	InitRechargerListener() error
	CloseListener()
}

type service struct{}

var repo repository.RechargeRepository
var msgbkr msgbroker.MessageBroker

// NewRechargeService creates a new recharge service
func NewRechargeService(
	repository repository.RechargeRepository,
	messageBroker msgbroker.MessageBroker,
) RechargeService {
	repo = repository
	msgbkr = messageBroker
	return &service{}
}

func (*service) Validate(recharge *entity.Recharge) error {
	v := validator.New()
	err := v.Struct(recharge)
	if err != nil {

		var valErrs string
		for _, err := range err.(validator.ValidationErrors) {
			valErrs += fmt.Sprintf("error: %v isn't aceptable value for %v \n", err.Value(), err.Field())
		}

		return fmt.Errorf(valErrs)
	}
	return nil
}

func (*service) Create(recharge *entity.Recharge) error {
	recharge.CreatedAt = time.Now().UTC().Unix()

	result, err := repo.Save(recharge)
	if err != nil {
		return errors.Wrap(err, "service.rechargeService.Create")
	}

	// format to common data type and send to queue
	err = addtoQueue(result)
	if err != nil {
		return errors.Wrap(err, "service.rechargeService.Create")
	}

	return nil
}

func (*service) List() ([]*entity.Recharge, error) {
	res, err := repo.FindAll()
	if err != nil {
		return nil, errors.Wrap(err, "service.rechargeService.List")
	}
	return res, nil
}

func addtoQueue(recharge *entity.Recharge) error {

	rmsg := entity.RechargeMsg{
		IDRecharge:  recharge.ID,
		PhoneNumber: recharge.PhoneNumber,
		Company:     recharge.Company,
		CardNumber:  recharge.CardNumber,
		CreatedAt:   recharge.CreatedAt,
	}

	body, err := json.Marshal(rmsg)
	if err != nil {
		return err
	}
	err = msgbkr.PublishOnQueue(body, "addRecharge")
	if err != nil {
		return err
	}
	return nil
}

// listener for recharger microservice
func (*service) InitRechargerListener() error {
	err := msgbkr.Subscribe("rechargeFromRechargerService", handleMessages)
	if err != nil {
		errors.Wrap(err, "services.rechargeService.InitRechargesListener")
	}
	return nil
}

func (*service) CloseListener() {
	msgbkr.Close()
}

//recibe back resolved recharges from recharger microservice
func handleMessages(data []byte) {
	fmt.Println("---------recive data back from microservice -------------------")
	fmt.Println(string(data))
	fmt.Println("---------end data back from microservice -------------------")
	rmsg := entity.RechargeMsg{}
	err := json.Unmarshal(data, &rmsg)
	if err != nil {
		log.Printf("%s : service.rechargeService.handleMessages", err.Error())
	}
	//update data base with news dates
	repo.UpdateRecharger(&rmsg)
	if err != nil {
		log.Printf("%s : service.rechargeService.handleMessages", err.Error())
	}
}
