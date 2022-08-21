package libs

import (
	"bytes"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"html/template"
	"log"
	"net"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type OpenVPN struct {
	Server   Network     `json:"server"`
	Config   OvpnOptions `json:"config"`
	IPAHosts []string    `json:"ipa_hosts,omitempty"`
}

type Network struct {
	IP   string `json:"ip"`
	Mask string `json:"mask"`
}

type OvpnOptions struct {
	// options with path
	Crl string `json:"crl"`

	CA      string `json:"ca"`
	TlsAuth string `json:"tls_auth"`
	DH      string `json:"dh"`

	Cert string `json:"cert"`
	Key  string `json:"key"`

	Ccd string `json:"ccd"`
	Ipp string `json:"ipp"`

	Status string `json:"status"`

	// options with key
	Local          string `json:"local"`
	Port           string `json:"port"`
	Proto          string `json:"proto"`
	Dev            string `json:"dev"`
	CompLzo        string `json:"comp_lzo"`
	TunMtu         string `json:"tun_mtu"`
	DataCiphers    string `json:"data_cipthers"`
	Auth           string `json:"auth"`
	Topology       string `json:"topology"`
	ClientToClient bool   `json:"client_to_client"`
}

func (ovpn *OpenVPN) Init(pathToConf string) error {
	b, err := os.ReadFile(pathToConf)
	if err != nil {
		log.Println(err)
		return err
	}

	for _, i := range strings.Split(string(b), "\n") {
		row := strings.Split(i, " ")

		switch row[0] {
		case "server":
			ovpn.Server.IP = strings.TrimSpace(row[1])
			ovpn.Server.Mask = strings.TrimSpace(row[2])

		case "ca":
			ovpn.Config.CA = strings.TrimSpace(row[1])
		case "tls-auth":
			ovpn.Config.TlsAuth = strings.TrimSpace(row[1])
		case "dh":
			ovpn.Config.DH = strings.TrimSpace(row[1])
		case "cert":
			ovpn.Config.Cert = strings.TrimSpace(row[1])
		case "key":
			ovpn.Config.Key = strings.TrimSpace(row[1])
		case "crl-verify":
			ovpn.Config.Crl = strings.TrimSpace(row[1])
		case "client-config-dir":
			ovpn.Config.Ccd = strings.TrimSpace(row[1])
		case "ifconfig-pool-persist":
			ovpn.Config.Ipp = strings.TrimSpace(row[1])
		case "status":
			ovpn.Config.Status = strings.TrimSpace(row[1])

		case "local":
			ovpn.Config.Local = strings.TrimSpace(row[1])
		case "port":
			ovpn.Config.Port = strings.TrimSpace(row[1])
		case "proto":
			ovpn.Config.Proto = strings.TrimSpace(row[1])
		case "dev":
			ovpn.Config.Dev = strings.TrimSpace(row[1])
		case "comp-lzo":
			ovpn.Config.CompLzo = strings.TrimSpace(row[1])
		case "tun-mtu":
			ovpn.Config.TunMtu = strings.TrimSpace(row[1])
		case "data-ciphers":
			ovpn.Config.DataCiphers = strings.TrimSpace(row[1])
		case "auth":
			ovpn.Config.Auth = strings.TrimSpace(row[1])
		case "topology":
			ovpn.Config.Topology = strings.TrimSpace(row[1])
		case "client-to-client":
			ovpn.Config.ClientToClient = true
		}
	}

	return nil
}

func (ovpn *OpenVPN) getIPAServer() string {
	timeout := 1 * time.Second

	for _, s := range ovpn.IPAHosts {
		_, err := net.DialTimeout("tcp", s, timeout)
		if err == nil {
			return s
		}
	}

	return ovpn.IPAHosts[0]
}

func (ovpn *OpenVPN) GetClientIP(client string) (string, error) {
	b, err := os.ReadFile(ovpn.Config.Ipp)
	if err != nil {
		return "", err
	}

	rows := strings.Split(string(b), "\n")
	for i := 0; i < len(rows); i++ {
		row := strings.Split(rows[i], ",")

		if row[0] == client {
			return row[1], nil
		}
	}

	return "", errors.New("ip not found")
}

func (ovpn *OpenVPN) UpdateClientIP(client, ip string) error {
	var rows []string
	changed := false

	b, err := os.ReadFile(ovpn.Config.Ipp)
	if err == nil {
		rows = strings.Split(strings.TrimSpace(string(b)), "\n")

		for i := 0; i < len(rows); i++ {
			row := strings.Split(rows[i], ",")

			if row[0] == client {
				rows[i] = client + "," + ip + "\n"
				changed = true
			}
		}
	}

	if !changed {
		rows = append(rows, client+","+ip+"\n")
	}

	data := strings.Join(rows, "\n")
	err = os.WriteFile(ovpn.Config.Ipp, []byte(data), 0644)
	if err != nil {
		return err
	}

	return nil
}

func (ovpn *OpenVPN) GetClientConfig(client string) (string, error) {
	path := filepath.Join(ovpn.Config.Ccd, client)

	b, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func (ovpn *OpenVPN) CreateServerConfig(path string) error {
	// Fill template
	data := new(bytes.Buffer)
	tmpl, err := template.ParseFiles("assets/server.tmpl")
	if err != nil {
		log.Println("[parseTemplate] [Error] ", err)
		return err
	}

	err = tmpl.Execute(data, map[string]interface{}{
		"network": ovpn.Server,
		"options": ovpn.Config,
	})
	if err != nil {
		log.Println("[tmpl.Execute] [Error] ", err)
		return err
	}

	err = os.WriteFile(path, data.Bytes(), 0644)
	if err != nil {
		log.Println("[WriteFile] [Error] ", err)
		return err
	}

	ovpn.UpdateCA()
	ovpn.UpdateCrl()

	return nil
}

func (ovpn *OpenVPN) UpdateClientConfig(client, data string) error {
	path := filepath.Join(ovpn.Config.Ccd, client)

	err := os.WriteFile(path, []byte(data), 0644)
	if err != nil {
		return err
	}

	return nil
}

func (ovpn *OpenVPN) UpdateCA() error {
	u := url.URL{
		Scheme: "https",
		Host:   ovpn.getIPAServer(),
		Path:   "/ipa/config/ca.crt",
	}

	_, err := DownloadFile(u.String(), ovpn.Config.CA)
	return err
}

func (ovpn *OpenVPN) GetCA() string {
	b, err := os.ReadFile(ovpn.Config.CA)
	if err != nil {
		log.Print(err)
		return ""
	}

	return string(b)
}

func (ovpn *OpenVPN) UpdateCrl() error {
	u := url.URL{
		Scheme: "https",
		Host:   ovpn.getIPAServer(),
		Path:   "/ipa/crl/MasterCRL.bin",
	}

	_, err := DownloadFile(u.String(), ovpn.Config.Crl)

	return err
}

func (ovpn *OpenVPN) GetStatusInfo() (string, error) {
	b, err := os.ReadFile(ovpn.Config.Status)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func (ovpn *OpenVPN) GetCrlInfo() (*pkix.TBSCertificateList, error) {
	b, err := os.ReadFile(ovpn.Config.Crl)
	if err != nil {
		return nil, err
	}

	crl, err := x509.ParseDERCRL(b)
	if err != nil {
		return nil, err
	}

	return &crl.TBSCertList, nil
}

func (ovpn *OpenVPN) GetServerConfig() *map[string]OvpnOptions {
	var name OvpnOptions
	var status OvpnOptions
	var message OvpnOptions

	// Check CRL
	name.Crl = "Certificate Revocation List"
	if b, err := os.ReadFile(ovpn.Config.Crl); err != nil {
		status.Crl = "false"
		message.Crl = "file not found"
	} else {
		status.Crl = "true"
		message.Crl = "Issuer: "

		if crl, err := x509.ParseDERCRL(b); err == nil {
			message.Crl += crl.TBSCertList.Issuer.String()
		}
	}

	// Check CA
	name.CA = "Certificate Authority"
	if b, err := os.ReadFile(ovpn.Config.CA); err != nil {
		status.CA = "false"
		message.CA = "file not found"
	} else {
		status.CA = "true"
		message.CA = "Issuer: "

		der, _ := pem.Decode(b)
		if ca, err := x509.ParseCertificate(der.Bytes); err == nil {
			message.CA += ca.Issuer.String()
		}
	}

	// Check tls-auth
	name.TlsAuth = "TLS Auth HMAC signature"
	if _, err := os.Stat(ovpn.Config.TlsAuth); err != nil {
		status.TlsAuth = "false"
		message.TlsAuth = "file not found"
	} else {
		status.TlsAuth = "true"
		message.TlsAuth = ""
	}

	// Check dh
	name.DH = "Diffie-Hellman"
	if _, err := os.Stat(ovpn.Config.DH); err != nil {
		status.DH = "false"
		message.DH = "file not found"
	} else {
		status.DH = "true"
		message.DH = ""
	}

	// Check cert
	name.Cert = "Server Certificate"
	if b, err := os.ReadFile(ovpn.Config.Cert); err != nil {
		status.Cert = "false"
		message.Cert = "file not found"
	} else {
		status.Cert = "true"
		message.Cert = "Issuer: "

		der, _ := pem.Decode(b)
		if cert, err := x509.ParseCertificate(der.Bytes); err == nil {
			message.Cert += cert.Issuer.String()
		}
	}

	// Check key
	name.Key = "Server Private Key"
	if _, err := os.ReadFile(ovpn.Config.Key); err != nil {
		status.Key = "false"
		message.Key = "file not found"
	} else {
		status.Key = "true"
		message.Key = ""
	}

	// Check ccd
	name.Ccd = "Client Config Dir"
	if _, err := os.Stat(ovpn.Config.Ccd); err != nil {
		status.Ccd = "false"
		message.Ccd = "folder not found"
	} else {
		status.Ccd = "true"
		message.Ccd = ""
	}

	// Check ipp
	name.Ipp = "Ifconfig Pool Persist"
	if _, err := os.Stat(ovpn.Config.Ipp); err != nil {
		status.Ipp = "false"
		message.Ipp = "file not found"
	} else {
		status.Ipp = "true"
		message.Ipp = ""
	}

	// Check status
	name.Status = "Status File"
	if _, err := os.Stat(ovpn.Config.Status); err != nil {
		status.Status = "false"
		message.Status = "file not found"
	} else {
		status.Status = "true"
		message.Status = ""
	}

	res := map[string]OvpnOptions{
		"name":    name,
		"value":   ovpn.Config,
		"status":  status,
		"message": message,
	}

	return &res
}
