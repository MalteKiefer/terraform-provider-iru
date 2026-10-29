package provider

import "github.com/hashicorp/terraform-plugin-framework/types"

// knownBool returns nil for a null or unknown value, so an attribute the
// configuration leaves to the server is not sent as false.
func knownBool(v types.Bool) *bool {
	if v.IsNull() || v.IsUnknown() {
		return nil
	}
	b := v.ValueBool()
	return &b
}

// knownString returns nil for a null or unknown value, so an attribute the
// configuration leaves to the server is not sent as "".
func knownString(v types.String) *string {
	if v.IsNull() || v.IsUnknown() {
		return nil
	}
	s := v.ValueString()
	return &s
}
