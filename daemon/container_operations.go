package daemon

import (
	"errors"
	"fmt"
	"os"

	derr "github.com/docker/docker/api/errors"
	"github.com/rancher/betsy/container"
)

var (
	// ErrRootFSReadOnly is returned when a container
	// rootfs is marked readonly.
	ErrRootFSReadOnly = errors.New("container rootfs is marked read-only")
)

func (daemon *Daemon) initializeNetworking(container *container.Container) error {
	var err error

	if container.HostConfig.NetworkMode.IsContainer() {
		// we need to get the hosts files from the container to join
		nc, err := daemon.getNetworkedContainer(container.ID, container.HostConfig.NetworkMode.ConnectedContainer())
		if err != nil {
			return err
		}
		initializeNetworkingPaths(container, nc)
		container.Config.Hostname = nc.Config.Hostname
		container.Config.Domainname = nc.Config.Domainname
		return nil
	}

	if container.HostConfig.NetworkMode.IsHost() {
		if container.Config.Hostname == "" {
			container.Config.Hostname, err = os.Hostname()
			if err != nil {
				return err
			}
		}
	}

	return container.BuildHostnameFile()
}

func (daemon *Daemon) getNetworkedContainer(containerID, connectedContainerID string) (*container.Container, error) {
	nc, err := daemon.GetContainer(connectedContainerID)
	if err != nil {
		return nil, err
	}
	if containerID == nc.ID {
		return nil, fmt.Errorf("cannot join own network")
	}
	if !nc.IsRunning() {
		err := fmt.Errorf("cannot join network of a non running container: %s", connectedContainerID)
		return nil, derr.NewRequestConflictError(err)
	}
	if nc.IsRestarting() {
		return nil, errContainerIsRestarting(connectedContainerID)
	}
	return nc, nil
}

func errRemovalContainer(containerID string) error {
	return fmt.Errorf("Container %s is marked for removal and cannot be connected or disconnected to the network", containerID)
}
