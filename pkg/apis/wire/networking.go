package wire

type NetworkingResponse struct {
	WifiSSID  string `json:"wifi_ssid"`
	WifiAddr  string `json:"wifi_addr"`
	WiredAddr string `json:"wired_addr"`
}
