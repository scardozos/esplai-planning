#!/bin/bash
export GRPC_SERVER_CLIENT_ADDR_SELF="0.0.0.0"
export GRPC_SERVER_CLIENT_PORT_SELF=9000
go run -race main.go