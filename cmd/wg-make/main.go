package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/tevino/log"
	"github.com/tevino/wg-make/config"
	"github.com/tevino/wg-make/example"
	"github.com/tevino/wg-make/rendering"
)

const (
	extConf             = ".conf"
	dirNetworks         = "networks"
	dirPeers            = "peers"
	filenameExampleConf = "example" + extConf
)

func setLoggerLevel(logLevel log.Level) {
	if logLevel != log.DEBUG {
		log.SetDefaultLogger(log.NewLogger(os.Stdout, 0))
	}
	log.SetOutputLevel(logLevel)
}

func main() {
	opt := new(opt).Parse()
	setLoggerLevel(opt.logLevel)
	if opt.needClean {
		cleanPeers()
	}
	if opt.needExample {
		createExampleNetwork()
	}
	renderNetworks()
}

func infoTitlef(f string, a ...interface{}) {
	log.Infof(fmt.Sprintf("==== %s ====", f), a...)
}

const fileModeSensitive = 0700

func createExampleNetwork() {
	files, err := ioutil.ReadDir(".")
	if err != nil {
		log.Fatalf("reading current directory: %v", err)
	}
	if len(files) > 0 {
		fmt.Print(`The current diectory is not empty, are you sure to create directory structure here?
You may want to cd into a new directory first (y/N): `)
		input := bufio.NewReader(os.Stdin)
		answer, _, err := input.ReadRune()
		if err != nil {
			log.Fatalf("reading input from stdin: %v", err)
		}
		if answer != 'y' && answer != 'Y' {
			return
		}
	}
	infoTitlef("Generating directory structure with examples")
	log.Info("Creating networks directory")
	err = os.MkdirAll(dirNetworks, fileModeSensitive)
	if err != nil {
		log.Fatalf("Creating directory(%s): %v", dirNetworks, err)
	}
	log.Info("Creating example network configuration")
	filename := path.Join(dirNetworks, filenameExampleConf)
	err = ioutil.WriteFile(filename, []byte(example.FileConfExample), fileModeSensitive)
	if err != nil {
		log.Fatalf("Creating file(%s): %v", filename, err)
	}
}

func cleanPeers() {
	log.Warn("Cleaning peers folder")
	err := os.RemoveAll(dirPeers)
	if err != nil {
		log.Fatalf("Removing folder(%s): %v", dirPeers, err)
	}
}

func renderNetworks() {
	files, err := ioutil.ReadDir(dirNetworks)
	if err != nil {
		log.Fatalf("Reading networks dir(%s): %v", dirNetworks, err)
	}

	for _, f := range files {
		if f.IsDir() || !strings.HasSuffix(f.Name(), extConf) {
			continue
		}

		pathNetworkConf := path.Join(dirNetworks, f.Name())
		conf, err := config.LoadConfigFromFile(pathNetworkConf)
		if err != nil {
			log.Fatalf("unexpected config file(%s): %v", pathNetworkConf, err)
		}
		// TODO: validateConf()
		infoTitlef("Found %d Peer(s) in network %s", len(conf.Peers), conf.Network.ID)
		err = rendering.RenderNetwork(conf, dirPeers)
		if err != nil {
			log.Fatalf("Rendering network %s: %v", conf.Network.ID, err)
		}
	}
}
