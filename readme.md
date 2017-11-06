# Dabanshan-GoLang Project

## dependency

* consul as discover service
* zipkin as trace service
* use grpc and protobuf
* use mongodb as database
* some go libs

## services

* svcs/product 
* svcs/user
* svcs/order

## how debug ?

* use "go get -v *" install go libs
* compile proto at "pb/Makefile"
* download and launch consul as default discover service.
* "go run cmd/productsvc/main.go" for launch product service
* "go run cmd/usersvc/main.go" for launch user service
* "go run cmd/gateway/main.go" fro launch gateway api

## debug example

* GET "http://localhost:8000/api/v1/products?userid=233&size=10"
* GET "http://localhost:8000/api/v1/users/59f05169668b9bcc7d442355"
