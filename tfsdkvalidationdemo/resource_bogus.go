package tfsdkvalidationdemo

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

const (
	minVal = 0
	maxVal = 2
)

type resourceBogusType struct{}

func (r resourceBogusType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			//"unvalidated": {
			//	Type:          types.Int64Type,
			//	Required:      false,
			//	PlanModifiers: tfsdk.AttributePlanModifiers{tfsdk.RequiresReplace()},
			//},
			"validated": {
				Type:          types.Int64Type,
				Required:      true,
				PlanModifiers: tfsdk.AttributePlanModifiers{tfsdk.RequiresReplace()},
				Validators:    []tfsdk.AttributeValidator{bogusValidator{}},
			},
		},
	}, nil
}

func (r resourceBogusType) NewResource(_ context.Context, p tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	return resourceBogus{
		p: *(p.(*provider)),
	}, nil
}

type resourceBogus struct {
	p provider
}

func (r resourceBogus) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	if !r.p.configured {
		resp.Diagnostics.AddError(
			"Provider not configured",
			"The provider hasn't been configured before apply, likely because it depends on an unknown value from another resource. This leads to weird stuff happening, so we'd prefer if you didn't do that. Thanks!",
		)
		return
	}

	// Retrieve values from plan
	var plan ResourceBogus
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set State
	diags = resp.State.Set(ctx, ResourceBogus{
		Unvalidated: types.Int64{Value: plan.Unvalidated.Value},
		Validated:   types.Int64{Value: plan.Validated.Value},
	})
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r resourceBogus) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
	// Get current state
	var state ResourceBogus
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Reset state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r resourceBogus) Update(_ context.Context, _ tfsdk.UpdateResourceRequest, _ *tfsdk.UpdateResourceResponse) {
	// No update method because Read() will never report a state change, only
	// resource existence (or not)
}

func (r resourceBogus) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	var state ResourceBogus
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Remove resource from state
	resp.State.RemoveResource(ctx)
}

func (r resourceBogus) ImportState(ctx context.Context, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateResponse) {
	//Save the import identifier in the id attribute
	tfsdk.ResourceImportStatePassthroughID(ctx, tftypes.NewAttributePath().WithAttributeName("id"), req, resp)
}

type bogusValidator struct {
	tfsdk.AttributeValidator
}

func (o bogusValidator) Description(_ context.Context) string {
	return fmt.Sprintf("Valid range is %d - %d", minVal, maxVal)
}

func (o bogusValidator) MarkdownDescription(_ context.Context) string {
	return fmt.Sprintf("Valid range is %d - %d", minVal, maxVal)
}

func (o bogusValidator) Validate(ctx context.Context, req tfsdk.ValidateAttributeRequest, resp *tfsdk.ValidateAttributeResponse) {
	var bogusAttr types.Int64
	diags := req.Config.GetAttribute(ctx, req.AttributePath, &bogusAttr)
	resp.Diagnostics.Append(diags...)

	if bogusAttr.Value < minVal || bogusAttr.Value > maxVal {
		resp.Diagnostics.AddError(
			"Value out of range",
			fmt.Sprintf("Valid range: %d - %d, value %d is out of range", minVal, maxVal, bogusAttr.Value),
		)
	}
}
