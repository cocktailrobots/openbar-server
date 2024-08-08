package networking

import (
	"fmt"
	"io"
	"net"
	"os/exec"
	"strings"
)

var ErrInterfaceNotFound = fmt.Errorf("not found")

func SSIDForHotspot() (string, error) {
	cmd := exec.Command("nmcli", "con", "show", "Hotspot")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return "", fmt.Errorf("failed to create stdout pipe: %w", err)
	}

	err = cmd.Start()
	if err != nil {
		return "", fmt.Errorf("failed to run nmcli: %w", err)
	}

	data, err := io.ReadAll(stdout)
	if err != nil {
		return "", fmt.Errorf("failed to read nmcli output: %w", err)
	}

	return processNmcliOutput(string(data))
}

func processNmcliOutput(data string) (string, error) {
	lines := strings.Split(data, "\n")
	for _, line := range lines {
		tokens := strings.Split(line, ":")
		if len(tokens) == 2 {
			key := strings.ToLower(strings.TrimSpace(tokens[0]))
			if key == "802-11-wireless.ssid" {
				return strings.TrimSpace(tokens[1]), nil
			}
		}
	}

	return "", fmt.Errorf("failed to find SSID in nmcli output")
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
