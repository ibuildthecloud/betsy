package plugin

import (
	"fmt"

	"github.com/docker/docker/pkg/plugingetter"
	"github.com/docker/docker/pkg/plugins"
)

/* defaultAPIVersion is the version of the plugin API for volume, network,
   IPAM and authz. This is a very stable API. When we update this API, then
   pluginType should include a version. eg "networkdriver/2.0".
*/
const defaultAPIVersion string = "1.0"

// ErrNotFound indicates that a plugin was not found locally.
type ErrNotFound string

func (name ErrNotFound) Error() string { return fmt.Sprintf("plugin %q not found", string(name)) }

// ErrAmbiguous indicates that a plugin was not found locally.
type ErrAmbiguous string

func (name ErrAmbiguous) Error() string {
	return fmt.Sprintf("multiple plugins found for %q", string(name))
}

// Get returns an enabled plugin matching the given name and capability.
func (ps *Store) Get(name, capability string, mode int) (plugingetter.CompatPlugin, error) {
	return plugins.Get(name, capability)
}

// GetAllManagedPluginsByCap returns a list of managed plugins matching the given capability.
func (ps *Store) GetAllManagedPluginsByCap(capability string) []plugingetter.CompatPlugin {
	return nil
}

// GetAllByCap returns a list of enabled plugins matching the given capability.
func (ps *Store) GetAllByCap(capability string) ([]plugingetter.CompatPlugin, error) {
	result := make([]plugingetter.CompatPlugin, 0, 1)
	pl, err := plugins.GetAll(capability)
	if err != nil {
		return nil, fmt.Errorf("legacy plugin: %v", err)
	}
	for _, p := range pl {
		result = append(result, p)
	}
	return result, nil
}

// Handle sets a callback for a given capability. It is only used by network
// and ipam drivers during plugin registration. The callback registers the
// driver with the subsystem (network, ipam).
func (ps *Store) Handle(capability string, callback func(string, *plugins.Client)) {
	plugins.Handle(capability, callback)
}
