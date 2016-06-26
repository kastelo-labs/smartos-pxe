package main

import (
	"flag"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/pin/tftp"
)

type state struct {
	Console    string
	RootShadow string
	Versions   []string
	Overlay    []string
	Timeout    int
}

var (
	console         = "ttyS2"
	listen          = ":69"
	rootShadow      = "$5$hEQ0l8d5$s0Jwt.oif76hVoQpzsgH2XVKhS8uCXnMlQhXltYgvaB"
	bootfile        = "grub/pxegrub"
	datadir         = "data"
	timeout         = 10 * time.Second
	upgradeServer   = "https://us-east.manta.joyent.com"
	upgradePath     = "/Joyent_Dev/public/SmartOS/latest"
	upgradeInterval = 24 * time.Hour
	debug           bool
)

func main() {
	flag.BoolVar(&debug, "debug", debug, "Debug output")
	flag.StringVar(&console, "console", console, "Console device")
	flag.StringVar(&listen, "listen", listen, "Listen address")
	flag.StringVar(&rootShadow, "root-pw", rootShadow, "Root password hash")
	flag.StringVar(&bootfile, "boot-file", bootfile, "Boot file")
	flag.StringVar(&datadir, "data-dir", datadir, "Data directory")
	flag.StringVar(&upgradeServer, "upgrade-server", upgradeServer, "Platform upgrade server")
	flag.StringVar(&upgradePath, "upgrade-path", upgradePath, "Path to latest platform file")
	flag.DurationVar(&timeout, "grub-timeout", timeout, "GRUB menu timeout")
	flag.DurationVar(&upgradeInterval, "upgrade-intv", upgradeInterval, "New platform upgrade interval (0 to disable)")
	flag.Parse()

	log.SetOutput(os.Stdout)
	if debug {
		log.SetFlags(log.Lshortfile | log.LUTC | log.Ltime)
	} else {
		log.SetFlags(0)
	}

	if upgradeInterval > 0 {
		go autoUpgrade(upgradeInterval)
	}

	if debug {
		log.Println("Listening on", listen)
	}
	s := tftp.NewServer(readHandler, nil)
	s.SetTimeout(60 * time.Second)
	if err := s.ListenAndServe(listen); err != nil {
		log.Fatal(err)
	}
}

func autoUpgrade(intv time.Duration) {
	for {
		if err := upgrade(); err != nil {
			log.Println("Downloading upgrade:", err)
		}
		time.Sleep(intv)
	}
}

func readHandler(filename string, rf io.ReaderFrom) error {
	if len(filename) > 0 && filename[0] == '/' {
		filename = filename[1:]
	}
	if debug {
		log.Println("Request for", filename)
	}

	if strings.HasPrefix(filename, "menu.lst") {
		if debug {
			log.Println("Redirect", filename, "-> menu.lst")
		}
		return menuLst(rf)
	}

	if filename == "bootfile" {
		if debug {
			log.Println("Redirect bootfile ->", bootfile)
		}
		filename = bootfile
	}

	if strings.HasPrefix(filename, "os/") {
		parts := strings.Split(filename, "/")
		newFilename := filepath.Join("platform-"+parts[1], filepath.Join(parts[3:]...))
		if debug {
			log.Println("Redirect", filename, "->", newFilename)
		}
		filename = newFilename
	}

	file, err := os.Open(filepath.Join(datadir, filename))
	if err != nil {
		if debug {
			log.Printf("Opening %s: %v", filename, err)
		}
		return err
	}
	defer file.Close()

	if tf, ok := rf.(tftp.OutgoingTransfer); ok {
		info, _ := file.Stat()
		if debug {
			log.Println("Size of", filename, info.Size())
		}
		tf.SetSize(info.Size())
	}

	if debug {
		log.Println("Sending", filename, "...")
	}
	n, err := rf.ReadFrom(file)
	if err != nil {
		if debug {
			log.Printf("Sending %s: %v", filename, err)
		}
		return err
	}
	if debug {
		log.Println("Sent", filename, n)
	}
	return nil
}
