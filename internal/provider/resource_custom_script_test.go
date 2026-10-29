package provider

import (
	"encoding/json"
	"testing"

	"github.com/dmcphearson/terraform-provider-iru/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// The API trims trailing whitespace from `script`. When the configured value differs
// from the server value only by trailing whitespace, we must keep the configured
// value in state so re-plan shows no diff.
func TestApplyCustomScriptResponse_PreservesTrailingNewline(t *testing.T) {
	prior := customScriptModel{Script: types.StringValue("echo hi\n")}
	out := &client.CustomScript{ID: "1", Name: "s", ExecutionFrequency: "once", Script: "echo hi"}

	var dst customScriptModel
	applyCustomScriptResponse(&dst, out, prior)

	if dst.Script.ValueString() != "echo hi\n" {
		t.Errorf("script = %q, want the configured value with trailing newline preserved", dst.Script.ValueString())
	}
}

// When the script body genuinely changed (not just trailing whitespace), adopt the
// server value.
func TestApplyCustomScriptResponse_AdoptsRealScriptChange(t *testing.T) {
	prior := customScriptModel{Script: types.StringValue("echo old\n")}
	out := &client.CustomScript{ID: "1", Script: "echo new"}

	var dst customScriptModel
	applyCustomScriptResponse(&dst, out, prior)

	if dst.Script.ValueString() != "echo new" {
		t.Errorf("script = %q, want server value adopted on real change", dst.Script.ValueString())
	}
}

// self_service_category_id / self_service_recommended are write-only (never returned
// by GET). The configured values must be carried forward, not wiped to empty.
func TestApplyCustomScriptResponse_PreservesWriteOnlySelfService(t *testing.T) {
	prior := customScriptModel{
		SelfServiceCategoryID:  types.StringValue("cat-123"),
		SelfServiceRecommended: types.BoolValue(true),
	}
	out := &client.CustomScript{ID: "1"} // API returns nothing for these

	var dst customScriptModel
	applyCustomScriptResponse(&dst, out, prior)

	if dst.SelfServiceCategoryID.ValueString() != "cat-123" {
		t.Errorf("self_service_category_id = %q, want carried forward cat-123", dst.SelfServiceCategoryID.ValueString())
	}
	if !dst.SelfServiceRecommended.ValueBool() {
		t.Errorf("self_service_recommended wiped; want carried forward true")
	}
}

// remediation_script is Optional+Computed: an unset config (null prior) adopts the
// API-returned "" so it settles on a known value (the API never returns null). This
// avoids a permanent "" -> null diff against state written with an empty string.
func TestApplyCustomScriptResponse_RemediationAdoptsEmptyString(t *testing.T) {
	prior := customScriptModel{RemediationScript: types.StringNull()}
	out := &client.CustomScript{ID: "1", RemediationScript: ""}

	var dst customScriptModel
	applyCustomScriptResponse(&dst, out, prior)

	if dst.RemediationScript.IsNull() || dst.RemediationScript.ValueString() != "" {
		t.Errorf("remediation_script = %v, want empty string adopted", dst.RemediationScript)
	}
}

// Unset (unknown) attributes must not be sent, so the server keeps its defaults, while
// known zero values must be sent so fields can be cleared or switched off.
func TestModelToCustomScript_SendsOnlyKnownValues(t *testing.T) {
	m := customScriptModel{
		Name:                   types.StringValue("s"),
		ExecutionFrequency:     types.StringValue("once"),
		Script:                 types.StringValue("echo"),
		RemediationScript:      types.StringValue(""),
		Active:                 types.BoolUnknown(),
		Restart:                types.BoolUnknown(),
		ShowInSelfService:      types.BoolValue(true),
		SelfServiceCategoryID:  types.StringUnknown(),
		SelfServiceRecommended: types.BoolValue(false),
	}
	b, err := json.Marshal(modelToCustomScript(m))
	if err != nil {
		t.Fatal(err)
	}
	var got map[string]any
	if err := json.Unmarshal(b, &got); err != nil {
		t.Fatal(err)
	}
	for _, k := range []string{"active", "restart", "self_service_category_id"} {
		if _, ok := got[k]; ok {
			t.Errorf("%s sent although unset: %s", k, b)
		}
	}
	if v, ok := got["remediation_script"]; !ok || v != "" {
		t.Errorf("remediation_script = %v, want \"\" sent to clear it: %s", v, b)
	}
	if v, ok := got["self_service_recommended"]; !ok || v != false {
		t.Errorf("self_service_recommended = %v, want false sent: %s", v, b)
	}
}

// The API rejects self_service_recommended unless Self Service is on, so turning Self
// Service off must not send it.
func TestModelToCustomScript_NoRecommendedWithoutSelfService(t *testing.T) {
	m := customScriptModel{
		ShowInSelfService:      types.BoolValue(false),
		SelfServiceRecommended: types.BoolValue(true),
	}
	if in := modelToCustomScript(m); in.SelfServiceRecommended != nil {
		t.Errorf("self_service_recommended = %v, want omitted when Self Service is off", *in.SelfServiceRecommended)
	}
}
