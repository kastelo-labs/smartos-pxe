smartos-tftp
============

A simple TFTP/PXE server tailor made for booting SmartOS.

Installation
------------

Download the [latest relase](https://github.com/calmh/smartos-pxe/releases).
Untar it into a suitable place, such as /opt/smartos-pxe.

Usage
-----

```bash
# cd /opt/smartos-pxe
# ./bin/smartos-pxe -debug
21:50:17 main.go:61: Listening on :69
21:50:19 upgrade.go:41: Fetching https://us-east.manta.joyent.com/Joyent_Dev/public/SmartOS/20160622T220759Z/platform-20160622T220759Z.tgz
21:50:19 upgrade.go:92: Unpacking platform-20160622T220759Z/i86pc/amd64/boot_archive.hash
21:50:19 upgrade.go:92: Unpacking platform-20160622T220759Z/i86pc/amd64/boot_archive.gitstatus
21:50:19 upgrade.go:92: Unpacking platform-20160622T220759Z/i86pc/amd64/boot_archive.manifest
...
```

You can modify a number of parameters and tweak or disable automatic
downloads of new platforms:

```bash
# smartos-pxe --help
Usage of smartos-pxe:
  -boot-file string
    	Boot file (default "grub/pxegrub")
  -console string
    	Console device (default "ttyS2")
  -data-dir string
    	Data directory (default "data")
  -debug
    	Debug output
  -grub-timeout duration
    	GRUB menu timeout (default 10s)
  -listen string
    	Listen address (default ":69")
  -root-pw string
    	Root password hash (default "$5$hEQ0l8d5$s0Jwt.oif76hVoQpzsgH2XVKhS8uCXnMlQhXltYgvaB")
  -upgrade-intv duration
    	New platform upgrade interval (0 to disable) (default 24h0m0s)
  -upgrade-path string
    	Path to latest platform file (default "/Joyent_Dev/public/SmartOS/latest")
  -upgrade-server string
    	Platform upgrade server (default "https://us-east.manta.joyent.com")
```

License
-------

MIT
