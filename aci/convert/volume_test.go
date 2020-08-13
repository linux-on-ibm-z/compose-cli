/*
   Copyright 2020 Docker, Inc.

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

package convert

import (
	"testing"

	"github.com/compose-spec/compose-go/types"
	"gotest.tools/v3/assert"
)

const (
	storageAccountNameKey = "storage_account_name"
	shareNameKey          = "share_name"
)

func TestGetRunVolumes(t *testing.T) {
	volumeStrings := []string{
		"myuser1@myshare1:/my/path/to/target1",
		"myuser2@myshare2:/my/path/to/target2",
		"myuser3@mydefaultsharename", // Use default placement at '/run/volumes/<share_name>'
	}
	var goldenVolumeConfigs = map[string]types.VolumeConfig{
		"volume-0": {
			Name:   "volume-0",
			Driver: "azure_file",
			DriverOpts: map[string]string{
				storageAccountNameKey: "myuser1",
				shareNameKey:          "myshare1",
			},
		},
		"volume-1": {
			Name:   "volume-1",
			Driver: "azure_file",
			DriverOpts: map[string]string{
				storageAccountNameKey: "myuser2",
				shareNameKey:          "myshare2",
			},
		},
		"volume-2": {
			Name:   "volume-2",
			Driver: "azure_file",
			DriverOpts: map[string]string{
				storageAccountNameKey: "myuser3",
				shareNameKey:          "mydefaultsharename",
			},
		},
	}
	goldenServiceVolumeConfigs := []types.ServiceVolumeConfig{
		{
			Type:   "azure_file",
			Source: "volume-0",
			Target: "/my/path/to/target1",
		},
		{
			Type:   "azure_file",
			Source: "volume-1",
			Target: "/my/path/to/target2",
		},
		{
			Type:   "azure_file",
			Source: "volume-2",
			Target: "/run/volumes/mydefaultsharename",
		},
	}

	volumeConfigs, serviceVolumeConfigs, err := GetRunVolumes(volumeStrings)
	assert.NilError(t, err)
	for k, v := range volumeConfigs {
		assert.DeepEqual(t, goldenVolumeConfigs[k], v)
	}
	for i, v := range serviceVolumeConfigs {
		assert.DeepEqual(t, goldenServiceVolumeConfigs[i], v)
	}
}

func TestGetRunVolumesMissingFileShare(t *testing.T) {
	_, _, err := GetRunVolumes([]string{"myaccount@"})
	assert.ErrorContains(t, err, "does not include a storage file share after '@'")
}

func TestGetRunVolumesMissingUser(t *testing.T) {
	_, _, err := GetRunVolumes([]string{"@myshare"})
	assert.ErrorContains(t, err, "does not include a storage account before '@'")
}

func TestGetRunVolumesNoShare(t *testing.T) {
	_, _, err := GetRunVolumes([]string{"noshare"})
	assert.ErrorContains(t, err, "does not include a storage account before '@'")
}
