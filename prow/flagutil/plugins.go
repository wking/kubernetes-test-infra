/*
Copyright 2018 The Kubernetes Authors.

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

package flagutil

import (
	"flag"

	"k8s.io/test-infra/prow/config"
	"k8s.io/test-infra/prow/plugins"
	"k8s.io/test-infra/prow/repoowners"
)

// PluginOptions holds options for interacting with Plugins.
type PluginOptions struct {
	configPath string
}

// AddFlags injects plugin options into the given FlagSet.
func (o *PluginOptions) AddFlags(fs *flag.FlagSet) {
	fs.StringVar(&o.configPath, "plugin-config", "/etc/plugins/plugins.yaml", "Path to plugin config file.")
}

// Validate validates plugin options.
func (o *PluginOptions) Validate(dryRun bool) error {
	return nil
}

// Agent returns a plugin agent.
func (o *PluginOptions) Agent(configAgent *config.Agent, pluginClient *plugins.PluginClient, start bool) (agent *plugins.PluginAgent, err error) {
	agent = &plugins.PluginAgent{}
	if pluginClient != nil {
		if configAgent != nil {
			pluginClient.OwnersClient = repoowners.NewClient(
				pluginClient.GitClient,
				pluginClient.GitHubClient,
				configAgent,
				agent.MDYAMLEnabled,
				agent.SkipCollaborators,
			)
		}
		agent.PluginClient = *pluginClient
	}

	if err := agent.Load(o.configPath); err != nil {
		return nil, err
	}

	if start {
		err = agent.Start(o.configPath)
		if err != nil {
			return nil, err
		}
	}
	return agent, err

	return agent, nil
}
