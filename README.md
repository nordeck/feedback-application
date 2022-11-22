<h1 align="center">
  <br>
  <a href="https://nordeck.net/"><img src="https://nordeck.net/wp-content/uploads/2020/05/NIC_logo_Nordeck-300x101.png" alt="Markdownify" width="300"></a>
  <br>
  Feedback-Application
  <br>
</h1>
<h4 align="center">A feedback tool for <a href="https://jitsi.org/jitsi-meet/" target="_blank">Jitsi Meet</a>.</h4>

### Key Features

This application is meant to run in conjunction with [Jitsi Meet](https://jitsi.org/jitsi-meet/) embedded as
an [Element](https://element.io/) widget to manage, view and persist feedback rating data.
The Jitsi instance must be hosted in the context of Element/[Matrix](https://matrix.org), as this application performs
authentication/authorization checks against
the [Matrix UVS](https://github.com/matrix-org/matrix-user-verification-service/) component.
The data is gathered from Jitsi's internal feedback dialogue and persisted into its own postgres database.
One is able to visualize the data with grafana.

### How To Use

This feedback application is made up of several components you need to install and configure.
You can find documentation on how to install, configure, and use each one in their respective folders.

1. [jitsi-feedback-plugin](./jitsi-feedback-plugin/) implements the frontend in form of a Jitsi Meet plugin
2. [backend](./backend/) is a REST service which verifies authorization and writes data to a database
3. [Grafana](https://grafana.com/) can be used to visualise the data using the provided [dashboard](./grafana/)

## Download

You can download the latest version of our feedback-application from this Github
repo (https://github.com/nordeck/feedback-application).

### How to Contribute

Please take a look at our [Contribution Guidelines](https://github.com/nordeck/.github/blob/main/docs/CONTRIBUTING.md).

### Related

[Jitsi Conference Mapper](https://github.com/nordeck/Jitsi-Conference-Mapper) - Jitsi Conference Mapper

### You may also like...

[Matrix Widget Toolkit](https://github.com/nordeck/matrix-widget-toolkit) - A widget toolkit
for [Matrix](https://matrix.org/)

[Matrix Poll](https://github.com/nordeck/matrix-poll) - A poll widget for [Matrix](https://matrix.org/)

### License

This project is licensed under [APACHE 2.0](./LICENSE).

### Sponsors

<p align="center">
   <a href="https://www.dataport.de/"><img src="./.docs/logos/dataportlogo.png" alt="Dataport" width="20%"></a>
   <a href="https://nordeck.net/"><img src="https://nordeck.net/wp-content/uploads/2020/05/NIC_logo_Nordeck-300x101.png" alt="Markdownify" width="300"></a>
</p>

> [nordeck.net](https://nordeck.net/) &nbsp;&middot;&nbsp;
> GitHub [Nordeck](https://github.com/nordeck/) 