package client

import (
	"bytes"
	"mime"
	"mime/multipart"
	"testing"
)

// Unset (nil) flags must be left out so an update does not switch a profile off.
func TestBuildProfileMultipart_OmitsNilFlags(t *testing.T) {
	on := true
	ct, body, err := buildProfileMultipart(CustomProfileInput{Name: "p", RunsOnMac: &on})
	if err != nil {
		t.Fatal(err)
	}
	_, params, err := mime.ParseMediaType(ct)
	if err != nil {
		t.Fatal(err)
	}
	form, err := multipart.NewReader(bytes.NewReader(body), params["boundary"]).ReadForm(1 << 20)
	if err != nil {
		t.Fatal(err)
	}
	if got := form.Value["runs_on_mac"]; len(got) != 1 || got[0] != "true" {
		t.Errorf("runs_on_mac = %v, want [true]", got)
	}
	for _, k := range []string{"active", "runs_on_iphone", "runs_on_ipad", "runs_on_tv"} {
		if v, ok := form.Value[k]; ok {
			t.Errorf("%s = %v sent although unset", k, v)
		}
	}
}
