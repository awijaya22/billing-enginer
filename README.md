# BILLING ENGINE

## Pre req 
- postgres
- go

## How to run
- prep database
  ```
  psql <insert ur postgres connection> -f database/user-loan.sql
  psql <insert ur postgres connection> -f database/user-billing.sql
  psql <insert ur postgres connection> -f database/index.sql
  ```
- run 
    ```bash
    go run main.go
    ```