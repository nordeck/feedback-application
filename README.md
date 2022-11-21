<h1 align="center">
  <br>
  <a href="https://nordeck.net/"><img src="https://nordeck.net/wp-content/uploads/2020/05/NIC_logo_Nordeck-300x101.png" alt="Markdownify" width="300"></a>
  <br>
  Feedback-Application
  <br>
</h1>
<h4 align="center">A feedback tool for <a href="https://jitsi.org/" target="_blank">Jitsi</a>.</h4>


### Key Features

This application is meant to run in parallel with a jitsi instance to manage, view and persist feedback rating data.
The data is gathered from jitsi's internal feedback dialogue and persisted into its own postgres database. One is able
to
visualize the data with grafana.

### How To Use

To clone and run this application, you'll need [Git](https://git-scm.com) as well as [Docker](https://docker.com/)
installed and
configured on your computer. 

1. Clone the repository
2. Create and run a postgres database
3. build and run the feedback-app image
4. run grafana (optional)

### Download

You can [TODO-download](https://nordeck.atlassian.net/browse/NEW-487) the latest installable version of our
feedback-application.

### Configuration

In order to run this application, you need to prepare your environment. You will need to set the following variables.
<div style="margin-left: auto;
            margin-right: auto;
            width: 70%">

| Environment variable name | Description                 | Example        |
|---------------------------|-----------------------------|----------------|
| DB_HOST                   | DB server's hostname        | localhost      |
| DB_PORT                   | DB server's port            | 5432           |
| DB_USER                   | DB server's username        | someUser       |
| DB_PASSWORD               | DB user's password          | somePassphrase |
| DB_NAME                   | Database name               | someDatabase   |
| SSL_MODE                  | Use SSL (enable or disable) | disable        |

</div>

### Grafana

To visualize the feedback data, we created a simple dashboard for grafana. You can import it from
/feedback-application/grafana/
and customize it to your needs.

You will need to configure your datasource beforehand. Use your grepped IP and your prepared database configuration.

### Development

The database is versioned using the goose plugin for go.

### Credits

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

### How to Contribute

Please take a look at our [Contribution Guidelines](https://github.com/nordeck/.github/blob/main/docs/CONTRIBUTING.md).


### Related

[Jitsi Conference Mapper](https://github.com/nordeck/Jitsi-Conference-Mapper) - Jitsi Conference Mapper

### You may also like...

[Matrix Widget Toolkit](https://github.com/nordeck/matrix-widget-toolkit) - A widget toolkit
for [Matrix](https://matrix.org/)

[Matrix Poll](https://github.com/nordeck/matrix-poll) - A poll widget for [Matrix](https://matrix.org/)

### License

- [Apache 2.0](TODO GITHUB LINK ZUR LICENSE)

### Sponsors

<p align="center">
   <a href="https://www.dataport.de/"><img src="./.docs/logos/dataportlogo.png" alt="Dataport" width="20%"></a>
   <a href="https://nordeck.net/"><img src="https://nordeck.net/wp-content/uploads/2020/05/NIC_logo_Nordeck-300x101.png" alt="Markdownify" width="300"></a>
</p>

> [nordeck.net](https://nordeck.net/) &nbsp;&middot;&nbsp;
> GitHub [Nordeck](https://github.com/nordeck/) 
