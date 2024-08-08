package networking

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestProcessNmcliOutput(t *testing.T) {
	sampleOut := `connection.id:                          Hotspot
connection.uuid:                        68d36582-9cc0-45ee-a38e-eb5a3d1ac5b2
connection.stable-id:                   --
connection.type:                        802-11-wireless
connection.interface-name:              wlan0
connection.autoconnect:                 yes
connection.autoconnect-priority:        0
connection.autoconnect-retries:         -1 (default)
connection.multi-connect:               0 (default)
connection.auth-retries:                -1
connection.timestamp:                   1722491329
connection.read-only:                   no
connection.permissions:                 --
connection.zone:                        --
connection.master:                      --
connection.slave-type:                  --
connection.autoconnect-slaves:          -1 (default)
connection.secondaries:                 --
connection.gateway-ping-timeout:        0
connection.metered:                     unknown
connection.lldp:                        default
connection.mdns:                        -1 (default)
connection.llmnr:                       -1 (default)
connection.dns-over-tls:                -1 (default)
connection.mptcp-flags:                 0x0 (default)
connection.wait-device-timeout:         -1
connection.wait-activation-delay:       -1
802-11-wireless.ssid:                   openbar-net-b5706c70
802-11-wireless.mode:                   ap
802-11-wireless.band:                   bg
802-11-wireless.channel:                0
802-11-wireless.bssid:                  --
802-11-wireless.rate:                   0
802-11-wireless.tx-power:               0
802-11-wireless.mac-address:            --
802-11-wireless.cloned-mac-address:     --
802-11-wireless.generate-mac-address-mask:--
802-11-wireless.mac-address-blacklist:  --
802-11-wireless.mac-address-randomization:default
802-11-wireless.mtu:                    auto
802-11-wireless.seen-bssids:            D8:3A:DD:66:D1:1A
802-11-wireless.hidden:                 no
802-11-wireless.powersave:              0 (default)
802-11-wireless.wake-on-wlan:           0x1 (default)
802-11-wireless.ap-isolation:           -1 (default)
802-11-wireless-security.key-mgmt:      wpa-psk
802-11-wireless-security.wep-tx-keyidx: 0
802-11-wireless-security.auth-alg:      --
802-11-wireless-security.proto:         --
802-11-wireless-security.pairwise:      --
802-11-wireless-security.group:         --
802-11-wireless-security.pmf:           0 (default)
802-11-wireless-security.leap-username: --
802-11-wireless-security.wep-key0:      <hidden>
802-11-wireless-security.wep-key1:      <hidden>
802-11-wireless-security.wep-key2:      <hidden>
802-11-wireless-security.wep-key3:      <hidden>`

	ssid, err := processNmcliOutput(sampleOut)
	require.NoError(t, err)
	require.Equal(t, "openbar-net-b5706c70", ssid)
}
