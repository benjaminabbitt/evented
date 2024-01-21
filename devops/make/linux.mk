# Consists of linux commands to manage the k8s cluster.  Run from inside the machine running the k8s cluster.

## Developer setup
setup: install-consul install-rabbit install-mongo install-deps

## Consul Shortcuts
install-consul:
	helm repo add hashicorp https://helm.releases.hashicorp.com
	helm repo update
	helm install evented-consul hashicorp/consul --wait --set global.name=consul --values ./devops/helm/consul/values.yaml

consul-ui-expose:
	kubectl port-forward svc/consul-ui 80

consul-service-expose:
	kubectl port-forward svc/consul-server 8500


## RabbitMQ Shortcuts
install-rabbit:
	helm repo add bitnami https://charts.bitnami.com/bitnami
	helm repo update
	helm install evented-rabbitmq bitnami/rabbitmq --wait --values ./devops/helm/rabbitmq/values.yaml

rabbit-ui-expose:
	kubectl port-forward svc/evented-rabbitmq 15672:15672

rabbit-service-expose:
	kubectl port-forward svc/evented-rabbitmq 5672:5672

rabbit-extract-cookie:
	@python devops/make/get-secret/get-secret.py --secret="rabbitmq-erlang-cookie"

rabbit-extract-pw:
	@python devops/make/get-secret/get-secret.py --secret="rabbitmq-password"

## Mongo Shortcuts
mongo-service-expose:
	kubectl port-forward --namespace default svc/evented-mongodb 27017:27017

install-mongo:
	helm repo add bitnami https://charts.bitnami.com/bitnami
	helm repo update
	helm install evented-mongodb bitnami/mongodb --wait --values ./devops/helm/mongodb/values.yaml

install-k8s-services: install-mongo install-consul install-rabbit

install-helm:
	curl https://baltocdn.com/helm/signing.asc | gpg --dearmor | sudo tee /usr/share/keyrings/helm.gpg > /dev/null
	sudo apt-get install apt-transport-https --yes
	echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/helm.gpg] https://baltocdn.com/helm/stable/debian/ all main" | sudo tee /etc/apt/sources.list.d/helm-stable-debian.list
	sudo apt-get update
	sudo apt-get install helm

install-minikube:
	curl -LO https://storage.googleapis.com/minikube/releases/latest/minikube_latest_amd64.deb
	sudo dpkg -i minikube_latest_amd64.deb

install-os-dependencies: install-minikube install-helm

## Minikube Shortcuts
docker:
	sudo service docker start

minikube: docker
	MINIKUBE_ROOTLESS=false minikube start --feature-gates=GRPCContainerProbe=true --memory=12288 --cpus 4

minikube_enable_lb:
	minikube tunnel


