# Latipe Promotion Service

## Introduction

This is a service that manages discount business include create, retrieve, apply voucher for Latipe's users . It is a
RESTful API and gRPC Application that provides the following features:

- Create a voucher (store/admin)
- Retrieve a voucher, list all vouchers of store/admin
- Apply a voucher to a user's order
- Retrieve all vouchers of a user can use
- Provide a gRPC service for the above features

## Technologies

- Golang (1.20)
- gRPC (v1.62)
- RESTful API (Fiber v2)
- MongoDB
- RabbitMQ
- Docker
- Prometheus

## Features

- API host: `localhost:5010`
- gRPC host: `localhost:6010`
- BasicAuth (username: `admin`, password: `123123`)
- Metrics (Prometheus): `localhost:5010/metrics`
- Fiberdash: `localhost:5010/fiber/dashboard`
- Health check: `localhost:5010/health`
- Readiness check: `localhost:5010/readiness`
- Liveness check: `localhost:5010/liveness`
- Swagger: `localhost:5010/swagger/index.html`

## Installation

- Change config file in `config/config.yml` to your own configuration
- Use Makefile to build and run the application or use docker file
- Run `make setup` to install all dependencies
- Run `make buildw` to build the application for windows (.exe)
- Run `make buildl` to build the application for linux

<hr>
<h4>Development by tdat.it</h4>