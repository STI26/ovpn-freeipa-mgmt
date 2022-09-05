# ovpn-freeipa-mgmt

[![License](https://img.shields.io/github/license/STI26/ovpn-freeipa-mgmt)](LICENSE)
[![Docker Image Size](https://badgen.net/docker/size/imagelist/ovpn_freeipa_mgmt/latest-ui?icon=docker&label=image%20size)](https://hub.docker.com/r/imagelist/ovpn_freeipa_mgmt)
[![Docker Image Size](https://badgen.net/docker/size/imagelist/ovpn_freeipa_mgmt/latest-api?icon=docker&label=image%20size)](https://hub.docker.com/r/imagelist/ovpn_freeipa_mgmt)

Web interface for openvpn uses Freeipa as Certificate Authority.

## Features

- Generate/Revoke user certificates
- Generate openvpn config
- Update certificate revocation lists

### Screenshot

![Screenshot](https://github.com/STI26/ovpn-freeipa-mgmt/blob/master/demo_ui.png?raw=true)

## Installation

Requirements: You need [openvpn](https://openvpn.net/community-downloads/) installed.

1. Install [Docker](https://docs.docker.com/get-docker/) and [docker-compose](https://docs.docker.com/compose/install/).
2. Download [docker-compose.yaml](https://github.com/STI26/ovpn-freeipa-mgmt/blob/master/docker-compose.yaml).
3. Change `https://ipa.example.com` to your FreeIPA server. See other options.

    ```bash
    command: >
        --ipa-server=https://ipa.example.com
    ```

4. Run docker container:

    ```bash
    sudo docker-compose up -d
    ```

5. Create openvpn server config - http://127.0.0.1:8080/config. To access the [user interface](http://127.0.0.1:8080), use the freeipa credentials.
6. Start systemd service:

    ```bash
    cd /etc/openvpn
    sudo systemctl start openvpn-server@server
    sudo systemctl enable openvpn-server@server
    ```

## Options

| Name | Default | Descriptions |
|------|---------|--------------|
| --addr | "0.0.0.0:8000" | Listening and serving address |
| --ipa-domain |  | Domain with IPA servers. Ignored if set `--ipa-server`. (search by SRV record) |
| --ipa-server |  | FreeIPA server with a scheme |
| --ipa-allowgroup | "admins" | IPA group with allowed access |
| --ipa-usergroup |  | Show users included in this ipa user group |
| --ipa-hostgroup |  | Show hosts included in this ipa host group |
| --ipa-cacn | "ipa" | Name of issuing CA |
| --ipa-ca-profile |  | IPA Certificate Profile to use |
| --ovpn-serverconf | "/etc/openvpn/server/server.conf" | Path to openvpn server.conf file |
| --ovpn-keys | "/etc/openvpn/keys" | Path to folder with user keys |
| --version |  | Show version. |
