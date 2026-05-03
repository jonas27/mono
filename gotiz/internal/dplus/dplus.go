// Package dplus downloads videos from Discovery+ using a headless Chromium browser.
package dplus

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
)

// Config holds Discovery+ credentials loaded from .discovery.yaml.
type Config struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Token    string `yaml:"token"` // unused; kept for config compatibility
}

// Client resolves Discovery+ video pages via a headless browser.
type Client struct {
	cfg Config
}

// New returns a Client using the given config.
func New(cfg Config) *Client {
	return &Client{cfg: cfg}
}

// VideoInfo holds the data needed to download a single Discovery+ video.
type VideoInfo struct {
	Title        string
	VideoID      string
	StreamURL    string
	Kind         string
	DRMProtected bool
}

// Resolve opens the Discovery+ video page in a headless browser, handles login,
// and intercepts the playback API response to return a VideoInfo.
func (c *Client) Resolve(ctx context.Context, pageURL string) (VideoInfo, error) {
	binPath := launcher.NewBrowser().BinPath()
	l := launcher.New().
		Bin(binPath).
		Headless(true).
		Set("no-sandbox").
		Set("disable-setuid-sandbox").
		Set("disable-dev-shm-usage").
		Set("no-first-run").
		Set("no-default-browser-check")

	controlURL, err := l.Launch()
	if err != nil {
		return VideoInfo{}, fmt.Errorf("launch browser: %w", err)
	}

	browser := rod.New().ControlURL(controlURL).NoDefaultDevice()
	if err := browser.Connect(); err != nil {
		return VideoInfo{}, fmt.Errorf("connect browser: %w", err)
	}
	defer browser.Close()

	page, err := browser.Page(proto.TargetCreateTarget{URL: "about:blank"})
	if err != nil {
		return VideoInfo{}, fmt.Errorf("open page: %w", err)
	}

	if _, err := page.EvalOnNewDocument(`Object.defineProperty(navigator,'webdriver',{get:()=>undefined})`); err != nil {
		return VideoInfo{}, fmt.Errorf("inject stealth script: %w", err)
	}

	type apiResult struct {
		manifestURL string
		err         error
	}
	resultCh := make(chan apiResult, 1)
	titleCh := make(chan string, 1)

	router := page.HijackRequests()

	// Intercept playback API (POST only, skip OPTIONS preflight).
	if err := router.Add("*/playback*", "", func(h *rod.Hijack) {
		if h.Request.Method() != "POST" {
			h.ContinueRequest(&proto.FetchContinueRequest{})
			return
		}
		h.MustLoadResponse()
		body := h.Response.Body()
		if len(body) < 10 {
			h.ContinueRequest(&proto.FetchContinueRequest{})
			return
		}
		manifestURL, err := extractManifestURL([]byte(body))
		select {
		case resultCh <- apiResult{manifestURL: manifestURL, err: err}:
		default:
		}
	}); err != nil {
		return VideoInfo{}, fmt.Errorf("setup playback intercept: %w", err)
	}

	// Intercept CMS routes to extract video title.
	if err := router.Add("*/cms/routes/video/*", "", func(h *rod.Hijack) {
		if h.Request.Method() != "GET" {
			h.ContinueRequest(&proto.FetchContinueRequest{})
			return
		}
		h.MustLoadResponse()
		body := h.Response.Body()
		if title := extractTitle([]byte(body)); title != "" {
			select {
			case titleCh <- title:
			default:
			}
		}
	}); err != nil {
		return VideoInfo{}, fmt.Errorf("setup title intercept: %w", err)
	}

	go router.Run()

	fmt.Printf("  opening browser for: %s\n", pageURL)

	// Login first, then navigate to video.
	if err := c.login(ctx, page); err != nil {
		return VideoInfo{}, fmt.Errorf("login: %w", err)
	}

	if err := page.Navigate(pageURL); err != nil {
		return VideoInfo{}, fmt.Errorf("navigate to video: %w", err)
	}

	// Wait up to 60 s for the playback API to fire.
	var manifestURL string
	select {
	case r := <-resultCh:
		router.Stop()
		if r.err != nil {
			return VideoInfo{}, r.err
		}
		manifestURL = r.manifestURL
	case <-time.After(60 * time.Second):
		return VideoInfo{}, fmt.Errorf("timed out waiting for playback API response")
	case <-ctx.Done():
		return VideoInfo{}, ctx.Err()
	}

	// Get title (best-effort: CMS response or page title).
	title := ""
	select {
	case t := <-titleCh:
		title = t
	default:
		if t := pageTitle(page); t != "" {
			title = t
		}
	}
	if title == "" {
		title = "discovery-plus-video"
	}

	return VideoInfo{
		Title:     title,
		StreamURL: manifestURL,
		Kind:      "dash",
	}, nil
}

// login navigates to the auth page and logs in using shadow DOM form filling.
func (c *Client) login(ctx context.Context, page *rod.Page) error {
	if err := page.Navigate("https://auth.discoveryplus.com/login?flow=login"); err != nil {
		return fmt.Errorf("navigate to login: %w", err)
	}
	page.WaitLoad() //nolint:errcheck
	time.Sleep(4 * time.Second)

	// Check if already logged in (redirected away from auth page).
	if info := page.MustInfo(); !strings.Contains(info.URL, "auth.discoveryplus.com") {
		return nil
	}

	if c.cfg.Username == "" || c.cfg.Password == "" {
		return fmt.Errorf("login required but no credentials; add 'username' and 'password' to .discovery.yaml")
	}

	fmt.Println("  logging in...")

	// Fill email via shadow DOM.
	if err := shadowFill(page, `input[type="email"]`, c.cfg.Username); err != nil {
		return fmt.Errorf("fill email: %w", err)
	}
	// Try both button click and Enter key press.
	shadowClick(page, `button[type="submit"]`)
	shadowPressEnter(page, `input[type="email"]`)
	time.Sleep(3 * time.Second)

	// Wait for password field (up to 20 s).
	passwordFound := false
	for i := 0; i < 20; i++ {
		// Check if we got redirected away from auth (SSO or already done).
		if info := page.MustInfo(); !strings.Contains(info.URL, "auth.discoveryplus.com") {
			fmt.Println("  login successful (redirected after email)")
			return nil
		}
		found, _ := shadowExists(page, `input[type="password"]`)
		if found {
			passwordFound = true
			break
		}
		fmt.Printf("  waiting for password field (%d/20)...\n", i+1)
		time.Sleep(1 * time.Second)
	}

	if !passwordFound {
		return fmt.Errorf("password field did not appear after email submission (URL: %s)", page.MustInfo().URL)
	}

	if err := shadowFill(page, `input[type="password"]`, c.cfg.Password); err != nil {
		return fmt.Errorf("fill password: %w", err)
	}
	shadowClick(page, `button[type="submit"]`)

	// Wait for redirect away from auth page (up to 30 s).
	deadline := time.Now().Add(30 * time.Second)
	for time.Now().Before(deadline) {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		time.Sleep(1 * time.Second)
		if info := page.MustInfo(); !strings.Contains(info.URL, "auth.discoveryplus.com") {
			fmt.Println("  login successful")
			return nil
		}
	}
	return fmt.Errorf("timed out waiting for login redirect")
}

// shadowFill fills a form field found inside shadow DOM subtrees.
func shadowFill(page *rod.Page, selector, value string) error {
	js := fmt.Sprintf(`() => {
		function findInShadow(root, sel) {
			let found = root.querySelector(sel);
			if (found) return found;
			for (let el of root.querySelectorAll("*")) {
				if (el.shadowRoot) {
					found = findInShadow(el.shadowRoot, sel);
					if (found) return found;
				}
			}
			return null;
		}
		let el = findInShadow(document, %q);
		if (!el) return "notfound";
		el.focus();
		let setter = Object.getOwnPropertyDescriptor(window.HTMLInputElement.prototype, "value").set;
		setter.call(el, %q);
		el.dispatchEvent(new Event("input", {bubbles: true}));
		el.dispatchEvent(new Event("change", {bubbles: true}));
		return "ok";
	}`, selector, value)
	r, err := page.Eval(js)
	if err != nil {
		return err
	}
	if r != nil && strings.Contains(r.Value.String(), "notfound") {
		return fmt.Errorf("element not found: %s", selector)
	}
	return nil
}

// shadowClick clicks the first element matching selector in shadow DOM subtrees.
func shadowClick(page *rod.Page, selector string) {
	js := fmt.Sprintf(`() => {
		function findInShadow(root, sel) {
			let found = root.querySelector(sel);
			if (found) return found;
			for (let el of root.querySelectorAll("*")) {
				if (el.shadowRoot) {
					found = findInShadow(el.shadowRoot, sel);
					if (found) return found;
				}
			}
			return null;
		}
		let el = findInShadow(document, %q);
		if (el) el.click();
	}`, selector)
	page.Eval(js) //nolint:errcheck
}

// shadowPressEnter dispatches Enter keydown/keyup events on a shadow DOM element.
func shadowPressEnter(page *rod.Page, selector string) {
	js := fmt.Sprintf(`() => {
		function findInShadow(root, sel) {
			let found = root.querySelector(sel);
			if (found) return found;
			for (let el of root.querySelectorAll("*")) {
				if (el.shadowRoot) {
					found = findInShadow(el.shadowRoot, sel);
					if (found) return found;
				}
			}
			return null;
		}
		let el = findInShadow(document, %q);
		if (el) {
			el.dispatchEvent(new KeyboardEvent("keydown", {bubbles:true, key:"Enter", keyCode:13}));
			el.dispatchEvent(new KeyboardEvent("keyup",  {bubbles:true, key:"Enter", keyCode:13}));
		}
	}`, selector)
	page.Eval(js) //nolint:errcheck
}

// shadowExists checks if a shadow DOM element exists.
func shadowExists(page *rod.Page, selector string) (bool, error) {
	js := fmt.Sprintf(`() => {
		function findInShadow(root, sel) {
			let found = root.querySelector(sel);
			if (found) return found;
			for (let el of root.querySelectorAll("*")) {
				if (el.shadowRoot) {
					found = findInShadow(el.shadowRoot, sel);
					if (found) return found;
				}
			}
			return null;
		}
		return findInShadow(document, %q) ? "found" : "notfound";
	}`, selector)
	r, err := page.Eval(js)
	if err != nil {
		return false, err
	}
	return r != nil && strings.Contains(r.Value.String(), "found"), nil
}

// extractManifestURL parses the new Bolt v2 playbackInfo response.
func extractManifestURL(body []byte) (string, error) {
	var pr struct {
		Manifest struct {
			URL    string `json:"url"`
			Format string `json:"format"`
		} `json:"manifest"`
		Errors []struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
		} `json:"errors"`
		Status  int    `json:"status"`
		Message string `json:"message"`
	}
	if err := json.Unmarshal(body, &pr); err != nil {
		return "", fmt.Errorf("decode playback response: %w", err)
	}

	if pr.Status != 0 && pr.Status != 200 {
		return "", fmt.Errorf("playback API error %d: %s", pr.Status, pr.Message)
	}

	if len(pr.Errors) > 0 {
		return "", fmt.Errorf("playback API error %d: %s", pr.Errors[0].Code, pr.Errors[0].Message)
	}

	if pr.Manifest.URL == "" {
		return "", fmt.Errorf("no manifest URL in playback response")
	}

	return pr.Manifest.URL, nil
}

// extractTitle extracts the video title from a CMS routes API response.
func extractTitle(body []byte) string {
	var root map[string]json.RawMessage
	if err := json.Unmarshal(body, &root); err != nil {
		return ""
	}

	// Try common paths: data.attributes.title or data[0].attributes.title
	type attrWrapper struct {
		Attributes struct {
			Title string `json:"title"`
			Name  string `json:"name"`
		} `json:"attributes"`
	}

	// data could be object or array
	if raw, ok := root["data"]; ok {
		var single attrWrapper
		if json.Unmarshal(raw, &single) == nil && single.Attributes.Title != "" {
			return single.Attributes.Title
		}
		if single.Attributes.Name != "" {
			return single.Attributes.Name
		}
		var arr []attrWrapper
		if json.Unmarshal(raw, &arr) == nil {
			for _, a := range arr {
				if a.Attributes.Title != "" {
					return a.Attributes.Title
				}
			}
		}
	}

	return ""
}

// pageTitle returns the cleaned page title (strips common suffixes).
func pageTitle(page *rod.Page) string {
	r, err := page.Eval(`() => document.title`)
	if err != nil || r == nil {
		return ""
	}
	title := strings.Trim(r.Value.String(), `"`)
	for _, suffix := range []string{" | Discovery+", " - Discovery+", " | Eurosport", " - Eurosport"} {
		if idx := strings.LastIndex(title, suffix); idx > 0 {
			title = title[:idx]
		}
	}
	return strings.TrimSpace(title)
}

// Download downloads the given VideoInfo to dir using ffmpeg.
func Download(ctx context.Context, info VideoInfo, dir string) error {
	filename := sanitize(info.Title) + ".mp4"
	dest := filepath.Join(dir, filename)

	if fi, err := os.Stat(dest); err == nil && fi.Size() > 0 {
		fmt.Printf("  already complete: %s\n", filename)
		return nil
	}

	fmt.Printf("downloading: %s\n", info.Title)

	args := []string{
		"-reconnect", "1",
		"-reconnect_streamed", "1",
		"-reconnect_delay_max", "5",
		"-i", info.StreamURL,
		"-c", "copy",
		"-y",
		dest,
	}

	cmd := exec.CommandContext(ctx, "ffmpeg", args...) //nolint:gosec
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		return fmt.Errorf("ffmpeg: %w", err)
	}

	fmt.Printf("  saved to %s\n", filename)
	return nil
}

func sanitize(s string) string {
	var b strings.Builder
	for _, r := range s {
		switch r {
		case '/', '\\', ':', '*', '?', '"', '<', '>', '|':
			b.WriteRune('_')
		default:
			b.WriteRune(r)
		}
	}
	return b.String()
}
