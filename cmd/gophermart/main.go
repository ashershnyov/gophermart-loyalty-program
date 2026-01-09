package main

import (
	"flag"
	"log"

	"github.com/ashershnyov/gophermart-loyalty-program/internal/server"
	"github.com/ashershnyov/gophermart-loyalty-program/internal/server/config"
	"github.com/caarlos0/env"
)

type params struct {
	Address        string `env:"RUN_ADDRESS" envDefault:""`
	DBAddress      string `env:"DATABASE_URI" envDefault:""`
	AccrualAddress string `env:"ACCRUAL_SYSTEM_ADDRESS" envDefault:""`
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

	if params.Address == "" {
		params.Address = *address
	}

	if params.DBAddress == "" {
		params.DBAddress = *dbAddress
	}

	if params.AccrualAddress == "" {
		params.AccrualAddress = *accrualAddress
	}

	return params
}

func main() {
	params := parseStartingParams()

	srv, err := server.New(
		config.SetAddress(params.Address),
		config.SetDBAddress(params.DBAddress),
	)

	if err != nil {
		log.Fatalf("an error when creating the server: %v", err)
	}

	if err := srv.Run(); err != nil {
		log.Fatal(err)
	}
}
