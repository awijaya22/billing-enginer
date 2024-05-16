package main

import (
	"billing-engine/loan"
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"
)

func main() {
	postgreConn := connectPostgre()
	defer postgreConn.Close()

	repo := loan.NewRepo(postgreConn)

	service := loan.NewService(repo)

	// create a loan
	startAt := time.Now()
	err := service.CreateLoan(context.Background(), "ANGEL", 5000000, &startAt)
	if err != nil {
		fmt.Println("Create Loan error = ", err)
	}

	// user make payment
	err = service.MakePayment(context.Background(), "ANGEL", 1, 550000)
	fmt.Println("Make Payment error = ", err)

	// user make wrong payment
	err = service.MakePayment(context.Background(), "ANGEL", 2, 500000)
	fmt.Println("user make wrong payment amount err = ", err)

	// user doble payment 
	err = service.MakePayment(context.Background(), "ANGEL", 1, 5500000)
	fmt.Println("user doble payment err =", err)

	// check outstanding
	outstanding, err := service.GetOutstanding(context.Background(), 3)
	if err != nil {
		fmt.Println("Get Outstanding error = ", err)
	}
	fmt.Println("outstanding = ", outstanding)

	// check is delinquent
	isDelinquent, err := service.IsDelinquent(context.Background(), "ANGEL")
	if err != nil {
		fmt.Println("Is Delinquent error = ", err)
	}
	fmt.Println("Is Delinquent = ", isDelinquent)
}

func connectPostgre() *sql.DB {
	// dbDsn := os.Getenv("DATABASE_URL")
	dbDsn := "postgres://fa-5327@localhost:5432/postgres?sslmode=disable"

	db, err := sql.Open("postgres", dbDsn)
	if err != nil {
		log.Fatalf("Error opening database: %v\n", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("Error pinging database: %v\n", err)
	}

	fmt.Println("Successfully connected to the database!")
	return db
}
