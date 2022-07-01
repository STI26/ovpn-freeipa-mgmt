package main

import (
	"context"
	"crypto/tls"
	"errors"
	"flag"
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

type IPASession struct {
	cookies []*http.Cookie
}

type GlobalConfig struct {
	Title         *string
	ListenAddress *string

	LdapDomain     *string
	LdapServer     *string
	LdapAllowGroup *string

	CSRF           *bool
	CSRFSecret     *string
	CookieSecret   *string
	SessionsMaxAge *int

	IPAServers []*net.SRV
}

func (cfg *GlobalConfig) init() {
	cfg.Title = flag.String("title", "OpenVPN", "Page Title.")
	cfg.ListenAddress = flag.String("addr", "127.0.0.1:8888", "Listening and serving address.")

	cfg.LdapDomain = flag.String("domain", "", "Domain with ldap servers. (search by SRV record)")
	cfg.LdapServer = flag.String("server", "ipa.local.net", "Ldap server.")
	cfg.LdapAllowGroup = flag.String("allowgroup", "admins", "LDAP group with allowed access.")

	cfg.CSRF = flag.Bool("csrf", false, "CSRF enable. (true/false)")
	cfg.CSRFSecret = flag.String("csrfsecret", "RT9I_oXQTU4hfGLPFAMIm-ve-dDUdiTeiYsN2NxgcPk=", "CSRF secret key. Should be 32-bytes long.")
	cfg.CookieSecret = flag.String("cookiesecret", "RT9I_oXQTU4hfGLPFAMIm-ve-dDUdiTeiYsN2NxgcPk=", "Cookies secret key. It is recommended to use an authentication key with 32 or 64 bytes.")
	cfg.SessionsMaxAge = flag.Int("sessionsmaxage", 24, "Time in hour. Cookies lifetime.")

	flag.Parse()
}

func (cfg *GlobalConfig) getIPAServers() ([]*net.SRV, error) {

	if cfg.IPAServers != nil {
		return cfg.IPAServers, nil
	}

	if *cfg.LdapServer != "" {
		cfg.IPAServers = []*net.SRV{{Target: *config.LdapServer, Port: 389, Priority: 0, Weight: 100}}
		return cfg.IPAServers, nil
	}

	// Use the ldap SRV record to detect freeipa hosts
	_, srvs, err := net.LookupSRV("ldap", "tcp", *cfg.LdapDomain)
	if err != nil {
		log.Printf("[SRV] [Error] %s", err.Error())

		return nil, err
	}

	sort.Slice(srvs, func(i, j int) bool {
		return srvs[i].Priority < srvs[j].Priority
	})

	cfg.IPAServers = srvs

	return srvs, nil
}

func (s *IPASession) getClient(addr ...string) *http.Client {
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

	jar, _ := cookiejar.New(nil)
	if len(addr) > 0 {
		url, _ := url.Parse(addr[0])
		jar.SetCookies(url, s.cookies)
	}

	client := &http.Client{Transport: tr, Jar: jar}

	return client
}

func (s *IPASession) post_auth(user, pass *string) error {
	var (
		err  error
		srvs []*net.SRV
		req  *http.Request
		resp *http.Response
	)

	srvs, err = config.getIPAServers()
	if err != nil {
		return err
	}

	data := url.Values{
		"user":     {*user},
		"password": {*pass},
	}

	client := s.getClient()

	for _, srv := range srvs {
		req, err = http.NewRequest("POST",
			"https://"+srv.Target+"/ipa/session/login_password",
			strings.NewReader(data.Encode()))
		if err != nil {
			log.Printf("[RPC] [Error] %s", err.Error())
			continue
		}

		req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=utf-8")
		req.Header.Add("referer", "https://"+srv.Target+"/ipa")
		req.Header.Add("Accept", "text/plain")

		resp, err = client.Do(req)
		if err != nil {
			log.Printf("[RPC] [Error] %s", err.Error())
			continue
		}

		defer resp.Body.Close()

		if resp.StatusCode == 200 {
			s.cookies = resp.Cookies()
			return nil
		} else if resp.StatusCode >= 500 {
			continue
		}

		return errors.New("invalid credentials")
	}

	return err
}

func (s *IPASession) rpc(ctx context.Context, method string, params ...interface{}) error {
	var (
		err  error
		srvs []*net.SRV
		resp *jsonrpc.RPCResponse
	)

	srvs, err = config.getIPAServers()
	if err != nil {
		return err
	}

	for _, srv := range srvs {

		baseUrl := "https://" + srv.Target + "/ipa"

		rpcClient := jsonrpc.NewClientWithOpts(baseUrl+"/session/json", &jsonrpc.RPCClientOpts{
			CustomHeaders:      map[string]string{"referer": baseUrl},
			HTTPClient:         s.getClient(baseUrl),
			AllowUnknownFields: true,
		})

		resp, err = rpcClient.Call(ctx, method, params...)
		if err != nil {
			log.Println(err.Error())
			continue
		}

		if resp.Error != nil {
			return fmt.Errorf("code: %v; data: %v", resp.Error.Code, resp.Error.Data)
		}

		type UserCertificate struct {
			Base64 string `json:"__base64__"`
		}
		type ResultUsers struct {
			DN               string            `json:"dn"`
			GidNumber        []string          `json:"gidnumber"`
			GivenName        []string          `json:"givenname"`
			HomeDirectory    []string          `json:"homedirectory"`
			KrbCanonicalName []string          `json:"krbcanonicalname"`
			KrbPrincipalName []string          `json:"krbprincipalname"`
			LoginShell       []string          `json:"loginshell"`
			Mail             []string          `json:"mail"`
			NsAccountLock    bool              `json:"nsaccountlock"`
			SN               []string          `json:"sn"`
			UID              []string          `json:"uid"`
			UidNumber        []string          `json:"uidnumber"`
			MemberOfGroup    []string          `json:"memberof_group"`
			UserCertificate  []UserCertificate `json:"usercertificate"`
		}
		type R struct {
			Count   int           `json:"count"`
			Message any           `json:"messages"`
			Result  []ResultUsers `json:"result"`
			Summary string        `json:"summary"`
		}

		var r R
		err := resp.GetObject(&r)
		if err != nil {
			fmt.Println("-----err------", err)
			return err
		}
		fmt.Println("-----m------", r.Message)
		fmt.Printf("----DN %v\n", r.Result[0].DN)
		fmt.Printf("----GidNumber %v\n", r.Result[0].GidNumber)
		fmt.Printf("----GivenName %v\n", r.Result[0].GivenName)
		fmt.Printf("----HomeDirectory %v\n", r.Result[0].HomeDirectory)
		fmt.Printf("----KrbCanonicalName %v\n", r.Result[0].KrbCanonicalName)
		fmt.Printf("----KrbPrincipalName %v\n", r.Result[0].KrbPrincipalName)
		fmt.Printf("----LoginShell %v\n", r.Result[0].LoginShell)
		fmt.Printf("----Mail %v\n", r.Result[0].Mail)
		fmt.Printf("----NsAccountLock %v\n", r.Result[0].NsAccountLock)
		fmt.Printf("----SN %v\n", r.Result[0].SN)
		fmt.Printf("----UID %v\n", r.Result[0].UID)
		fmt.Printf("----UidNumber %v\n", r.Result[0].UidNumber)
		fmt.Printf("----MemberOfGroup %v\n", r.Result[0].MemberOfGroup)
		fmt.Printf("----UserCertificate %v\n", r.Result[0].UserCertificate[0].Base64)

		// var res map[string]interface{}

		// json.NewDecoder(resp.Body).Decode(&res)

		// fmt.Println("--------resp.StatusCode---------")
		// fmt.Println(resp.StatusCode)
		// fmt.Println(res)

		// if resp.StatusCode == 200 {
		// 	return nil
		// } else if resp.StatusCode >= 500 {
		// 	continue
		// }

		return errors.New("invalid credentials")
	}

	return err
}
