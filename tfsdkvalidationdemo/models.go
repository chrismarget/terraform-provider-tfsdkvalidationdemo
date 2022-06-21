package tfsdkvalidationdemo

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ResourceBogus struct {
	Unvalidated types.Int64 `tfsdk:"unvalidated"`
	Validated   types.Int64 `tfsdk:"validated"`
}
