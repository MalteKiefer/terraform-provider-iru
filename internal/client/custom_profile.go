package client

import (
	"bytes"
	"context"
	"fmt"
	"mime/multipart"
	"strconv"
)

const customProfilesPath = "/api/v1/library/custom-profiles"

// CustomProfileInput carries the fields for create/update. ProfileXML is the raw
// .mobileconfig content, uploaded as the multipart `file` field. A nil bool is left
// out of the request so the server keeps its value.
type CustomProfileInput struct {
	Name         string
	ProfileXML   string
	Active       *bool
	RunsOnMac    *bool
	RunsOnIPhone *bool
	RunsOnIPad   *bool
	RunsOnTV     *bool
	IncludeFile  bool // when false (update without a file change), omit the file part
}

func buildProfileMultipart(in CustomProfileInput) (string, []byte, error) {
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)

	_ = w.WriteField("name", in.Name)
	for _, f := range []struct {
		name string
		v    *bool
	}{
		{"active", in.Active},
		{"runs_on_mac", in.RunsOnMac},
		{"runs_on_iphone", in.RunsOnIPhone},
		{"runs_on_ipad", in.RunsOnIPad},
		{"runs_on_tv", in.RunsOnTV},
	} {
		if f.v != nil {
			_ = w.WriteField(f.name, strconv.FormatBool(*f.v))
		}
	}

	if in.IncludeFile {
		fw, err := w.CreateFormFile("file", in.Name+".mobileconfig")
		if err != nil {
			return "", nil, err
		}
		if _, err := fw.Write([]byte(in.ProfileXML)); err != nil {
			return "", nil, err
		}
	}

	if err := w.Close(); err != nil {
		return "", nil, err
	}
	return w.FormDataContentType(), buf.Bytes(), nil
}

// CreateCustomProfile uploads a new custom profile.
func (c *Client) CreateCustomProfile(ctx context.Context, in CustomProfileInput) (*CustomProfile, error) {
	in.IncludeFile = true
	ct, body, err := buildProfileMultipart(in)
	if err != nil {
		return nil, fmt.Errorf("building multipart body: %w", err)
	}
	var out CustomProfile
	if err := c.DoRaw(ctx, "POST", customProfilesPath, ct, body, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetCustomProfile fetches a profile by id. Returns NotFoundError on 404.
func (c *Client) GetCustomProfile(ctx context.Context, id string) (*CustomProfile, error) {
	var out CustomProfile
	if err := c.DoJSON(ctx, "GET", customProfilesPath+"/"+id, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdateCustomProfile patches a profile. Pass includeFile=true to replace the file.
func (c *Client) UpdateCustomProfile(ctx context.Context, id string, in CustomProfileInput) (*CustomProfile, error) {
	ct, body, err := buildProfileMultipart(in)
	if err != nil {
		return nil, fmt.Errorf("building multipart body: %w", err)
	}
	var out CustomProfile
	if err := c.DoRaw(ctx, "PATCH", customProfilesPath+"/"+id, ct, body, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteCustomProfile removes a profile by id.
func (c *Client) DeleteCustomProfile(ctx context.Context, id string) error {
	return c.DoJSON(ctx, "DELETE", customProfilesPath+"/"+id, nil, nil)
}

// ListCustomProfiles returns all custom profiles (paginated).
func (c *Client) ListCustomProfiles(ctx context.Context) ([]CustomProfile, error) {
	return listAll[CustomProfile](ctx, c, customProfilesPath)
}
