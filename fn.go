package main

import (
	"context"

	"github.com/pkg/errors"
	awsv1beta1 "github.com/upbound/provider-aws/v2/apis/namespaced/s3/v1beta1"

	"github.com/crossplane/function-sdk-go/logging"
	fnv1 "github.com/crossplane/function-sdk-go/proto/v1"
	"github.com/crossplane/function-sdk-go/request"
	"github.com/crossplane/function-sdk-go/resource"
	"github.com/crossplane/function-sdk-go/resource/composed"
	"github.com/crossplane/function-sdk-go/response"
	"github.com/crossplane/learn/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type Function struct {
	fnv1.UnimplementedFunctionRunnerServiceServer

	log logging.Logger
}

func (f *Function) RunFunction(_ context.Context, req *fnv1.RunFunctionRequest) (*fnv1.RunFunctionResponse, error) {

	rsp := response.To(req, response.DefaultTTL)
	xr, err := request.GetObservedCompositeResource(req)
	if err != nil {
		f.log.Info(err.Error())
		response.Fatal(rsp, errors.Wrapf(err, "failed to get observed XR"))
		return rsp, nil
	}

	// bucket data safely stored using typed struct
	var buckets v1alpha1.XBuckets
	err = runtime.DefaultUnstructuredConverter.FromUnstructured(
		xr.Resource.Object,
		&buckets,
	)
	if err != nil {
		f.log.Info(err.Error())
		response.Fatal(rsp, errors.Wrapf(err, "failed to convert the object into type"))
		return rsp, nil
	}

	// get current desired composed
	desired, err := request.GetDesiredComposedResources(req)
	if err != nil {
		f.log.Debug(err.Error())
	}

	// register the type into known scheme, allowing composed method to transform the golang types into the equivalent unstructured object, which the resource desired composed can only accept
	awsv1beta1.AddToScheme(composed.Scheme)

	// create new bucket for each bucket name
	for _, name := range buckets.Spec.Names {
		newBucket := awsv1beta1.Bucket{}
		newBucket.ObjectMeta = metav1.ObjectMeta{Annotations: map[string]string{
			// set external name for the bucket name
			"crossplane.io/external-name": name,
		}}
		newBucket.Name = name
		newBucket.Spec.ForProvider.Region = &buckets.Spec.Region

		// create unstructured object of the s3 managed resources
		composedAwsBucket, err := composed.From(&newBucket)
		if err != nil {
			f.log.Info(err.Error())
		}

		// store it into the desired composed resource
		desired[resource.Name(name)] = &resource.DesiredComposed{Resource: composedAwsBucket}
	}

	// return the desired composed
	if err = response.SetDesiredComposedResources(rsp, desired); err != nil {
		f.log.Debug(err.Error())
	}

	// return crossplane managed resources
	response.ConditionTrue(rsp, "FunctionSuccess", "Success").
		TargetCompositeAndClaim()

	return rsp, nil
}
