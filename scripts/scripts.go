package scripts

import _ "embed"

//go:embed install.txt.tmpl
var Text []byte

//go:embed install.sh.tmpl
var LinuxShell []byte

//go:embed install.ps1.tmpl
var WindowsShell []byte
