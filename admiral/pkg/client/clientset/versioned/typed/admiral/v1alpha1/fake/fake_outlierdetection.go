/*
Copyright The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	"context"
	json "encoding/json"
	"fmt"

	v1alpha1 "github.com/istio-ecosystem/admiral/admiral/pkg/apis/admiral/v1alpha1"
	admiralv1alpha1 "github.com/istio-ecosystem/admiral/admiral/pkg/client/applyconfiguration/admiral/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeOutlierDetections implements OutlierDetectionInterface
type FakeOutlierDetections struct {
	Fake *FakeAdmiralV1alpha1
	ns   string
}

var outlierdetectionsResource = v1alpha1.SchemeGroupVersion.WithResource("outlierdetections")

var outlierdetectionsKind = v1alpha1.SchemeGroupVersion.WithKind("OutlierDetection")

// Get takes name of the outlierDetection, and returns the corresponding outlierDetection object, and an error if there is any.
func (c *FakeOutlierDetections) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.OutlierDetection, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(outlierdetectionsResource, c.ns, name), &v1alpha1.OutlierDetection{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.OutlierDetection), err
}

// List takes label and field selectors, and returns the list of OutlierDetections that match those selectors.
func (c *FakeOutlierDetections) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.OutlierDetectionList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(outlierdetectionsResource, outlierdetectionsKind, c.ns, opts), &v1alpha1.OutlierDetectionList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.OutlierDetectionList{ListMeta: obj.(*v1alpha1.OutlierDetectionList).ListMeta}
	for _, item := range obj.(*v1alpha1.OutlierDetectionList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested outlierDetections.
func (c *FakeOutlierDetections) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(outlierdetectionsResource, c.ns, opts))

}

// Create takes the representation of a outlierDetection and creates it.  Returns the server's representation of the outlierDetection, and an error, if there is any.
func (c *FakeOutlierDetections) Create(ctx context.Context, outlierDetection *v1alpha1.OutlierDetection, opts v1.CreateOptions) (result *v1alpha1.OutlierDetection, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(outlierdetectionsResource, c.ns, outlierDetection), &v1alpha1.OutlierDetection{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.OutlierDetection), err
}

// Update takes the representation of a outlierDetection and updates it. Returns the server's representation of the outlierDetection, and an error, if there is any.
func (c *FakeOutlierDetections) Update(ctx context.Context, outlierDetection *v1alpha1.OutlierDetection, opts v1.UpdateOptions) (result *v1alpha1.OutlierDetection, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(outlierdetectionsResource, c.ns, outlierDetection), &v1alpha1.OutlierDetection{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.OutlierDetection), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeOutlierDetections) UpdateStatus(ctx context.Context, outlierDetection *v1alpha1.OutlierDetection, opts v1.UpdateOptions) (*v1alpha1.OutlierDetection, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(outlierdetectionsResource, "status", c.ns, outlierDetection), &v1alpha1.OutlierDetection{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.OutlierDetection), err
}

// Delete takes name of the outlierDetection and deletes it. Returns an error if one occurs.
func (c *FakeOutlierDetections) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteActionWithOptions(outlierdetectionsResource, c.ns, name, opts), &v1alpha1.OutlierDetection{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeOutlierDetections) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(outlierdetectionsResource, c.ns, listOpts)

	_, err := c.Fake.Invokes(action, &v1alpha1.OutlierDetectionList{})
	return err
}

// Patch applies the patch and returns the patched outlierDetection.
func (c *FakeOutlierDetections) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.OutlierDetection, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(outlierdetectionsResource, c.ns, name, pt, data, subresources...), &v1alpha1.OutlierDetection{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.OutlierDetection), err
}

// Apply takes the given apply declarative configuration, applies it and returns the applied outlierDetection.
func (c *FakeOutlierDetections) Apply(ctx context.Context, outlierDetection *admiralv1alpha1.OutlierDetectionApplyConfiguration, opts v1.ApplyOptions) (result *v1alpha1.OutlierDetection, err error) {
	if outlierDetection == nil {
		return nil, fmt.Errorf("outlierDetection provided to Apply must not be nil")
	}
	data, err := json.Marshal(outlierDetection)
	if err != nil {
		return nil, err
	}
	name := outlierDetection.Name
	if name == nil {
		return nil, fmt.Errorf("outlierDetection.Name must be provided to Apply")
	}
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(outlierdetectionsResource, c.ns, *name, types.ApplyPatchType, data), &v1alpha1.OutlierDetection{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.OutlierDetection), err
}

// ApplyStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating ApplyStatus().
func (c *FakeOutlierDetections) ApplyStatus(ctx context.Context, outlierDetection *admiralv1alpha1.OutlierDetectionApplyConfiguration, opts v1.ApplyOptions) (result *v1alpha1.OutlierDetection, err error) {
	if outlierDetection == nil {
		return nil, fmt.Errorf("outlierDetection provided to Apply must not be nil")
	}
	data, err := json.Marshal(outlierDetection)
	if err != nil {
		return nil, err
	}
	name := outlierDetection.Name
	if name == nil {
		return nil, fmt.Errorf("outlierDetection.Name must be provided to Apply")
	}
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(outlierdetectionsResource, c.ns, *name, types.ApplyPatchType, data, "status"), &v1alpha1.OutlierDetection{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.OutlierDetection), err
}
