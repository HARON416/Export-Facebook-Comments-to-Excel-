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
	XFbDebug                   string
	CSPNonce                   string
	LSDToken                   string
	Cookies                    string
	ExpansionTokens            map[string]string
	LastUpdated                time.Time
	CommentsUserID             string
	CommentsRequestID          string
	CommentsHash               string
	CommentsRevision           string
	CommentsSession            string
	CommentsHashSessionID      string
	CommentsDynamic            string
	CommentsCSR                string
	CommentsHSDP               string
	CommentsHBLP               string
	CommentsSJSP               string
	CommentsFbDtsg             string
	CommentsJazoest            string
	CommentsSpinR              string
	CommentsSpinT              string
	RepliesUserID              string
	RepliesRequestID           string
	RepliesHash                string
	RepliesRevision            string
	RepliesSession             string
	RepliesHashSessionID       string
	RepliesDynamic             string
	RepliesCSR                 string
	RepliesHSDP                string
	RepliesHBLP                string
	RepliesSJSP                string
	RepliesFbDtsg              string
	RepliesJazoest             string
	RepliesSpinR               string
	RepliesSpinT               string
	NestedRepliesUserID        string
	NestedRepliesRequestID     string
	NestedRepliesHash          string
	NestedRepliesRevision      string
	NestedRepliesSession       string
	NestedRepliesHashSessionID string
	NestedRepliesDynamic       string
	NestedRepliesCSR           string
	NestedRepliesHSDP          string
	NestedRepliesHBLP          string
	NestedRepliesSJSP          string
	NestedRepliesFbDtsg        string
	NestedRepliesJazoest       string
	NestedRepliesSpinR         string
	NestedRepliesSpinT         string
}

func CountFacebookCommentsAndReplies(comments []FacebookComment) (int, int, int) {
	mainComments := 0
	totalReplies := 0

	for _, comment := range comments {
		if comment.CommentParent == nil && comment.CommentDirectParent == nil {
			mainComments++
		} else {
			totalReplies++
		}
	}

	totalAll := mainComments + totalReplies
	return mainComments, totalReplies, totalAll
}

func getDefaultFacebookConfig() *FacebookConfig {
	return &FacebookConfig{
		XFbDebug:                   "",
		CSPNonce:                   "",
		LSDToken:                   "dhL62Eyc-Qfy6VJ1VqAUGF",
		Cookies:                    "datr=W-nLZxRTui39vfU0vdb0SHEr; sb=XOnLZzjZUznCGJQE89_qDE7G; ps_l=1; ps_n=1; c_user=61575555323184; wd=1366x681; fr=1yNe7T5XzIJ1KN5oo.AWfU6GqGgl-LvHr5svl6aAikgx8NlmsbLCWfszOtaaA-W_Yitvk.BohXK-..AAA.0.0.BohXK-.AWewc7qC6fmh_M9g1vOY8yjs75E; xs=18%3Ay5FhAH33c2RlRw%3A2%3A1753204356%3A-1%3A-1%3A%3AAcXKc2wtnhLk8A5un026fuGGsti5sPd4ZS01W8_nKpk; presence=C%7B%22t3%22%3A%5B%5D%2C%22utc3%22%3A1753576569657%2C%22v%22%3A1%7D",
		ExpansionTokens:            make(map[string]string),
		LastUpdated:                time.Now(),
		CommentsUserID:             "61575555323184",
		CommentsRequestID:          "17",
		CommentsHash:               "20296.HYP%3Acomet_pkg.2.1...0",
		CommentsRevision:           "1025191199",
		CommentsSession:            "39tg6c%3Ahv39n2%3Aw8szb4",
		CommentsHashSessionID:      "7531553972170399721",
		CommentsDynamic:            "7xeUjGU5a5Q1ryaxG4Vp41twWwIxu13wFwhUKbgS3q2ibwNw9G2Saw8i2S1DwUx60GE3Qwb-q7oc81EEc87m221Fwgo9oO0-E4a3a4oaEnxO0Bo7O2l2Utwqo31wiE4u9x-3m1mzXw8W58jwGzEaE5e3ym2SU4i5oe8464-5pUfEe88o4Wm7-2K1yw9q2-awLyESE2KwwwOg2cwMwhEkxebwHwKG4UrwFg2fwxyo566k1FwgUjz89oeE4WVU-4FqwIK6E4-mEbUaU2wwgo620XEaUcGy8qxG",
		CommentsCSR:                "g9Y74bhLf3sh6l5gD5NYsGnllbEJiiHW5N5jR229eB5R4RYDexaFbv9sWVlZ9l-RPfVlv8ikybriAWmBuDqDWmFTFyES45UJDVmCAFlyGBKqEGh5BQFHGdh9VFpqp4UDxB28SF8GbHDzebimFaBKmlkFpKeDBGbl1O8x24u48CeLKjjAAyGCxWEpVE-5oSq5m6WBDwMyoK3CucG8wQzVV8hy89VUnJa6o6q8xTwIBwXwUzoqwCXgO27UnBCwQz-13xOcCK2uU-dwzwNU98K2e5EcoK2249FazU8Egxa7EdUizUak1swgEW3S3aE8UaEf8mw9qi6Eeo989Uhggxydx2rU4S0Pki1Fw2cEy0PE6W1dw55wYw5LwWw5Vw67CCigPa0iSexd1nzk2203_l00g6U2gw0Bzw0e0-0pm1fgey07wo3mw2aE2kw3jE35Hw15S0op04-y85ap06dwf-08qwLw64cdw3A829wzCw7WG0k-1Kw1q608_w0zUw1Sh03aU0GF05cK0btS0n50FwiU1rU1so0l8gfp2w3NE",
		CommentsHSDP:               "g4olEG4Fq8wCzAi2itaDagqAGF8hEoAxW7EyBiEF7f2OiQwz7I43F-CsokgljNT2fch4jHOhh2q649azkiyNrdi8NqiIR5AzF8EIqyliq5IGHh1chh78bhnFP4A5EyA22NA4KKZaPkMYB5h25EdSMVkH3ZEoCrphD9pGa9g-VWiBcuFpAfozQ4pJeXq8yP3yxKv8bpqoCbgh8UYiIiVLg-VdQigJa8mzJAlGyd8N9pD8RahavfkraR5Vtq5O6uq9amy35gyBwCglzp58mKaaZAJ7Eiy4AKtJh1WLwDzrG8glymgxF9Ru4p8HAqEwyUKbA8IHGAVQ5k9xl3p9qKENAfpUsAxea8mtpbJCADy_seyU8QiiUcHaeizppy3kO8q9Koxk2KAAXc9xxDmcKi8F1e9XyofU9EKEepEuAt2EhAleuaDUy1uAwYxN4zkh1SewiER2iwSgjo9CF8boiwdK5oKbwDwAz86q4o6C1nwmEb4362WbxuFE8i0iE566j192A5UfEmwkHwwwIyNMfotwnk2Zc80WwBxwd0ookwkZ0Nada0_E1co8E4GOwMw4izU522e1zwfefwq82Pg2Aw4nwJg1080EXxa18wZw4Hw4Twno22w6twv8aU1vUsw55wSw4sw8G0GE-1Aw49w41waO0HEdU0Gu084w8m0afwe214wfq0liq",
		CommentsHBLP:               "08fgG5Vo2no-1Awgk48eUy6o5i2G1GwhUvw4AwMwh9E429wc23y8wRxS2S4bwUwHx25k2y8zEsBzouwb62W2Ja5VE20Cwt8twpo6im1cg4GEvwkU2LwG-bxW4FU6W1vzrwwwq84m5850wK1ywaSbG1mwpEhw4swJwjEe9k1Lxi3a9xO1AwGwkE36wdS0x82tU4jo1m81houxC0P85m1gxefwde0wU-11xWu58lw960lWdx-223h0eG08-wbG11AyE1881dU5um0wE1nUfoao51wHwf210wjoswoE3ywSw4swvEK1GzrwJzXwo8lwvU7e0g60H8f8y1EyEK22ey80E-8w5SCwsEG1kwMG3m3GEaoowKgco490uoowGxC1nwe248W2i1TZk1oyE4m0VU42",
		CommentsSJSP:               "g4olEG4Fq8wCzAi2itaDagqAGF8hEoAyNa7EyBiEF7f2OiQpNcuMNeWhbhOAptuljPP5Wl_FV4heKQx9mlECgB6p9av4BcG8hOqiAbhp8Wicz8FaAtd4mElxibDBt1iniyChwTyVQ2NpoKcigS9Bp8_wGgO4bxqEy488UxDy8oyUrww8489IwJyk2K5EuQ68miyaHuQ45x6dh9O129GGF0wAG8kxDx9y23agym2u5oS48ppbhW4EwSu5VK2u3i9g-45x616xiqbglgC5k4EyUa88ocEy1lwKwSIwZ1qtzkZ8q9Koxk2KbCoC66togx11e9U23wVCwAgvAmcxi0Gotxx049qAw4sxm3-0mO0KEbEK5UaS0PQ1ywbKb70Zwv8bQMw1wxwd03qo4UE1gU3Ew5Iwbd01nLxa18w",
		CommentsFbDtsg:             "NAfvDFxvSqCXNn87L7j47AEjxnv2v_FAZqEjdx_D3XOEjeJUy3ov3_g%3A18%3A1753204356",
		CommentsJazoest:            "25520",
		CommentsSpinR:              "1025191199",
		CommentsSpinT:              "1753576559",
		RepliesUserID:              "61575555323184",
		RepliesRequestID:           "1u",
		RepliesHash:                "20296.HYP%3Acomet_pkg.2.1...0",
		RepliesRevision:            "1025191199",
		RepliesSession:             "x4pz4j%3Ahv39n2%3Aw8szb4",
		RepliesHashSessionID:       "7531553972170399721",
		RepliesDynamic:             "7xeUjGU5a5Q1ryaxG4Vp41twWwIxu13wFwhUKbgS3q2ibwNw9G2Saw8i2S1DwUx60GE3Qwb-q7oc81EEc87m221Fwgo9oO0-E4a3a4oaEnxO0Bo7O2l2Utwqo31wiE4u9x-3m1mzXw8W58jwGzEaE5e3ym2SU4i5oe8464-5pUfEe88o4Wm7-2K1yw9q2-awLyESE2KwwwOg2cwMwhEkxebwHwKG4UrwFg2fwxyo566k1FwgUjz89oeE4WVU-4FqwIK6E4-mEbUaU2wwgo620XEaUcGy8qxG",
		RepliesCSR:                 "g9Y74bhL174hYl2sn7NObORbNcAGZhshkl229eB5R4RPsW4GAJYBPHWvTlvIAP-lncikybviAWmBtnqDRmFtWoGum9A9XhCvBqtGloGHrGWyF4GGJaHKAuh4DCBBFAt2ueGt28Th8GirmiFGHCIBqHGmVpliBCUWu-GVBgsy8gx7x29zHXAQQV8GFFrDKVaCG9XpEGiuqWzqF6zdxiqFkAEhxW9yVbyA2iucG8wQzVV9EG8gGUnDyEOQEWazU46LBDAKi6e2Om2Sdx2eDx2dxqVojzeQp4Gez_UlBCAwyK8z-13xOcGXwDLzoSm7olxLwAyU8UmwNyU88gCAGfwyx24EsUdUizUak1rg4aewZwOG2e2G3O5E2mAxGi3i2i2u4k48ozogC-1dwcR4wqo3Vyog80hd5y87J0ko6W1dw55wYw5LwWxu0m61nw4LCCigPa1X81PwWzEjglUR0ww0_Rg041K0A809oU03wfw6lwjQ3Ew1U60RE0yG0B80QW0NqU4W043E1xA0jW8wkFA0lK360_U0xG2-0ogMS08zw5Iw8C2eq0vGE1jU6W08ew2no1d80z-05aU0YO07p40cHw2GA0kOU0JTo1sk2C1bw5Lw8S0UU0l8gfp2w3NE",
		RepliesHSDP:                "g4okkxamy89EV4wADiFOA6FaGi4q6984OQAx17f2OiQwz7I4mDWpNxh1lf7gDN2i8WsAkgCxx2iER4EImPkycmAHdhp8Wiab6EBkCxraGQgj4khO2QlWsN91q8F0wI6mKZaPkMYB5h25EdSMVkH3ZEoCrphD9peyykfJ6FakNAGQp3SxckkGiQBXq8yP4gwErDO2SmC9yQ4ief4H4KrQfKjt4Abiy5EXp5qeQz4BCszkF4FYZhIHknBREn8pVEAFq8cl2am2p1mdAkxqUEHSiQuxa8iiVSR47G-2udKFS5oBA8qitnx6iaV6G88KbyV2baWFet1l2olgSimHGcp3Su798jyy5DmiXpF9ULT3EK2d4AK3aOzAESmowRcy6yrC8l0HF9eP2oopRzbAyagjyuUC3-2qbG3Cq7F7gG4p5jDyF-8wnF8f8sh8R4gtzE4GdgAEdA4S2pGi2S4E3rxmbyU9U98O1Cx61FwlU5G2N0NwKyUnGq24w4G1hxAMigF1u3W5E5aU88b8Is3S7o5R0Lj20eE9oo3g66585fgciziwfW0j62a1aIEc814E-1gwzwoU3PzU6y0IQ0F815Ubk0g20aeUiwi8fo1aU1dU5S0wE1Do7O2K0n-781hodE1782awaGfwp812o10o2IwaW3u04H825w7Nw9S0U84i0ZE1mU",
		RepliesHBLP:                "08fgG5Vo2no-1AwgkfgeUy6o5i2G1GwhUvw4AwMwh9E429wc23y8wRxS2S4bwUwHx25k2y8zEsBzouwb62W2Ja5VE20Cwt8twpo6im1cgdEjG7U5e0HUaLyUuxau1KwnUSU886y4ocUkyE4kwK1ywaSbG1mwpEhw9a2C1vwJg8Uaoe9k1Lxi3a9xO1AwGwkE4W1TwdS0x82tU4jo1m81houxC0BoaEK1lwk8jzU3jw8efwgouDxi5o27yo3wxe1FzovxGm3h0eG08-wbG11AyE1883awr85um0wE1nUfoao51wwyEeU2Qwg84S786a0UEdE1787WbwqESUbo-U625o7-1Pwf_waO3O8wq8GbwwzEzw2z48w5SCwrWyE5i32EdoeGwFxy2V0NwgA1Vxt0GxDwlo3wx2ewAwt_l0m8G3-m0VU42",
		RepliesSJSP:                "g4okkxamy89EV4wADiFOA6FaGi4q6984OQAx17f2OiQAD4NX15F4FD65RVlffcGWl_FSAyeCQx9mgy9A9hCiiDN9jay4sCAF2QmieAz8OaiF7jh5G5okyVVngkBQEFAogCx-bDgb5ByUN93oClAz-2F38gK5GCUgxmdy6p5Q6bdj458QFqxsxk9wCO2S9gaUmxXgoxpa8GJXgUC4oR4D848CGGA22iExi6u4C88cF29o9UlzogxBAJ7Eiy3pUnCU9Usxu9g-45x616xiqbglgC5k4EyUa88ocEy1lwKwSIwZ1qtzkZ8q9Koxk2KbCoC66togx11e9U23wVCwAgvAmcxi0Gotxx049qAw4sxm3-0mO0KEbEK5UaS0PQ1ywbKb70Zwv8bQMw1wxwd03qo4UE1gU3Ew5Iwbd01nLxa18w",
		RepliesFbDtsg:              "NAfvDFxvSqCXNn87L7j47AEjxnv2v_FAZqEjdx_D3XOEjeJUy3ov3_g%3A18%3A1753204356",
		RepliesJazoest:             "25520",
		RepliesSpinR:               "1025191199",
		RepliesSpinT:               "1753576559",
		NestedRepliesUserID:        "61575555323184",
		NestedRepliesRequestID:     "w",
		NestedRepliesHash:          "20296.HYP%3Acomet_pkg.2.1...0",
		NestedRepliesRevision:      "1025191199",
		NestedRepliesSession:       "rcnfri%3A98dq38%3A41yfi5",
		NestedRepliesHashSessionID: "7531562106898030442",
		NestedRepliesDynamic:       "7xeUjGU5a5Q1ryaxG4Vp41twWwIxu13wFwhUKbgS3q2ibwNw9G2Saw8i2S1DwUx60GE5O0BU2_CxS320qa321Rwwwqo462mcwfG12wOx62G5Usw9m1YwBgK7o6C0Mo4G17yovwRwlE-U2exi4UaEW2G1jwUBwJK14xm3y11xfxmu3W3y261eBx_wHwoE2mwLyEbUGdG0HE88cA0z8c84q58jyUaUbGxe6Uak0zU8oC1hxB0qo4e4UO2m3G1eKufxamEbbxG1fBG2-2K0E8461wweW2K3aEy6Eqw",
		NestedRepliesCSR:           "g9Y74bhL1p3k9Nsf8IlbEJiiHR5N5jMCyjFhthDOsW4GBlYBPHmBZRnXnc_AtYx9i8JJ9vFqlWeF-BGtWoGdx1uVSvBqqiBl2qBKqEGhaFtaqWy94iuWBBFAjyu6k8zqAyEKKucUJ9qAGmVpliBCUWiFqyRgsy8gx7x29gS-VddeiaGq7GxDCzUlzpElokCGmXwMyoK2a5FUOEy3ifDAx68wDDxuQEpwNwQy8tUb9oeUvxydxG2rJ388-5VpEd8_wgUsz9HwDKfzo8Ucu2ibwzxq36bwwx2qiE-3e4EuwTxafwFg5J0gEW3S3aE8UaEf8mw9qi6Eeo989Uhggxydx2rU4S0Pki1Fw2cEy0PE6W1dw55wYw5LwWw5Vw67CCigPa0iSexd1nzk2203_l00g6U2gw0Bzw0e0-0pm1fgey07wo3mw2aE2kw3jE35Hw15S0op04-y85ap06dwf-08qwLw64cdw3A829wzCw7WG0k-1Kw16O0ji08_w0zUw1FC0OA0cHw2GA0kOU0JTo1sk2C1bw5Lw8S0UU0l8gfp2w3NE",
		NestedRepliesHSDP:          "g4p78aE9EV4wLid2OiGAAyq6984O8Ax17f3iQwz2215F-CsokgkcMDN2h6GsAkga458R4EIoxa2i1my8p6EBkCxvazA2kj8b8TFP4A5EyH227giWX8Hdj3Okl2G3tIelaM_q69CSkpOmqyykfJ6FalGuFpAfozQ4pJeXq8yP3yxKv8bpqoCbgh8UYiIiVLg-VdQigJa8mzJAlEXicimpOdiAiDPR6OJhunmxsxDCyiBEwNk8Fo9A5oikxqUEHSiQuxa8iiVSR47G-2udKFS5oBA8qitnx6iaV6G88KbyV2baWFet1l2olgSimHGcp3Su798jyy5DmiXpF9ULT3EK2d4AK3aOzAESmowRcy6yrC8l0HF9eP2oopRzbAyagjyuUC3-2qbG3Cq7F7gG4p5jDyF-8wnF8f8sh8R4gtzE4GdgAEdA4S2pGi2S4E3rxmbyU9U98O1Cx61FwlU5G2N0NwKyUnGq24w4G1hxAMigF1u3W5E5aU88b8Is3S7o5R0Lj20eE9oo3g66585fgciziwfW0j62a1aIEc814E-1gwzwoU3PzU6y0IQ0F815Ubk0g20aeUiwi8fo1aU1dU5S0wE1Do7O2K0n-781hodE1782awaGfwp812o10o2IwaW3u04H825w7Nw9S0U84i0ZE1mU",
		NestedRepliesHBLP:          "08fgG5-U2no-fwl8453Q3K8xC1kwGwqE9Uvx-0ii3214Cwg8C0M8e8y3m7obogK3y7KFUgxl0Ey8W79oS7E2NwKwHixuq0w9E7i7o6m1ABwj41aG7U5e0HUaLyUuxau1KwnUSU886y15xiawhi2U6a0HoKE5q1Cx60AE29wJwjEe9kU6K58cEC786i2G1iwjE7u0RVo24w9Twhdw5ow55xW6o30yU5m1gxefwde0wU-11xWu58lw960lWdx-223h0eG08-wbG11AyE1881dU5um0wE1nUfoao51wHwf210wjoswoE3ywSw4swvEK1GzrwJzXwo8lwvU7e0g60H8f8y1EyEK22ey80EW9w5SCwrWyE5i32EdoeGwFxy2V0NwgA1Vxt0GxDwlo3wx2ewAwt_l0m8G3-m0VU42",
		NestedRepliesSJSP:          "g4p78aE9EV4wLid2OiGAAyq6984O8Ax17f3iQAD4NX15F4FD65RVlffcJql_FSAhGCQx9ilq9A9hCiiDN9jay4sCAF2QmieAz8OaiF7jmBFU_yoGbQilQ59taFShwTyut0Imm9oN93oClAz-2F38gK5Gy8gxmdyFUy68K6U82122r8boB0Hxq7J1y5AEyGTJ3yohzkiswgyqGGg89ay58pUiowwOA8BwDxm7EppbhW4EwSu5VK2u78nykfx1ohwhEkCyQ5k9xl1a8K2y263a8wlobEdH8fgmDoRfi6yrC8l0HyVC9xxDm48ggjyu0wUepE947V5z8kwaC7oog12mF8178lw_w5IwbG2Wbxu2JwcZ0oE2XyNMfo7O2Zc80o8o3g0SC1ea0ke0W81r82Pg0lXUiwi8",
		NestedRepliesFbDtsg:        "NAfujEplFPkTHueyZIkbQl3PP6Mog3gByHaPJj-n6ZA_Yp_0VUitZkw%3A18%3A1753204356",
		NestedRepliesJazoest:       "25640",
		NestedRepliesSpinR:         "1025191199",
		NestedRepliesSpinT:         "1753578453",
	}
}

func updateFacebookConfigFromResponse(config *FacebookConfig, responseHeaders http.Header, responseBody string, apiType string) {
	if xFbDebug := responseHeaders.Get("X-Fb-Debug"); xFbDebug != "" {
		fmt.Printf("🔄 Updating X-Fb-Debug token: %s...\n", xFbDebug[:20])
		config.XFbDebug = xFbDebug
	}
	if csp := responseHeaders.Get("Content-Security-Policy"); csp != "" && strings.Contains(csp, "nonce-") {
		nonceRegex := regexp.MustCompile(`'nonce-([^']+)'`)
		if matches := nonceRegex.FindStringSubmatch(csp); len(matches) > 1 {
			newNonce := matches[1]
			fmt.Printf("🔄 Updating CSP nonce: %s\n", newNonce)
			config.CSPNonce = newNonce
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
			if apiType == "replies" {
				reqSequence := []string{"b", "c", "d", "e", "f", "g", "h", "i", "j"}
				currentIndex := -1
				for i, req := range reqSequence {
					if req == config.RepliesRequestID {
						currentIndex = i
						break
					}
				}
				if currentIndex >= 0 {
					nextIndex := (currentIndex + 1) % len(reqSequence)
					config.RepliesRequestID = reqSequence[nextIndex]
				} else {
					config.RepliesRequestID = newReq
				}
			} else if apiType == "nested_replies" {
				reqSequence := []string{"w", "x", "y", "z", "10", "11", "12", "13", "14"}
				currentIndex := -1
				for i, req := range reqSequence {
					if req == config.NestedRepliesRequestID {
						currentIndex = i
						break
					}
				}
				if currentIndex >= 0 {
					nextIndex := (currentIndex + 1) % len(reqSequence)
					config.NestedRepliesRequestID = reqSequence[nextIndex]
				} else {
					config.NestedRepliesRequestID = newReq
				}
			} else {
				reqSequence := []string{"17", "18", "19", "1a", "1b", "1c", "1d", "1e", "1f"}
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
				} else {
					config.CommentsRequestID = newReq
				}
			}
		}
	}

	if strings.Contains(responseBody, "__s") {
		sessionRegex := regexp.MustCompile(`"__s":"([^"]+)"`)
		if matches := sessionRegex.FindStringSubmatch(responseBody); len(matches) > 1 {
			newSession := matches[1]
			switch apiType {
			case "replies":
				if newSession != config.RepliesSession {
					sessionPreview := newSession
					if len(newSession) > 20 {
						sessionPreview = newSession[:20]
					}
					fmt.Printf("🔄 Updating Replies Session: %s...\n", sessionPreview)
					config.RepliesSession = newSession
				}
			case "nested_replies":
				if newSession != config.NestedRepliesSession {
					sessionPreview := newSession
					if len(newSession) > 20 {
						sessionPreview = newSession[:20]
					}
					fmt.Printf("🔄 Updating Nested Replies Session: %s...\n", sessionPreview)
					config.NestedRepliesSession = newSession
				}
			default:
				if newSession != config.CommentsSession {
					sessionPreview := newSession
					if len(newSession) > 20 {
						sessionPreview = newSession[:20]
					}
					fmt.Printf("🔄 Updating Comments Session: %s...\n", sessionPreview)
					config.CommentsSession = newSession
				}
			}
		}
	}

	paramMappings := map[string]struct {
		commentsField      *string
		repliesField       *string
		nestedRepliesField *string
		paramName          string
	}{
		"__hs":     {&config.CommentsHash, &config.RepliesHash, &config.NestedRepliesHash, "Hash"},
		"__rev":    {&config.CommentsRevision, &config.RepliesRevision, &config.NestedRepliesRevision, "Revision"},
		"__hsi":    {&config.CommentsHashSessionID, &config.RepliesHashSessionID, &config.NestedRepliesHashSessionID, "HashSessionID"},
		"__dyn":    {&config.CommentsDynamic, &config.RepliesDynamic, &config.NestedRepliesDynamic, "Dynamic"},
		"__csr":    {&config.CommentsCSR, &config.RepliesCSR, &config.NestedRepliesCSR, "CSR"},
		"__hsdp":   {&config.CommentsHSDP, &config.RepliesHSDP, &config.NestedRepliesHSDP, "HSDP"},
		"__hblp":   {&config.CommentsHBLP, &config.RepliesHBLP, &config.NestedRepliesHBLP, "HBLP"},
		"__sjsp":   {&config.CommentsSJSP, &config.RepliesSJSP, &config.NestedRepliesSJSP, "SJSP"},
		"fb_dtsg":  {&config.CommentsFbDtsg, &config.RepliesFbDtsg, &config.NestedRepliesFbDtsg, "FbDtsg"},
		"jazoest":  {&config.CommentsJazoest, &config.RepliesJazoest, &config.NestedRepliesJazoest, "Jazoest"},
		"__spin_r": {&config.CommentsSpinR, &config.RepliesSpinR, &config.NestedRepliesSpinR, "SpinR"},
		"__spin_t": {&config.CommentsSpinT, &config.RepliesSpinT, &config.NestedRepliesSpinT, "SpinT"},
	}

	for param, mapping := range paramMappings {
		if strings.Contains(responseBody, param) {
			regex := regexp.MustCompile(fmt.Sprintf(`"%s":"([^"]+)"`, param))
			if matches := regex.FindStringSubmatch(responseBody); len(matches) > 1 {
				newValue := matches[1]
				switch apiType {
				case "replies":
					if newValue != *mapping.repliesField {
						preview := newValue
						if len(newValue) > 20 {
							preview = newValue[:20]
						}
						fmt.Printf("🔄 Updating Replies %s: %s...\n", mapping.paramName, preview)
						*mapping.repliesField = newValue
					}
				case "nested_replies":
					if newValue != *mapping.nestedRepliesField {
						preview := newValue
						if len(newValue) > 20 {
							preview = newValue[:20]
						}
						fmt.Printf("🔄 Updating Nested Replies %s: %s...\n", mapping.paramName, preview)
						*mapping.nestedRepliesField = newValue
					}
				default:
					if newValue != *mapping.commentsField {
						preview := newValue
						if len(newValue) > 20 {
							preview = newValue[:20]
						}
						fmt.Printf("🔄 Updating Comments %s: %s...\n", mapping.paramName, preview)
						*mapping.commentsField = newValue
					}
				}
			}
		}
	}

	switch apiType {
	case "replies":
		fmt.Printf("🔄 Facebook Replies Token Rotation Summary:\n")
		sessionPreview := config.RepliesSession
		if len(config.RepliesSession) > 20 {
			sessionPreview = config.RepliesSession[:20]
		}
		fmt.Printf("   RequestID: %s, Session: %s...\n", config.RepliesRequestID, sessionPreview)
	case "nested_replies":
		fmt.Printf("🔄 Facebook Nested Replies Token Rotation Summary:\n")
		sessionPreview := config.NestedRepliesSession
		if len(config.NestedRepliesSession) > 20 {
			sessionPreview = config.NestedRepliesSession[:20]
		}
		fmt.Printf("   RequestID: %s, Session: %s...\n", config.NestedRepliesRequestID, sessionPreview)
	default:
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
	req.Header.Set("sec-ch-ua-full-version-list", "\"Not)A;Brand\";v=\"8.0.0.0\", \"Chromium\";v=\"138.0.7204.157\", \"Google Chrome\";v=\"138.0.7204.157\"")
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

	fmt.Printf("🔍 === Facebook Initial Comments Response ===\n")
	fmt.Printf("📡 Response Status: %s\n", resp.Status)
	fmt.Printf("📋 Response Headers:\n")
	for key, values := range resp.Header {
		for _, value := range values {
			fmt.Printf("    %s: %s\n", key, value)
		}
	}
	fmt.Printf("📄 Response Body Length: %d bytes\n", len(body))
	if len(body) < 2000 {
		fmt.Printf("📝 Response Body: %s\n", string(body))
	} else {
		fmt.Printf("📝 Response Body Sample (first 2000 chars): %s...\n", string(body[:2000]))
	}
	fmt.Printf("🔍 === End Facebook Initial Comments Response ===\n\n")

	updateFacebookConfigFromResponse(config, resp.Header, string(body), "comments")

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
	req.Header.Set("sec-ch-ua", "\"Google Chrome\";v=\"137\", \"Chromium\";v=\"137\", \"Not/A)Brand\";v=\"24\"")
	req.Header.Set("sec-ch-ua-full-version-list", "\"Google Chrome\";v=\"137.0.7151.119\", \"Chromium\";v=\"137.0.7151.119\", \"Not/A)Brand\";v=\"24.0.0.0\"")
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

	fmt.Printf("🔍 === Facebook Paginated Comments Response ===\n")
	fmt.Printf("📡 Response Status: %s\n", resp.Status)
	fmt.Printf("📋 Response Headers:\n")
	for key, values := range resp.Header {
		for _, value := range values {
			fmt.Printf("    %s: %s\n", key, value)
		}
	}
	fmt.Printf("📄 Response Body Length: %d bytes\n", len(body))
	if len(body) < 2000 {
		fmt.Printf("📝 Response Body: %s\n", string(body))
	} else {
		fmt.Printf("📝 Response Body Sample (first 2000 chars): %s...\n", string(body[:2000]))
	}
	fmt.Printf("🔍 === End Facebook Paginated Comments Response ===\n\n")

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
	fmt.Printf("🔍 Starting to fetch all main comments from post: %s\n", postID)

	var allFacebookComments []FacebookComment
	var currentCursor string
	pageCount := 1
	hasNextPage := true
	retryCount := 0
	maxRetries := 5
	mxc := 500

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

	allFacebookComments = append(allFacebookComments, pageFacebookComments...)
	fmt.Printf("✅ Page %d: Found %d comments (Total: %d)\n", pageCount, len(pageFacebookComments), len(allFacebookComments))

	if len(allFacebookComments) >= mxc {
		allFacebookComments = allFacebookComments[:mxc]
		fmt.Printf("🛑 Reached max comments limit (%d)\n", mxc)
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

		allFacebookComments = append(allFacebookComments, pageFacebookComments...)
		fmt.Printf("✅ Page %d: Found %d comments (Total: %d)\n", pageCount, len(pageFacebookComments), len(allFacebookComments))

		if len(allFacebookComments) >= mxc {
			allFacebookComments = allFacebookComments[:mxc]
			fmt.Printf("🛑 Reached max comments limit (%d)\n", mxc)
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

	mainComments, totalReplies, totalAll := CountFacebookCommentsAndReplies(allFacebookComments)
	fmt.Printf("📊 Facebook Main Comments Summary:\n")
	fmt.Printf("   📄 Pages fetched: %d\n", pageCount-1)
	fmt.Printf("   📝 Main comments: %d\n", mainComments)
	fmt.Printf("   💬 Total replies: %d\n", totalReplies)
	fmt.Printf("   🔢 Total items: %d\n", totalAll)
	fmt.Printf("📊 Total comments fetched: %d items from %d pages\n", totalAll, pageCount-1)

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

	var urlArg string

	if len(os.Args) > 1 {
		urlArg = os.Args[1]
	} else {

		fmt.Println("🚀 Facebook Comments & Replies Scraper")
		fmt.Println("=====================================")
		fmt.Println("Please enter a Facebook post URL or share link:")
		fmt.Print("➡️ ")

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

	fmt.Println("📥 Fetching replies for comments...")
	var allComments []FacebookComment
	allComments = append(allComments, comments...)

	mainComments, totalReplies, totalAll := CountFacebookCommentsAndReplies(allComments)
	fmt.Printf("📊 Facebook Extraction Summary:\n")
	fmt.Printf("   📝 Main comments: %d\n", mainComments)
	fmt.Printf("   💬 Total replies: %d\n", totalReplies)
	fmt.Printf("   🔢 Total items: %d\n", totalAll)
	fmt.Printf("📊 Total comments and replies: %d\n", totalAll)

	excelPath, err := exportFacebookCommentsToExcel(allComments, urlArg)
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
	fmt.Printf("Need more comments and replies?\nContact the developer:\nEmail: haronkibetrutoh@gmail.com\nWhatsApp: +254718448461")
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

	commentsByParent := make(map[string][]int)
	mainComments := []int{}

	for i, comment := range comments {
		if comment.CommentParent != nil || comment.CommentDirectParent != nil {
			var parentID string
			if comment.CommentParent != nil {
				parentID = comment.CommentParent.ID
			} else {
				parentID = comment.CommentDirectParent.ID
			}
			commentsByParent[parentID] = append(commentsByParent[parentID], i)
		} else {

			mainComments = append(mainComments, i)
		}
	}

	row := 2

	var writeCommentWithReplies func(commentIdx int, depth int) int
	writeCommentWithReplies = func(commentIdx int, depth int) int {
		comment := comments[commentIdx]

		createdTime := time.Unix(comment.CreatedTime, 0).Format("2006-01-02 15:04:05")

		isReply := comment.CommentParent != nil || comment.CommentDirectParent != nil

		var parentID, parentAuthor string
		if comment.CommentParent != nil {
			parentID = comment.CommentParent.ID
			parentAuthor = comment.CommentParent.Author.Name
		} else if comment.CommentDirectParent != nil {
			parentID = comment.CommentDirectParent.ID
			parentAuthor = comment.CommentDirectParent.Author.Name
		}

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
			isReply,
			parentID,
			parentAuthor,
			comment.Feedback.URL,
		}

		for i, value := range rowData {
			cell := fmt.Sprintf("%s%d", string(rune('A'+i)), row)
			f.SetCellValue(sheetName, cell, value)
		}

		if comment.Depth > 1 {
			fmt.Printf("📊 Excel: Writing depth=%d for comment by %s (ID: %s)\n",
				comment.Depth, comment.Author.Name, comment.ID)
		}

		if depth > 0 {

			indentStyle, _ := f.NewStyle(&excelize.Style{
				Fill: excelize.Fill{
					Type:    "pattern",
					Color:   []string{fmt.Sprintf("#F8F8F%d", 8-depth)},
					Pattern: 1,
				},
			})
			if indentStyle != 0 {
				f.SetCellStyle(sheetName, fmt.Sprintf("A%d", row), fmt.Sprintf("L%d", row), indentStyle)
			}
		}

		row++

		if replies, ok := commentsByParent[comment.ID]; ok {
			for _, replyIdx := range replies {
				row = writeCommentWithReplies(replyIdx, depth+1)
			}
		}

		return row
	}

	for _, commentIdx := range mainComments {
		row = writeCommentWithReplies(commentIdx, 0)
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
		f.SetCellValue(metaSheetName, "B4", len(mainComments))

		f.SetCellValue(metaSheetName, "A5", "Reply Comments")
		f.SetCellValue(metaSheetName, "B5", len(comments)-len(mainComments))
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
