package wireguard

import "gopkg.in/ini.v1"

// Config reflects a wireguard configuration file.
type Config struct {
	Interface `ini:"Interface"`
	Peers     []Peer `ini:"Peer,,,nonunique"`
}

// Interface reflects the Interface section in a wireguard configuration file.
type Interface struct {
	PrivateKey string   `ini:"PrivateKey,omitempty"`
	Address    string   `ini:"Address,omitempty"`
	ListenPort int      `ini:"ListenPort,omitempty"`
	DNS        string   `ini:"DNS,omitempty"`
	MTU        int      `ini:"MTU,omitempty"`
	PreUps     []string `ini:"PreUp,omitempty,allowshadow"`
	PreDowns   []string `ini:"PreDown,omitempty,allowshadow"`
	PostUps    []string `ini:"PostUp,omitempty,allowshadow"`
	PostDowns  []string `ini:"PostDown,omitempty,allowshadow"`
}

// Peer reflects the Peer section in a wireguard configuration file.
type Peer struct {
	Endpoint            string `ini:"Endpoint,omitempty"`
	PublicKey           string `ini:"PublicKey,omitempty"`
	AllowedIPs          string `ini:"AllowedIPs,omitempty"`
	PersistentKeepalive int    `ini:"PersistentKeepalive,omitempty"`
}

// LoadOptions contains the options to load the config correctly.
var LoadOptions = ini.LoadOptions{
	Insensitive:            true,
	AllowNonUniqueSections: true,
	AllowShadows:           true,
}
