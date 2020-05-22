package config

import (
	"fmt"
	"net"
	"strings"

	"github.com/tevino/wg-make/config/wireguard"
	"gopkg.in/ini.v1"
)

// Config reflects a network configuration file.
type Config struct {
	Network `ini:"Network"`
	Peers   []Peer `ini:"Peer,,,nonunique"`
}

// GetPeerByID returns Peer of given ID.
func (p *Config) GetPeerByID(id string) (*Peer, bool) {
	for i, peer := range p.Peers {
		if peer.ID == id {
			return &p.Peers[i], true
		}
	}
	return nil, false
}

// Network reflects the Network section within a network configuration file.
type Network struct {
	ID     string `ini:"ID"`
	Subnet string `ini:"Subnet"`
}

// Peer reflects a Peer section within a network configuration file.
type Peer struct {
	wireguard.Interface `ini:"Peer"`
	wireguard.Peer      `ini:"Peer"`
	ID                  string `ini:"ID"`
	LocalSubnets        string `int:"LocalSubnets"`
	PublicInterface     string `ini:"PublicInterface,omitempty"`
	OS                  string `ini:"OS,omitempty"`
}

// IsBounceServer returns true if the peer is capable of traffic relaying.
func (p *Peer) IsBounceServer() bool {
	return p.Endpoint != "" && p.PublicInterface != ""
}

// All OS types, currently only used to distinguish Linux.
const (
	OSLinux = "Linux"
)

// IsLinux returns true if OS is Linux.
func (p *Peer) IsLinux() bool {
	return p.OS == OSLinux
}

const ipv4Bits = 32

func isIPInSubnets(address string, subnets []string) bool {
	ip, _, err := net.ParseCIDR(address)
	if err != nil {
		panic(fmt.Sprintf("unexpected address(%s): %v", address, err))
	}
	for _, subnet := range subnets {
		if subnet == "" {
			continue
		}
		_, subnet, err := net.ParseCIDR(subnet)
		if err != nil {
			panic(fmt.Sprintf("unexpected subnet address(%s): %v", subnet, err))
		}
		if subnet.Contains(ip) {
			return true
		}
	}
	return false
}

// AllowedIPsForPeer returns the computed AllowedIPs at the perspective of given peer.
func (p *Peer) AllowedIPsForPeer(peer *Peer) string {
	allowedIPs := []string{}
	if p.Address != "" {
		allowedIPs = append(allowedIPs, p.Address)
	}
	//localSubnets := strings.Split(peer.LocalSubnets+peer.AllowedIPs, ",")
	localSubnets := strings.Split(peer.LocalSubnets+","+peer.AllowedIPs, ",")
	for _, ip := range strings.Split(p.AllowedIPs, ",") {
		if ip == "" {
			continue
		}
		if !isIPInSubnets(ip, localSubnets) {
			allowedIPs = append(allowedIPs, ip)
		}
	}
	return strings.Trim(strings.Join(allowedIPs, ","), ",")
}

// Validate returns the first error when validating the Peer.
func (p *Peer) Validate() error {
	if _, _, err := net.ParseCIDR(p.Address); err != nil && p.Address != "" {
		return fmt.Errorf("invalid address(%s): %w", p.Address, err)
	}
	if _, _, err := net.ParseCIDR(p.AllowedIPs); err != nil && p.AllowedIPs != "" {
		return fmt.Errorf("invalid AllowedIPs(%s): %w", p.Address, err)
	}
	return nil
}

// LoadOptions contains the options to load the config correctly.
var LoadOptions = ini.LoadOptions{
	Insensitive:            true,
	AllowNonUniqueSections: true,
	AllowShadows:           true,
}

// LoadConfigFromFile reads Config from given filePath.
func LoadConfigFromFile(filePath string) (*Config, error) {
	confFile, err := ini.LoadSources(LoadOptions, filePath)
	if err != nil {
		return nil, fmt.Errorf("loading source(%s): %w", filePath, err)
	}

	conf := new(Config)
	err = confFile.MapTo(conf)
	if err != nil {
		return nil, fmt.Errorf("mapping config(%s) to struct: %w", filePath, err)
	}
	return conf, nil
}
