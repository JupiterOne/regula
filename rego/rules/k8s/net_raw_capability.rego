# Copyright 2020-2021 Fugue, Inc.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

package rules.k8s_net_raw_capability

import data.fugue

__rego__metadoc__ := {
	"custom": {
		"controls": {"CIS-Kubernetes_v1.6.1": ["CIS-Kubernetes_v1.6.1_5.2.7"]},
		"severity": "Medium",
	},
	"description": "",
	"id": "FG_R00513",
	"title": "Minimize the admission of containers with the NET_RAW capability",
}

input_type = "k8s"

resource_type = "MULTIPLE"

resources = fugue.resources("Pod")

is_valid(resource) {
    true
}

policy[j] {
	resource := resources[_]
	is_valid(resource)
	j = fugue.allow_resource(resource)
}

policy[j] {
	resource := resources[_]
	not is_valid(resource)
	j = fugue.deny_resource(resource)
}
