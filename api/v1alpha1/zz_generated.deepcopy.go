//go:build !ignore_autogenerated
// +build !ignore_autogenerated

/*
Copyright 2023.

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

// Code generated by controller-gen. DO NOT EDIT.

package v1alpha1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CommandAction) DeepCopyInto(out *CommandAction) {
	*out = *in
	if in.Args != nil {
		in, out := &in.Args, &out.Args
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CommandAction.
func (in *CommandAction) DeepCopy() *CommandAction {
	if in == nil {
		return nil
	}
	out := new(CommandAction)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MatchedTagStatus) DeepCopyInto(out *MatchedTagStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MatchedTagStatus.
func (in *MatchedTagStatus) DeepCopy() *MatchedTagStatus {
	if in == nil {
		return nil
	}
	out := new(MatchedTagStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TagCopyPlugin) DeepCopyInto(out *TagCopyPlugin) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TagCopyPlugin.
func (in *TagCopyPlugin) DeepCopy() *TagCopyPlugin {
	if in == nil {
		return nil
	}
	out := new(TagCopyPlugin)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TagDockerBuildPlugin) DeepCopyInto(out *TagDockerBuildPlugin) {
	*out = *in
	if in.Commands != nil {
		in, out := &in.Commands, &out.Commands
		*out = make([]CommandAction, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TagDockerBuildPlugin.
func (in *TagDockerBuildPlugin) DeepCopy() *TagDockerBuildPlugin {
	if in == nil {
		return nil
	}
	out := new(TagDockerBuildPlugin)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TagReflector) DeepCopyInto(out *TagReflector) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TagReflector.
func (in *TagReflector) DeepCopy() *TagReflector {
	if in == nil {
		return nil
	}
	out := new(TagReflector)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *TagReflector) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TagReflectorList) DeepCopyInto(out *TagReflectorList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]TagReflector, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TagReflectorList.
func (in *TagReflectorList) DeepCopy() *TagReflectorList {
	if in == nil {
		return nil
	}
	out := new(TagReflectorList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *TagReflectorList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TagReflectorPlugin) DeepCopyInto(out *TagReflectorPlugin) {
	*out = *in
	if in.Copy != nil {
		in, out := &in.Copy, &out.Copy
		*out = new(TagCopyPlugin)
		**out = **in
	}
	if in.DockerBuild != nil {
		in, out := &in.DockerBuild, &out.DockerBuild
		*out = new(TagDockerBuildPlugin)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TagReflectorPlugin.
func (in *TagReflectorPlugin) DeepCopy() *TagReflectorPlugin {
	if in == nil {
		return nil
	}
	out := new(TagReflectorPlugin)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TagReflectorSpec) DeepCopyInto(out *TagReflectorSpec) {
	*out = *in
	out.Regex = in.Regex
	in.Action.DeepCopyInto(&out.Action)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TagReflectorSpec.
func (in *TagReflectorSpec) DeepCopy() *TagReflectorSpec {
	if in == nil {
		return nil
	}
	out := new(TagReflectorSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TagReflectorStatus) DeepCopyInto(out *TagReflectorStatus) {
	*out = *in
	if in.MatchedTags != nil {
		in, out := &in.MatchedTags, &out.MatchedTags
		*out = make(map[string]*MatchedTagStatus, len(*in))
		for key, val := range *in {
			var outVal *MatchedTagStatus
			if val == nil {
				(*out)[key] = nil
			} else {
				in, out := &val, &outVal
				*out = new(MatchedTagStatus)
				**out = **in
			}
			(*out)[key] = outVal
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TagReflectorStatus.
func (in *TagReflectorStatus) DeepCopy() *TagReflectorStatus {
	if in == nil {
		return nil
	}
	out := new(TagReflectorStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TagRegex) DeepCopyInto(out *TagRegex) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TagRegex.
func (in *TagRegex) DeepCopy() *TagRegex {
	if in == nil {
		return nil
	}
	out := new(TagRegex)
	in.DeepCopyInto(out)
	return out
}
