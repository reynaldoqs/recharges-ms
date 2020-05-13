package main

import (
	"log"
	"time"

	"audiman/projects/recharger-x/controller"
	router "audiman/projects/recharger-x/http"
	msgbroker "audiman/projects/recharger-x/messageBroker"
	"audiman/projects/recharger-x/repository"
	"audiman/projects/recharger-x/service"
)

func main() {

	rrepo := setupMongoRepo()
	mb := setupRabbitMQ()
	rs := service.NewRechargeService(rrepo, mb)
	go rs.InitRechargerListener()
	rc := controller.NewRechargeController(rs)
	ht := router.NewChiRouter()

	const port string = ":8000"

	ht.POST("/posts", rc.AddRecharge)
	ht.GET("/posts", rc.GetRecharges)
	ht.SERVE(port)
	defer rs.CloseListener()

}

func setupMongoRepo() repository.RechargeRepository {
	//mongoURL := os.Getenv("MONGO_URL")
	//mongodb := os.Getenv("MONGO_DB")
	//mongoTimeout, _ := strconv.Atoi(os.Getenv("MONGO_TIMEOUT"))
	mongoURL := "mongodb://localhost:27017"
	mongodb := "recharger"
	mongoTimeout := 30
	repo, err := repository.NewMongoRepository(mongoURL, mongodb, mongoTimeout)
	if err != nil {
		log.Fatal(err)
	}
	return repo
}
func setupRabbitMQ() msgbroker.MessageBroker {
	//rabbitURL := os.Getenv("RABBIT_URL")
	rabbitURL := "amqp://guest:guest@localhost:5672/"
	brk, err := msgbroker.NewRabbitMqBroker(rabbitURL)
	if err != nil {
		time.Sleep(time.Second * 2)
		log.Printf("Trying to connect: %s\n", rabbitURL)
		return setupRabbitMQ()
	}
	return brk
}
