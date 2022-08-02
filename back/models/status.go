package models

type Status struct {
	Title      string   `json:"title"`
	ClientList []Client `json:"client_list"`
}

type Client struct {
	CommonName         string `json:"common_name"`
	RealAddress        string `json:"real_address"`
	VirtualAddress     string `json:"virtual_address"`
	VirtualIPv6Address string `json:"virtual_ipv6_address"`
	BytesReceived      string `json:"bytes_received"`
	BytesSent          string `json:"bytes_send"`
	ConnectedSince     string `json:"connected_since"`
	Username           string `json:"user_name"`
	ClientID           string `json:"client_id"`
	PeerID             string `json:"peer_id"`
}
