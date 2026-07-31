// SPDX-FileCopyrightText: 2026 Xquik contributors
//
// SPDX-License-Identifier: Apache-2.0

package releaseconfig

import "github.com/goccy/go-yaml"

type goreleaserConfig struct {
	Release struct {
		UseExistingDraft bool `yaml:"use_existing_draft"`
	} `yaml:"release"`
}

func usesExistingDraft(contents []byte) (bool, error) {
	var config goreleaserConfig
	if err := yaml.Unmarshal(contents, &config); err != nil {
		return false, err
	}
	return config.Release.UseExistingDraft, nil
}
