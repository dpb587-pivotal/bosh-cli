package localvm

import (
	"encoding/json"
	"io/ioutil"
	"os/exec"
	"strings"

	"github.com/cloudfoundry/bosh-agent/settings"
	"github.com/cloudfoundry/bosh-cli/cloud"
	"github.com/cloudfoundry/bosh-cli/cloud/fakes"
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	biproperty "github.com/cloudfoundry/bosh-utils/property"
)

type Cloud struct {
	fakes.FakeCloud // TODO actually implement interface to avoid unexpected calls and quiet errors

	Mbus string
}

var _ cloud.Cloud = &Cloud{}

// TODO use system.FileSystem and system.CmdRunner
func (c Cloud) CreateVM(
	agentID string,
	_ string,
	_ biproperty.Map,
	networksInterfaces map[string]biproperty.Map,
	env biproperty.Map,
) (string, error) {
	var networksInterfacesDst settings.Networks
	var envDst settings.Env

	err := c.remarshal(networksInterfaces, &networksInterfacesDst)
	if err != nil {
		return "", bosherr.WrapError(err, "remarshalling networks interfaces")
	}

	err = c.remarshal(env, &envDst)
	if err != nil {
		return "", bosherr.WrapError(err, "remarshalling env")
	}

	{ // write static settings
		s := settings.Settings{
			AgentID: agentID,
			Mbus:    c.Mbus,
			Blobstore: settings.Blobstore{
				Type: "local",
				Options: map[string]interface{}{
					"blobstore_path": "/var/vcap/micro_bosh/data/cache",
				},
			},
			Networks: networksInterfacesDst,
			VM: settings.VM{
				Name: "vm-#{agent_id}",
			},
			Env: envDst,
		}

		sb, err := json.Marshal(s)
		if err != nil {
			return "", bosherr.WrapError(err, "marshalling localvm-settings")
		}

		// for testing with docker on macos
		sb = []byte(strings.Replace(string(sb), `"preconfigured":false`, `"preconfigured":true`, -1))

		err = ioutil.WriteFile("/var/vcap/bosh/localvm-settings.json", sb, 0400)
		if err != nil {
			return "", bosherr.WrapError(err, "writing localvm-settings.json")
		}
	}

	{ // override infrastructure settings to read from local file
		infraBytes, err := ioutil.ReadFile("/var/vcap/bosh/agent.json")
		if err != nil {
			return "", bosherr.WrapError(err, "reading agent.json")
		}

		var infra map[string]interface{}

		err = json.Unmarshal(infraBytes, &infra)
		if err != nil {
			return "", bosherr.WrapError(err, "unmarshalling agent.json")
		}

		infra["Infrastructure"] = map[string]interface{}{
			"Settings": map[string]interface{}{
				"Sources": []map[string]string{
					{
						"Type":         "File",
						"SettingsPath": "/var/vcap/bosh/localvm-settings.json",
					},
				},
				"UseRegistry": false,
			},
		}

		infraBytes, err = json.Marshal(infra)
		if err != nil {
			return "", bosherr.WrapError(err, "marshalling agent.json")
		}

		err = ioutil.WriteFile("/var/vcap/bosh/agent.json", infraBytes, 0400)
		if err != nil {
			return "", bosherr.WrapError(err, "writing agent.json")
		}
	}

	cmd := exec.Command("sv", "restart", "agent")
	err = cmd.Run()
	if err != nil {
		return "", bosherr.WrapError(err, "restarting agent")
	}

	return "localhost", nil
}

func (c Cloud) remarshal(in interface{}, out interface{}) error {
	tmp, err := json.Marshal(in)
	if err != nil {
		return bosherr.WrapError(err, "marshalling")
	}

	err = json.Unmarshal(tmp, out)
	if err != nil {
		return bosherr.WrapError(err, "unmarshalling")
	}

	return nil
}
