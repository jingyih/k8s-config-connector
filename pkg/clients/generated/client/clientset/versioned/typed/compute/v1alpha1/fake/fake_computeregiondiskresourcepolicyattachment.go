// Copyright 2020 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// *** DISCLAIMER ***
// Config Connector's go-client for CRDs is currently in ALPHA, which means
// that future versions of the go-client may include breaking changes.
// Please try it out and give us feedback!

// Code generated by main. DO NOT EDIT.

package fake

import (
	"context"

	v1alpha1 "github.com/GoogleCloudPlatform/k8s-config-connector/pkg/clients/generated/apis/compute/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeComputeRegionDiskResourcePolicyAttachments implements ComputeRegionDiskResourcePolicyAttachmentInterface
type FakeComputeRegionDiskResourcePolicyAttachments struct {
	Fake *FakeComputeV1alpha1
	ns   string
}

var computeregiondiskresourcepolicyattachmentsResource = schema.GroupVersionResource{Group: "compute.cnrm.cloud.google.com", Version: "v1alpha1", Resource: "computeregiondiskresourcepolicyattachments"}

var computeregiondiskresourcepolicyattachmentsKind = schema.GroupVersionKind{Group: "compute.cnrm.cloud.google.com", Version: "v1alpha1", Kind: "ComputeRegionDiskResourcePolicyAttachment"}

// Get takes name of the computeRegionDiskResourcePolicyAttachment, and returns the corresponding computeRegionDiskResourcePolicyAttachment object, and an error if there is any.
func (c *FakeComputeRegionDiskResourcePolicyAttachments) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.ComputeRegionDiskResourcePolicyAttachment, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(computeregiondiskresourcepolicyattachmentsResource, c.ns, name), &v1alpha1.ComputeRegionDiskResourcePolicyAttachment{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.ComputeRegionDiskResourcePolicyAttachment), err
}

// List takes label and field selectors, and returns the list of ComputeRegionDiskResourcePolicyAttachments that match those selectors.
func (c *FakeComputeRegionDiskResourcePolicyAttachments) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.ComputeRegionDiskResourcePolicyAttachmentList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(computeregiondiskresourcepolicyattachmentsResource, computeregiondiskresourcepolicyattachmentsKind, c.ns, opts), &v1alpha1.ComputeRegionDiskResourcePolicyAttachmentList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.ComputeRegionDiskResourcePolicyAttachmentList{ListMeta: obj.(*v1alpha1.ComputeRegionDiskResourcePolicyAttachmentList).ListMeta}
	for _, item := range obj.(*v1alpha1.ComputeRegionDiskResourcePolicyAttachmentList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested computeRegionDiskResourcePolicyAttachments.
func (c *FakeComputeRegionDiskResourcePolicyAttachments) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(computeregiondiskresourcepolicyattachmentsResource, c.ns, opts))

}

// Create takes the representation of a computeRegionDiskResourcePolicyAttachment and creates it.  Returns the server's representation of the computeRegionDiskResourcePolicyAttachment, and an error, if there is any.
func (c *FakeComputeRegionDiskResourcePolicyAttachments) Create(ctx context.Context, computeRegionDiskResourcePolicyAttachment *v1alpha1.ComputeRegionDiskResourcePolicyAttachment, opts v1.CreateOptions) (result *v1alpha1.ComputeRegionDiskResourcePolicyAttachment, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(computeregiondiskresourcepolicyattachmentsResource, c.ns, computeRegionDiskResourcePolicyAttachment), &v1alpha1.ComputeRegionDiskResourcePolicyAttachment{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.ComputeRegionDiskResourcePolicyAttachment), err
}

// Update takes the representation of a computeRegionDiskResourcePolicyAttachment and updates it. Returns the server's representation of the computeRegionDiskResourcePolicyAttachment, and an error, if there is any.
func (c *FakeComputeRegionDiskResourcePolicyAttachments) Update(ctx context.Context, computeRegionDiskResourcePolicyAttachment *v1alpha1.ComputeRegionDiskResourcePolicyAttachment, opts v1.UpdateOptions) (result *v1alpha1.ComputeRegionDiskResourcePolicyAttachment, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(computeregiondiskresourcepolicyattachmentsResource, c.ns, computeRegionDiskResourcePolicyAttachment), &v1alpha1.ComputeRegionDiskResourcePolicyAttachment{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.ComputeRegionDiskResourcePolicyAttachment), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeComputeRegionDiskResourcePolicyAttachments) UpdateStatus(ctx context.Context, computeRegionDiskResourcePolicyAttachment *v1alpha1.ComputeRegionDiskResourcePolicyAttachment, opts v1.UpdateOptions) (*v1alpha1.ComputeRegionDiskResourcePolicyAttachment, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(computeregiondiskresourcepolicyattachmentsResource, "status", c.ns, computeRegionDiskResourcePolicyAttachment), &v1alpha1.ComputeRegionDiskResourcePolicyAttachment{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.ComputeRegionDiskResourcePolicyAttachment), err
}

// Delete takes name of the computeRegionDiskResourcePolicyAttachment and deletes it. Returns an error if one occurs.
func (c *FakeComputeRegionDiskResourcePolicyAttachments) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteActionWithOptions(computeregiondiskresourcepolicyattachmentsResource, c.ns, name, opts), &v1alpha1.ComputeRegionDiskResourcePolicyAttachment{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeComputeRegionDiskResourcePolicyAttachments) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(computeregiondiskresourcepolicyattachmentsResource, c.ns, listOpts)

	_, err := c.Fake.Invokes(action, &v1alpha1.ComputeRegionDiskResourcePolicyAttachmentList{})
	return err
}

// Patch applies the patch and returns the patched computeRegionDiskResourcePolicyAttachment.
func (c *FakeComputeRegionDiskResourcePolicyAttachments) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.ComputeRegionDiskResourcePolicyAttachment, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(computeregiondiskresourcepolicyattachmentsResource, c.ns, name, pt, data, subresources...), &v1alpha1.ComputeRegionDiskResourcePolicyAttachment{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.ComputeRegionDiskResourcePolicyAttachment), err
}
