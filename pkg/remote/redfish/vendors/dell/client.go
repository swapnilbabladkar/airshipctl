// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package dell wraps the standard Redfish client in order to provide additional functionality required to perform
// actions on iDRAC servers.
package dell

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	redfishAPI "opendev.org/airship/go-redfish/api"
	redfishClient "opendev.org/airship/go-redfish/client"

	"opendev.org/airship/airshipctl/pkg/log"
	"opendev.org/airship/airshipctl/pkg/remote/ifc"
	"opendev.org/airship/airshipctl/pkg/remote/redfish"
)

const (
	// ClientType is used by other packages as the identifier of the Redfish client.
	ClientType           = "redfish-dell"
	endpointImportSysCFG = "%s/redfish/v1/Managers/%s/Actions/Oem/EID_674_Manager.ImportSystemConfiguration"
	vCDBootRequestBody   = `{
	    "ShareParameters": {
	        "Target": "ALL"
	    },
	    "ShutdownType": "NoReboot",
	    "ImportBuffer": "<SystemConfiguration>
	                       <Component FQDD=\"iDRAC.Embedded.1\">
	                         <Attribute Name=\"ServerBoot.1#BootOnce\">Enabled</Attribute>
	                         <Attribute Name=\"ServerBoot.1#FirstBootDevice\">VCD-DVD</Attribute>
	                       </Component>
	                     </SystemConfiguration>"
	}`
)

// Client is a wrapper around the standard airshipctl Redfish client. This allows vendor specific Redfish clients to
// override methods without duplicating the entire client.
type Client struct {
	username   string
	password   string
	redfishURL string
	redfish.Client
	RedfishAPI redfishAPI.RedfishAPI
	RedfishCFG *redfishClient.Configuration
}

type iDRACAPIRespErr struct {
	Err iDRACAPIErr `json:"error"`
}

type iDRACAPIErr struct {
	ExtendedInfo []iDRACAPIExtendedInfo `json:"@Message.ExtendedInfo"`
	Code         string                 `json:"code"`
	Message      string                 `json:"message"`
}

type iDRACAPIExtendedInfo struct {
	Message    string `json:"Message"`
	Resolution string `json:"Resolution,omitempty"`
}

// SetBootSourceByType sets the boot source of the ephemeral node to a virtual CD, "VCD-DVD".
func (c *Client) SetBootSourceByType(ctx context.Context) error {
	log.Debug("Setting boot device to 'VCD-DVD'.")
	managerID, err := redfish.GetManagerID(
		redfish.SetAuth(ctx, c.username, c.password),
		c.RedfishAPI, c.NodeID())
	if err != nil {
		log.Debugf("Failed to retrieve manager ID for node '%s'.", c.NodeID())
		return err
	}

	// NOTE(drewwalters96): Setting the boot device to a virtual media type requires an API request to the iDRAC
	// actions API. The request is made below using the same HTTP client used by the Redfish API and exposed by the
	// standard airshipctl Redfish client. Only iDRAC 9 >= 3.3 is supports this endpoint.
	url := fmt.Sprintf(endpointImportSysCFG, c.RedfishCFG.BasePath, managerID)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBufferString(vCDBootRequestBody))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	if len(c.password+c.username) != 0 {
		req.SetBasicAuth(c.username, c.password)
	}

	httpResp, err := c.RedfishCFG.HTTPClient.Do(req)
	if httpResp.StatusCode != http.StatusAccepted {
		body, ok := ioutil.ReadAll(httpResp.Body)
		if ok != nil {
			log.Debugf("Malformed iDRAC response: %s", body)
			return redfish.ErrRedfishClient{Message: "Unable to set boot device. Malformed iDRAC response."}
		}

		var iDRACResp iDRACAPIRespErr
		ok = json.Unmarshal(body, &iDRACResp)
		if ok != nil {
			log.Debugf("Malformed iDRAC response: %s", body)
			return redfish.ErrRedfishClient{Message: "Unable to set boot device. Malformed iDrac response."}
		}

		return redfish.ErrRedfishClient{
			Message: fmt.Sprintf("Unable to set boot device. %s", iDRACResp.Err.ExtendedInfo[0]),
		}
	} else if err != nil {
		return redfish.ErrRedfishClient{Message: fmt.Sprintf("Unable to set boot device. %v", err)}
	}

	log.Debug("Successfully set boot device.")
	defer httpResp.Body.Close()

	return nil
}

// RemoteDirect implements remote direct interface
func (c *Client) RemoteDirect(ctx context.Context, isoURL string) error {
	return redfish.RemoteDirect(ctx, isoURL, c.redfishURL, c)
}

// newClient returns a client with the capability to make Redfish requests.
func newClient(redfishURL string,
	insecure bool,
	useProxy bool,
	username string,
	password string,
	systemActionRetries int,
	systemRebootDelay int) (*Client, error) {
	genericClient, err := redfish.NewClient(redfishURL, insecure, useProxy, username, password,
		systemActionRetries, systemRebootDelay)
	if err != nil {
		return nil, err
	}

	c := &Client{username, password, redfishURL, *genericClient, genericClient.RedfishAPI, genericClient.RedfishCFG}

	return c, nil
}

// ClientFactory is a constructor for redfish ifc.Client implementation
var ClientFactory ifc.ClientFactory = func(redfishURL string,
	insecure bool,
	useProxy bool,
	username string,
	password string,
	systemActionRetries int,
	systemRebootDelay int) (ifc.Client, error) {
	return newClient(redfishURL, insecure, useProxy,
		username, password, systemActionRetries, systemRebootDelay)
}
