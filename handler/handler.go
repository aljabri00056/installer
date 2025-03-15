package handler

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"text/template"
	"time"

	"github.com/divyam234/installer/logger"

	"github.com/divyam234/installer/handler/provider"
	"github.com/divyam234/installer/scripts"
)

const (
	cacheTTL = time.Hour

	ErrInvalidPath     = "Invalid path - must specify program name"
	ErrUnknownType     = "Unknown response type requested"
	ErrUnknownProvider = "Unknown provider specified"
	ErrProviderURL     = "Provider URL is required for Forgejo"
)

var (
	isTermRe = regexp.MustCompile(`(?i)^(curl|wget|.+WindowsPowerShell)\/`)
	errMsgRe = regexp.MustCompile(`[^A-Za-z0-9\ :\/\.]`)
)

type Query struct {
	User, Program, AsProgram, Release, Include, Arch, Token, Platform, ProviderURL string
	MoveToPath, Insecure, Private                                                  bool
}

type Result struct {
	Query
	Timestamp time.Time
	Assets    []provider.Asset
	Version   string
	M1Asset   bool
}

func (q Query) cacheKey() string {
	hw := sha256.New()
	jw := json.NewEncoder(hw)
	if err := jw.Encode(q); err != nil {
		panic(err)
	}
	return base64.StdEncoding.EncodeToString(hw.Sum(nil))
}

// Handler serves install scripts using Github releases
type Handler struct {
	Config
	cacheMut sync.Mutex
	cache    map[string]Result
}

func (h *Handler) detectProvider(path string) (provider, user string) {
	if path == "" {
		return "github", ""
	}

	parts := strings.SplitN(path, "/", 2)
	first := strings.ToLower(parts[0])

	switch first {
	case "github", "codeberg", "forgejo":
		if len(parts) > 1 {
			return first, parts[1]
		}
		return first, ""
	default:
		if len(parts) > 1 {
			repoPath := path
			if mappedProvider, ok := h.Config.RepoProviderMap[repoPath]; ok {
				return mappedProvider, path
			}
		}
		if h.Config.Provider != "" {
			return h.Config.Provider, path
		}
		return "github", path
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// calculate response type
	ext := ""
	script := ""
	qtype := r.URL.Query().Get("type")
	if qtype == "" {
		ua := r.Header.Get("User-Agent")
		switch {
		case isTermRe.MatchString(ua):
			qtype = "script"
		default:
			qtype = "text"
		}
	}
	// type specific error response
	showError := func(msg string, code int) {
		// prevent shell injection
		cleaned := errMsgRe.ReplaceAllString(msg, "")
		if qtype == "script" {
			cleaned = fmt.Sprintf("echo '%s'", cleaned)
		}
		http.Error(w, cleaned, code)
	}

	q := Query{
		User:      "",
		Program:   "",
		Release:   "",
		Insecure:  r.URL.Query().Get("insecure") == "1",
		AsProgram: r.URL.Query().Get("as"),
		Include:   r.URL.Query().Get("include"),
		Arch:      r.URL.Query().Get("arch"),
		Platform:  r.URL.Query().Get("platform"),
	}
	if q.Platform == "" {
		q.Platform = "linux"
	}
	if r.URL.Query().Get("move") == "" {
		q.MoveToPath = true
	} else {
		q.MoveToPath = r.URL.Query().Get("move") == "1"
	}
	switch qtype {
	case "script":
		if q.Platform == "windows" {
			w.Header().Set("Content-Type", "text/x-powershell")
			ext = "ps1"
			script = string(scripts.WindowsShell)
		} else {
			w.Header().Set("Content-Type", "text/x-shellscript")
			ext = "sh"
			script = string(scripts.LinuxShell)
		}
	case "text":
		w.Header().Set("Content-Type", "text/plain")
		ext = "txt"
		script = string(scripts.Text)
	default:
		showError("Unknown type", http.StatusInternalServerError)
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/")

	detectedProvider, remainingPath := h.detectProvider(path)

	switch detectedProvider {
	case "github":
		q.ProviderURL = "https://github.com"
	case "codeberg":
		q.ProviderURL = "https://codeberg.org"
	case "forgejo":
		q.ProviderURL = h.Config.ProviderURL
	default:
		showError("Unknown provider", http.StatusBadRequest)
		return
	}

	var rest string
	q.User, rest = splitHalf(remainingPath, "/")
	q.Program, q.Release = splitHalf(rest, "@")

	// no program? treat first part as program, use default user
	if q.Program == "" {
		q.Program = q.User
		q.User = h.Config.User
	}

	if h.Config.Provider != "" {
		detectedProvider = h.Config.Provider
	}
	if q.Release == "" {
		q.Release = "latest"
	}

	valid := q.Program != ""
	if !valid {
		if path == "" {
			http.Redirect(w, r, "https://github.com/divyam234/installer", http.StatusMovedPermanently)
			return
		}
		logger.Debug("invalid path: query: %+v", q)
		showError("Invalid path - must specify program name", http.StatusBadRequest)
		return
	}

	split := strings.Split(r.Header.Get("Authorization"), "Bearer ")
	if len(split) > 1 {
		q.Token = split[1]
	}
	// if q.Token == "" {
	// 	q.Token = os.Getenv("GITHUB_TOKEN")
	// }
	provider, err := provider.NewProvider(detectedProvider, h.Config.ProviderURL)
	if err != nil {
		showError(err.Error(), http.StatusBadRequest)
		return
	}
	res, err := provider.GetRepo(q.User, q.Program, q.Token)
	if err != nil {
		showError(err.Error(), http.StatusBadRequest)
		return
	}
	q.Private = res.Private
	result, err := h.execute(provider, q)
	if err != nil {
		showError(err.Error(), http.StatusBadGateway)
		return
	}
	t, err := template.New("installer").Parse(script)
	if err != nil {
		showError("installer BUG: "+err.Error(), http.StatusInternalServerError)
		return
	}
	buff := bytes.Buffer{}
	if err := t.Execute(&buff, result); err != nil {
		showError("Template error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	logger.Info("serving script %s/%s@%s (%s)", q.User, q.Program, q.Release, ext)
	w.Write(buff.Bytes())
}
