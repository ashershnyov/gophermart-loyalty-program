package main

import (
	"flag"
	"log"

	"github.com/ashershnyov/gophermart-loyalty-program/internal/server"
	"github.com/ashershnyov/gophermart-loyalty-program/internal/server/config"
	"github.com/caarlos0/env"
)

type params struct {
	address        string `env:"RUN_ADDRESS" envDefault:""`
	dbAddress      string `env:"DATABASE_URI" envDefault:""`
	accrualAddress string `env:"ACCRUAL_SYSTEM_ADDRESS" envDefault:""`
}

func parseStartingParams() params {
	var params params
	err := env.Parse(&params)
	if err != nil {
		log.Fatal(err)
	}

	address := flag.String("a", "localhost:8080", "Specifies the address for the server to start on. Can be overriden with RUN_ADDRESS env variable.")
	dbAddress := flag.String("d", "", "Specifies the address of the DB. Can be overriden with DATABASE_URI env variable.")
	accrualAddress := flag.String("r", "", "Specifies the address of the accrual service. Can be overriden with ACCRUAL_SYSTEM_ADDRESS env variable.")
	flag.Parse()

	if params.address == "" {
		params.address = *address
	}

	if params.dbAddress == "" {
		params.dbAddress = *dbAddress
	}

	if params.accrualAddress == "" {
		params.accrualAddress = *accrualAddress
	}

	return params
}

func main() {
	params := parseStartingParams()

	srv, err := server.New(
		config.SetAddress(params.address),
		config.SetDBAddress(params.dbAddress),
	)

	if err != nil {
		log.Fatalf("an error when creating the server: %v", err)
	}

	if err := srv.Run(); err != nil {
		log.Fatal(err)
	}
}
