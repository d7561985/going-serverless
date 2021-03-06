package function

import (
	"path/filepath"

	"github.com/pulumi/pulumi-gcp/sdk/v3/go/gcp/cloudfunctions"
	"github.com/pulumi/pulumi-gcp/sdk/v3/go/gcp/storage"
	"github.com/pulumi/pulumi/sdk/v2/go/pulumi"
)

// GCPFunctionConfig holds the various bits that we'll need to correctly
// configure a GCP Cloud Function
type GCPFunctionConfig struct {
	Name       string
	Path       string
	TalkPhase  string
	CodeBucket *storage.Bucket
	BlobBucket *storage.Bucket
}

// NewGCP takes a pulumi context, and a GCPFunctionConfig, and returns a
// GCP Cloud Function configured to serve the toy app
func NewGCP(ctx *pulumi.Context, cfg GCPFunctionConfig) (*cloudfunctions.Function, error) {
	codeObjectArgs := &storage.BucketObjectArgs{
		Bucket: cfg.CodeBucket.Name,
		Source: pulumi.NewFileArchive(cfg.Path),
	}
	codeObject, err := storage.NewBucketObject(ctx, filepath.Base(cfg.Path), codeObjectArgs)
	if err != nil {
		return nil, err
	}

	functionArgs := &cloudfunctions.FunctionArgs{
		SourceArchiveBucket: cfg.CodeBucket.Name,
		SourceArchiveObject: codeObject.Name,
		Runtime:             pulumi.String("go113"),
		EntryPoint:          pulumi.String("Handle"),
		TriggerHttp:         pulumi.Bool(true),
		AvailableMemoryMb:   pulumi.Int(128),
		Labels: pulumi.Map{
			"project":    pulumi.String("going-serverless-talk"),
			"talk-phase": pulumi.String(cfg.TalkPhase),
		},
		EnvironmentVariables: pulumi.Map{
			"BUCKET_NAME": cfg.BlobBucket.Name,
		},
	}

	function, err := cloudfunctions.NewFunction(ctx, cfg.Name, functionArgs)
	if err != nil {
		return nil, err
	}

	_, err = cloudfunctions.NewFunctionIamMember(ctx, "invoker", &cloudfunctions.FunctionIamMemberArgs{
		Project:       function.Project,
		Region:        function.Region,
		CloudFunction: function.Name,
		Role:          pulumi.String("roles/cloudfunctions.invoker"),
		Member:        pulumi.String("allUsers"),
	})
	if err != nil {
		return nil, err
	}
	return function, nil
}
