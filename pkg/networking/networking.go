package networking

import (
	"fmt"
	"net"

	wifiname "github.com/gar-r/wifi-name"
)

var ErrInterfaceNotFound = fmt.Errorf("not found")

func SSIDForInterface(name string) (string, error) {
	has, err := hasInterface(name)
	if err == nil && !has {
		return "", fmt.Errorf("interface '%s' : %w", name, ErrInterfaceNotFound)
	}

	return wifiname.GetSSID(name)
}

func hasInterface(name string) (bool, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return false, fmt.Errorf("failed to get interfaces: %w", err)
	}

	for _, i := range ifaces {
		if i.Name == name {
			return true, nil
		}
	}

	return false, nil
}

func AddrForInterface(name string) (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", fmt.Errorf("failed to get interfaces: %w", err)
	}

	for _, i := range ifaces {
		if i.Name == name {

			addrs, err := i.Addrs()
			if err != nil {
				return "", fmt.Errorf("failed to get addresses: %w", err)
			}

			for _, a := range addrs {
				v, ok := a.(*net.IPNet)
				if ok {
					ip4 := v.IP.To4()
					if ip4 != nil {
						return ip4.String(), nil
					}
				}
			}

			return "", fmt.Errorf("no ipv4 address found for interface '%s'", name)
		}
	}

	return "", fmt.Errorf("interface '%s' : %w", name, ErrInterfaceNotFound)
}
