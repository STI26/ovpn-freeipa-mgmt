package libs

import (
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

type IfconfigPoolPersist struct {
	Path    string
	Network Network
}

func (ipp IfconfigPoolPersist) getAllIP() []string {
	ips := []string{}

	b, err := os.ReadFile(ipp.Path)
	if err != nil {
		return ips
	}

	rows := strings.Split(string(b), "\n")
	for i := 0; i < len(rows); i++ {
		if strings.HasPrefix(rows[i], "#") {
			continue
		}

		row := strings.Split(rows[i], ",")

		if len(row) > 0 {
			ips = append(ips, strings.TrimSpace(row[1]))
		}
	}

	return ips
}

func (ipp IfconfigPoolPersist) isFreeIP(ip net.IP, pool []string) bool {
	for _, i := range pool {
		if ip.String() == i {
			return false
		}
	}

	return true
}

func (ipp IfconfigPoolPersist) getFreeIP() string {
	_mask := net.ParseIP(ipp.Network.Mask)
	ip := net.ParseIP(ipp.Network.IP)

	mask := net.IPMask(_mask)
	// Get first ip
	ip = ip.Mask(mask)

	mSize, size32 := mask.Size()
	// Get count ip's
	ips := 1 << (size32 - mSize)

	pool := ipp.getAllIP()

	for i := 0; i < ips; i++ {
		if ipp.isFreeIP(ip, pool) {
			return ip.String()
		}
		nextIP(ip)
	}

	return ""
}

func (ipp IfconfigPoolPersist) GetIP(name string) (string, error) {
	rows, err := os.ReadFile(ipp.Path)
	if err != nil {
		return "", err
	}

	for _, i := range strings.Split(string(rows), "\n") {
		row := strings.Split(string(i), ",")

		if row[0] == name {
			return row[1], nil
		}
	}

	return "", errors.New("not found ip")
}

func (ipp IfconfigPoolPersist) AddIP(name string) (string, error) {
	createFile(ipp.Path)

	rows, err := os.ReadFile(ipp.Path)
	if err != nil {
		log.Println("[ReadFile] [Warn] ", err)
	}

	ip := ipp.getFreeIP()
	if ip == "" {
		return "", errors.New("unable to allocate new ip")
	}

	newRow := fmt.Sprintf("%s,%s\n", name, ip)
	rows = append(rows, []byte(newRow)...)

	err = os.WriteFile(ipp.Path, rows, 0644)
	if err != nil {
		return "", err
	}

	return "", nil
}

func (ipp IfconfigPoolPersist) UpdateIP(name, ip string) error {
	rows, err := os.ReadFile(ipp.Path)
	if err != nil {
		return err
	}

	ips := []string{}
	for _, i := range strings.Split(string(rows), "\n") {
		row := strings.Split(string(i), ",")

		if row[0] == name {
			continue
		}

		ips = append(ips, i)
	}

	newRow := fmt.Sprintf("%s,%s\n", name, ip)
	ips = append(ips, newRow)

	err = os.WriteFile(ipp.Path, []byte(strings.Join(ips, "\n")), 0644)
	if err != nil {
		return err
	}

	return nil
}

func (ipp IfconfigPoolPersist) DeleteIP(name string) error {
	rows, err := os.ReadFile(ipp.Path)
	if err != nil {
		return err
	}

	ips := []string{}
	for _, i := range strings.Split(string(rows), "\n") {
		row := strings.Split(string(i), ",")

		if row[0] == name {
			continue
		}

		ips = append(ips, i)
	}

	err = os.WriteFile(ipp.Path, []byte(strings.Join(ips, "\n")), 0644)
	if err != nil {
		return err
	}

	return nil
}
