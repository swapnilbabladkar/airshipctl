/*
 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

     https://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

package pull

import (
	"opendev.org/airship/airshipctl/pkg/config"
	"opendev.org/airship/airshipctl/pkg/document/repo"
	"opendev.org/airship/airshipctl/pkg/log"
)

// Pull clones repositories
func Pull(cfgFactory config.Factory, noCheckout bool) error {
	cfg, err := cfgFactory()
	if err != nil {
		return err
	}

	return cloneRepositories(cfg, noCheckout)
}

func cloneRepositories(cfg *config.Config, noCheckout bool) error {
	// Clone main repository
	currentManifest, err := cfg.CurrentContextManifest()
	log.Debugf("Reading current context manifest information from %s", cfg.LoadedConfigPath())
	if err != nil {
		return err
	}

	// Clone repositories
	for repoName, extraRepoConfig := range currentManifest.Repositories {
		err := extraRepoConfig.Validate()
		if err != nil {
			return err
		}
		repository, err := repo.NewRepository(currentManifest.TargetPath, extraRepoConfig)
		if err != nil {
			return err
		}
		log.Printf("Downloading %s repository %s from %s into %s",
			repoName, repository.Name, extraRepoConfig.URL(), currentManifest.TargetPath)
		err = repository.Download(noCheckout)
		if err != nil {
			return err
		}
		repository.Driver.Close()
	}

	return nil
}
