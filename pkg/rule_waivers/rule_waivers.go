// Copyright 2022 Fugue, Inc.
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

package rule_waivers

import (
	"strings"

	"github.com/fugue/regula/v2/pkg/loader"
	"github.com/fugue/regula/v2/pkg/reporter"
)

func ExactMatchOrWildcards(waiverElem string, resourceElem string) bool {
	var matchElem bool
	// if elem is fully escaped do a exact match
	if strings.HasPrefix(waiverElem, "`") && strings.HasSuffix(waiverElem, "`") {
		matchElem = strings.Trim(waiverElem, "`") == resourceElem
	} else {
		matchElem = Match(waiverElem, resourceElem)
	}
	return matchElem
}

type RuleWaiver struct {
	ID               string
	ResourceID       string
	ResourceProvider string
	ResourceTag      string
	ResourceType     string
	RuleID           string
}

// TODO: Add an interface for results/resources so we can use this both at
// runtime as in IaC.  Move out configs and replace it by a method in this
// interface.
func (waiver RuleWaiver) Match(
	configs loader.LoadedConfigurations,
	result reporter.RuleResult,
) bool {
	configFilepath := result.Filepath
	if strp := configs.ConfigurationPath(result.Filepath); strp != nil {
		configFilepath = *strp
	}

	return ExactMatchOrWildcards(waiver.ResourceID, result.ResourceID) &&
		ExactMatchOrWildcards(waiver.ResourceProvider, configFilepath) &&
		ExactMatchOrWildcards(waiver.ResourceType, result.ResourceType) &&
		ExactMatchOrWildcards(waiver.RuleID, result.RuleID) &&
		waiver.matchTags(result)
}

func (waiver RuleWaiver) matchTags(result reporter.RuleResult) bool {
	if waiver.ResourceTag == "*" {
		return true
	}

	escapeTag := func(str string) string {
		return strings.Replace(str, ":", "\\:", -1)
	}

	tags := []string{}
	for key, val := range result.ResourceTags {
		if str, ok := val.(string); ok {
			tags = append(tags, escapeTag(key)+":"+escapeTag(str))
		} else if val == nil {
			tags = append(tags, escapeTag(key))
		}
	}

	for _, tag := range tags {
		if ExactMatchOrWildcards(waiver.ResourceTag, tag) {
			return true
		}
	}

	return false
}

func ApplyRuleWaivers(
	configs loader.LoadedConfigurations,
	report *reporter.RegulaReport,
	waivers []RuleWaiver,
) {
	for i := range report.RuleResults {
		for _, waiver := range waivers {
			if waiver.Match(configs, report.RuleResults[i]) {
				report.RuleResults[i].RuleResult = "WAIVED"

				if waiver.ID != "" {
					report.RuleResults[i].ActiveWaivers = append(
						report.RuleResults[i].ActiveWaivers,
						waiver.ID,
					)
				}
			}
		}
	}

	report.RecomputeSummary()
}
