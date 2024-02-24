# Latipe Promotion Service

### Introduction

This is a service that manages discount business include create, retrieve, apply voucher for Latipe's users . It is a
RESTful API and gRPC Application that provides the following features:

- Create a voucher (store/admin)
- Retrieve a voucher, list all vouchers of store/admin
- Apply a voucher to a user's order
- Retrieve all vouchers of a user can use
- Provide a gRPC service for the above features

### Technologies

> - Golang (1.20)
> - gRPC (v1.62)
> - RESTful API (Fiber v2)
> - gRPC v1.62
> - MongoDB
> - RabbitMQ
> - Prometheus
> - Docker
### API Documentation

- **Base URL:** http://localhost:5010
- **Swagger:** [API Documentation](http://localhost:5010/swagger/index.html)
- **API Endpoints:**
    - **Admin Endpoints:**
        - `GET /api/v1/vouchers/admin`: Get all vouchers of store/admin.
        - `POST /api/v1/vouchers/admin`: Create a new voucher for administrators.
        - `GET /api/v1/vouchers/admin/:id`: Get a voucher by ID for administrators.
        - `GET /api/v1/vouchers/admin/code/:code`: Get a voucher by code for administrators.
        - `PATCH /api/v1/vouchers/admin/code/:code`: Update the status of a voucher for administrators.
    - **User Endpoints**
        - `GET /api/v1/vouchers/user/foryou`: Get vouchers available for the user.
        - `GET /api/v1/vouchers/user/code/:code`: Get a voucher by code for the user.

    - **Store Endpoints**
        - `GET /api/v1/vouchers/store`: Get all vouchers for the store.
        - `POST /api/v1/vouchers/store`: Create a new voucher for the store.
        - `GET /api/v1/vouchers/store/code/:code`: Get a voucher by code for the store.
        - `PATCH /api/v1/vouchers/store/cancel`: Cancel the status of a voucher for the store.
    - **Additional Endpoint**
        - `POST /api/v1/vouchers/checking` Check the validity of a voucher.
- **Metrics:**
    - `GET /metrics`: Get metrics data for monitoring the application, you can also use prometheus.
    - `GET /health`: Get health check data for monitoring the application.
    - `GET /readiness`: Get ready check.
    - `GET /liveness`: Get live check.
    - `GET /fiber/dashboard`: Get fiber dashboard for monitoring the application.
- **gRPC Service:**
    - `ApplyVoucherToPurchase`: Apply a voucher to a user's order.
    - `CheckUsingVouchersForCheckout`: Get voucher data for purchase checkout.

### Installation

- Change config file in `config/config.yml` to your own configuration
- Use Makefile to build and run the application
    - Run `make setup` to install all dependencies
    - Run `make buildw` to build the application for windows (.exe)
    - Run `make buildl` to build the application for linux
    - Run `make runw` to run the application for windows (.exe) or run the binary file in /build folder
- You can also use docker file to build and run the application:
  ```bash
    docker build -t latipe-promotion-service .
    docker run -p 5010:5010 latipe-promotion-service
  ```

<hr>
<h4>Development by tdat.it</h4>