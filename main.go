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
	bootfile           = "grub/pxegrub"
	console            = "ttyS2"
	datadir            = "./data"
	listen             = ":69"
	rootShadow         = "$5$5x85uZWD$AQUMEs1UiMwXcjWjYopG2cMUm/eAoFxtjWiHokw7SL."
	timeout            = 10 * time.Second
	downloadInterval   = 24 * time.Hour
	downloadLatestPath = "/Joyent_Dev/public/SmartOS/latest"
	downloadServer     = "https://us-east.manta.joyent.com"
	verbose            bool
)

func main() {
	flag.StringVar(&bootfile, "boot-file", bootfile, "Boot file (within data-dir)")
	flag.StringVar(&console, "grub-console", console, "GRUB os_console device")
	flag.StringVar(&datadir, "data-dir", datadir, "Data directory")
	flag.StringVar(&listen, "listen", listen, "TFTP listen address")
	flag.StringVar(&rootShadow, "root-pw", rootShadow, "Root password hash")
	flag.DurationVar(&timeout, "grub-timeout", timeout, "GRUB menu timeout")
	flag.DurationVar(&downloadInterval, "download-intv", downloadInterval, "New platform download interval (0 to disable)")
	flag.StringVar(&downloadLatestPath, "download-latest-path", downloadLatestPath, "Path to latest platform indicator file")
	flag.StringVar(&downloadServer, "download-server", downloadServer, "Platform download server")
	flag.BoolVar(&verbose, "verbose", verbose, "Verbose output")
	flag.Parse()

	log.SetOutput(os.Stdout)
	if verbose {
		log.SetFlags(log.Lshortfile | log.LUTC | log.Ltime)
	} else {
		log.SetFlags(0)
	}

	if downloadInterval > 0 {
		go autoDownload(downloadInterval)
	}

	if verbose {
		log.Println("Listening on", listen)
	}
	s := tftp.NewServer(readHandler, nil)
	s.SetTimeout(60 * time.Second)
	if err := s.ListenAndServe(listen); err != nil {
		log.Fatal(err)
	}
}

func autoDownload(intv time.Duration) {
	for {
		if err := downloadPlatform(); err != nil {
			log.Println("Downloading new platform:", err)
		}
		time.Sleep(intv)
	}
}

func readHandler(filename string, rf io.ReaderFrom) error {
	if len(filename) > 0 && filename[0] == '/' {
		filename = filename[1:]
	}
	if verbose {
		log.Println("Request for", filename)
	}

	if strings.HasPrefix(filename, "menu.lst") {
		if verbose {
			log.Println("Redirect", filename, "-> menu.lst")
		}
		return menuLst(rf)
	}

	if filename == "bootfile" {
		if verbose {
			log.Println("Redirect bootfile ->", bootfile)
		}
		filename = bootfile
	}

	if strings.HasPrefix(filename, "os/") {
		parts := strings.Split(filename, "/")
		newFilename := filepath.Join("platform-"+parts[1], filepath.Join(parts[3:]...))
		if verbose {
			log.Println("Redirect", filename, "->", newFilename)
		}
		filename = newFilename
	}

	file, err := os.Open(filepath.Join(datadir, filename))
	if err != nil {
		if verbose {
			log.Printf("Opening %s: %v", filename, err)
		}
		return err
	}
	defer file.Close()

	if tf, ok := rf.(tftp.OutgoingTransfer); ok {
		info, _ := file.Stat()
		if verbose {
			log.Println("Size of", filename, info.Size())
		}
		tf.SetSize(info.Size())
	}

	if verbose {
		log.Println("Sending", filename, "...")
	}
	n, err := rf.ReadFrom(file)
	if err != nil {
		if verbose {
			log.Printf("Sending %s: %v", filename, err)
		}
		return err
	}
	if verbose {
		log.Println("Sent", filename, n)
	}
	return nil
}
