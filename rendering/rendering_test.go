package rendering

import (
	"bytes"
	"regexp"
	"strings"
	"testing"

	"github.com/flexi-cache/pkg/testutil"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/tevino/wg-make/config"
	"github.com/tevino/wg-make/example"
)

func TestRenderPeerConfig(t *testing.T) {
	Convey("Render example config", t, func() {
		var (
			conf *config.Config
			err  error
		)
		testutil.WithTempFile(t, example.FileConfExample, func(filename string) {
			conf, err = config.LoadConfigFromFile(filename)
		})
		So(err, ShouldBeNil)
		So(conf, ShouldNotBeNil)

		var buf bytes.Buffer

		err = renderPeerConfig(&buf, conf, "Tento")
		So(err, ShouldBeNil)
		confTento := buf.String()
		Convey("Config of Tento should contain expected contents", func() {
			So(confTento, ShouldContainSubstring, "[Interface]")
			So(confTento, ShouldContainSubstring, "# ID = Tento")
			So(confTento, ShouldContainSubstring, "PrivateKey = private-key-of-tento")
			So(confTento, ShouldContainSubstring, "Address = 192.168.25.55/32")

			So(confTento, ShouldContainSubstring, "[Peer]")
			So(confTento, ShouldContainSubstring, "# ID = Pata")
			So(confTento, ShouldContainSubstring, "Endpoint = pata.example.com:49736")
			So(confTento, ShouldContainSubstring, "PublicKey = public-key-of-pata")
			So(confTento, ShouldContainSubstring, "AllowedIPs = 192.168.25.1/32")
		})
		Convey("Config of Tento should not contain unexpected contents", func() {
			So(confTento, ShouldNotContainSubstring, "public-key-of-tento")
			So(confTento, ShouldNotContainSubstring, "private-key-of-pata")
			So(regexp.MustCompile("(?m)= ?$").MatchString(confTento), ShouldBeFalse)
		})

		buf.Reset()
		err = renderPeerConfig(&buf, conf, "Pata")
		So(err, ShouldBeNil)
		confPata := buf.String()
		Convey("Config of Pata should contain expected contents", func() {
			So(confPata, ShouldContainSubstring, "[Interface]")
			So(confPata, ShouldContainSubstring, "# ID = Pata")
			So(confPata, ShouldContainSubstring, "PrivateKey = private-key-of-pata")
			So(confPata, ShouldContainSubstring, "Address = 192.168.25.1/32")

			So(confPata, ShouldContainSubstring, "[Peer]")
			So(confPata, ShouldContainSubstring, "# ID = Tento")
			So(confPata, ShouldContainSubstring, "PublicKey = public-key-of-tento")
			So(confPata, ShouldContainSubstring, "AllowedIPs = 192.168.25.55/32")
		})
		Convey("Config of Pata should not contain unexpected contents", func() {
			So(confPata, ShouldNotContainSubstring, "public-key-of-pata")
			So(confPata, ShouldNotContainSubstring, "private-key-of-tento")
		})
		Convey("All peers should be included for a server", func() {
			So(confPata, ShouldContainSubstring, "# ID = Agu")
			So(confPata, ShouldContainSubstring, "# ID = Tento")
		})
		Convey("Only peers acting as servers should be included for a client", func() {
			So(confTento, ShouldContainSubstring, "# ID = Pata")
			So(confTento, ShouldNotContainSubstring, "# ID = Agu")
			So(strings.Count(confTento, "[Peer]"), ShouldEqual, strings.Count(confTento, "Endpoint = "))
		})
		Convey("Bounce servers should relay for the WireGuard subnet", func() {
			So(confTento, ShouldContainSubstring, conf.Network.Subnet)
			So(strings.Count(confTento, conf.Network.Subnet), ShouldEqual, strings.Count(confTento, "[Peer]"))
		})
		Convey("PersistentKeepAlive for peers should be read from the Interface section", func() {
			So(confTento, ShouldContainSubstring, "PersistentKeepalive = 25")
			So(confTento, ShouldNotContainSubstring, "PersistentKeepalive = 5")

			So(confPata, ShouldNotContainSubstring, "PersistentKeepalive")
		})

		// general validations
		Convey("Sections are at the begin of lines", func() {
			rePeer := regexp.MustCompile(`(?m)^\[Peer\]$`)
			So(rePeer.MatchString(confTento), ShouldBeTrue)
			So(rePeer.MatchString(confPata), ShouldBeTrue)

			reInterface := regexp.MustCompile(`(?m)^\[Interface\]$`)
			So(reInterface.MatchString(confTento), ShouldBeTrue)
			So(reInterface.MatchString(confPata), ShouldBeTrue)
		})
		Convey("Only one field should be found per line", func() {
			reMultiField := regexp.MustCompile("(?m)^.+ = .+ = (.+)?$")
			So(reMultiField.MatchString(confTento), ShouldBeFalse)
			So(reMultiField.MatchString(confPata), ShouldBeFalse)
		})
		Convey("Fields must be finished", func() {
			reUnfinishedField := regexp.MustCompile("(?m)= ?$")
			So(reUnfinishedField.MatchString(confPata), ShouldBeFalse)
			So(reUnfinishedField.MatchString(confTento), ShouldBeFalse)
		})
	})
}
