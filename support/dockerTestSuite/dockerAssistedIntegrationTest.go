package dockerTestSuite

import (
	"context"
	"errors"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

type DockerAssistedIntegrationTest struct {
	Id    string
	Ports []types.Port
}

func (o *DockerAssistedIntegrationTest) GetPortMapping(servicePort uint16) (hostPort uint16, err error) {
	for _, port := range o.Ports {
		if servicePort == port.PrivatePort {
			return port.PublicPort, nil
		}
	}
	return 0, errors.New(fmt.Sprintf("Service port not found: %d", servicePort))
}

func (o *DockerAssistedIntegrationTest) PublicPort() (port uint16) {
	return o.Ports[0].PublicPort
}

func (o *DockerAssistedIntegrationTest) getPorts() (ports []types.Port, err error) {
	cli, err := client.NewEnvClient()
	if err != nil {
		fmt.Println("Unable to create docker client")
		panic(err)
	}

	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}

	for _, cont := range containers {
		if cont.ID == o.Id {
			ports = cont.Ports
			return ports, nil
		}
	}
	return nil, errors.New(fmt.Sprintf("No running cont found with id %s", o.Id))
}

func (o *DockerAssistedIntegrationTest) CreateNewContainer(image string, internalPorts []uint16) error {
	cli, err := client.NewEnvClient()
	if err != nil {
		fmt.Println("Unable to create docker client")
		panic(err)
	}

	hostBinding := nat.PortBinding{
		HostIP: "0.0.0.0",
	}

	var portBinding = make(nat.PortMap)
	for _, port := range internalPorts {
		strPort := fmt.Sprintf("%d", port)
		containerPort, err := nat.NewPort("tcp", strPort)
		if err != nil {
			panic("Unable to get the port")
		}
		portBinding[containerPort] = []nat.PortBinding{hostBinding}
	}

	cont, err := cli.ContainerCreate(
		context.Background(),
		&container.Config{
			Image: image,
		},
		&container.HostConfig{
			PortBindings: portBinding,
		}, nil, "")
	if err != nil {
		panic(err)
	}

	_ = cli.ContainerStart(context.Background(), cont.ID, types.ContainerStartOptions{})
	o.Id = cont.ID
	o.Ports, _ = o.getPorts()
	return nil
}

func (o *DockerAssistedIntegrationTest) StopContainer() error {
	cli, err := client.NewEnvClient()
	if err != nil {
		return err
	}

	err = cli.ContainerStop(context.Background(), o.Id, nil)
	return err
}
