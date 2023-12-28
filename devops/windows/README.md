# Container Setup
This suite of applications is container based.  As the primary development is performed under Windows, which is not container native, some configuration is required.

We endeavor to make this as simple as possible and welcome feedback and suggestions.

## The current approach:

### WSL2

[Setting up WSL2](https://learn.microsoft.com/en-us/windows/wsl/install)

TL,DR:

`wsl --install`


`wsl -s Debian`

Debian is chosen due to its fairly zealous licensing, ensuring that the instructions do not become encumbered.

### Docker

The docker version in use is Docker Engine, *not* Docker Desktop.  This is viable for large enterprises and is not encumbered by licenses.  [Docker Engine is licensed under Apache License, v2.0](https://docs.docker.com/engine/)

[Install Docker](https://docs.docker.com/engine/install/debian/#install-using-the-repository)

[Post-install](https://docs.docker.com/engine/install/linux-postinstall/)

TL,DR:

`wsl`

```sh
# Add Docker's official GPG key:
sudo apt-get update
sudo apt-get install ca-certificates curl gnupg
sudo install -m 0755 -d /etc/apt/keyrings
curl -fsSL https://download.docker.com/linux/debian/gpg | sudo gpg --dearmor -o /etc/apt/keyrings/docker.gpg
sudo chmod a+r /etc/apt/keyrings/docker.gpg

# Add the repository to Apt sources:
echo \
  "deb [arch="$(dpkg --print-architecture)" signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/debian \
  "$(. /etc/os-release && echo "$VERSION_CODENAME")" stable" | \
  sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
sudo apt-get update
```

```shell
sudo apt-get install docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin
```

```shell
sudo docker run hello-world
```

[Post-Install](https://docs.docker.com/engine/install/linux-postinstall/)

```shell
sudo groupadd docker
```

Note that in the below command $USER is pre-populated by the distro.  This can be copy+pasted directly.
```shell
sudo usermod -aG docker $USER
```

```shell
newgrp docker
```

```shell
docker run hello-world
```


### Minikube

[Installing Minikube](https://minikube.sigs.k8s.io/docs/start/)

TL,DR:

`curl -LO https://storage.googleapis.com/minikube/releases/latest/minikube_latest_amd64.deb`

`sudo dpkg -i minikube_latest_amd64.deb`

#### Fixing DNS

```shell
mkdir -p ~/.minikube/files/etc
echo nameserver 8.8.8.8 > ~/.minikube/files/etc/resolv.conf
minikube stop
minikube start
```

[Source](https://rpi4cluster.com/awx/dns/#)


### Kubectl

```shell
sudo apt-get update
# apt-transport-https may be a dummy package; if so, you can skip that package
sudo apt-get install -y apt-transport-https ca-certificates curl
```

```shell
curl -fsSL https://pkgs.k8s.io/core:/stable:/v1.28/deb/Release.key | sudo gpg --dearmor -o /etc/apt/keyrings/kubernetes-apt-keyring.gpg
```

```shell
echo 'deb [signed-by=/etc/apt/keyrings/kubernetes-apt-keyring.gpg] https://pkgs.k8s.io/core:/stable:/v1.28/deb/ /' | sudo tee /etc/apt/sources.list.d/kubernetes.list
```

```shell
sudo apt-get update
sudo apt-get install -y kubectl
```

### Golang
```shell
sudo apt install golang
```


# Dependency Setup

TODO: Make a script vs a long list of installables?

## Helm

Helm is used to install both this framework and its dependencies (for dev, at least)

Follow instructions here: https://helm.sh/docs/intro/install/

```
curl https://baltocdn.com/helm/signing.asc | gpg --dearmor | sudo tee /usr/share/keyrings/helm.gpg > /dev/null
sudo apt-get install apt-transport-https --yes
echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/helm.gpg] https://baltocdn.com/helm/stable/debian/ all main" | sudo tee /etc/apt/sources.list.d/helm-stable-debian.list
sudo apt-get update
sudo apt-get install helm
```





## Consul

