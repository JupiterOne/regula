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

package rules.k8s_seccomp_profile

import data.fugue
import data.k8s

__rego__metadoc__ := {
	"custom": {
		"controls": {"CIS-Kubernetes_v1.6.1": ["CIS-Kubernetes_v1.6.1_5.7.2"]},
		"severity": "Medium",
	},
	"description": "",
	"id": "FG_R00522",
	"title": "Ensure that the seccomp profile is set to docker/default in your pod definitions",
}

input_type = "k8s"

resource_type = "MULTIPLE"

resources = k8s.resources_with_pod_templates

seccomp_set(template) {
    annotations := template.metadata.annotations
	annotations["seccomp.security.alpha.kubernetes.io/pod"] == "docker/default"
}

seccomp_set(template) {
	annotations := template.metadata.annotations
	annotations["seccomp.security.alpha.kubernetes.io/pod"] == "runtime/default"
}

policy[j] {
	resource := resources[_]
	template := k8s.pod_template(resource)
	seccomp_set(template)
	j = fugue.allow_resource(resource)
}

policy[j] {
	resource := resources[_]
	template := k8s.pod_template(resource)
	not seccomp_set(template)
	j = fugue.deny_resource(resource)
}
