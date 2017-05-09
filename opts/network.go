package opts

import (
	"encoding/csv"
	"fmt"
	"github.com/docker/docker/api/types/network"
	"strings"
)

const (
	networkOptName  = "name"
	networkOptAlias = "alias"
	driverOpt       = "driver-opt"
)

// NetworkOpt represents a network config in swarm mode.
type NetworkOpt struct {
	options []network.Options
}

// Set networkopts value
func (n *NetworkOpt) Set(value string) error {
	csvReader := csv.NewReader(strings.NewReader(value))
	fields, err := csvReader.Read()
	if err != nil {
		return err
	}
	var netOpt network.Options

	if len(fields) == 1 && !strings.Contains(fields[0], "=") {
		//Support legacy non csv format
		netOpt.Target = fields[0]
		goto attach
	}
	netOpt.DriverOpt = make(map[string]string)
	for _, field := range fields {
		parts := strings.SplitN(field, "=", 2)

		if len(parts) < 2 {
			return fmt.Errorf("invalid field %s", field)
		}

		key := strings.TrimSpace(strings.ToLower(parts[0]))
		value := strings.TrimSpace(strings.ToLower(parts[1]))

		switch key {
		case networkOptName:
			netOpt.Target = value
		case networkOptAlias:
			netOpt.Aliases = append(netOpt.Aliases, value)
		case driverOpt:
			key, value, err = parseDriverOpt(value)
			if err == nil {
				netOpt.DriverOpt[key] = value
			} else {
				return err
			}
		default:
			return fmt.Errorf("invalid field key %s", key)
		}
	}
	if len(netOpt.Target) == 0 {
		return fmt.Errorf("network name/id is not specified")
	}
attach:
	n.options = append(n.options, netOpt)

	return nil
}

// Type returns the type of this option
func (n *NetworkOpt) Type() string {
	return "network"
}

// Value returns the networkopts
func (n *NetworkOpt) Value() []network.Options {
	return n.options
}

// String returns the network opts as a string
func (n *NetworkOpt) String() string {
	networks := []string{}
	for _, network := range n.options {
		str := fmt.Sprintf("%s %s", network.Target, strings.Join(network.Aliases, "/"))
		networks = append(networks, str)
	}
	return strings.Join(networks, ",")
}

func parseDriverOpt(driverOpt string) (key string, value string, err error) {

	parts := strings.SplitN(driverOpt, "=", 2)
	if len(parts) != 2 {
		err = fmt.Errorf("invalid key value pair format in driver options")
	}
	key = strings.TrimSpace(strings.ToLower(parts[0]))
	value = strings.TrimSpace(strings.ToLower(parts[1]))
	return
}
