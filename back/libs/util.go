package libs

import (
	"bytes"
	b64 "encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func contains(source *[]int, target int) bool {
	for _, v := range *source {
		if v == target {
			return true
		}
	}

	return false
}

func KeyContains(source *[]fs.DirEntry, target string) bool {
	for _, v := range *source {
		if v.Name() == target {
			return true
		}
	}

	return false
}

func ParseToken(token, key string) (string, string, error) {

	if len(token) < 15 {
		return "", "", errors.New("error parse token - not valid size")
	}

	// cut "MagBearerToken="
	sDec, _ := b64.StdEncoding.DecodeString(token[15:])

	t := strings.Split(string(sDec), key)

	if len(t) != 2 {
		return "", "", errors.New("error parse token - IPA host not found")
	}

	// Return Host, MagBearerToken, Error
	return t[0], "MagBearerToken=" + t[1], nil
}

func UpdateToken(token, url, key string) string {
	sEnc := b64.StdEncoding.EncodeToString([]byte(url + key + token[15:]))

	return "MagBearerToken=" + string(sEnc)
}

func ParseResponse(result any, toType any) error {
	js, err := json.Marshal(result)
	if err != nil {
		return err
	}

	err = json.Unmarshal(js, toType)
	if err != nil {
		return err
	}

	return nil
}

func DownloadFile(url string, filepath string) (*bytes.Buffer, error) {

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if filepath == "" {
		buf := new(bytes.Buffer)

		_, err = io.Copy(buf, resp.Body)
		if err != nil {
			return nil, err
		}

		return buf, nil
	}

	// Create the file
	file, err := os.Create(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Write the body to file
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func GetPrincipal(subject, server string) string {
	domain := strings.SplitAfterN(server, ".", 2)[1]
	name := strings.SplitAfterN(subject, ".", 2)

	principal := subject + "@" + strings.ToUpper(domain)

	if len(name) > 1 && name[1] == domain {
		principal = "host/" + principal
	}

	return principal
}

func createFile(p string) (*os.File, error) {
	if err := os.MkdirAll(filepath.Dir(p), 0644); err != nil {
		return nil, err
	}
	return os.Create(p)
}
