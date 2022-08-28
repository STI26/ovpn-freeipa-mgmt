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

	"github.com/sti26/ovpn_freeipa_mgmt/models"
)

type OpenVPN struct {
	Server   Network             `json:"server"`
	Optons   []models.OvpnOption `json:"options"`
	IPAHosts []string            `json:"ipa_hosts,omitempty"`
}

type Network struct {
	IP   string `json:"ip"`
	Mask string `json:"mask"`
}

func (ovpn *OpenVPN) Init(pathToConf string) error {
	// Clear options
	ovpn.Optons = []models.OvpnOption{}

	b, err := os.ReadFile(pathToConf)
	if err != nil {
		log.Println(err)
		return err
	}

	for _, i := range strings.Split(string(b), "\n") {
		row := strings.SplitAfterN(i, " ", 2)

		key := strings.TrimSpace(row[0])
		option, exist := models.OvpnAvailableOptions[key]
		if exist {
			if len(row) > 1 {
				option.Value = strings.TrimSpace(row[1])
			} else {
				option.Value = ""
			}
			ovpn.Optons = append(ovpn.Optons, option)

			if option.Key == "server" {
				v := strings.Split(row[1], " ")
				ovpn.Server.IP = strings.TrimSpace(v[0])
				ovpn.Server.Mask = strings.TrimSpace(v[1])
			}
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

func (ovpn *OpenVPN) GetOptionsByKey(key string) []*models.OvpnOption {
	opts := []*models.OvpnOption{}

	for _, opt := range ovpn.Optons {
		if opt.Key == key {
			opts = append(opts, &opt)
		}
	}

	return opts
}

func (ovpn *OpenVPN) GetOptionByKey(key string) (*models.OvpnOption, int) {
	for idx, opt := range ovpn.Optons {
		if opt.Key == key {
			return &opt, idx
		}
	}

	return nil, 0
}

func (ovpn *OpenVPN) GetClientIP(client string) (string, error) {
	ipp, _ := ovpn.GetOptionByKey("ifconfig-pool-persist")
	if ipp == nil {
		return "", errors.New("ifconfig-pool-persist not found")
	}

	b, err := os.ReadFile(ipp.Value)
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

	ipp, _ := ovpn.GetOptionByKey("ifconfig-pool-persist")
	if ipp == nil {
		return errors.New("ifconfig-pool-persist not found")
	}

	b, err := os.ReadFile(ipp.Value)
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
	err = os.WriteFile(ipp.Value, []byte(data), 0644)
	if err != nil {
		return err
	}

	return nil
}

func (ovpn *OpenVPN) GetClientConfig(client string) (string, error) {
	ccd, _ := ovpn.GetOptionByKey("client-config-dir")
	if ccd == nil {
		return "", errors.New("client-config-dir not found")
	}

	path := filepath.Join(ccd.Value, client)

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

	local, _ := ovpn.GetOptionByKey("local")
	err = tmpl.Execute(data, map[string]interface{}{
		"network": ovpn.Server,
		"local":   local.Value,
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

	return nil
}

func (ovpn *OpenVPN) UpdateClientConfig(client, data string) error {
	ccd, _ := ovpn.GetOptionByKey("client-config-dir")
	if ccd == nil {
		return errors.New("client-config-dir not found")
	}

	path := filepath.Join(ccd.Value, client)

	err := os.WriteFile(path, []byte(data), 0644)
	if err != nil {
		return err
	}

	return nil
}

func (ovpn *OpenVPN) UpdateCrl() error {
	u := url.URL{
		Scheme: "https",
		Host:   ovpn.getIPAServer(),
		Path:   "/ipa/crl/MasterCRL.bin",
	}

	crl, _ := ovpn.GetOptionByKey("crl-verify")
	BackupFile(crl.Value)

	_, err := DownloadFile(u.String(), crl.Value)

	return err
}

func (ovpn *OpenVPN) GetStatusInfo() (string, error) {
	status, _ := ovpn.GetOptionByKey("status")
	b, err := os.ReadFile(status.Value)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func (ovpn *OpenVPN) GetCrlInfo() (*pkix.TBSCertificateList, error) {
	_crl, _ := ovpn.GetOptionByKey("crl-verify")
	b, err := os.ReadFile(_crl.Value)
	if err != nil {
		return nil, err
	}

	crl, err := x509.ParseDERCRL(b)
	if err != nil {
		return nil, err
	}

	return &crl.TBSCertList, nil
}

func (ovpn *OpenVPN) GetServerConfig() []models.OvpnOption {

	// Check CRL
	crl, idx := ovpn.GetOptionByKey("crl-verify")
	if crl != nil {
		if b, err := os.ReadFile(crl.Value); err != nil {
			ovpn.Optons[idx].Status = false
			ovpn.Optons[idx].Note = "file not found"
		} else {
			ovpn.Optons[idx].Status = true
			ovpn.Optons[idx].Note = "Issuer: "

			if _crl, err := x509.ParseDERCRL(b); err == nil {
				ovpn.Optons[idx].Note += _crl.TBSCertList.Issuer.String()
			}
		}
	}

	// Check CA
	ca, idx := ovpn.GetOptionByKey("ca")
	if ca != nil {
		if b, err := os.ReadFile(ca.Value); err != nil {
			ovpn.Optons[idx].Status = false
			ovpn.Optons[idx].Note = "file not found"
		} else {
			ovpn.Optons[idx].Status = true
			ovpn.Optons[idx].Note = "Issuer: "

			der, _ := pem.Decode(b)
			if _ca, err := x509.ParseCertificate(der.Bytes); err == nil {
				ovpn.Optons[idx].Note += _ca.Issuer.String()
			}
		}
	}

	// Check tls-auth
	tlsAuth, idx := ovpn.GetOptionByKey("tls-auth")
	if tlsAuth != nil {
		_tlsAuth := strings.Split(tlsAuth.Value, " ")
		if _, err := os.Stat(_tlsAuth[0]); err != nil {
			ovpn.Optons[idx].Status = false
			ovpn.Optons[idx].Note = "file not found"
		} else {
			ovpn.Optons[idx].Status = true
		}
	}

	// Check dh
	dh, idx := ovpn.GetOptionByKey("dh")
	if dh != nil {
		if _, err := os.Stat(dh.Value); err != nil {
			ovpn.Optons[idx].Status = false
			ovpn.Optons[idx].Note = "file not found"
		} else {
			ovpn.Optons[idx].Status = true
		}
	}

	// Check cert
	cert, idx := ovpn.GetOptionByKey("cert")
	if cert != nil {
		if b, err := os.ReadFile(cert.Value); err != nil {
			ovpn.Optons[idx].Status = false
			ovpn.Optons[idx].Note = "file not found"
		} else {
			ovpn.Optons[idx].Status = true
			ovpn.Optons[idx].Note = "Issuer: "

			der, _ := pem.Decode(b)
			if _cert, err := x509.ParseCertificate(der.Bytes); err == nil {
				ovpn.Optons[idx].Note += _cert.Issuer.String()
			}
		}
	}

	// Check key
	key, idx := ovpn.GetOptionByKey("key")
	if key != nil {
		if _, err := os.Stat(key.Value); err != nil {
			ovpn.Optons[idx].Status = false
			ovpn.Optons[idx].Note = "file not found"
		} else {
			ovpn.Optons[idx].Status = true
		}
	}

	// Check ccd
	ccd, idx := ovpn.GetOptionByKey("client-config-dir")
	if ccd != nil {
		if _, err := os.Stat(ccd.Value); err != nil {
			ovpn.Optons[idx].Status = false
			ovpn.Optons[idx].Note = "folder not found"
		} else {
			ovpn.Optons[idx].Status = true
		}
	}

	// Check ipp
	ipp, idx := ovpn.GetOptionByKey("ifconfig-pool-persist")
	if ipp != nil {
		if _, err := os.Stat(ipp.Value); err != nil {
			ovpn.Optons[idx].Status = false
			ovpn.Optons[idx].Note = "file not found"
		} else {
			ovpn.Optons[idx].Status = true
		}
	}

	// Check status
	status, idx := ovpn.GetOptionByKey("ifconfig-pool-persist")
	if status != nil {
		if _, err := os.Stat(status.Value); err != nil {
			ovpn.Optons[idx].Status = false
			ovpn.Optons[idx].Note = "file not found"
		} else {
			ovpn.Optons[idx].Status = true
		}
	}

	return ovpn.Optons
}
