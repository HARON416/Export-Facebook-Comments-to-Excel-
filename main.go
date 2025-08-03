package main

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
)

var (
	globalPostID  string
	globalPostURL string
)

type FacebookConfig struct {
	XFbDebug              string
	CSPNonce              string
	LSDToken              string
	Cookies               string
	ExpansionTokens       map[string]string
	LastUpdated           time.Time
	CommentsUserID        string
	CommentsRequestID     string
	CommentsHash          string
	CommentsRevision      string
	CommentsSession       string
	CommentsHashSessionID string
	CommentsDynamic       string
	CommentsCSR           string
	CommentsHSDP          string
	CommentsHBLP          string
	CommentsSJSP          string
	CommentsFbDtsg        string
	CommentsJazoest       string
	CommentsSpinR         string
	CommentsSpinT         string
}

func CountFacebookComments(comments []FacebookComment) int {
	mainComments := 0
	for _, comment := range comments {
		if comment.CommentParent == nil && comment.CommentDirectParent == nil {
			mainComments++
		}
	}
	return mainComments
}

func getDefaultFacebookConfig() *FacebookConfig {
	return &FacebookConfig{
		XFbDebug:              "",
		CSPNonce:              "",
		LSDToken:              "3bL9I7RJqQAmvrzk8ckc3m",
		Cookies:               "sb=XOnLZzjZUznCGJQE89_qDE7G; ps_l=1; ps_n=1; datr=7qiMaDwDu83IF3sgoxhRvmsx; locale=en_GB; c_user=61555590462485; dpr=1.5; wd=1366x681; fr=1PIUJUKzs5V8kf9PV.AWcHvGKPZ2BAaYwrgQbDVnZAaiOzGv8j5uVVJ3TPW31Q6Tcm3pY.Boj3Er..AAA.0.0.Boj3Er.AWcy7vshu5ahtNVjDqYHtlh5rfM; xs=19%3AD6hlVJsfNlxziw%3A2%3A1754048966%3A-1%3A-1%3A%3AAcUQ5p0cKnQfIkbPMjDSgSEcRHBKqdVqRcsqM8SovA; presence=C%7B%22t3%22%3A%5B%5D%2C%22utc3%22%3A1754231368846%2C%22v%22%3A1%7D",
		ExpansionTokens:       make(map[string]string),
		LastUpdated:           time.Now(),
		CommentsUserID:        "61555590462485",
		CommentsRequestID:     "15",
		CommentsHash:          "20303.HYP%3Acomet_pkg.2.1...0",
		CommentsRevision:      "1025463513",
		CommentsSession:       "2a2sph%3A9ep293%3A64cc7d",
		CommentsHashSessionID: "7534366326965440049",
		CommentsDynamic:       "7xeUjGU5a5Q1ryaxG4Vp41twWwIxu13wFwhUngS3q2ibwNw9G2Saw8i2S1DwUx60GE3Qwb-q7oc81EEc87m221Fwgo9oO0-E4a3a4oaEnxO0Bo7O2l2Utwqo31wiE4u9x-3m1mzXw8W58jwGzEaE5e3ym2SU4i5oe8464-5pUfEe88o4Wm7-2K0-obUG2-azqwaW223908O3216xi4UK2K2WEjxK2B08-269wkopg6C13xecwBwWwjHDzUiBG2OUqwjVqwLwHwa211wo83KwHwOG8xG6E",
		CommentsCSR:           "g4Zstk4RnsAdMkbb5ED25hneOIbNibfOjEl4nRsDTV0FP8BFh2QWRFVmjnjinHC8KyKATcLihkEKjDGN6Fahm-VlAhAim8SiVmuim88GBUOXhkFEWhdd7BLh8GkyBUCVkegy7p-XFVVqVVFoOF-bXz9ooXU9EhzoyjDFpbDG8xG4oOECum6UWdK5VoiWjGAjwPx25p8Cm6rzVUm-6Vbxa4o-4E-78ao8k4oG2S4V8dE8oSEaUaUyeyFbwBwAGAm7Vp8C2yqvxW5842UgG5oW3udzEgyo9UbFV8Tz8pwjaxu0zo9ecwRBgbUkU8U5q4Upxy3e1Zxi1iwkUxe48nwt99oiyqU4i0XEvJ5m3S1gg8o7q0KU2WxN2o3fxWah6Q5E4i0Xk4V88m6o8UC8wJxt5w9G0S8dk2y22E6a0I8Ojwu8W0xUmwpEcF8lyQU9EbEjwJw0Nxw0Fow1te0ZQ04mE1Io6i05IU02oVw0Raw2ko9U3ZU1uE0Ye0Ho0gVwmJ01lV022i0Ew3l8O085w7MwJw1E6t04pg0DK09-w4YIM0dAi02zE3fw76w1g-gE0l_w61w0DrK0DUG0aAFw1e-0g8wow",
		CommentsHSDP:          "g4pC32agf85ai8F1Yy78y2WkV1n2GcFOPE4kiB4bC494QgyqaN4gyAf4qsgyvrkkzFtEPtEOfqi8AyiAl8gH34y3QlkxWsmxsy48ImwCxsO6jMjGb5hy16h8A9E9Ekh9MKVpp1hMWyW7sikLeYDitBGoFAghAyO6ye7YExy4BLWBgAwlz93tKaPyFoKaK4GHFQm9tn9bbhet2p91522qcldEFegzaGsPAEhPHojrEwGiqAichGe69rh4xkEhgyXCyGK2y8ACjc6my66AQuXyyGpDxmaJEN8WdypyFdoB6mHAxoF28ZAzEGEqc2d3WDwlUynxt7yEgwzgMxqUbOFy9FuF1oM8XyUjwBszCxCEB3FQUW3C6UG4AiQ7ogCiySUFt5HCw8aVV8no8qVU462miF8hyoF2Oxe8K2h0CwZo4hwCyAWK3G1xwio4C0OUy0iG5A5UnG5FHBCc2q12zQ3Z7xu1g7d0ywjWgcqBgmyyFtQaCy2wEwpA3i262l4oexwMAu7Cu0xEfbwiofawhU5G3C17wcy3u1jweG0SkJ0Zw7Twbi1ywEg1y81ko4S540_o3Ww5wwyxW0wUfpo1IHwhU4m0IE6m0Bo3Dwde0SE3OwvE0gywho560PU0xO0W83tx60gO0NE",
		CommentsHBLP:          "0rUqhogjwiUeomoy1lwJwMwiGxS0BU767U5m1kwaq222aq1vCwau3a1RUrxa3F3E8obUozEjgCawxwk87i2WE3izp8bEhximi1WU4ym3q2O6EdEnwc65EmBw9au10z84mi7eGxq5898yU2awhU9Yw1Po428g1dE5K4WzUlwbi1xwzwwwNwuFEdoW1nwvF6E2bwcu0youzoqwmE1iHwPyE1b9Efo8o2nwCwKw48wba3C2GbwRxG1yw5TwWwaq3K1KwoA0VEmxWawd20zo3iwyxR08e3Sm0U8lx-1VwxK17who2Owprw960VU7-1jwrEe89UGdwd24UK1ixZ6xy040o4m8wTx28BUG0G8G0gO0gyq17wlU4x2o3ix6483Xw9i3ucwmrw8Sq8xS1xwnoy6o10982YBlwCwNw",
		CommentsSJSP:          "g4pC32agf85ai8F1Yy78y92ikV1n2GcFOEO10j9FjLCch6hdd2y9n8DAZ9KGYKSLrGzFS_qlJFbX9jGO98y8uOBBjjbtqGaxBqGcugFvgF3RlBxaUJebDyonCyShbAgb8ScFyq89hkA79poO8oSAEvV44649Q8ykc4og836Q2e4WxmcAttKaPx2bzElGh2oBRJcIKnDgWghgLzkhcwJrhqwES4Sow-fzpEy69rgakcK68a8yibO1BExxFEgxxwNyod4USp4Ki5yBxeeyE8k2d1G1wBwJwgy0YxqnGgQwMvByUjwF8VEpG9gmg6y7prgtx2pacjyGQ0B8ak2q19wOwam3-1UjGU2swio4C08mgnxuE8po3fwq1M2MBg8iFtQ0V88o6woc97xVDw4hwYw4Dw36kJ0Zw7Tw50g",
		CommentsFbDtsg:        "NAftgH6NimkTA7FpjPVKRJ4TP7Stn69LiFKB5fdo8eY40pPGVslHb5A%3A19%3A1754048966",
		CommentsJazoest:       "25285",
		CommentsSpinR:         "1025463513",
		CommentsSpinT:         "1754231361",
	}
}

func updateFacebookConfigFromResponse(config *FacebookConfig, responseHeaders http.Header, responseBody string) {
	// Track if any tokens were updated
	tokensUpdated := false

	if xFbDebug := responseHeaders.Get("X-Fb-Debug"); xFbDebug != "" {
		fmt.Printf("🔄 Updating X-Fb-Debug token: %s...\n", xFbDebug[:20])
		config.XFbDebug = xFbDebug
		tokensUpdated = true
	}
	if csp := responseHeaders.Get("Content-Security-Policy"); csp != "" && strings.Contains(csp, "nonce-") {
		nonceRegex := regexp.MustCompile(`'nonce-([^']+)'`)
		if matches := nonceRegex.FindStringSubmatch(csp); len(matches) > 1 {
			newNonce := matches[1]
			fmt.Printf("🔄 Updating CSP nonce: %s\n", newNonce)
			config.CSPNonce = newNonce
			tokensUpdated = true
		}
	}

	if strings.Contains(responseBody, "expansion_token") {
		if jsonStart := strings.Index(responseBody, "{"); jsonStart != -1 {
			jsonData := responseBody[jsonStart:]
			if endPos := findJSONEnd(jsonData); endPos > 0 {
				jsonData = jsonData[:endPos]
				tokenRegex := regexp.MustCompile(`"expansion_token":"([^"]+)"`)
				matches := tokenRegex.FindAllStringSubmatch(jsonData, -1)
				for _, match := range matches {
					if len(match) > 1 {
						token := match[1]
						config.ExpansionTokens["latest"] = token
						fmt.Printf("🔄 Updated expansion token: %s...\n", token[:30])
						tokensUpdated = true
						break
					}
				}
			}
		}
	}

	if strings.Contains(responseBody, "__req") {
		reqRegex := regexp.MustCompile(`"__req":"([^"]+)"`)
		if matches := reqRegex.FindStringSubmatch(responseBody); len(matches) > 1 {
			newReq := matches[1]
			reqSequence := []string{"15", "16", "17", "18", "19", "1a", "1b", "1c", "1d", "1e", "1f"}
			currentIndex := -1
			for i, req := range reqSequence {
				if req == config.CommentsRequestID {
					currentIndex = i
					break
				}
			}
			if currentIndex >= 0 {
				nextIndex := (currentIndex + 1) % len(reqSequence)
				config.CommentsRequestID = reqSequence[nextIndex]
				tokensUpdated = true
			} else {
				config.CommentsRequestID = newReq
				tokensUpdated = true
			}
		}
	}

	if strings.Contains(responseBody, "__s") {
		sessionRegex := regexp.MustCompile(`"__s":"([^"]+)"`)
		if matches := sessionRegex.FindStringSubmatch(responseBody); len(matches) > 1 {
			newSession := matches[1]
			if newSession != config.CommentsSession {
				sessionPreview := newSession
				if len(newSession) > 20 {
					sessionPreview = newSession[:20]
				}
				fmt.Printf("🔄 Updating Comments Session: %s...\n", sessionPreview)
				config.CommentsSession = newSession
				tokensUpdated = true
			}
		}
	}

	paramMappings := map[string]struct {
		commentsField *string
		paramName     string
	}{
		"__hs":     {&config.CommentsHash, "Hash"},
		"__rev":    {&config.CommentsRevision, "Revision"},
		"__hsi":    {&config.CommentsHashSessionID, "HashSessionID"},
		"__dyn":    {&config.CommentsDynamic, "Dynamic"},
		"__csr":    {&config.CommentsCSR, "CSR"},
		"__hsdp":   {&config.CommentsHSDP, "HSDP"},
		"__hblp":   {&config.CommentsHBLP, "HBLP"},
		"__sjsp":   {&config.CommentsSJSP, "SJSP"},
		"fb_dtsg":  {&config.CommentsFbDtsg, "FbDtsg"},
		"jazoest":  {&config.CommentsJazoest, "Jazoest"},
		"__spin_r": {&config.CommentsSpinR, "SpinR"},
		"__spin_t": {&config.CommentsSpinT, "SpinT"},
	}

	for param, mapping := range paramMappings {
		if strings.Contains(responseBody, param) {
			regex := regexp.MustCompile(fmt.Sprintf(`"%s":"([^"]+)"`, param))
			if matches := regex.FindStringSubmatch(responseBody); len(matches) > 1 {
				newValue := matches[1]
				if newValue != *mapping.commentsField {
					preview := newValue
					if len(newValue) > 20 {
						preview = newValue[:20]
					}
					fmt.Printf("🔄 Updating Comments %s: %s...\n", mapping.paramName, preview)
					*mapping.commentsField = newValue
					tokensUpdated = true
				}
			}
		}
	}

	// Always show the summary if any tokens were updated
	if tokensUpdated {
		fmt.Printf("🔄 Facebook Comments Token Rotation Summary:\n")
		sessionPreview := config.CommentsSession
		if len(config.CommentsSession) > 20 {
			sessionPreview = config.CommentsSession[:20]
		}
		fmt.Printf("   RequestID: %s, Session: %s...\n", config.CommentsRequestID, sessionPreview)
	}

	config.LastUpdated = time.Now()
}

func findJSONEnd(jsonData string) int {
	depth := 0
	for i, c := range jsonData {
		switch c {
		case '{':
			depth++
		case '}':
			depth--
			if depth == 0 {
				return i + 1
			}
		}
	}
	return -1
}

func fetchInitialComments(postID string, config *FacebookConfig) (string, error) {
	url := "https://web.facebook.com/api/graphql/"

	variablesPart := fmt.Sprintf(`variables=%%7B%%22commentsAfterCount%%22%%3A100%%2C%%22commentsIntentToken%%22%%3A%%22RANKED_UNFILTERED_CHRONOLOGICAL_REPLIES_INTENT_V1%%22%%2C%%22feedLocation%%22%%3A%%22POST_PERMALINK_DIALOG%%22%%2C%%22feedbackSource%%22%%3A2%%2C%%22focusCommentID%%22%%3Anull%%2C%%22scale%%22%%3A1%%2C%%22useDefaultActor%%22%%3Afalse%%2C%%22id%%22%%3A%%22%s%%22%%2C%%22__relay_internal__pv__IsWorkUserrelayprovider%%22%%3Afalse%%7D`, postID)

	payload := fmt.Sprintf("av=%s&__aaid=0&__user=%s&__a=1&__req=%s&__hs=%s&dpr=1&__ccg=MODERATE&__rev=%s&__s=%s&__hsi=%s&__dyn=%s&__csr=%s&__hsdp=%s&__hblp=%s&__sjsp=%s&__comet_req=15&fb_dtsg=%s&jazoest=%s&lsd=%s&__spin_r=%s&__spin_b=trunk&__spin_t=%s&__crn=comet.fbweb.CometSinglePostDialogRoute&fb_api_caller_class=RelayModern&fb_api_req_friendly_name=CommentListComponentsRootQuery&%s&server_timestamps=true&doc_id=24509177942033202",
		config.CommentsUserID,
		config.CommentsUserID,
		config.CommentsRequestID,
		config.CommentsHash,
		config.CommentsRevision,
		config.CommentsSession,
		config.CommentsHashSessionID,
		config.CommentsDynamic,
		config.CommentsCSR,
		config.CommentsHSDP,
		config.CommentsHBLP,
		config.CommentsSJSP,
		config.CommentsFbDtsg,
		config.CommentsJazoest,
		config.LSDToken,
		config.CommentsSpinR,
		config.CommentsSpinT,
		variablesPart)

	req, err := http.NewRequest("POST", url, strings.NewReader(payload))
	if err != nil {
		return "", fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("accept", "*/*")
	req.Header.Set("accept-language", "en-US,en;q=0.9")
	req.Header.Set("cache-control", "no-cache")
	req.Header.Set("content-type", "application/x-www-form-urlencoded")
	req.Header.Set("origin", "https://web.facebook.com")
	req.Header.Set("pragma", "no-cache")
	req.Header.Set("priority", "u=1, i")
	req.Header.Set("referer", globalPostURL)
	req.Header.Set("sec-ch-prefers-color-scheme", "light")
	req.Header.Set("sec-ch-ua", "\"Not)A;Brand\";v=\"8\", \"Chromium\";v=\"138\", \"Google Chrome\";v=\"138\"")
	req.Header.Set("sec-ch-ua-full-version-list", "\"Not)A;Brand\";v=\"8.0.0.0\", \"Chromium\";v=\"138.0.7204.168\", \"Google Chrome\";v=\"138.0.7204.168\"")
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-model", "\"\"")
	req.Header.Set("sec-ch-ua-platform", "\"Linux\"")
	req.Header.Set("sec-ch-ua-platform-version", "\"6.8.0\"")
	req.Header.Set("sec-fetch-dest", "empty")
	req.Header.Set("sec-fetch-mode", "cors")
	req.Header.Set("sec-fetch-site", "same-origin")
	req.Header.Set("user-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/138.0.0.0 Safari/537.36")
	req.Header.Set("x-asbd-id", "359341")
	req.Header.Set("x-fb-friendly-name", "CommentListComponentsRootQuery")
	req.Header.Set("x-fb-lsd", config.LSDToken)

	cookieStr := config.Cookies
	for cookie := range strings.SplitSeq(cookieStr, "; ") {
		parts := strings.SplitN(cookie, "=", 2)
		if len(parts) == 2 {
			req.AddCookie(&http.Cookie{
				Name:  parts[0],
				Value: parts[1],
			})
		}
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response: %w", err)
	}

	// fmt.Printf("🔍 === Facebook Initial Comments Response ===\n")
	// fmt.Printf("📡 Response Status: %s\n", resp.Status)
	// fmt.Printf("📋 Response Headers:\n")
	// for key, values := range resp.Header {
	// 	for _, value := range values {
	// 		fmt.Printf("    %s: %s\n", key, value)
	// 	}
	// }
	// fmt.Printf("📄 Response Body Length: %d bytes\n", len(body))
	// if len(body) < 2000 {
	// 	fmt.Printf("📝 Response Body: %s\n", string(body))
	// } else {
	// 	fmt.Printf("📝 Response Body Sample (first 2000 chars): %s...\n", string(body[:2000]))
	// }
	// fmt.Printf("🔍 === End Facebook Initial Comments Response ===\n\n")

	updateFacebookConfigFromResponse(config, resp.Header, string(body))

	return string(body), nil
}

func fetchPaginatedComments(cursor, postID string, config *FacebookConfig) (string, error) {
	url := "https://web.facebook.com/api/graphql/"

	basePayload := fmt.Sprintf("av=%s&__aaid=0&__user=%s&__a=1&__req=%s&__hs=%s&dpr=1&__ccg=MODERATE&__rev=%s&__s=%s&__hsi=%s&__dyn=%s&__csr=%s&__hsdp=%s&__hblp=%s&__sjsp=%s&__comet_req=15&fb_dtsg=%s&jazoest=%s&lsd=%s&__spin_r=%s&__spin_b=trunk&__spin_t=%s&__crn=comet.fbweb.CometSinglePostDialogRoute&fb_api_caller_class=RelayModern&fb_api_req_friendly_name=CommentsListComponentsPaginationQuery",
		config.CommentsUserID,
		config.CommentsUserID,
		config.CommentsRequestID,
		config.CommentsHash,
		config.CommentsRevision,
		config.CommentsSession,
		config.CommentsHashSessionID,
		config.CommentsDynamic,
		config.CommentsCSR,
		config.CommentsHSDP,
		config.CommentsHBLP,
		config.CommentsSJSP,
		config.CommentsFbDtsg,
		config.CommentsJazoest,
		config.LSDToken,
		config.CommentsSpinR,
		config.CommentsSpinT)

	cursorVariables := fmt.Sprintf("&variables=%%7B%%22commentsAfterCount%%22%%3A100%%2C%%22commentsAfterCursor%%22%%3A%%22%s%%22%%2C%%22commentsBeforeCount%%22%%3Anull%%2C%%22commentsBeforeCursor%%22%%3Anull%%2C%%22commentsIntentToken%%22%%3A%%22RANKED_UNFILTERED_CHRONOLOGICAL_REPLIES_INTENT_V1%%22%%2C%%22feedLocation%%22%%3A%%22POST_PERMALINK_DIALOG%%22%%2C%%22focusCommentID%%22%%3Anull%%2C%%22scale%%22%%3A1%%2C%%22useDefaultActor%%22%%3Afalse%%2C%%22id%%22%%3A%%22%s%%22%%2C%%22__relay_internal__pv__IsWorkUserrelayprovider%%22%%3Afalse%%7D&server_timestamps=true&doc_id=9994312660685367", cursor, postID)

	payload := basePayload + cursorVariables

	req, err := http.NewRequest("POST", url, strings.NewReader(payload))
	if err != nil {
		return "", fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("accept", "*/*")
	req.Header.Set("accept-language", "en-US,en;q=0.9")
	req.Header.Set("cache-control", "no-cache")
	req.Header.Set("content-type", "application/x-www-form-urlencoded")
	req.Header.Set("origin", "https://web.facebook.com")
	req.Header.Set("pragma", "no-cache")
	req.Header.Set("priority", "u=1, i")
	req.Header.Set("referer", globalPostURL)
	req.Header.Set("sec-ch-prefers-color-scheme", "light")
	req.Header.Set("sec-ch-ua", "\"Not)A;Brand\";v=\"8\", \"Chromium\";v=\"138\", \"Google Chrome\";v=\"138\"")
	req.Header.Set("sec-ch-ua-full-version-list", "\"Not)A;Brand\";v=\"8.0.0.0\", \"Chromium\";v=\"138.0.7204.168\", \"Google Chrome\";v=\"138.0.7204.168\"")
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-model", "\"\"")
	req.Header.Set("sec-ch-ua-platform", "\"Linux\"")
	req.Header.Set("sec-ch-ua-platform-version", "\"6.8.0\"")
	req.Header.Set("sec-fetch-dest", "empty")
	req.Header.Set("sec-fetch-mode", "cors")
	req.Header.Set("sec-fetch-site", "same-origin")
	req.Header.Set("user-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36")
	req.Header.Set("x-asbd-id", "359341")
	req.Header.Set("x-fb-friendly-name", "CommentsListComponentsPaginationQuery")
	req.Header.Set("x-fb-lsd", config.LSDToken)

	cookieStr := config.Cookies
	for cookie := range strings.SplitSeq(cookieStr, "; ") {
		parts := strings.SplitN(cookie, "=", 2)
		if len(parts) == 2 {
			req.AddCookie(&http.Cookie{
				Name:  parts[0],
				Value: parts[1],
			})
		}
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response: %w", err)
	}

	// Update tokens from response
	updateFacebookConfigFromResponse(config, resp.Header, string(body))

	// fmt.Printf("🔍 === Facebook Paginated Comments Response ===\n")
	// fmt.Printf("📡 Response Status: %s\n", resp.Status)
	// fmt.Printf("📋 Response Headers:\n")
	// for key, values := range resp.Header {
	// 	for _, value := range values {
	// 		fmt.Printf("    %s: %s\n", key, value)
	// 	}
	// }
	// fmt.Printf("📄 Response Body Length: %d bytes\n", len(body))
	// if len(body) < 2000 {
	// 	fmt.Printf("📝 Response Body: %s\n", string(body))
	// } else {
	// 	fmt.Printf("📝 Response Body Sample (first 2000 chars): %s...\n", string(body[:2000]))
	// }
	// fmt.Printf("🔍 === End Facebook Paginated Comments Response ===\n\n")

	return string(body), nil
}

func extractDataFromFacebookResponse(body string) (map[string]any, error) {

	jsonStart := strings.Index(body, "{")
	if jsonStart == -1 {
		return nil, fmt.Errorf("no JSON object found in response")
	}

	jsonData := body[jsonStart:]

	depth := 0
	endPos := -1

	for i, c := range jsonData {
		if c == '{' {
			depth++
		} else if c == '}' {
			depth--
			if depth == 0 {
				endPos = i + 1
				break
			}
		}
	}

	if endPos > 0 {
		jsonData = jsonData[:endPos]
	}

	var result map[string]any
	if err := json.Unmarshal([]byte(jsonData), &result); err != nil {
		return nil, fmt.Errorf("error parsing JSON: %w", err)
	}

	return result, nil
}

func extractComments(data map[string]any) (any, error) {

	if data == nil {
		return nil, fmt.Errorf("no data provided")
	}

	dataObj, ok := data["data"].(map[string]any)
	if !ok {
		return nil, fmt.Errorf("data.data not found or not an object")
	}

	nodeObj, ok := dataObj["node"].(map[string]any)
	if !ok {
		return nil, fmt.Errorf("data.data.node not found or not an object")
	}

	renderingObj, ok := nodeObj["comment_rendering_instance_for_feed_location"].(map[string]any)
	if !ok {
		return nil, fmt.Errorf("comment_rendering_instance_for_feed_location not found or not an object")
	}

	commentsObj, ok := renderingObj["comments"].(map[string]any)
	if !ok {
		return nil, fmt.Errorf("comments not found or not an object")
	}

	return commentsObj, nil
}

func extractEndCursor(commentsObj map[string]any) (string, bool) {
	pageInfo, ok := commentsObj["page_info"].(map[string]any)
	if !ok {
		return "", false
	}

	endCursor, ok := pageInfo["end_cursor"].(string)
	if !ok {
		return "", false
	}

	hasNextPage, ok := pageInfo["has_next_page"].(bool)
	if !ok {
		hasNextPage = false
	}

	return endCursor, hasNextPage
}

type FacebookComment struct {
	ID                                                string                      `json:"id"`
	IsSubreplyParentDeleted                           bool                        `json:"is_subreply_parent_deleted"`
	Author                                            Author                      `json:"author"`
	Feedback                                          Feedback                    `json:"feedback"`
	LegacyFbid                                        string                      `json:"legacy_fbid"`
	Depth                                             int                         `json:"depth"`
	Body                                              TextWithEntities            `json:"body"`
	Attachments                                       []any                       `json:"attachments"`
	IsMarkdownEnabled                                 bool                        `json:"is_markdown_enabled"`
	CommunityCommentSignalRenderer                    any                         `json:"community_comment_signal_renderer"`
	CommentMenuTooltip                                any                         `json:"comment_menu_tooltip"`
	ShouldShowCommentMenu                             bool                        `json:"should_show_comment_menu"`
	IsAuthorWeakReference                             bool                        `json:"is_author_weak_reference"`
	CommentActionLinks                                []CommentActionLink         `json:"comment_action_links"`
	PreferredBody                                     TextWithEntities            `json:"preferred_body"`
	BodyRenderer                                      BodyRenderer                `json:"body_renderer"`
	CommentParent                                     *CommentParent              `json:"comment_parent"`
	IsDeclinedByGroupAdminAssistant                   bool                        `json:"is_declined_by_group_admin_assistant"`
	IsGamingVideoComment                              bool                        `json:"is_gaming_video_comment"`
	TimestampInVideo                                  any                         `json:"timestamp_in_video"`
	TranslatabilityForViewer                          TranslatabilityForViewer    `json:"translatability_for_viewer"`
	WrittenWhileVideoWasLive                          bool                        `json:"written_while_video_was_live"`
	GroupCommentInfo                                  any                         `json:"group_comment_info"`
	BizwebCommentInfo                                 any                         `json:"bizweb_comment_info"`
	HasConstituentBadge                               bool                        `json:"has_constituent_badge"`
	CanViewerSeeSubsribeButton                        bool                        `json:"can_viewer_see_subsribe_button"`
	CanSeeConstituentBadgeUpsell                      bool                        `json:"can_see_constituent_badge_upsell"`
	LegacyToken                                       string                      `json:"legacy_token"`
	ParentFeedback                                    ParentFeedback              `json:"parent_feedback"`
	QuestionAndAnswerType                             any                         `json:"question_and_answer_type"`
	AuthorUserSignalsRenderer                         any                         `json:"author_user_signals_renderer"`
	AuthorBadgeRenderers                              []any                       `json:"author_badge_renderers"`
	IdentityBadgesWeb                                 []any                       `json:"identity_badges_web"`
	CanShowMultipleIdentityBadges                     bool                        `json:"can_show_multiple_identity_badges"`
	DiscoverableIdentityBadgesWeb                     []DiscoverableIdentityBadge `json:"discoverable_identity_badges_web"`
	User                                              User                        `json:"user"`
	IsViewerCommentPoster                             bool                        `json:"is_viewer_comment_poster"`
	ParentPostStory                                   ParentPostStory             `json:"parent_post_story"`
	GenAiContentTransparencyLabelRenderer             any                         `json:"gen_ai_content_transparency_label_renderer"`
	WorkAmaAnswerStatus                               any                         `json:"work_ama_answer_status"`
	WorkKnowledgeInlineAnnotationCommentBadgeRenderer any                         `json:"work_knowledge_inline_annotation_comment_badge_renderer"`
	BusinessCommentAttributes                         []any                       `json:"business_comment_attributes"`
	IsLiveVideoComment                                bool                        `json:"is_live_video_comment"`
	CreatedTime                                       int64                       `json:"created_time"`
	TranslationAvailableForViewer                     bool                        `json:"translation_available_for_viewer"`
	InlineSurveyConfig                                any                         `json:"inline_survey_config"`
	SpamDisplayMode                                   string                      `json:"spam_display_mode"`
	AttachedStory                                     any                         `json:"attached_story"`
	CommentDirectParent                               *CommentDirectParent        `json:"comment_direct_parent"`
	IfViewerCanSeeMemberPageTooltip                   any                         `json:"if_viewer_can_see_member_page_tooltip"`
	IsDisabled                                        bool                        `json:"is_disabled"`
	WorkAnsweredEventCommentRenderer                  any                         `json:"work_answered_event_comment_renderer"`
	CommentUpperBadgeRenderer                         any                         `json:"comment_upper_badge_renderer"`
	ElevatedCommentData                               any                         `json:"elevated_comment_data"`
	Typename                                          string                      `json:"__typename"`
}

type Author struct {
	Typename             string         `json:"__typename"`
	IsActor              string         `json:"__isActor"`
	Name                 string         `json:"name"`
	ID                   string         `json:"id"`
	ProfilePictureDepth0 ProfilePicture `json:"profile_picture_depth_0"`
	ProfilePictureDepth1 ProfilePicture `json:"profile_picture_depth_1"`
	Gender               string         `json:"gender"`
	IsEntity             string         `json:"__isEntity"`
	URL                  string         `json:"url"`
	WorkInfo             any            `json:"work_info"`
	IsVerified           bool           `json:"is_verified"`
	ShortName            string         `json:"short_name"`
	SubscribeStatus      string         `json:"subscribe_status"`
}

type ProfilePicture struct {
	URI string `json:"uri"`
}

type Feedback struct {
	ID                             string                         `json:"id"`
	ExpansionInfo                  ExpansionInfo                  `json:"expansion_info"`
	RepliesFields                  RepliesFields                  `json:"replies_fields"`
	ViewerActor                    any                            `json:"viewer_actor"`
	ActorProvider                  ActorProvider                  `json:"actor_provider"`
	URL                            string                         `json:"url"`
	Typename                       string                         `json:"__typename"`
	IfViewerCanCommentAnonymously  any                            `json:"if_viewer_can_comment_anonymously"`
	Plugins                        []Plugin                       `json:"plugins"`
	CommentComposerPlaceholder     string                         `json:"comment_composer_placeholder"`
	ConstituentBadgeBannerRenderer any                            `json:"constituent_badge_banner_renderer"`
	AssociatedGroup                any                            `json:"associated_group"`
	HaveCommentsBeenDisabled       bool                           `json:"have_comments_been_disabled"`
	AreLiveVideoCommentsDisabled   bool                           `json:"are_live_video_comments_disabled"`
	IsViewerMuted                  bool                           `json:"is_viewer_muted"`
	CommentRenderingInstance       any                            `json:"comment_rendering_instance"`
	CommentsDisabledNoticeRenderer CommentsDisabledNoticeRenderer `json:"comments_disabled_notice_renderer"`
	RepliesConnection              RepliesConnection              `json:"replies_connection"`
	ParentObjectEnt                ParentObjectEnt                `json:"parent_object_ent"`
	ViewerFeedbackReactionInfo     any                            `json:"viewer_feedback_reaction_info"`
	TopReactions                   TopReactions                   `json:"top_reactions"`
	Reactors                       Reactors                       `json:"reactors"`
}

type ExpansionInfo struct {
	ExpansionToken       string `json:"expansion_token"`
	ShouldShowReplyCount bool   `json:"should_show_reply_count"`
}

type RepliesFields struct {
	Count      int `json:"count"`
	TotalCount int `json:"total_count"`
}

type ActorProvider struct {
	Typename     string `json:"__typename"`
	CurrentActor any    `json:"current_actor"`
	ID           string `json:"id"`
}

type Plugin struct {
	Typename                                          string          `json:"__typename"`
	ContextID                                         any             `json:"context_id,omitempty"`
	PostID                                            string          `json:"post_id,omitempty"`
	ModuleOperationUseCometUFIComposerPluginsFeedback ModuleOperation `json:"__module_operation_useCometUFIComposerPlugins_feedback"`
	ModuleComponentUseCometUFIComposerPluginsFeedback ModuleComponent `json:"__module_component_useCometUFIComposerPlugins_feedback"`
	EmojiSize                                         int             `json:"emoji_size,omitempty"`
}

type ModuleOperation struct {
	Dr string `json:"__dr"`
}

type ModuleComponent struct {
	Dr string `json:"__dr"`
}

type CommentsDisabledNoticeRenderer struct {
	Typename                                             string           `json:"__typename"`
	NoticeMessage                                        TextWithEntities `json:"notice_message"`
	ModuleOperationCometUFICommentDisabledNoticeFeedback ModuleOperation  `json:"__module_operation_CometUFICommentDisabledNotice_feedback"`
	ModuleComponentCometUFICommentDisabledNoticeFeedback ModuleComponent  `json:"__module_component_CometUFICommentDisabledNotice_feedback"`
}

type RepliesConnection struct {
	Edges    []any    `json:"edges"`
	PageInfo PageInfo `json:"page_info"`
}

type PageInfo struct {
	EndCursor       any  `json:"end_cursor"`
	HasNextPage     bool `json:"has_next_page"`
	HasPreviousPage bool `json:"has_previous_page"`
	StartCursor     any  `json:"start_cursor"`
}

type ParentObjectEnt struct {
	Typename                      string         `json:"__typename"`
	Feedback                      SimpleFeedback `json:"feedback"`
	InlineRepliesExpanderRenderer any            `json:"inline_replies_expander_renderer"`
	ID                            string         `json:"id"`
}

type SimpleFeedback struct {
	ID string `json:"id"`
}

type TopReactions struct {
	Edges []ReactionEdge `json:"edges"`
}

type ReactionEdge struct {
	VisibleInBlingBar bool         `json:"visible_in_bling_bar"`
	Node              ReactionNode `json:"node"`
	ReactionCount     int          `json:"reaction_count"`
}

type ReactionNode struct {
	ID string `json:"id"`
}

type Reactors struct {
	CountReduced string `json:"count_reduced"`
}

type TextWithEntities struct {
	Typename          string `json:"__typename,omitempty"`
	DelightRanges     []any  `json:"delight_ranges"`
	ImageRanges       []any  `json:"image_ranges"`
	InlineStyleRanges []any  `json:"inline_style_ranges"`
	AggregatedRanges  []any  `json:"aggregated_ranges"`
	Ranges            []any  `json:"ranges"`
	ColorRanges       []any  `json:"color_ranges"`
	Text              string `json:"text"`
	TranslationType   string `json:"translation_type,omitempty"`
}

type CommentActionLink struct {
	Typename                                         string            `json:"__typename"`
	Comment                                          ActionLinkComment `json:"comment"`
	ModuleOperationCometUFICommentActionLinksComment ModuleOperation   `json:"__module_operation_CometUFICommentActionLinks_comment"`
	ModuleComponentCometUFICommentActionLinksComment ModuleComponent   `json:"__module_component_CometUFICommentActionLinks_comment"`
}

type ActionLinkComment struct {
	ID          string `json:"id"`
	CreatedTime int64  `json:"created_time"`
	URL         string `json:"url"`
}

type BodyRenderer struct {
	Typename                                              string          `json:"__typename"`
	DelightRanges                                         []any           `json:"delight_ranges"`
	ImageRanges                                           []any           `json:"image_ranges"`
	InlineStyleRanges                                     []any           `json:"inline_style_ranges"`
	AggregatedRanges                                      []any           `json:"aggregated_ranges"`
	Ranges                                                []any           `json:"ranges"`
	ColorRanges                                           []any           `json:"color_ranges"`
	Text                                                  string          `json:"text"`
	ModuleOperationCometUFICommentTextBodyRendererComment ModuleOperation `json:"__module_operation_CometUFICommentTextBodyRenderer_comment"`
	ModuleComponentCometUFICommentTextBodyRendererComment ModuleComponent `json:"__module_component_CometUFICommentTextBodyRenderer_comment"`
}

type CommentParent struct {
	Author Author `json:"author"`
	ID     string `json:"id"`
}

type TranslatabilityForViewer struct {
	SourceDialect string `json:"source_dialect"`
}

type ParentFeedback struct {
	ID                  string        `json:"id"`
	ShareFbid           string        `json:"share_fbid"`
	PoliticalFigureData any           `json:"political_figure_data"`
	OwningProfile       OwningProfile `json:"owning_profile"`
}

type OwningProfile struct {
	Typename string `json:"__typename"`
	Name     string `json:"name"`
	ID       string `json:"id"`
}

type DiscoverableIdentityBadge struct {
	GreyBadgeAsset           string `json:"grey_badge_asset"`
	DarkModeBadgeAsset       string `json:"dark_mode_badge_asset"`
	LightModeBadgeAsset      string `json:"light_mode_badge_asset"`
	IsEarned                 bool   `json:"is_earned"`
	InformationTitle         string `json:"information_title"`
	InformationDescription   string `json:"information_description"`
	IsEnabled                bool   `json:"is_enabled"`
	IsManageable             bool   `json:"is_manageable"`
	Serialized               string `json:"serialized"`
	IdentityBadgeType        string `json:"identity_badge_type"`
	InformationButtonEnabled bool   `json:"information_button_enabled"`
	InformationButtonURI     string `json:"information_button_uri"`
	InformationButtonText    string `json:"information_button_text"`
	TierInfo                 any    `json:"tier_info"`
}

type User struct {
	Name           string         `json:"name"`
	ProfilePicture ProfilePicture `json:"profile_picture"`
	ID             string         `json:"id"`
}

type ParentPostStory struct {
	Attachments []PostAttachment `json:"attachments"`
	ID          string           `json:"id"`
}

type PostAttachment struct {
	Media PostMedia `json:"media"`
}

type PostMedia struct {
	Typename string `json:"__typename"`
	ID       string `json:"id"`
}

type CommentDirectParent struct {
	Author DirectParentAuthor `json:"author"`
	ID     string             `json:"id"`
}

type DirectParentAuthor struct {
	Typename string `json:"__typename"`
	Name     string `json:"name"`
	Gender   string `json:"gender"`
	ID       string `json:"id"`
}

func fetchAllPostComments(postID string, config *FacebookConfig) ([]FacebookComment, error) {
	fmt.Printf("🔍 Starting to fetch main comments from post: %s\n", postID)

	var allFacebookComments []FacebookComment
	var currentCursor string
	pageCount := 1
	hasNextPage := true
	retryCount := 0
	maxRetries := 5
	const maxComments = 500

	fmt.Printf("⚠️  Comment limit: %d main comments maximum\n", maxComments)
	fmt.Printf("📄 Fetching comments page %d...\n", pageCount)
	var response string
	var err error

	for retryCount <= maxRetries {
		response, err = fetchInitialComments(postID, config)
		if err != nil {
			retryCount++
			if retryCount <= maxRetries {
				fmt.Printf("⚠️ Error fetching initial comments. Retry %d/%d: %v\n", retryCount, maxRetries, err)
				continue
			} else {
				return nil, fmt.Errorf("error fetching initial comments after %d retries: %w", maxRetries, err)
			}
		}

		if len(response) > 100 && strings.Contains(response, "comments") {
			break
		} else {
			retryCount++
			if retryCount <= maxRetries {
				fmt.Printf("⚠️ Received invalid response. Retry %d/%d (response length: %d)\n", retryCount, maxRetries, len(response))
				continue
			} else {
				return nil, fmt.Errorf("received invalid response after %d retries", maxRetries)
			}
		}
	}

	data, err := extractDataFromFacebookResponse(response)
	if err != nil {
		return nil, fmt.Errorf("error parsing initial response: %w", err)
	}

	pageFacebookComments, err := extractFacebookComments(data)
	if err != nil {
		return nil, fmt.Errorf("error extracting initial comments: %w", err)
	}

	// Filter only main comments (no replies)
	var mainComments []FacebookComment
	for _, comment := range pageFacebookComments {
		if comment.CommentParent == nil && comment.CommentDirectParent == nil {
			mainComments = append(mainComments, comment)
		}
	}

	allFacebookComments = append(allFacebookComments, mainComments...)
	fmt.Printf("✅ Page %d: Found %d main comments (Total: %d)\n", pageCount, len(mainComments), len(allFacebookComments))

	if len(allFacebookComments) >= maxComments {
		allFacebookComments = allFacebookComments[:maxComments]
		fmt.Printf("🛑 Reached comment limit of %d comments. Stopping extraction.\n", maxComments)
		return allFacebookComments, nil
	}

	if commentsObj, err := extractComments(data); err == nil {
		if commentsMap, ok := commentsObj.(map[string]any); ok {
			currentCursor, hasNextPage = extractEndCursor(commentsMap)
		}
	}

	pageCount++

	for hasNextPage {
		fmt.Printf("📄 Fetching comments page %d...\n", pageCount)

		retryCount = 0
		var paginatedResponse string
		var paginatedErr error

		for retryCount <= maxRetries {
			paginatedResponse, paginatedErr = fetchPaginatedComments(currentCursor, postID, config)
			if paginatedErr != nil {
				retryCount++
				if retryCount <= maxRetries {
					fmt.Printf("⚠️ Error fetching page %d. Retry %d/%d: %v\n", pageCount, retryCount, maxRetries, paginatedErr)
					continue
				} else {
					fmt.Printf("❌ Failed to fetch page %d after %d retries: %v\n", pageCount, maxRetries, paginatedErr)
					break
				}
			}

			if len(paginatedResponse) > 100 {
				break
			} else {
				retryCount++
				if retryCount <= maxRetries {
					fmt.Printf("⚠️ Invalid response for page %d. Retry %d/%d (length: %d)\n", pageCount, retryCount, maxRetries, len(paginatedResponse))
					continue
				} else {
					fmt.Printf("❌ Invalid response for page %d after %d retries\n", pageCount, maxRetries)
					break
				}
			}
		}

		if paginatedErr != nil {
			fmt.Printf("⚠️ Skipping page %d due to persistent errors\n", pageCount)
			break
		}

		response = paginatedResponse

		data, err := extractDataFromFacebookResponse(response)
		if err != nil {
			fmt.Printf("⚠️ Error parsing page %d: %v\n", pageCount, err)
			break
		}

		pageFacebookComments, err := extractFacebookComments(data)
		if err != nil {
			fmt.Printf("⚠️ Error extracting comments from page %d: %v\n", pageCount, err)
			break
		}

		if len(pageFacebookComments) == 0 {
			fmt.Printf("🛑 No more comments found on page %d\n", pageCount)
			break
		}

		// Filter only main comments (no replies)
		var mainComments []FacebookComment
		for _, comment := range pageFacebookComments {
			if comment.CommentParent == nil && comment.CommentDirectParent == nil {
				mainComments = append(mainComments, comment)
			}
		}

		allFacebookComments = append(allFacebookComments, mainComments...)
		fmt.Printf("✅ Page %d: Found %d main comments (Total: %d)\n", pageCount, len(mainComments), len(allFacebookComments))

		if len(allFacebookComments) >= maxComments {
			allFacebookComments = allFacebookComments[:maxComments]
			fmt.Printf("🛑 Reached comment limit of %d comments. Stopping extraction.\n", maxComments)
			break
		}

		if commentsObj, err := extractComments(data); err == nil {
			if commentsMap, ok := commentsObj.(map[string]any); ok {
				currentCursor, hasNextPage = extractEndCursor(commentsMap)
				if !hasNextPage {
					fmt.Printf("🏁 Reached end of comments pagination\n")
				}
			}
		}

		pageCount++
	}

	mainCommentsCount := CountFacebookComments(allFacebookComments)
	fmt.Printf("📊 Facebook Main Comments Summary:\n")
	fmt.Printf("   📄 Pages fetched: %d\n", pageCount-1)
	fmt.Printf("   📝 Main comments: %d\n", mainCommentsCount)
	fmt.Printf("📊 Total main comments fetched: %d items from %d pages\n", mainCommentsCount, pageCount-1)

	return allFacebookComments, nil
}

func extractPostIDFromURL(facebookURL string) (string, error) {

	globalPostURL = facebookURL

	originalURL := facebookURL

	parsedURL, err := url.Parse(facebookURL)
	if err != nil {
		return "", fmt.Errorf("invalid URL: %w", err)
	}

	if strings.Contains(parsedURL.Path, "permalink.php") {

		queryParams := parsedURL.Query()
		storyFbid := queryParams.Get("story_fbid")

		if storyFbid != "" {
			fmt.Printf("🔍 Detected permalink.php URL with story_fbid: %s\n", storyFbid)

			if strings.HasPrefix(storyFbid, "pfbid") || regexp.MustCompile(`^\d+$`).MatchString(storyFbid) {
				postID := base64.URLEncoding.EncodeToString([]byte("feedback:" + storyFbid))
				fmt.Printf("🔄 Converted story_fbid to base64 post ID: %s\n", postID)
				return postID, nil
			}

			return storyFbid, nil
		}
	}

	pathSegments := strings.Split(strings.Trim(parsedURL.Path, "/"), "/")

	if len(pathSegments) >= 3 && pathSegments[0] == "share" && (pathSegments[1] == "p" || pathSegments[1] == "r") {
		shareType := pathSegments[1]
		shareID := pathSegments[2]
		fmt.Printf("🔍 Detected Facebook share link (type: %s) with ID: %s\n", shareType, shareID)
		fmt.Printf("🔄 Following redirect to get actual post URL...\n")

		client := &http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		}

		resp, err := client.Head(facebookURL)
		if err != nil {
			return "", fmt.Errorf("error following share link: %w", err)
		}
		defer resp.Body.Close()

		location := resp.Header.Get("Location")
		if location == "" {
			return "", fmt.Errorf("share link did not provide a redirect location")
		}

		fmt.Printf("✅ Share link redirected to: %s\n", location)

		return extractPostIDFromURL(location)
	}

	if len(pathSegments) >= 2 && pathSegments[0] == "reel" {
		reelID := pathSegments[1]

		if queryIndex := strings.Index(reelID, "?"); queryIndex != -1 {
			reelID = reelID[:queryIndex]
		}

		fmt.Printf("🔍 Detected Facebook Reel with ID: %s\n", reelID)

		if matched, _ := regexp.MatchString(`^\d+$`, reelID); matched {
			postID := base64.URLEncoding.EncodeToString([]byte("feedback:" + reelID))
			fmt.Printf("🔄 Converted reel ID to base64 post ID: %s\n", postID)
			return postID, nil
		}

		return reelID, nil
	}

	if len(pathSegments) >= 4 && pathSegments[0] == "groups" && pathSegments[2] == "permalink" {
		postID := pathSegments[3]
		fmt.Printf("🔍 Detected Facebook group post with ID: %s\n", postID)

		encodedID := base64.URLEncoding.EncodeToString([]byte("feedback:" + postID))
		fmt.Printf("🔄 Converted group post ID to base64: %s\n", encodedID)
		return encodedID, nil
	}

	for i, segment := range pathSegments {
		if segment == "posts" && i+1 < len(pathSegments) {
			postSlug := pathSegments[i+1]

			if queryIndex := strings.Index(postSlug, "?"); queryIndex != -1 {
				postSlug = postSlug[:queryIndex]
			}
			if hashIndex := strings.Index(postSlug, "#"); hashIndex != -1 {
				postSlug = postSlug[:hashIndex]
			}

			fmt.Printf("🔍 Found post slug: %s\n", postSlug)

			if strings.HasPrefix(postSlug, "pfbid") {
				fmt.Printf("🔍 pfbid format detected: %s\n", postSlug)
				postID := base64.URLEncoding.EncodeToString([]byte("feedback:" + postSlug))
				fmt.Printf("🔄 Converted pfbid to base64 post ID: %s\n", postID)
				return postID, nil
			}

			if matched, _ := regexp.MatchString(`^\d+$`, postSlug); matched {
				postID := base64.URLEncoding.EncodeToString([]byte("feedback:" + postSlug))
				fmt.Printf("🔄 Converted numeric ID to post ID: %s\n", postID)
				return postID, nil
			}

			return postSlug, nil
		}
	}

	return "", fmt.Errorf("could not extract post ID from URL: %s", originalURL)
}

func showUsage() {
	fmt.Println("🚀 Facebook Comments & Replies Scraper")
	fmt.Println("=====================================")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Printf("  %s <facebook_post_url_or_share_link>\n", os.Args[0])
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  # Use Facebook share link (recommended)")
	fmt.Printf("  %s \"https://web.facebook.com/share/p/1AguUnrRzz/\"\n", os.Args[0])
	fmt.Println()
	fmt.Println("  # Use direct post URL")
	fmt.Printf("  %s \"https://web.facebook.com/username/posts/pfbid123...\"\n", os.Args[0])
	fmt.Println()
	fmt.Println("  # Use base64 encoded post ID (advanced)")
	fmt.Printf("  %s \"ZmVlZGJhY2s6cGZiaWQxMjM...\"\n", os.Args[0])
	fmt.Println()
	fmt.Println("📝 Note: The scraper will automatically:")
	fmt.Println("   • Follow share link redirects to get the actual post URL")
	fmt.Println("   • Extract the post ID from the URL")
	fmt.Println("   • Encode it properly for Facebook's GraphQL API")
	fmt.Println("   • Scrape all comments and replies from the post")
	fmt.Println()
}

func main() {

	startTime := time.Now()

	fmt.Println("📱 Facebook Comment Scraper (Free Limited Version)")
	fmt.Println("==================================================")
	fmt.Println("This tool extracts main comments from Facebook posts, reels, and shares.")
	fmt.Println("⚠️  Limited to 500 main comments maximum per post")
	fmt.Println()
	fmt.Println("💎 Need unlimited comments with replies?")
	fmt.Println("📧 Email: haronkibetrutoh@gmail.com")
	fmt.Println("📱 WhatsApp: +254718448461")
	fmt.Println()

	var urlArg string

	if len(os.Args) > 1 {
		urlArg = os.Args[1]
	} else {
		fmt.Print("🔗 Enter Facebook URL (post/reel/share): ")
		scanner := bufio.NewScanner(os.Stdin)
		if scanner.Scan() {
			urlArg = scanner.Text()
		} else {
			if err := scanner.Err(); err != nil {
				fmt.Printf("❌ Error reading input: %v\n", err)
			} else {
				fmt.Println("❌ No URL provided")
			}
			os.Exit(1)
		}

		urlArg = strings.TrimSpace(urlArg)

		if urlArg == "" {
			fmt.Println("❌ No URL provided")
			showUsage()
			os.Exit(1)
		}
	}

	postID, err := extractPostIDFromURL(urlArg)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✅ Successfully extracted post ID: %s\n", postID)

	globalPostID = postID
	globalPostURL = urlArg

	fmt.Printf("⏱️ Extraction started at: %s\n", startTime.Format("15:04:05"))

	config := getDefaultFacebookConfig()
	fmt.Printf("🔑 Initialized Facebook config with token rotation\n")

	fmt.Println("📥 Fetching comments...")
	comments, err := fetchAllPostComments(postID, config)
	if err != nil {
		fmt.Printf("❌ Error fetching comments: %v\n", err)
		return
	}

	fmt.Printf("✅ Successfully fetched %d main comments\n", len(comments))

	mainCommentsCount := CountFacebookComments(comments)
	fmt.Printf("📊 Facebook Extraction Summary:\n")
	fmt.Printf("   📝 Main comments: %d\n", mainCommentsCount)
	fmt.Printf("📊 Total main comments: %d\n", mainCommentsCount)

	excelPath, err := exportFacebookCommentsToExcel(comments, urlArg)
	if err != nil {
		fmt.Printf("❌ Error exporting to Excel: %v\n", err)
		return
	}

	endTime := time.Now()
	actualDuration := endTime.Sub(startTime)
	minutes := int(actualDuration.Minutes())
	seconds := int(actualDuration.Seconds()) % 60

	fmt.Printf("⏱️ Extraction completed at: %s\n", endTime.Format("15:04:05"))
	fmt.Printf("🎯 Actual processing time: %d minutes %d seconds (%.1f seconds total)\n",
		minutes, seconds, actualDuration.Seconds())
	fmt.Printf("📂 Exported comments to Excel: %s\n", excelPath)
	fmt.Printf("💎 Need unlimited comments with replies?\n📧 Email: haronkibetrutoh@gmail.com\n📱 WhatsApp: +254718448461")
	fmt.Println()
}

func exportFacebookCommentsToExcel(comments []FacebookComment, sourceURL string) (string, error) {
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println("Error closing Excel file:", err)
		}
	}()

	sheetName := "Comments"
	index, err := f.NewSheet(sheetName)
	if err != nil {
		return "", fmt.Errorf("error creating sheet: %w", err)
	}
	f.SetActiveSheet(index)

	f.DeleteSheet("Sheet1")

	headers := []string{
		"Comment ID", "Author Name", "Author ID", "Comment Text", "Created Time",
		"Likes Count", "Reply Count", "Depth", "Is Reply", "Parent Comment ID", "Parent Author",
		"URL",
	}

	columnWidths := map[string]float64{
		"A": 20,
		"B": 25,
		"C": 20,
		"D": 60,
		"E": 20,
		"F": 12,
		"G": 12,
		"H": 10,
		"I": 10,
		"J": 20,
		"K": 25,
		"L": 40,
	}

	for col, width := range columnWidths {
		f.SetColWidth(sheetName, col, col, width)
	}

	headerStyle, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
			Size: 12,
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#E0E0E0"},
			Pattern: 1,
		},
		Border: []excelize.Border{
			{Type: "bottom", Color: "#000000", Style: 1},
		},
	})
	if err != nil {
		fmt.Println("Warning: Error creating header style:", err)
	}

	for i, header := range headers {
		cell := fmt.Sprintf("%s1", string(rune('A'+i)))
		f.SetCellValue(sheetName, cell, header)
		if headerStyle != 0 {
			f.SetCellStyle(sheetName, cell, cell, headerStyle)
		}
	}

	row := 2

	for _, comment := range comments {
		createdTime := time.Unix(comment.CreatedTime, 0).Format("2006-01-02 15:04:05")

		likesCount := 0
		if comment.Feedback.Reactors.CountReduced != "" {
			countStr := comment.Feedback.Reactors.CountReduced
			if strings.HasSuffix(countStr, "K") {
				baseStr := strings.TrimSuffix(countStr, "K")
				if base, err := strconv.ParseFloat(baseStr, 64); err == nil {
					likesCount = int(base * 1000)
				}
			} else {
				if count, err := strconv.Atoi(countStr); err == nil {
					likesCount = count
				}
			}
		}

		rowData := []any{
			comment.ID,
			comment.Author.Name,
			comment.Author.ID,
			comment.Body.Text,
			createdTime,
			likesCount,
			comment.Feedback.RepliesFields.TotalCount,
			comment.Depth,
			false, // Is Reply - always false for main comments
			"",    // Parent Comment ID - empty for main comments
			"",    // Parent Author - empty for main comments
			comment.Feedback.URL,
		}

		for i, value := range rowData {
			cell := fmt.Sprintf("%s%d", string(rune('A'+i)), row)
			f.SetCellValue(sheetName, cell, value)
		}

		row++
	}

	metaSheetName := "Metadata"
	_, err = f.NewSheet(metaSheetName)
	if err != nil {
		fmt.Println("Warning: Error creating metadata sheet:", err)
	} else {

		f.SetCellValue(metaSheetName, "A1", "Source URL")
		f.SetCellValue(metaSheetName, "B1", sourceURL)

		f.SetCellValue(metaSheetName, "A2", "Extraction Date")
		f.SetCellValue(metaSheetName, "B2", time.Now().Format("2006-01-02 15:04:05"))

		f.SetCellValue(metaSheetName, "A3", "Total Comments")
		f.SetCellValue(metaSheetName, "B3", len(comments))

		f.SetCellValue(metaSheetName, "A4", "Main Comments")
		f.SetCellValue(metaSheetName, "B4", len(comments))

		f.SetCellValue(metaSheetName, "A5", "Reply Comments")
		f.SetCellValue(metaSheetName, "B5", 0)
	}

	exportDir := "exports"
	if _, err := os.Stat(exportDir); os.IsNotExist(err) {
		if err := os.Mkdir(exportDir, 0755); err != nil {
			return "", fmt.Errorf("error creating exports directory: %w", err)
		}
	}

	timestamp := time.Now().Format("20060102_150405")
	filename := filepath.Join(exportDir, fmt.Sprintf("facebook_comments_%s.xlsx", timestamp))

	if err := f.SaveAs(filename); err != nil {
		return "", fmt.Errorf("error saving Excel file: %w", err)
	}

	fmt.Printf("✅ Exported %d comments to Excel file: %s\n", len(comments), filename)
	return filename, nil
}

func extractFacebookComments(data map[string]any) ([]FacebookComment, error) {

	dataObj, ok := data["data"].(map[string]any)
	if !ok {
		return nil, fmt.Errorf("data.data not found")
	}

	nodeObj, ok := dataObj["node"].(map[string]any)
	if !ok {
		return nil, fmt.Errorf("data.data.node not found")
	}

	renderingObj, ok := nodeObj["comment_rendering_instance_for_feed_location"].(map[string]any)
	if !ok {
		return nil, fmt.Errorf("comment_rendering_instance_for_feed_location not found")
	}

	commentsObj, ok := renderingObj["comments"].(map[string]any)
	if !ok {
		return nil, fmt.Errorf("comments not found")
	}

	edges, ok := commentsObj["edges"].([]any)
	if !ok {
		return nil, fmt.Errorf("edges not found")
	}

	var comments []FacebookComment

	for _, edge := range edges {
		edgeMap, ok := edge.(map[string]any)
		if !ok {
			continue
		}

		node, ok := edgeMap["node"].(map[string]any)
		if !ok {
			continue
		}

		var comment FacebookComment
		jsonBytes, err := json.Marshal(node)
		if err != nil {
			fmt.Printf("⚠️ Warning: Failed to marshal comment: %v\n", err)
			continue
		}

		err = json.Unmarshal(jsonBytes, &comment)
		if err != nil {
			fmt.Printf("⚠️ Warning: Failed to unmarshal comment: %v\n", err)
			continue
		}

		comments = append(comments, comment)
	}

	return comments, nil
}
