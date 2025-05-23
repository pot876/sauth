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

		Issuer:    "issuer_id",
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
	// {"exp":1136217845,"iat":1136214245,"iss":"issuer_id","sid":"sidb45f63a4-ed43-4605-a4fb-23a3af022d64","nonce":"2cd0bac4-4483-4224-9405-6b604c319858","siat":1136214245,"info":{"realm_id":"7d2a5354-657c-4b17-af49-ec59523050a8","user_id":"dd0dfa9e-e7ae-4526-ab45-cdb8306b701c","role":"staff"}}
	// SyxtHr4DGx85weyzmHDMomuEETcezJhgROpBrnyFcu0btlh5y5FefpFOwzrdbOKI6uvJtHdAeNcUuovfvFkO4r4BYoibcLlIedTJo8vGSdAHWywz3j1lbGardyIn0KAmY30lkJ0hZmwd8UkpKjPLpXu0MikNNGbBeXBadL-zg5GqGxy4lgdW2-Q0NojCBF1Z_6OPDDOOXUemUdFuOkH02qFbxMXHTdQJ0feOUlEKnJ_M_pNEELypj3GB9DN-FvCFoUL7quisAhJU_Ou22lVOiyiOvqrzOaMvteU7NAChvZ4bW-NcxJ2CMVi2a_deMh_vPwNEdtywuyjZtGK2QJpUxA
	// ---
	// {"alg":"RS256","kid":"ref590aa1ac-b776-4c53-bcc9-99656842419b","typ":"JWT"}
	// {"exp":1136217845,"iat":1136214245,"iss":"issuer_id","sid":"sidb45f63a4-ed43-4605-a4fb-23a3af022d64","nonce":"2cd0bac4-4483-4224-9405-6b604c319858","siat":1136214245,"info":{"realm_id":"7d2a5354-657c-4b17-af49-ec59523050a8","user_id":"dd0dfa9e-e7ae-4526-ab45-cdb8306b701c","role":"staff"}}
	// eZ0ks8RZbQfl9Gmvc0vF4OofnuBH5qB3nUgwT6lDUvgiiU58cuxVzRyNdlKD5YA2U_8sZd603joV6GIKeRpBja9fIx9pxegS885cPBBRcuLtNRvvI8_utGgGevB0fcJbHK7DUiJDrAUPEzqnNm-6deGEPMy6qWlrKDwyE2bN4edV1EpWEDQH-qIyxkDrwxieKLp29-8P4QdD25RCXnqwg_AYGvls2YQU7R4wzTak4QP0wjMqFk2GyKrtd_eDeqWWPamOQgUHUq8ARar5QhIOOHwmigyegKty2yYny33tfQXY1yDW4xe7vceJcWdd5TrsjWXxsxcyeZjdV_vaaIUZXA
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
