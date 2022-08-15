package libs

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type IfconfigPoolPersist struct {
	Path    string
	Network Network
}

func (ipp IfconfigPoolPersist) getFreeIP(b *[]byte) string {
	ips := []int{}

	for _, i := range strings.Split(string(*b), "\n") {
		row := strings.Split(string(i), ".")

		num, err := strconv.Atoi(row[len(row)-1])
		if err == nil {
			ips = append(ips, num)
		}
	}

	for ip := 2; ip < 255; ip++ {
		if contains(&ips, ip) {
			continue
		}

		s := strings.Split(ipp.Network.IP, ".")
		s[3] = fmt.Sprint(ip)

		return strings.Join(s, ".")
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

	ip := ipp.getFreeIP(&rows)
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
