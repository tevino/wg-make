package config

import (
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/tevino/wg-make/example"

	"gopkg.in/ini.v1"
)

func TestMapping(t *testing.T) {
	Convey("Load a correct config", t, func() {
		file, err := ini.LoadSources(LoadOptions, strings.NewReader(example.FileConfExample))
		So(err, ShouldBeNil)

		Convey("Mapping to Config", func() {
			conf := new(Config)
			err := file.MapTo(conf)
			So(err, ShouldBeNil)

			// Network
			So(conf.Network.ID, ShouldEqual, "example")
			So(conf.Network.Subnet, ShouldEqual, "192.168.25.0/24")
			So(conf.Peers, ShouldHaveLength, 3)
			// Peer 1
			So(conf.Peers[0].ID, ShouldEqual, "Tento")
			So(conf.Peers[0].Address, ShouldEqual, "192.168.25.55/32")
			So(conf.Peers[0].PrivateKey, ShouldEqual, "private-key-of-tento")
			So(conf.Peers[0].PublicKey, ShouldEqual, "public-key-of-tento")
			So(conf.Peers[0].PersistentKeepalive, ShouldEqual, 25)
			// Peer 2
			So(conf.Peers[1].ID, ShouldEqual, "Pata")
			So(conf.Peers[1].Address, ShouldEqual, "192.168.25.1/32")
			So(conf.Peers[1].PrivateKey, ShouldEqual, "private-key-of-pata")
			So(conf.Peers[1].PublicKey, ShouldEqual, "public-key-of-pata")
			So(conf.Peers[1].ListenPort, ShouldEqual, 49736)
			So(conf.Peers[1].Endpoint, ShouldEqual, "pata.example.com:49736")
			So(conf.Peers[1].PublicInterface, ShouldEqual, "eth0")
			So(conf.Peers[1].OS, ShouldEqual, "Linux")
		})
	})
}

func TestIsLinux(t *testing.T) {
	Convey("Create Peers with different OSes", t, func() {
		empty := Peer{}
		iPhone := Peer{OS: "iOS"}
		freeBSD := Peer{OS: "FreeBSD"}
		linux := Peer{OS: "Linux"}
		So(empty.IsLinux(), ShouldBeFalse)
		So(iPhone.IsLinux(), ShouldBeFalse)
		So(freeBSD.IsLinux(), ShouldBeFalse)
		So(linux.IsLinux(), ShouldBeTrue)
	})
}

func TestIsBounceServer(t *testing.T) {
	Convey("Create Peers with different properties", t, func() {
		empty := Peer{}

		endpointOnly := new(Peer)
		endpointOnly.Endpoint = "something"

		publicInterfaceOnly := new(Peer)
		publicInterfaceOnly.PublicInterface = "FreeBSD"

		both := new(Peer)
		both.Endpoint = "a"
		both.PublicInterface = "b"

		So(empty.IsBounceServer(), ShouldBeFalse)
		So(endpointOnly.IsBounceServer(), ShouldBeFalse)
		So(publicInterfaceOnly.IsBounceServer(), ShouldBeFalse)
		So(both.IsBounceServer(), ShouldBeTrue)
	})
}

func TestAllowedIPsForPeer(t *testing.T) {
	Convey("Create Peers with different properties", t, func() {
		empty := new(Peer)

		addressOnly := new(Peer)
		addressOnly.Address = "10.1.1.0/24"

		allowedIPsOnly := new(Peer)
		allowedIPsOnly.AllowedIPs = "20.1.1.0/24"

		both := new(Peer)
		both.Address = "10.1.1.0/24"
		both.AllowedIPs = "20.1.1.0/24"

		Convey("Both Address and AllowedIPs should be included", func() {
			So(empty.AllowedIPsForPeer(empty), ShouldBeEmpty)
			So(addressOnly.AllowedIPsForPeer(empty), ShouldEqual, addressOnly.Address)
			So(allowedIPsOnly.AllowedIPsForPeer(empty), ShouldEqual, allowedIPsOnly.AllowedIPs)
			So(both.AllowedIPsForPeer(empty), ShouldEqual, both.Address+","+both.AllowedIPs)
		})

		subnetB := new(Peer)
		subnetB.LocalSubnets = "20.1.1.0/24"

		subnetC := new(Peer)
		subnetC.LocalSubnets = "30.1.1.0/24"
		Convey("Peer's perspective should be applied", func() {
			So(both.AllowedIPsForPeer(subnetB), ShouldNotContainSubstring, subnetB.LocalSubnets)
			So(both.AllowedIPsForPeer(subnetC), ShouldEqual, both.AllowedIPsForPeer(empty))
		})
		Convey("allowedIPs should be treated as local subnet", func() {
			So(both.AllowedIPsForPeer(allowedIPsOnly), ShouldNotContainSubstring, allowedIPsOnly.AllowedIPs)
		})
	})
}

func TestIsIPInSubnet(t *testing.T) {
	Convey("", t, func() {
		So(isIPInSubnets("10.1.1.0/32", []string{"10.1.1.0/24"}), ShouldBeTrue)
	})
}
