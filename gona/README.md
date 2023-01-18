# gona-dev
`import "github.com/netactuate/gona/gona"`

* [Overview](#pkg-overview)
* [Index](#pkg-index)

## <a name="pkg-overview">Overview</a>
Package gona provides a simple golang interface to the NetActuate
Rest API at <a href="https://vapi2.netactuate.com/">https://vapi2.netactuate.com/</a>




## <a name="pkg-index">Index</a>
* [Constants](#pkg-constants)
* [func GetKeyFromEnv() string](#GetKeyFromEnv)
* [type BuildServerRequest](#BuildServerRequest)
* [type Client](#Client)
  * [func NewClient(apikey string) *Client](#NewClient)
  * [func NewClientCustom(apikey string, apiurl string) *Client](#NewClientCustom)
  * [func (c *Client) BuildServer(id int, r *BuildServerRequest) (b ServerBuild, err error)](#Client.BuildServer)
  * [func (c *Client) CreateSSHKey(name, key string) (sshkey SSHKey, err error)](#Client.CreateSSHKey)
  * [func (c *Client) CreateServer(r *CreateServerRequest) (b ServerBuild, err error)](#Client.CreateServer)
  * [func (c *Client) DeleteSSHKey(id int) error](#Client.DeleteSSHKey)
  * [func (c *Client) DeleteServer(id int, cancelBilling bool) error](#Client.DeleteServer)
  * [func (c *Client) GetLocations() ([]Location, error)](#Client.GetLocations)
  * [func (c *Client) GetOSs() ([]OS, error)](#Client.GetOSs)
  * [func (c *Client) GetSSHKey(id int) (sshkey SSHKey, err error)](#Client.GetSSHKey)
  * [func (c *Client) GetSSHKeys() (keys []SSHKey, err error)](#Client.GetSSHKeys)
  * [func (c *Client) GetServer(id int) (server Server, err error)](#Client.GetServer)
  * [func (c *Client) GetServers() ([]Server, error)](#Client.GetServers)
  * [func (c *Client) UnlinkServer(id int) error](#Client.UnlinkServer)
* [type CreateServerRequest](#CreateServerRequest)
* [type Location](#Location)
* [type OS](#OS)
* [type SSHKey](#SSHKey)
* [type Server](#Server)
* [type ServerBuild](#ServerBuild)


#### <a name="pkg-files">Package files</a>
[client.go](/src/target/client.go) [locations.go](/src/target/locations.go) [os.go](/src/target/os.go) [servers.go](/src/target/servers.go) [sshkeys.go](/src/target/sshkeys.go) 


## <a name="pkg-constants">Constants</a>
``` go
const (
    Version      = "0.2.0"
    BaseEndpoint = "https://vapi2.netactuate.com/api/"
    ContentType  = "application/json"
)
```
Version, BaseEndpoint, ContentType constants




## <a name="GetKeyFromEnv">func</a> [GetKeyFromEnv](/src/target/client.go?s=775:802#L38)
``` go
func GetKeyFromEnv() string
```
GetKeyFromEnv is a simple function to try to yank the value for
"NA_API_KEY" from the environment




## <a name="BuildServerRequest">type</a> [BuildServerRequest](/src/target/servers.go?s=2907:3337#L86)
``` go
type BuildServerRequest struct {
    Location      int    `url:"location,omitempty"`
    Image         int    `url:"image,omitempty"`
    FQDN          string `url:"fqdn,omitempty"`
    SSHKey        string `url:"ssh_key,omitempty"`
    SSHKeyID      int    `url:"ssh_key_id,omitempty"`
    Password      string `url:"password,omitempty"`
    CloudConfig   string `url:"cloud_config,omitempty"`
    ScriptContent string `url:"script_content,omitempty"`
}

```
BuildServerRequest is a set of parameters for a server re-building call.










## <a name="Client">type</a> [Client](/src/target/client.go?s=567:669#L29)
``` go
type Client struct {
    // contains filtered or unexported fields
}

```
Client is the main object (struct) to which we attach most
methods/functions.
It has the following fields:
(client, userAgent, endPoint, apiKey)







### <a name="NewClient">func</a> [NewClient](/src/target/client.go?s=1547:1584#L66)
``` go
func NewClient(apikey string) *Client
```
NewClient takes an apikey and calls NewClientCustom with the hardcoded
BaseEndpoint constant API URL


### <a name="NewClientCustom">func</a> [NewClientCustom](/src/target/client.go?s=1022:1080#L45)
``` go
func NewClientCustom(apikey string, apiurl string) *Client
```
NewClientCustom is the main entrypoint for instantiating a Client struct.
It takes your API Key as it's sole argument
and returns the Client struct ready to talk to the API





### <a name="Client.BuildServer">func</a> (\*Client) [BuildServer](/src/target/servers.go?s=3404:3490#L98)
``` go
func (c *Client) BuildServer(id int, r *BuildServerRequest) (b ServerBuild, err error)
```
BuildServer external method on Client to re-build an instance




### <a name="Client.CreateSSHKey">func</a> (\*Client) [CreateSSHKey](/src/target/sshkeys.go?s=713:787#L34)
``` go
func (c *Client) CreateSSHKey(name, key string) (sshkey SSHKey, err error)
```
CreateSSHKey creates a key




### <a name="Client.CreateServer">func</a> (\*Client) [CreateServer](/src/target/servers.go?s=2475:2555#L69)
``` go
func (c *Client) CreateServer(r *CreateServerRequest) (b ServerBuild, err error)
```
CreateServer external method on Client to buy and build a new instance.




### <a name="Client.DeleteSSHKey">func</a> (\*Client) [DeleteSSHKey](/src/target/sshkeys.go?s=1034:1077#L47)
``` go
func (c *Client) DeleteSSHKey(id int) error
```
DeleteSSHKey deletes a key




### <a name="Client.DeleteServer">func</a> (\*Client) [DeleteServer](/src/target/servers.go?s=3846:3909#L115)
``` go
func (c *Client) DeleteServer(id int, cancelBilling bool) error
```
DeleteServer external method on Client to destroy an instance.




### <a name="Client.GetLocations">func</a> (\*Client) [GetLocations](/src/target/locations.go?s=384:435#L14)
``` go
func (c *Client) GetLocations() ([]Location, error)
```
GetLocations public method on Client to get a list of locations




### <a name="Client.GetOSs">func</a> (\*Client) [GetOSs](/src/target/os.go?s=356:395#L15)
``` go
func (c *Client) GetOSs() ([]OS, error)
```
GetOSs returns a list of OS objects from the api




### <a name="Client.GetSSHKey">func</a> (\*Client) [GetSSHKey](/src/target/sshkeys.go?s=493:554#L26)
``` go
func (c *Client) GetSSHKey(id int) (sshkey SSHKey, err error)
```
GetSSHKey as in one key




### <a name="Client.GetSSHKeys">func</a> (\*Client) [GetSSHKeys](/src/target/sshkeys.go?s=270:326#L17)
``` go
func (c *Client) GetSSHKeys() (keys []SSHKey, err error)
```
GetSSHKeys as in many keys




### <a name="Client.GetServer">func</a> (\*Client) [GetServer](/src/target/servers.go?s=1227:1288#L39)
``` go
func (c *Client) GetServer(id int) (server Server, err error)
```
GetServer external method on Client to get an instance




### <a name="Client.GetServers">func</a> (\*Client) [GetServers](/src/target/servers.go?s=985:1032#L30)
``` go
func (c *Client) GetServers() ([]Server, error)
```
GetServers external method on Client to list your instances




### <a name="Client.UnlinkServer">func</a> (\*Client) [UnlinkServer](/src/target/servers.go?s=4178:4221#L124)
``` go
func (c *Client) UnlinkServer(id int) error
```
UnlinkServer external method on Client to unlink a billing package from a location




## <a name="CreateServerRequest">type</a> [CreateServerRequest](/src/target/servers.go?s=1495:2216#L47)
``` go
type CreateServerRequest struct {
    Plan                     string `url:"plan,omitempty"`
    Location                 int    `url:"location,omitempty"`
    Image                    int    `url:"image,omitempty"`
    FQDN                     string `url:"fqdn,omitempty"`
    SSHKey                   string `url:"ssh_key,omitempty"`
    SSHKeyID                 int    `url:"ssh_key_id,omitempty"`
    Password                 string `url:"password,omitempty"`
    PackageBilling           string `url:"package_billing,omitempty"`
    PackageBillingContractId string `url:"package_billing_contract_id,omitempty"`
    CloudConfig              string `url:"cloud_config,omitempty"`
    ScriptContent            string `url:"script_content,omitempty"`
}

```
CreateServerRequest is as set of parameters for a server creation call.










## <a name="Location">type</a> [Location](/src/target/locations.go?s=87:315#L4)
``` go
type Location struct {
    ID        int    `json:"id"`
    Name      string `json:"name"`
    IATACode  string `json:"iata_code"`
    Continent string `json:"continent"`
    Flag      string `json:"flat"`
    Disabled  int    `json:"disabled"`
}

```
Location is an API response message identifyin a particular location.










## <a name="OS">type</a> [OS](/src/target/os.go?s=68:302#L4)
``` go
type OS struct {
    ID      int    `json:"id,string"`
    Os      string `json:"os"`
    Type    string `json:"type"`
    Subtype string `json:"subtype"`
    Size    string `json:"size"`
    Bits    string `json:"bits"`
    Tech    string `json:"tech"`
}

```
OS is a struct for storing the attributes of an OS










## <a name="SSHKey">type</a> [SSHKey](/src/target/sshkeys.go?s=72:238#L9)
``` go
type SSHKey struct {
    ID          int    `json:"id"`
    Name        string `json:"name"`
    Key         string `json:"ssh_key"`
    Fingerprint string `json:"fingerprint"`
}

```
SSHKey is what it is










## <a name="Server">type</a> [Server](/src/target/servers.go?s=138:920#L11)
``` go
type Server struct {
    Name                     string `json:"fqdn"`
    ID                       int    `json:"mbpkgid"`
    OS                       string `json:"os"`
    OSID                     int    `json:"os_id"`
    PrimaryIPv4              string `json:"ip"`
    PrimaryIPv6              string `json:"ipv6"`
    PlanID                   int    `json:"plan_id"`
    Package                  string `json:"package"`
    PackageBilling           string `json:"package_billing"`
    PackageBillingContractId string `json:"package_billing_contract_id"`
    Location                 string `json:"city"`
    LocationID               int    `json:"location_id"`
    ServerStatus             string `json:"status"`
    PowerStatus              string `json:"state"`
    Installed                int    `json:"installed"`
}

```
Server struct defines what a VPS looks like










## <a name="ServerBuild">type</a> [ServerBuild](/src/target/servers.go?s=2272:2398#L62)
``` go
type ServerBuild struct {
    ServerID int    `json:"mbpkgid"`
    Status   string `json:"status"`
    Build    int    `json:"build"`
}

```
ServerBuild is a server creation response message.














- - -
Generated by [godoc2md](http://godoc.org/github.com/davecheney/godoc2md)
