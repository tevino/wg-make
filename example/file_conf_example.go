package example

// FileConfExample is the content of the example configuration file.
const FileConfExample = `# This is an example of the Network Description File, it defines a network you are going to create with WireGuard.
# With this single description file, wg-make generates configurations for all peers
# so you only need to install WireGuard and copy the configuration file(s) to the peer and everything just works.
#
# NOTE: Settings without special instructions are mandatory.

[Network]
# The ID of the network will be the suffix of the WireGuard interface name.
ID = "example"
# The subnet that contains all peers' addresses, this will be added into AllowedIPs for bounce servers.
Subnet = 192.168.25.0/24

# The Peer section is the combination of [Interface] and [Peer] in the WireGuard configuration file plus some extended settings.

# NOTE: Customizing number of the routing table is not supported for the moment.

# This is a client peer.
[Peer]
# The name of the peer, must be unique across networks.
ID = Tento
# The WireGuard IP Address of the peer.
Address = 192.168.25.55/32
# The subnets in which the peer already resides.
# This is useful when a peer is at the same subnet with a bounce server who's relaying the traffic(See AllowedIPs for a bounce server) to the subnet,
# in this case, setting this can avoid local subnet from being routed to the WireGuard interface, optional.
LocalSubnets = 10.1.1.0/24
# PrivateKey of the peer, could be generated with:
# prik=$(wg genkey); pubk=$(echo "$prik" | wg pubkey); echo -e "PrivateKey = $prik\nPublicKey = $pubk"
PrivateKey = private-key-of-tento
# PublicKey of the peer, could be generated like above.
PublicKey = public-key-of-tento
# Add this if THIS PEER is behind a NAT(no public IP), optional.
PersistentKeepalive = 25


# The peer acting as a server, relaying traffic for client peers.
[Peer]
ID = Pata
Address = 192.168.25.1/32
PrivateKey = private-key-of-pata
PublicKey = public-key-of-pata
#
# This Peer is a bounce server.
# The following settings are only for bounce servers, all optional.
#
# Port to listen.
ListenPort = 49736
# A publicly accessible address for other remote peers.
# The port number is mandatory, it could be different from the ListenPort.
Endpoint = pata.example.com:49736
# A range of the IPs or subnets that the bounce server is capable of routing traffic for.
AllowedIPs = 10.1.1.0/24
#
# NOTE: Omit the following two settings if you don't want WireGuard to change packet forwarding rules. (e.g. sysctl and iptables)
#
# Name of the network interface connecting to the Internet, used for adding packet forwarding rules.
PublicInterface = eth0
# Operating System, used to decide how to enable packet forwarding.
OS = Linux


# Another client behind NAT.
[Peer]
ID = Agu
Address = 192.168.25.15/32
LocalSubnets = 192.168.1.0/24
PrivateKey = private-key-of-agu
PublicKey = public-key-of-agu
PersistentKeepalive = 5
`
