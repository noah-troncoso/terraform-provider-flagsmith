package flagsmith

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/Flagsmith/flagsmith-go-api-client"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ resource.Resource = &featureResource{}
var _ resource.ResourceWithImportState = &featureResource{}

func newFeatureResource() resource.Resource {
	return &featureResource{}
}

type featureResource struct {
	client *flagsmithapi.Client
}

func (r *featureResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_feature"
}


func (r *featureResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*flagsmithapi.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *flagsmithapi.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = client
}
func (t *featureResource) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Flagsmith Feature/ Remote config",

		Attributes: map[string]tfsdk.Attribute{
			"id": {
				Computed:            true,
				MarkdownDescription: "ID of the feature",
				PlanModifiers: tfsdk.AttributePlanModifiers{
					resource.UseStateForUnknown(),
				},
				Type: types.Int64Type,
			},
			"uuid": {
				Computed:            true,
				MarkdownDescription: "UUID of the feature",
				PlanModifiers: tfsdk.AttributePlanModifiers{
					resource.UseStateForUnknown(),
				},
				Type: types.StringType,
			},
			"project_id": {
				Computed:            true,
				MarkdownDescription: "ID of the project",
				PlanModifiers: tfsdk.AttributePlanModifiers{
					resource.UseStateForUnknown(),
				},
				Type: types.Int64Type,
			},
			"feature_name": {
				Required:            true,
				MarkdownDescription: "Name of the feature",
				Type:                types.StringType,
			},
			"type": {
				Required:            true,
				MarkdownDescription: "Type of the feature, can be STANDARD, or MULTIVARIATE",
				Type:                types.StringType,
			},
			"default_enabled": {
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Determines if the feature is enabled by default. If unspecified, it will default to false",
				Type:                types.BoolType,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					resource.UseStateForUnknown(),
					BoolDefaultModifier{Default: false},
				},
			},
			"initial_value": {
				Computed:            true,
				Optional:            true,
				MarkdownDescription: "Determines the initial value of the feature.",
				Type:                types.StringType,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					resource.UseStateForUnknown(),
				},
			},
			"description": {
				Optional:            true,
				MarkdownDescription: "Description of the feature",
				Type:                types.StringType,
			},
			"is_archived": {
				Computed:            true,
				Optional:            true,
				MarkdownDescription: "Can be used to archive/unarchive a feature. If unspecified, it will default to false",
				Type:                types.BoolType,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					resource.UseStateForUnknown(),
					BoolDefaultModifier{Default: false},
				},
			},
			"owners": {
				Optional:            true,
				Type:                types.SetType{ElemType: types.Int64Type},
				MarkdownDescription: "List of user IDs who are owners of the feature",
			},
			"project_uuid": {
				MarkdownDescription: "UUID of project the feature belongs to",
				Required:            true,
				Type:                types.StringType,
			},
		},
	}, nil
}



func (r *featureResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data FeatureResourceData

	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}
	clientFeature := data.ToClientFeature()

	// Create the feature
	err := r.client.CreateFeature(clientFeature)

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create feature, got error: %s", err))
		return
	}
	resourceData := MakeFeatureResourceDataFromClientFeature(clientFeature)

	diags = resp.State.Set(ctx, &resourceData)
	resp.Diagnostics.Append(diags...)
}

func (r *featureResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data FeatureResourceData
	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	// Early return if the state is wrong
	if diags.HasError() {
		return
	}

	feature, err := r.client.GetFeature(data.UUID.ValueString())
	if err != nil {
		panic(err)
	}
	resourceData := MakeFeatureResourceDataFromClientFeature(feature)

	diags = resp.State.Set(ctx, &resourceData)
	resp.Diagnostics.Append(diags...)

}

func (r *featureResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	//Get plan values
	var plan FeatureResourceData
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		tflog.Error(ctx, "Update: Error reading plan data")
		return
	}

	// Get current state
	var state FeatureResourceData
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		tflog.Error(ctx, "Update: Error reading state data")
		return
	}
	// Generate API request body from plan
	clientFeature := plan.ToClientFeature()

	err := r.client.UpdateFeature(clientFeature)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update feature, got error: %s", err))
		return
	}

	resourceData := MakeFeatureResourceDataFromClientFeature(clientFeature)

	// Update the state with the new values
	diags = resp.State.Set(ctx, &resourceData)

	resp.Diagnostics.Append(diags...)

}

func (r *featureResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Get current state
	var state FeatureResourceData
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		tflog.Error(ctx, "Delete: Error reading state data")
		return
	}
	// Generate API request body from plan
	clientFeature := state.ToClientFeature()

	err := r.client.DeleteFeature(*clientFeature.ProjectID, *clientFeature.ID)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete feature, got error: %s", err))
		return
	}
	resp.State.RemoveResource(ctx)

}
func (r *featureResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("uuid"), req, resp)
}
