package main

import (
	"github.com/authorizer/internal/core/service"
	"github.com/authorizer/internal/driven/database"
	"github.com/authorizer/internal/driven/repository"
	"github.com/authorizer/internal/driver/cli"
	"log"
	"os"
)

func main() {
	db := database.NewInMemoryDB()

	accountRepo := repository.NewAccountRepository(db)

	as := service.NewAccount(accountRepo)
	ts := service.NewTransaction(accountRepo)

	handler := cli.NewHandler(as, ts)

	log.SetOutput(os.Stdout)

	err := handler.Handle(os.Stdin, os.Stdout)
	if err != nil {
		log.Fatal(err)
	}
}
