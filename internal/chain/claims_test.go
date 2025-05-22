package chain

import (
	"encoding/base64"
	"fmt"
	"strings"
	"time"
)

func ExampleClaims() {
	authService, _ := prepareAuthService(nil)

	now := parseTime("2006-01-02T15:04:05Z")
	sessionID := "sid" + "b45f63a4-ed43-4605-a4fb-23a3af022d64"

	claims := &Claims{
		ExpiresAt: now.Add(time.Hour).Unix(),
		IssuedAt:  now.Unix(),

		Issuer:    "TODO",
		SessionID: sessionID,
		Nonce:     "2cd0bac4-4483-4224-9405-6b604c319858",
		Seq:       0,
		Siat:      now.Unix(),

		Info: struct {
			RealmID string `json:"realm_id"`
			UserID  string `json:"user_id"`
			Role    string `json:"role"`
		}{
			RealmID: "7d2a5354-657c-4b17-af49-ec59523050a8",
			UserID:  "dd0dfa9e-e7ae-4526-ab45-cdb8306b701c",
			Role:    "staff",
		},
	}

	printToken(authService.issueAccessToken(claims))
	printToken(authService.issueRefreshToken(claims))

	// Output:
	// {"alg":"RS256","kid":"acc54cd3d7f-9538-4663-9533-3085bf47c735","typ":"JWT"}
	// {"exp":1136217845,"iat":1136214245,"iss":"TODO","sid":"sidb45f63a4-ed43-4605-a4fb-23a3af022d64","nonce":"2cd0bac4-4483-4224-9405-6b604c319858","siat":1136214245,"info":{"realm_id":"7d2a5354-657c-4b17-af49-ec59523050a8","user_id":"dd0dfa9e-e7ae-4526-ab45-cdb8306b701c","role":"staff"}}
	// gk4rH_RwEB0BpiY-zd6C1g-qQwXxzMayHPkKQBfmg5jazzh1HSmSD_5Vt87FJnJbmMJuSZ-G9WzfjgpKP1_FKMTvZ491ISWtx3qQYDrrrKtyBLOA96N9QTXHPA5tb6WGYxZP-jv2E0p7BXB6z8Qny34fsAeKf-U5jVXyqKjgvBh6qd14OeGHTdxHYpFx002lRDqY_h3GDaHjgfAAUna0DsSqY_PiyL8v9Uj3skUFJxalugqDCbOU1kc5NTrHHB44kPG_ZJ8ts19cKH7lkJm2m2PgmOLX3LyHKLA1-ysQGgVUoffzgtCutY1PVtN8LSJI_UskzV9WYqE_8c2fofQ18g
	// ---
	// {"alg":"RS256","kid":"ref590aa1ac-b776-4c53-bcc9-99656842419b","typ":"JWT"}
	// {"exp":1136217845,"iat":1136214245,"iss":"TODO","sid":"sidb45f63a4-ed43-4605-a4fb-23a3af022d64","nonce":"2cd0bac4-4483-4224-9405-6b604c319858","siat":1136214245,"info":{"realm_id":"7d2a5354-657c-4b17-af49-ec59523050a8","user_id":"dd0dfa9e-e7ae-4526-ab45-cdb8306b701c","role":"staff"}}
	// bQgU7gSo1qul6tMfT0qX8gS6uxM9cVjDRtj0ToGQpfG9lx7JHOntPqzxlFMJSFDmYS1mnW6Xdn34cFB0zPPBvgqDZyjIBnzIk-OUH96H6pLuiwVY4I7xoUYrBeCfcXNaI8mWrN4oIfzrccenGZEL4ooJJmdLQgzXn8bUKY_Nd8gBc1PJaECB4tM8jc1Qt-STxart8BTs72WaQfjrpSlXoDqGRVbcCd_AHqTNN3-z0sgL70hyjEEYEmbsHrAAg6lrylImSVqaeDkqLKDnaQ6_NI1wOoqacjdy8BnWendCPB2542I4ftGDKL1qEHsESi69Nh51MV664-lLd_TVcc0abA
	// ---
}

func printToken(token string, err error) {
	if err != nil {
		panic(err)
	}

	tokenSplitted := strings.Split(token, ".")
	p0, _ := base64.RawStdEncoding.DecodeString(tokenSplitted[0])
	p1, _ := base64.RawStdEncoding.DecodeString(tokenSplitted[1])
	p2 := tokenSplitted[2]

	fmt.Printf("%s\n", p0)
	fmt.Printf("%s\n", p1)
	fmt.Printf("%s\n", p2)
	fmt.Printf("---\n")
}

func parseTime(value string) time.Time {
	parsedTime, err := time.Parse(time.RFC3339, value)
	if err != nil {
		panic(err)
	}

	return parsedTime
}
