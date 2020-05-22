package rendering

import (
	"fmt"
	"io"
	"os"
	"path"
	"time"

	"github.com/tevino/log"
	"github.com/tevino/wg-make/config"
)

const (
	fileModeSensitive = 0700
	wgInterfacePrefix = "wg-"
)

// RenderNetwork render configurations of peers of a network described by networkConfFile into dirPeers.
func RenderNetwork(conf *config.Config, dirPeers string) error {
	for _, p := range conf.Peers {
		log.Infof("Rendering config for peer: %s\n", p.ID)
		dirPeer, err := ensurePeersDir(dirPeers, p.ID)
		if err != nil {
			return fmt.Errorf("ensuring folder for peer(%s): %w", p.ID, err)
		}
		confName := wgInterfacePrefix + conf.Network.ID + ".conf"
		confPath := path.Join(dirPeer, confName)
		flPeerConf, err := os.OpenFile(confPath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, fileModeSensitive)
		if err != nil {
			return fmt.Errorf("opening peer config(%s): %w", confPath, err)
		}

		err = renderPeerConfig(flPeerConf, conf, p.ID)
		flPeerConf.Close()
		if err != nil {
			return fmt.Errorf("rendering peer config: %w", err)
		}
	}
	return nil
}

func ensurePeersDir(dirPeers string, peerID string) (string, error) {
	dirPeer := path.Join(dirPeers, peerID)
	if _, err := os.Stat(dirPeer); os.IsNotExist(err) {
		err := os.MkdirAll(dirPeer, fileModeSensitive)
		if err != nil {
			return "", fmt.Errorf("creating config folder for peer (%s): %w", peerID, err)
		}
	} else if err != nil {
		return "", fmt.Errorf("reading config folder for peer (%s): %w", dirPeer, err)
	}
	return dirPeer, nil
}

func renderPeerConfig(dst io.Writer, conf *config.Config, peerID string) error {
	peers := []config.Peer{}
	targetPeer, ok := conf.GetPeerByID(peerID)
	if !ok {
		return fmt.Errorf("peer(%s) not found", peerID)
	}
	for _, p := range conf.Peers {
		if p.ID == peerID {
			continue
		}
		// Bounce servers should have all peers in its config.
		// Client peer only need the bounce server peers.
		if targetPeer.IsBounceServer() || p.IsBounceServer() {
			peers = append(peers, p)
		}
	}
	ctx := &PeerConfigTplContext{
		Network:     &conf.Network,
		Interface:   targetPeer,
		Peers:       peers,
		GeneratedAt: time.Now().Local(),
	}
	err := tplPeerConfig.Execute(dst, ctx)
	if err != nil {
		err = fmt.Errorf("rendering config for Peer(%s): %w", ctx.Interface.ID, err)
	}
	return err
}
