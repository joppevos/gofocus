package main

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

func TestSplitContent(t *testing.T) {
	content := []byte("www.youtube.com\nwww.google.com")

	got := splitContent(content)
	want := []string{"www.youtube.com", "www.google.com"}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %s, want %s", got, want)
	}
}

func TestFormatHostConfig(t *testing.T) {
	h := []string{"www.youtube.com", "www.google.com"}

	got := formatHostsConfig(h)
	want := fmt.Sprintf(`
%s www.youtube.com
%s www.google.com
`, IPAddress, IPAddress)

	AssertEqual(t, got, want)
}

func TestFormatMinutes(t *testing.T) {
	d := time.Duration(60) * time.Minute

	got := formatMinutes(d)
	want := "60:00"

	AssertEqual(t, got, want)
}

func AssertEqual(t *testing.T, got, want string) {
	if got != want {
		t.Errorf("got %s, want %s", got, want)
	}

}
