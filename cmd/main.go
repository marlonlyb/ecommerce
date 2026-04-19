package main

import (
	"log"

	"github.com/mlbautomation/Ecommmerce_MLB/application"
	"github.com/mlbautomation/Ecommmerce_MLB/domain/services"
	"github.com/mlbautomation/Ecommmerce_MLB/infrastructure/handlers"
	paypaladapter "github.com/mlbautomation/Ecommmerce_MLB/infrastructure/paypal"
	"github.com/mlbautomation/Ecommmerce_MLB/infrastructure/postgres"
)

func main() {

	err := loadEnv()
	if err != nil {
		log.Fatal(err)
	}

	err = validateEnvironments()
	if err != nil {
		log.Fatal(err)
	}

	dbPool, err := NewDBConnection()
	if err != nil {
		log.Fatal(err)
	}

	uRepository := postgres.NewUser(dbPool)
	uService := services.NewUser(uRepository)
	uHandlers := handlers.NewUser(uService)

	pRepository := postgres.NewProduct(dbPool)
	pService := services.NewProduct(pRepository)
	pHandlers := handlers.NewProduct(pService)

	poRepository := postgres.NewPurchaseOrder(dbPool)
	poService := services.NewPurchaseOrder(poRepository, pRepository)
	poHandlers := handlers.NewPurchaseOrder(poService)

	oRepository := postgres.NewOrder(dbPool)
	ppOrdersClient := paypaladapter.NewOrdersClient(nil)
	oService := services.NewOrder(oRepository, pRepository, ppOrdersClient)
	oHandlers := handlers.NewOrder(oService)

	lService := services.NewLogin(uService)
	lHandlers := handlers.NewLogin(lService)

	iRepository := postgres.NewInvoice(dbPool)
	irRepository := postgres.NewInvoiceReport(dbPool)
	iService := services.NewInvoice(iRepository, irRepository, poService)
	iHandlers := handlers.NewInvoice(iService)

	ppVerifier := paypaladapter.NewVerifier(nil)
	ppProcessor := application.NewPaymentFlow(ppVerifier, poService, iService)
	ppHandlers := handlers.NewPaypal(ppProcessor)

	httpServer := NewServer(uHandlers, pHandlers, poHandlers, oHandlers, lHandlers, ppHandlers, iHandlers)
	httpServer.Initialize()

}
