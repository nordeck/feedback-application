# Grafana Dashboard 

To visualize the feedback data, we created a simple dashboard for grafana.
You can import it from `/feedback-application/grafana/` and customize it to your needs.

You will need to configure your datasource beforehand to be the [backend](../backend/)'s postgres database.
The exact way to do this depends on your setup, for example you could use a docker network and its internal DNS, or publish the backend container's ports, etc.