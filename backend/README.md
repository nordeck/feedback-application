# Feedback Backend

## Key Features

This application implements the backend component of the feedback-application.
It provides a REST API on which feedback may be submitted.
An authorization mechanism is implemented in conjunction with the frontend, which ensures that feedback can only be
submitted by authorized [Matrix](https://matrix.org) users
through [Matrix UVS](https://github.com/matrix-org/matrix-user-verification-service/).

## How To Use

To clone and run this application, you'll need [Git](https://git-scm.com) as well as [Docker](https://docker.com/)
installed and configured on your computer.

1. Clone this repository
2. Create and run a postgres database
3. Build and run the image with Docker
    1. `cd backend`
    2. `docker build --tag=nordeck/feedback-app .`
    3. `docker run nordeck/feedback-app` with the fitting environment and port publishing parameters for your setup
4. Run Grafana with the [provided dashboard](../grafana) (optional)

## Configuration

In order to run this application, you need to prepare your environment.
You will need to set the following variables.

<div style="margin-left: auto;
            margin-right: auto;
            width: 70%">

| Environment variable name | Description                                                   | Example                      |
|---------------------------|---------------------------------------------------------------|------------------------------|
| DB_HOST                   | DB server's hostname                                          | localhost                    |
| DB_PORT                   | DB server's port                                              | 5432                         |
| DB_USER                   | DB server's username                                          | someUser                     |
| DB_PASSWORD               | DB user's password                                            | somePassphrase               |
| DB_NAME                   | Database name                                                 | someDatabase                 |
| SSL_MODE                  | Use SSL (enable or disable)                                   | disable                      |
| OIDC_VALIDATION_URL       | the URL of the MVS the OIDC Token has to be validated against | https://some.url/verify/user |
| JWT_SIGNATURE             | Some unique String the JWT will get signed with               | someArbitraryString          |
| MATRIX_SERVER_NAME        | The server name which the OIDC token is validated against     | domain.tld                   |

</div>

## Development

The database is versioned using the goose plugin for go.

## Credits

This software uses the following open source packages:

- [github.com/dariubs/gorm-jsonb](https://github.com/dariubs/gorm-jsonb) v0.1.5
- [github.com/gorilla/mux](https://github.com/gorilla/mux) v1.8.0
- [github.com/lib/pq](https://github.com/lib/pq) v1.10.7
- [github.com/pressly/goose/v3](https://github.com/pressly/goose/v3) v3.7.0
- [github.com/stretchr/testify](https://github.com/stretchr/testify) v1.8.1
- [github.com/testcontainers/testcontainers-go](https://github.com/estcontainers/testcontainers-go) v0.15.0
- [go.uber.org/zap](https://pkg.go.dev/go.uber.org/zap) v1.23.0
- [gorm.io/driver/postgres](https://pkg.go.dev/gorm.io/driver/postgres) v1.4.5
- [gorm.io/gorm](https://pkg.go.dev/gorm.io/gorm) v1.24.1-0.20221019064659-5dd2bb482755
