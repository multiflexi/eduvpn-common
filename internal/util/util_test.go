package util

import (
	"bytes"
	"testing"
	"time"
)

func Test_EnsureValidURL(t *testing.T) {
	_, validErr := EnsureValidURL("%notvalid%")

	if validErr == nil {
		t.Fatal("Got nil error, want: non-nil")
	}

	valid, validErr := EnsureValidURL("valid.com")
	if validErr != nil {
		t.Fatalf("Got: %v, want: nil", validErr)
	}

	afterValid := "https://valid.com"
	if valid != afterValid {
		t.Fatalf("Got: %v, want: %v", valid, afterValid)
	}

	valid, validErr = EnsureValidURL("http://valid.com")
	if validErr != nil {
		t.Fatalf("Got: %v, want: nil", validErr)
	}

	afterValid = "http://valid.com"
	if valid != afterValid {
		t.Fatalf("Got: %v, want: %v", valid, afterValid)
	}
}

func Test_MakeRandomByteSlice(t *testing.T) {
	random, randomErr := MakeRandomByteSlice(32)
	if randomErr != nil {
		t.Fatalf("Got: %v, want: nil", randomErr)
	}
	if len(random) != 32 {
		t.Fatalf("Got length: %d, want length: 32", len(random))
	}

	random2, randomErr2 := MakeRandomByteSlice(32)
	if randomErr2 != nil {
		t.Fatalf("2, Got: %v, want: nil", randomErr)
	}

	if bytes.Equal(random2, random) {
		t.Fatalf("Two random byteslices are the same: %v, %v", random2, random)
	}
}

func Test_GetCurrentTime(t *testing.T) {
	time_now := GetCurrentTime()

	time.Sleep(1 * time.Second)

	time_after_1_second := GetCurrentTime()

	if !time_after_1_second.After(time_now) {
		t.Fatal("Time is not after previous time")
	}
}

func Test_WAYFEncode(t *testing.T) {
	// AuthTemplate
	returnTo := "127.0.0.1:8000/test123bla/#wow   "

	// URL encoding but with spaces replace as + instead of %20
	wantReturnTo := "127.0.0.1%3A8000%2Ftest123bla%2F%23wow+++"
	encode := WAYFEncode(returnTo)
	if encode != wantReturnTo {
		t.Fatalf("Got: %s, want: %s", encode, wantReturnTo)
	}
}

func Test_ReplaceWAYF(t *testing.T) {
	// We expect url encoding but the spaces to be correctly replace with a + instead of a %20
	// And we expect that the return to and org_id are correctly replaced
	replaced := ReplaceWAYF(
		"@RETURN_TO@@ORG_ID@",
		"127.0.0.1:8000/&%$3#kM_-            ",
		"idp-test.nl.org/",
	)
	wantReplaced := "127.0.0.1%3A8000%2F%26%25%243%23kM_-++++++++++++idp-test.nl.org%2F"
	if replaced != wantReplaced {
		t.Fatalf("Got: %s, want: %s", replaced, wantReplaced)
	}

	// No RETURN_TO in template
	replaced = ReplaceWAYF("@ORG_ID@", "127.0.0.1:8000", "idp-test.nl.org/")
	wantReplaced = "127.0.0.1:8000"
	if replaced != wantReplaced {
		t.Fatalf("Got: %s, want: %s", replaced, wantReplaced)
	}

	// NO ORG_ID in template
	replaced = ReplaceWAYF("@RETURN_TO@", "127.0.0.1:8000", "idp-test.nl.org")
	wantReplaced = "127.0.0.1:8000"
	if replaced != wantReplaced {
		t.Fatalf("Got: %s, want: %s", replaced, wantReplaced)
	}

	// Template is empty
	replaced = ReplaceWAYF("", "127.0.0.1:8000", "idp-test.nl.org")
	wantReplaced = "127.0.0.1:8000"
	if replaced != wantReplaced {
		t.Fatalf("Got: %s, want: %s", replaced, wantReplaced)
	}

	// Template contains both @RETURN_TO@ and @ORG_ID@ but there is not enough to replace both
	replaced = ReplaceWAYF("@RETURN_TO@ORG_ID@", "127.0.0.1:8000", "idp-test.nl.org")
	wantReplaced = "127.0.0.1:8000"
	if replaced != wantReplaced {
		t.Fatalf("Got: %s, want: %s", replaced, wantReplaced)
	}
}

func Test_GetLanguageMatched(t *testing.T) {
	// func GetLanguageMatched(languageMap map[string]string, languageTag string) string {

	// exact match
	returned := GetLanguageMatched(map[string]string{"en": "test", "de": "test2"}, "en")
	if returned != "test" {
		t.Fatalf("Got: %s, want: %s", returned, "test")
	}

	// starts with language tag
	returned = GetLanguageMatched(map[string]string{"en-US-test": "test", "de": "test2"}, "en-US")
	if returned != "test" {
		t.Fatalf("Got: %s, want: %s", returned, "test")
	}

	// starts with en-
	returned = GetLanguageMatched(map[string]string{"en-UK": "test", "en": "test2"}, "en-US")
	if returned != "test" {
		t.Fatalf("Got: %s, want: %s", returned, "test")
	}

	// exact match for en
	returned = GetLanguageMatched(map[string]string{"de": "test", "en": "test2"}, "en-US")
	if returned != "test2" {
		t.Fatalf("Got: %s, want: %s", returned, "test2")
	}

	// We default to english
	returned = GetLanguageMatched(map[string]string{"es": "test", "en": "test2"}, "nl-NL")
	if returned != "test2" {
		t.Fatalf("Got: %s, want: %s", returned, "test2")
	}

	// We default to english with a - as well
	returned = GetLanguageMatched(map[string]string{"est": "test", "en-": "test2"}, "en-US")
	if returned != "test2" {
		t.Fatalf("Got: %s, want: %s", returned, "test2")
	}

	// None found just return one
	returned = GetLanguageMatched(map[string]string{"es": "test"}, "en-US")
	if returned != "test" {
		t.Fatalf("Got: %s, want: %s", returned, "test")
	}
}