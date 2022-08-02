package libs

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
	"log"
	"net"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type OpenVPN struct {
	Server   Network
	Config   OvpnOptions
	IPAHosts []string
}

type Network struct {
	IP   string
	Mask string
}

type OvpnOptions struct {
	// options with path
	Crl string

	CA      string
	TlsAuth string

	Cert string
	Key  string

	Ccd string
	Ipp string

	Status string

	// options with key
	Local       string
	Port        string
	Proto       string
	Dev         string
	CompLzo     string
	TunMtu      string
	DataCiphers string
	Auth        string
}

func (ovpn *OpenVPN) Init(pathToConf string) {
	b, err := os.ReadFile(pathToConf)
	if err != nil {
		log.Fatal(err)
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
		}
	}
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

func (ovpn *OpenVPN) GetFreeIP() string {
	return ""
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

	// TODO: ...
	// fmt.Println("-----------------------")
	// crl, _ := x509.ParseDERCRL(b.Bytes())
	// fmt.Println(crl.TBSCertList.RevokedCertificates)
	// os.Exit(0)
	// ...
	return err
}

func (ovpn *OpenVPN) GetStatusInfo() (string, error) {
	b, err := os.ReadFile(ovpn.Config.Status)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func (ovpn *OpenVPN) GetConfig() *map[string]OvpnOptions {
	var status OvpnOptions

	// Check CRL
	if b, err := os.ReadFile(ovpn.Config.Crl); err != nil {
		status.Crl = "file not found"
	} else {
		status.Crl = "exist, Issuer: "

		if crl, err := x509.ParseDERCRL(b); err == nil {
			status.Crl += crl.TBSCertList.Issuer.String()
		}
	}

	// Check CA
	if b, err := os.ReadFile(ovpn.Config.CA); err != nil {
		status.CA = "file not found"
	} else {
		status.CA = "exist, Issuer: "

		der, _ := pem.Decode(b)
		if ca, err := x509.ParseCertificate(der.Bytes); err == nil {
			status.CA += ca.Issuer.String()
		}
	}

	// Check tls-auth
	if _, err := os.Stat(ovpn.Config.TlsAuth); err != nil {
		status.TlsAuth = "file not found"
	} else {
		status.TlsAuth = "exist"
	}

	// Check cert
	if _, err := os.Stat(ovpn.Config.Cert); err != nil {
		status.Cert = "file not found"
	} else {
		status.Cert = "exist"
	}

	// Check key
	if _, err := os.Stat(ovpn.Config.Key); err != nil {
		status.Key = "file not found"
	} else {
		status.Key = "exist"
	}

	// Check ccd
	if _, err := os.Stat(ovpn.Config.Ccd); err != nil {
		status.Ccd = "file not found"
	} else {
		status.Ccd = "exist"
	}

	// Check ipp
	if _, err := os.Stat(ovpn.Config.Ipp); err != nil {
		status.Ipp = "file not found"
	} else {
		status.Ipp = "exist"
	}

	// Check status
	if _, err := os.Stat(ovpn.Config.Status); err != nil {
		status.Status = "file not found"
	} else {
		status.Status = "exist"
	}

	res := map[string]OvpnOptions{
		"value":  ovpn.Config,
		"status": status,
	}

	return &res
}
