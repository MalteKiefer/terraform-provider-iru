package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestCheckBlueprintType(t *testing.T) {
	cases := []struct {
		name    string
		prior   types.String
		planned types.String
		wantErr bool
	}{
		{"create map", types.StringNull(), types.StringValue("map"), false},
		{"create unset", types.StringNull(), types.StringUnknown(), false},
		{"create classic", types.StringNull(), types.StringValue("classic"), true},
		{"keep imported classic", types.StringValue("classic"), types.StringValue("classic"), false},
		{"replace map with classic", types.StringValue("map"), types.StringValue("classic"), true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var diags diag.Diagnostics
			checkBlueprintType(tc.prior, tc.planned, &diags)
			if diags.HasError() != tc.wantErr {
				t.Errorf("error = %v, want %v: %v", diags.HasError(), tc.wantErr, diags)
			}
		})
	}
}

// An enrollment_code block without is_active must not send is_active=false and so
// deactivate the code.
func TestBlueprintFromModel_UnsetIsActiveIsOmitted(t *testing.T) {
	for _, isActive := range []types.Bool{types.BoolUnknown(), types.BoolNull()} {
		ec, d := types.ObjectValue(enrollmentCodeAttrTypes, map[string]attr.Value{
			"code":      types.StringValue("1234"),
			"is_active": isActive,
		})
		if d.HasError() {
			t.Fatal(d)
		}
		b, diags := blueprintFromModel(context.Background(), blueprintModel{EnrollmentCode: ec})
		if diags.HasError() {
			t.Fatal(diags)
		}
		if b.EnrollmentCode == nil || !b.EnrollmentCode.OmitIsActive {
			t.Errorf("is_active %v: OmitIsActive not set: %+v", isActive, b.EnrollmentCode)
		}
	}
}
