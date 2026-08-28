package app

import (
	"micro-warehouse/transaction-service/configs"
	"micro-warehouse/transaction-service/controller"
	"micro-warehouse/transaction-service/database"
	"micro-warehouse/transaction-service/pkg/httpclient"
	"micro-warehouse/transaction-service/pkg/midtrans"
	"micro-warehouse/transaction-service/pkg/rabbitmq"
	"micro-warehouse/transaction-service/repository"
	"micro-warehouse/transaction-service/usecase"

	"github.com/gofiber/fiber/v2/log"
)

type Container struct {
	TransactionController controller.TransactionControllerInterface
}

func BuildContainer() *Container {
	cfg := configs.NewConfig()

	db, err := database.ConnectionPostgres(*cfg)
	if err != nil {
		log.Errorf("failed to connect database: %v", err)
	}

	transactionRepo := repository.NewTransactionRepository(db.DB)

	// http client
	merchantClient := httpclient.NewMerchantClient(*cfg)
	userClient := httpclient.NewUserClient(*cfg)
	productClient := httpclient.NewProductClient(*cfg)

	// RabbitMQ Client
	rabbitMQService, err := rabbitmq.NewRabbitMQService(cfg.RabbitMQ.URL())
	if err != nil {
		log.Fatalf("failed to connect RabbitMQ: %v", err)
	}

	transactionUsecase := usecase.NewTransactionUsecase(transactionRepo, merchantClient, rabbitMQService, productClient, userClient)

	midtransSerice := midtrans.NewMidtransService(*cfg)
	transactionController := controller.NewTransactionController(transactionUsecase, midtransSerice)

	return &Container{
		TransactionController: transactionController,
	}
}
