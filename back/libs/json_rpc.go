package libs

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/ybbus/jsonrpc/v3"
)

type FreeIPA struct {
	Domain     string
	Server     string
	ipaServers []string
	Secret     string
}

func (ipa *FreeIPA) GetIPAServers() ([]string, error) {

	if len(ipa.ipaServers) > 0 {
		return ipa.ipaServers, nil
	}

	if ipa.Server != "" {
		u, _ := url.Parse(ipa.Server)
		ipa.ipaServers = []string{u.Host}
		return ipa.ipaServers, nil
	}

	// Use the ldap SRV record to detect freeipa hosts
	_, srvs, err := net.LookupSRV("ldap", "tcp", ipa.Domain)
	if err != nil {
		log.Printf("[SRV] [Error] %s", err.Error())

		return nil, err
	}

	sort.Slice(srvs, func(i, j int) bool {
		return srvs[i].Priority < srvs[j].Priority
	})

	for _, srv := range srvs {
		ipa.ipaServers = append(ipa.ipaServers, srv.Target)
	}

	return ipa.ipaServers, nil
}

func (ipa *FreeIPA) createCookie(token string) *cookiejar.Jar {
	jar, _ := cookiejar.New(nil)
	host, bt, err := ParseToken(token, ipa.Secret)

	if err == nil {
		u, _ := url.Parse("https://" + host)
		c := &http.Cookie{
			Name:     "ipa_session",
			Value:    bt,
			Domain:   host,
			Path:     "/ipa",
			HttpOnly: true,
			Secure:   true,
		}
		jar.SetCookies(u, []*http.Cookie{c})
	}

	return jar
}

func (ipa *FreeIPA) getClient(token string) *http.Client {
	tr := &http.Transport{
		MaxIdleConns: 10,
		DialContext: (&net.Dialer{
			Timeout:   5 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		IdleConnTimeout:     30 * time.Second,
		DisableCompression:  true,
		TLSHandshakeTimeout: 3 * time.Second,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	cookie := ipa.createCookie(token)
	client := &http.Client{Transport: tr, Jar: cookie}

	return client
}

func (ipa FreeIPA) PostAuth(user, pass *string) (string, error) {
	var (
		err  error
		srvs []string
		req  *http.Request
		resp *http.Response
	)

	srvs, err = ipa.GetIPAServers()
	if err != nil {
		return "", err
	}

	data := url.Values{
		"user":     {*user},
		"password": {*pass},
	}

	client := ipa.getClient("")

	for _, srv := range srvs {
		req, err = http.NewRequest("POST",
			"https://"+srv+"/ipa/session/login_password",
			strings.NewReader(data.Encode()))
		if err != nil {
			log.Printf("[JSON_RPC] [Error] %s", err.Error())
			continue
		}

		req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=utf-8")
		req.Header.Add("referer", "https://"+srv+"/ipa")
		req.Header.Add("Accept", "text/plain")

		resp, err = client.Do(req)
		if err != nil {
			log.Printf("[JSON_RPC] [Error] %s", err.Error())
			continue
		}

		defer resp.Body.Close()

		if resp.StatusCode == 200 {
			token := UpdateToken(resp.Cookies()[0].Value, srv, ipa.Secret)
			return token, nil
		} else if resp.StatusCode >= 500 {
			continue
		}

		return "", errors.New("invalid credentials")
	}

	return "", err
}

func (ipa FreeIPA) Jrpc(ctx context.Context, token, method string, params ...interface{}) (any, int, error) {

	host, _, err := ParseToken(token, ipa.Secret)
	if err != nil {
		log.Println("[Warn] ", err)
		return nil, 401, errors.New("401, Unauthorized")
	}

	baseUrl := "https://" + host + "/ipa"

	rpcClient := jsonrpc.NewClientWithOpts(baseUrl+"/session/json", &jsonrpc.RPCClientOpts{
		CustomHeaders:      map[string]string{"referer": baseUrl},
		HTTPClient:         ipa.getClient(token),
		AllowUnknownFields: true,
	})

	resp, err := rpcClient.Call(ctx, method, params...)
	if err != nil {
		return nil, resp.Error.Code, err
	}

	if resp.Error != nil {
		return nil, resp.Error.Code, fmt.Errorf("code: %v; data: %v", resp.Error.Code, resp.Error.Data)
	}

	return resp.Result, 200, nil
}
