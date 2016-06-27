package main

import (
	"bytes"
	"html/template"
	"io"
	"os"
	"path/filepath"
	"sort"
)

const menuLstTpl = `default 0
timeout {{.Timeout}}
serial --speed=115200 --unit=1 --word=8 --parity=no --stop=1
terminal composite
variable os_console {{.Console}}

{{range .Versions}}
title SmartOS ({{.}})
    kernel$ /os/{{.}}/platform/i86pc/kernel/amd64/unix -B console=${os_console},${os_console}-mode="115200,8,n,1,-",smartos=true,root_shadow={{$.RootShadow}}
    module$ /os/{{.}}/platform/i86pc/amd64/boot_archive type=rootfs name=ramdisk
    module$ /os/{{.}}/platform/i86pc/amd64/boot_archive.hash type=hash name=ramdisk
{{- range $.Overlay}}
    module$ /overlay/{{.}} type=file name={{.}}
{{- end}}

title SmartOS ({{.}}) noinstall
    kernel$ /os/{{.}}/platform/i86pc/kernel/amd64/unix -B console=${os_console},${os_console}-mode="115200,8,n,1,-",noimport=true,root_shadow={{$.RootShadow}}
    module$ /os/{{.}}/platform/i86pc/amd64/boot_archive type=rootfs name=ramdisk
    module$ /os/{{.}}/platform/i86pc/amd64/boot_archive.hash type=hash name=ramdisk
{{- range $.Overlay}}
    module$ /overlay/{{.}} type=file name={{.}}
{{- end}}
{{end}}
`

func menuLst(rf io.ReaderFrom) error {
	versions, err := versions()
	if err != nil {
		return err
	}
	overlay, err := overlay()
	if err != nil {
		return err
	}
	st := state{
		Console:    console,
		RootShadow: rootShadow,
		Versions:   versions,
		Overlay:    overlay,
		Timeout:    int(timeout.Seconds()),
	}

	buf := new(bytes.Buffer)
	tpl := template.Must(template.New("menu.lst").Parse(menuLstTpl))
	if err := tpl.Execute(buf, st); err != nil {
		return err
	}
	_, err = rf.ReadFrom(buf)
	return err
}

func versions() ([]string, error) {
	pattern := filepath.Join(datadir, "platform-*")
	files, err := filepath.Glob(pattern)
	if err != nil {
		return nil, err
	}
	for i := range files {
		files[i] = files[i][len(pattern)-1:]
	}
	sort.Sort(versionList(files))
	return files, nil
}

func overlay() ([]string, error) {
	dir := filepath.Join(datadir, "overlay")
	var files []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		files = append(files, path[len(dir)+1:])
		return nil
	})
	if err != nil {
		return nil, err
	}
	return files, nil
}

type versionList []string

func (v versionList) Len() int {
	return len(v)
}

func (v versionList) Less(a, b int) bool {
	if v[a] == preferredVersion {
		return true
	}
	if v[b] == preferredVersion {
		return false
	}
	return b < a // Reverse sorting to get highest version at top
}

func (v versionList) Swap(a, b int) {
	v[a], v[b] = v[b], v[a]
}
