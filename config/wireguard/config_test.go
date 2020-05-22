package wireguard

import (
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"gopkg.in/ini.v1"
)

const fullConf = `
[Interface]
Address = 192.0.2.3/32
ListenPort = 51820
PrivateKey = priKey
DNS = 1.1.1.1,8.8.8.8
Table = 12345
MTU = 1500
PreUp = /pre/up 1 %i
PreUp = /pre/up 2 %i
PostUp = /pos/up 1 %i
PostUp = /pos/up 2 %i
PreDown = /pre/down 1 %i
PreDown = /pre/down 2 %i
PostDown = /pos/down 1 %i
PostDown = /pos/down 2 %i

[Peer]
AllowedIPs = 192.0.2.1/24
Endpoint = node1.example.tld:1
PublicKey = pubKey1
PersistentKeepalive = 25

[Peer]
AllowedIPs = 192.0.2.2/24
Endpoint = node2.example.tld:2
PublicKey = pubKey2
PersistentKeepalive = 25
`

func TestMapping(t *testing.T) {
	Convey("Load a correct config file", t, func() {
		file, err := ini.LoadSources(LoadOptions, strings.NewReader(fullConf))
		So(err, ShouldBeNil)

		Convey("Map to config", func() {
			conf := new(Config)
			err := file.StrictMapTo(conf)
			So(err, ShouldBeNil)

			Convey("Expected values should be filled", func() {
				// Interface
				sIF := conf.Interface
				So(sIF.Address, ShouldEqual, "192.0.2.3/32")
				So(sIF.ListenPort, ShouldEqual, 51820)
				So(sIF.PrivateKey, ShouldEqual, "priKey")
				So(sIF.DNS, ShouldEqual, "1.1.1.1,8.8.8.8")
				So(sIF.MTU, ShouldEqual, 1500)
				So(sIF.PostUps, ShouldResemble, []string{"/pos/up 1 %i", "/pos/up 2 %i"})
				So(sIF.PostDowns, ShouldResemble, []string{"/pos/down 1 %i", "/pos/down 2 %i"})
				So(sIF.PreUps, ShouldResemble, []string{"/pre/up 1 %i", "/pre/up 2 %i"})
				So(sIF.PreDowns, ShouldResemble, []string{"/pre/down 1 %i", "/pre/down 2 %i"})
				// Peers
				sPS := conf.Peers
				So(sPS, ShouldHaveLength, 2)
				// Peer 1
				So(sPS[0].AllowedIPs, ShouldEqual, "192.0.2.1/24")
				So(sPS[0].Endpoint, ShouldEqual, "node1.example.tld:1")
				So(sPS[0].PublicKey, ShouldEqual, "pubKey1")
				So(sPS[0].PersistentKeepalive, ShouldEqual, 25)
				// Peer 2
				So(sPS[1].AllowedIPs, ShouldEqual, "192.0.2.2/24")
				So(sPS[1].Endpoint, ShouldEqual, "node2.example.tld:2")
				So(sPS[1].PublicKey, ShouldEqual, "pubKey2")
				So(sPS[1].PersistentKeepalive, ShouldEqual, 25)
			})
		})
	})
}
