package main

import "testing"

func TestFormatHostFile(t *testing.T) {
	file := HostsFile{
		IpAddress: "127.0.0.1",
		HostNames: []HostNames{
			{HostName: "www.youtube.com"},
			{HostName: "www.google.com"},
		},
	}

	got := FormatHostFile(file)
	want := `
127.0.0.1 www.youtube.com
127.0.0.1 www.google.com
`
	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}