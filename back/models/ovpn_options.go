package models

type OvpnOption struct {
	Key         string `json:"key"`
	Value       string `json:"value"`
	Description string `json:"description"`
	Note        string `json:"note"`
	Status      bool   `json:"status,omitempty"`
}

var OvpnAvailableOptions = map[string]OvpnOption{
	"local": {
		"local",
		"",
		`Local host name or IP address for bind.
If specified, OpenVPN will bind to this address only.
If unspecified, OpenVPN will bind to all interfaces.`,
		"",
		false,
	},
	"port": {
		"port",
		"1194",
		`Which TCP/UDP port should OpenVPN listen on?
If you want to run multiple OpenVPN instances
on the same machine, use a different port
number for each one.  You will need to
open up this port on your firewall.`,
		"",
		false,
	},
	"proto": {
		"proto",
		"udp",
		`TCP or UDP server?`,
		"",
		false,
	},
	"dev": {
		"dev",
		"tun",
		`"dev tun" will create a routed IP tunnel,
"dev tap" will create an ethernet tunnel.
Use "dev tap0" if you are ethernet bridging
and have precreated a tap0 virtual interface
and bridged it with your ethernet interface.
If you want to control access policies
over the VPN, you must create firewall
rules for the the TUN/TAP interface.
On non-Windows systems, you can give
an explicit unit number, such as tun0.
On Windows, use "dev-node" for this.
On most systems, the VPN will not function
unless you partially or fully disable
the firewall for the TUN/TAP interface.`,
		"",
		false,
	},
	"dev-node": {
		"dev-node",
		"",
		`Windows needs the TAP-Win32 adapter name
from the Network Connections panel if you
have more than one.  On XP SP2 or higher,
you may need to selectively disable the
Windows firewall for the TAP adapter.
Non-Windows systems usually don't need this.`,
		"",
		false,
	},
	"ca": {
		"ca",
		"ca.crt",
		`Path to SSL/TLS root certificate (ca).`,
		"",
		false,
	},
	"cert": {
		"cert",
		"server.crt",
		`Path to server certificate.`,
		"",
		false,
	},
	"key": {
		"key",
		"server.key",
		`Path to private key.`,
		"",
		false,
	},
	"dh": {
		"dh",
		"dh2048.pem",
		`Path to Diffie hellman parameters.`,
		"",
		false,
	},
	"crl-verify": {
		"crl-verify",
		"/etc/openvpn/crl.pem",
		`This directive names a Certificate Revocation List file,
described below in the Revoking Certificates section.
The CRL file can be modified on the fly, and changes will
take effect immediately for new connections, or existing
connections which are renegotiating their SSL/TLS channel
(occurs once per hour by default). If you would like to kill
a currently connected client whose certificate has just
been added to the CRL, use the management interface (described below).`,
		"",
		false,
	},
	"topology": {
		"topology",
		"subnet",
		`Network topology
Should be subnet (addressing via IP)
unless Windows clients v2.0.9 and lower have to
be supported (then net30, i.e. a /30 per client)
Defaults to net30 (not recommended).`,
		"",
		false,
	},
	"remote-cert-tls": {
		"remote-cert-tls",
		"client",
		`For a server to verify that only hosts
with a client certificate can connect.`,
		"",
		false,
	},
	"persist-key": {
		"persist-key",
		"",
		`This option can be combined with --user nobody
to allow restarts triggered by the SIGUSR1 signal. 
Normally if you drop root privileges in OpenVPN.`,
		"",
		false,
	},
	"persist-tun": {
		"persist-tun",
		"",
		`Don't close and reopen TUN/TAP device or run up/down
scripts across SIGUSR1 or --ping-restart restarts.`,
		"",
		false,
	},
	"tun-mtu": {
		"tun-mtu",
		"1500",
		`Take the TUN device MTU.`,
		"",
		false,
	},
	"server": {
		"server",
		"10.8.0.0 255.255.255.0",
		`Configure server mode and supply a VPN subnet
for OpenVPN to draw client addresses from.
The server will take 10.8.0.1 for itself,
the rest will be made available to clients.
Each client will be able to reach the server
on 10.8.0.1. Comment this line out if you are
ethernet bridging. See the man page for more info.`,
		"",
		false,
	},
	"ifconfig-pool-persist": {
		"ifconfig-pool-persist",
		"/etc/openvpn/ipp.txt",
		`Maintain a record of client <-> virtual IP address
associations in this file.  If OpenVPN goes down or
is restarted, reconnecting clients can be assigned
the same virtual IP address from the pool that was
previously assigned.`,
		"",
		false,
	},
	"server-bridge": {
		"server-bridge",
		"10.8.0.4 255.255.255.0 10.8.0.50 10.8.0.100",
		`Configure server mode for ethernet bridging.
You must first use your OS's bridging capability
to bridge the TAP interface with the ethernet
NIC interface.  Then you must manually set the
IP/netmask on the bridge interface, here we
assume 10.8.0.4/255.255.255.0.  Finally we
must set aside an IP range in this subnet
(start=10.8.0.50 end=10.8.0.100) to allocate
to connecting clients.  Leave this line commented
out unless you are ethernet bridging.`,
		"",
		false,
	},
	"push": {
		"push",
		"route 192.168.10.0 255.255.255.0",
		`Push routes to the client to allow it
to reach other private subnets behind
the server.  Remember that these
private subnets will also need
to know to route the OpenVPN client
address pool (10.8.0.0/255.255.255.0)
back to the OpenVPN server.`,
		"",
		false,
	},
	"client-config-dir": {
		"client-config-dir",
		"/etc/openvpn/ccd",
		`To assign specific IP addresses to specific
clients or if a connecting client has a private
subnet behind it that should also have VPN access,
use the subdirectory "ccd" for client-specific
configuration files.`,
		"",
		false,
	},
	"client-to-client": {
		"client-to-client",
		"",
		`By default, clients will only see the server.
To force clients to only see the server, you
will also need to appropriately firewall the
server's TUN/TAP interface.`,
		"",
		false,
	},
	"duplicate-cn": {
		"duplicate-cn",
		"",
		`IF YOU HAVE NOT GENERATED INDIVIDUAL
CERTIFICATE/KEY PAIRS FOR EACH CLIENT,
EACH HAVING ITS OWN UNIQUE "COMMON NAME",
UNCOMMENT THIS LINE OUT.`,
		"",
		false,
	},
	"keepalive": {
		"keepalive",
		"10 120",
		`The keepalive directive causes ping-like
messages to be sent back and forth over
the link so that each side knows when
the other side has gone down.
Ping every 10 seconds, assume that remote
peer is down if no ping received during
a 120 second time period.`,
		"",
		false,
	},
	"tls-auth": {
		"tls-auth",
		"/etc/openvpn/ta.key 0",
		`For extra security beyond that provided
by SSL/TLS, create an "HMAC firewall"
to help block DoS attacks and UDP port flooding.
The server and each client must have
a copy of this key.
The second parameter should be '0'
on the server and '1' on the clients.`,
		"",
		false,
	},
	"auth": {
		"auth",
		"SHA256",
		`The OpenVPN data channel protocol uses encrypt-then-mac
i.e. first encrypt a packet then HMAC the resulting ciphertext),
which prevents padding oracle attacks.`,
		"",
		false,
	},
	"data-ciphers": {
		"data-ciphers",
		"AES-256-GCM:AES-128-GCM:AES-256-CBC:BF-CBC",
		`Data channel ciphers.`,
		"",
		false,
	},
	"cipher": {
		"cipher",
		"AES-256-CBC",
		`Select a cryptographic cipher.
This config item must be copied to
the client config file as well.
Note that v2.4 client/server will automatically
negotiate AES-256-GCM in TLS mode.
See also the ncp-cipher option in the manpage.`,
		"",
		false,
	},
	"comp-lzo": {
		"comp-lzo",
		"no",
		`For compression compatible with older clients use comp-lzo
If you enable it here, you must also
enable it in the client config file.`,
		"",
		false,
	},
	"max-clients": {
		"max-clients",
		"100",
		`The maximum number of concurrently connected
clients we want to allow.`,
		"",
		false,
	},
	"user": {
		"user",
		"nobody",
		`It's a good idea to reduce the OpenVPN
daemon's privileges after initialization.
You can uncomment this out on
non-Windows systems.`,
		"",
		false,
	},
	"group": {
		"group",
		"nobody",
		`It's a good idea to reduce the OpenVPN
daemon's privileges after initialization.
You can uncomment this out on
non-Windows systems.`,
		"",
		false,
	},
	"status": {
		"status",
		"/etc/openvpn/openvpn-status.log",
		`Output a short status file showing
current connections, truncated
and rewritten every minute.`,
		"",
		false,
	},
	"log": {
		"log",
		"/var/log/openvpn/openvpn.log",
		`By default, log messages will go to the syslog (or
on Windows, if running as a service, they will go to
the "\Program Files\OpenVPN\log" directory).
Use log or log-append to override this default.
"log" will truncate the log file on OpenVPN startup,
while "log-append" will append to it.  Use one
or the other (but not both).`,
		"",
		false,
	},
	"log-append": {
		"log-append",
		"/var/log/openvpn/openvpn.log",
		`By default, log messages will go to the syslog (or
on Windows, if running as a service, they will go to
the "\Program Files\OpenVPN\log" directory).
Use log or log-append to override this default.
"log" will truncate the log file on OpenVPN startup,
while "log-append" will append to it.  Use one
or the other (but not both).`,
		"",
		false,
	},
	"verb": {
		"verb",
		"3",
		`Set the appropriate level of log file verbosity.
0 is silent, except for fatal errors
4 is reasonable for general usage
5 and 6 can help to debug connection problems
9 is extremely verbose.`,
		"",
		false,
	},
	"mute": {
		"mute",
		"20",
		`Silence repeating messages.  At most 20
sequential messages of the same message
category will be output to the log.`,
		"",
		false,
	},
	"daemon": {
		"daemon",
		"",
		`Become a daemon after all initialization functions are completed.
This option will cause all message and error output
to be sent to the syslog file (such as /var/log/messages),
except for the output of scripts and ifconfig commands,
which will go to /dev/null unless otherwise redirected.
The syslog redirection occurs immediately at the point
that --daemon is parsed on the command line even though
the daemonization point occurs later.
If one of the --logoptions is present, it will supercede syslog redirection.
The optional progname parameter will cause OpenVPN
to report its program name to the system logger as progname.
This can be useful in linking OpenVPN messages in
the syslog file with specific tunnels.
When unspecified,progname defaults to "openvpn".`,
		"",
		false,
	},
}
