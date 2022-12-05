# Feedback Backend Helm Chart

## Key Features

This Helm chart will install and configure for you:

- Our feedback backend component
- A PostgreSQL database as persistent storage for the backend, using the [bitnami/postgresql](https://github.com/bitnami/charts/tree/main/bitnami/postgresql) Helm chart
- Optionally, a non-persistent Grafana preconfigured with the provided [dashboard](/grafana) and the mentioned database as datasource

This Helm chart **will not** set up and configure all the other components of a Matrix + Jitsi deployment. You need to do this yourself, in particular the following are required:

- Synapse, the Matrix homeserver
  - [Matrix User Verification Service](https://github.com/matrix-org/matrix-user-verification-service/), a Synapse plugin
- Jitsi
  - [Prosody Auth Matrix User Verification](https://github.com/matrix-org/prosody-mod-auth-matrix-user-verification) module
  - Our (jitsi-feedback-plugin)[/jitsi-feedback-plugin] installed
- Element, the Matrix client, configured to use the above Jitsi deployment

## Configuration

While this Helm chart tries to stick to reasonable defaults, you still need to configure a few settings to match your existing deployment:

- `service.oidcValidationUrl`
- `service.matrixServerName`

Further notable values are:

- `image.tag` to deploy an image other than `appVersion`
- `image.repository` if you build your own images of the feedback-backend

For detailed documentation of all possible values, refer to the comments in the [`values.yaml`](values.yaml) file.

