smartos-tftp
============

A simple TFTP/PXE server tailor made for booting SmartOS.

Installation
------------

Download the [latest relase](https://github.com/calmh/smartos-pxe/releases).
Untar it and place somewhere suitable, such as /opt/smartos-pxe.

Usage
-----

Start the `smartos-pxe` binary from the installation dir, or set `-data-dir`
to point at the `data` directory from the distribution. By default
`smartos-pxe` automatically downloads and serves the latest platform image
from Joyent, although this can be disabled. You can also download an unpack
platofrm images manually into the `data` dir.

```text
# cd /opt/smartos-pxe
# ./bin/smartos-pxe -verbose
21:50:17 main.go:61: Listening on :69
21:50:19 upgrade.go:41: Fetching https://us-east.manta.joyent.com/Joyent_Dev/public/SmartOS/20160622T220759Z/platform-20160622T220759Z.tgz
21:50:19 upgrade.go:92: Unpacking platform-20160622T220759Z/i86pc/amd64/boot_archive.hash
21:50:19 upgrade.go:92: Unpacking platform-20160622T220759Z/i86pc/amd64/boot_archive.gitstatus
...
```

Use `-help` to see the available options:

```text
# ./bin/smartos-pxe -help
Usage of smartos-pxe:
  -boot-file string
      Boot file (within data-dir) (default "grub/pxegrub")
  -data-dir string
      Data directory (default "./data")
  -download-intv duration
      New platform download interval (0 to disable) (default 24h0m0s)
  -download-latest-path string
      Path to latest platform indicator file (default "/Joyent_Dev/public/SmartOS/latest")
  -download-server string
      Platform download server (default "https://us-east.manta.joyent.com")
  -grub-console string
      GRUB os_console device (default "ttyS2")
  -grub-timeout duration
      GRUB menu timeout (default 10s)
  -listen string
      TFTP listen address (default ":69")
  -prefer string
      Preferred platform version (default "latest")
  -root-pw string
      Root password hash (default "$5$5x85uZWD$AQUMEs1UiMwXcjWjYopG2cMUm/eAoFxtjWiHokw7SL.")
  -verbose
      Verbose output
```

By default the latest available platform version will be preferred. To lock
it down to a specific version, set it using `-prefer`:

```text
# ./bin/smartos-pxe -prefer 20160622T220759Z
```


DHCP Configuration
------------------

The DHCP server must be configured to point at your PXE server. In ISC dhcpd
the relevant options are:

```
next-server 192.0.2.0;
filename "bootfile";
```

where `192.0.2.0` is the IP where `smartos-pxe` is running. Other servers
will have different but similar configuration options.

License
-------

MIT
