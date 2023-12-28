# Mongo Installation

[Bitnami's MongoDB Helm Chart](https://github.com/bitnami/charts/tree/master/bitnami/mongodb)


This may require a values file in the future (particularly to change domain to non-default in the event we're in a shared cluster), but works with defaults for now.


```shell
helm install mongodb  oci://registry-1.docker.io/bitnamicharts/mongodb --values=values.yaml
```