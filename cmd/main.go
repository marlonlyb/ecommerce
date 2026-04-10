package main

import (
	"log"

	"github.com/mlbautomation/ProyectoEMLB/application"
	"github.com/mlbautomation/ProyectoEMLB/domain/services"
	"github.com/mlbautomation/ProyectoEMLB/infrastructure/handlers"
	paypaladapter "github.com/mlbautomation/ProyectoEMLB/infrastructure/paypal"
	"github.com/mlbautomation/ProyectoEMLB/infrastructure/postgres"
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

	lService := services.NewLogin(uService)
	lHandlers := handlers.NewLogin(lService)

	iRepository := postgres.NewInvoice(dbPool)
	irRepository := postgres.NewInvoiceReport(dbPool)
	iService := services.NewInvoice(iRepository, irRepository, poService)
	iHandlers := handlers.NewInvoice(iService)

	ppVerifier := paypaladapter.NewVerifier(nil)
	ppProcessor := application.NewPaymentFlow(ppVerifier, poService, iService)
	ppHandlers := handlers.NewPaypal(ppProcessor)

	httpServer := NewServer(uHandlers, pHandlers, poHandlers, lHandlers, ppHandlers, iHandlers)
	httpServer.Initialize()

}
