package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
)

func shadowFill(page *rod.Page, selector, value string) {
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
		let nativeInputValueSetter = Object.getOwnPropertyDescriptor(window.HTMLInputElement.prototype, "value").set;
		nativeInputValueSetter.call(el, %q);
		el.dispatchEvent(new Event("input", {bubbles: true}));
		el.dispatchEvent(new Event("change", {bubbles: true}));
	}`, selector, value)
	page.Eval(js) //nolint:errcheck
}

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

func main() {
	binPath := launcher.NewBrowser().BinPath()
	l := launcher.New().Bin(binPath).Headless(true).
		Set("no-sandbox").Set("disable-setuid-sandbox").Set("disable-dev-shm-usage")

	controlURL, _ := l.Launch()
	browser := rod.New().ControlURL(controlURL).NoDefaultDevice()
	browser.Connect() //nolint:errcheck
	defer browser.Close()

	page, _ := browser.Page(proto.TargetCreateTarget{URL: "about:blank"})
	page.EvalOnNewDocument(`Object.defineProperty(navigator,'webdriver',{get:()=>undefined})`) //nolint:errcheck

	router := page.HijackRequests()
	manifestCh := make(chan string, 1)

	router.Add("*", "", func(h *rod.Hijack) { //nolint:errcheck
		url := h.Request.URL().String()
		host := h.Request.URL().Hostname()

		if (strings.Contains(host, "discoveryplus.com") || strings.Contains(host, "dplus") || strings.Contains(host, "discomax")) &&
			(strings.Contains(url, "playbackInfo") || strings.Contains(url, "playback/v1")) &&
			h.Request.Method() == "POST" {
			h.MustLoadResponse()
			body := h.Response.Body()
			// Extract manifest URL
			manifestURL := ""
			if idx := strings.Index(body, `"manifest":{"availEnd"`); idx >= 0 {
				rest := body[idx:]
				if uidx := strings.Index(rest, `"url":"`); uidx >= 0 {
					rest2 := rest[uidx+7:]
					if end := strings.Index(rest2, `"`); end >= 0 {
						manifestURL = rest2[:end]
					}
				}
			}
			if manifestURL == "" {
				// Try simpler extraction
				if idx := strings.Index(body, `"url":"`); idx >= 0 {
					rest := body[idx+7:]
					if end := strings.Index(rest, `"`); end >= 0 {
						u := rest[:end]
						if strings.Contains(u, "media.max.com") || strings.Contains(u, "akamai") {
							manifestURL = u
						}
					}
				}
			}
			select {
			case manifestCh <- manifestURL:
			default:
			}
			return
		}
		h.ContinueRequest(&proto.FetchContinueRequest{})
	})
	go router.Run()

	// Login
	page.Navigate("https://auth.discoveryplus.com/login?flow=login") //nolint:errcheck
	page.WaitLoad()                                                   //nolint:errcheck
	time.Sleep(4 * time.Second)

	shadowFill(page, `input[type="email"]`, "jonas.burster@gmail.com")
	shadowClick(page, `button[type="submit"]`)
	time.Sleep(3 * time.Second)

	for i := 0; i < 8; i++ {
		js := `() => {
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
			let el = findInShadow(document, 'input[type="password"]');
			return el ? el.tagName : "notfound";
		}`
		r, _ := page.Eval(js)
		if r != nil && !strings.Contains(r.Value.String(), "notfound") {
			break
		}
		time.Sleep(1 * time.Second)
	}

	shadowFill(page, `input[type="password"]`, "uS1Ov5Ow5GplpseQoWV")
	shadowClick(page, `button[type="submit"]`)
	time.Sleep(8 * time.Second)

	fmt.Println("Navigating to video...")
	page.Navigate("https://play.discoveryplus.com/video/watch-sport/79fa53e3-2fa0-5e7e-80bb-12809aae497f/b20eb0ce-ca83-4684-b6fa-c90dfb111ec8") //nolint:errcheck

	select {
	case manifestURL := <-manifestCh:
		router.Stop()
		fmt.Printf("Manifest URL: %s\n\n", manifestURL)

		if manifestURL == "" {
			fmt.Println("No manifest URL captured")
			return
		}

		// Try downloading with ffmpeg - 30 second test
		fmt.Println("Attempting ffmpeg download (30 seconds)...")
		outFile := "/tmp/test_dplus.mp4"
		os.Remove(outFile)

		cmd := exec.Command("ffmpeg",
			"-reconnect", "1",
			"-reconnect_streamed", "1",
			"-reconnect_delay_max", "5",
			"-t", "30", // only 30 seconds
			"-i", manifestURL,
			"-c", "copy",
			"-y",
			outFile,
		)
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			fmt.Printf("ffmpeg error: %v\n", err)
		}

		// Check if output file has content
		fi, statErr := os.Stat(outFile)
		if statErr == nil && fi.Size() > 0 {
			fmt.Printf("SUCCESS! Output file: %d bytes\n", fi.Size())
			// Check if it's valid
			probe := exec.Command("ffprobe", "-v", "error", "-show_entries", "format=duration", "-of", "csv=p=0", outFile)
			probe.Stdout = os.Stdout
			probe.Stderr = os.Stderr
			probe.Run() //nolint:errcheck
		} else {
			fmt.Println("No output file - content is DRM encrypted")
		}
	case <-time.After(60 * time.Second):
		fmt.Println("Timeout waiting for playback API")
	}
}
