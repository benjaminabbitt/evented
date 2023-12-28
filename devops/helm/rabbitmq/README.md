# RabbitMQ Installation

[Bitnami's Rabbit Helm Chart](https://github.com/bitnami/charts/tree/master/bitnami/rabbitmq)

This may require a values file in the future (particularly to change domain to non-default in the event we're in a shared cluster), but works with defaults for now.

```shell
helm install rabbitmq oci://registry-1.docker.io/bitnamicharts/rabbitmq --values=./values.yaml
```
