package main

import (
	"testing"
)

func TestSetUrl(t *testing.T) {
	c := CreateCrawler()
	c.setURL("ts4.travian.com")
	expected := "https://ts4.travian.com"

	if c.url != expected {
		t.Error("HTTPS not prepended to URL")
	}
}

func TestSetUrl1(t *testing.T) {
	c := CreateCrawler()
	c.setURL("https://ts4.travian.com")
	expected := "https://ts4.travian.com"

	if c.url != expected {
		t.Error("Existing HTTPS not in URL anymore")
	}
}
