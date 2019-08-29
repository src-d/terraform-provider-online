package online

type Server struct {
	ID                  int         `json:"id"`
	Offer               string      `json:"offer"`
	Hostname            string      `json:"hostname"`
	OS                  interface{} `json:"os"`
	Power               string      `json:"power"`
	BootMode            string      `json:"boot_mode"`
	InstallStatus       string      `json:"install_status"`
	LastReboot          string      `json:"last_reboot"`
	AntiDDOS            bool        `json:"anti_ddos"`
	HardwareWatch       bool        `json:"hardware_watch"`
	ProactiveMonitoring bool        `json:"proactive_monitoring"`
	Support             string      `json:"support"`
	Abuse               string      `json:"abuse"`
	Location            *Location   `json:"location"`
	Network             struct {
		Public  []string      `json:"ip"`
		Private []string      `json:"private"`
		Ipfo    []interface{} `json:"ipfo"`
	} `json:"network"`
	IP       []*Interface `json:"ip"`
	Contacts struct {
		Owner string `json:"owner"`
		Tech  string `json:"tech"`
	} `json:"contacts"`
	Disks []struct {
		Ref string `json:"$ref"`
	} `json:"disks"`
	DriveArrays []struct {
		Disks []struct {
			Ref string `json:"$ref"`
		} `json:"disks"`
		RaidController struct {
			Ref string `json:"$ref"`
		} `json:"raid_controller"`
		RaidLevel string `json:"raid_level"`
	} `json:"drive_arrays"`
	RaidControllers []struct {
		Ref string `json:"$ref"`
	} `json:"raid_controllers"`
	BMC struct {
		SessionKey interface{} `json:"session_key"`
	} `json:"bmc"`
}

type ServerInstall struct {
	Hostname                string   `json:"hostname"`
	OS_ID                   string   `json:"os_id"`
	UserLogin               string   `json:"user_login"`
	UserPassword            string   `json:"user_password"`
	RootPassword            string   `json:"root_password"`
	PartitioningTemplateRef string   `json:"partitioning_template_ref"`
	SSHKeys                 []string `json:"ssh_keys"`
}

type OS struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type OperatingSystem struct {
	OS
	ID        int    `json:"id"`
	Type      string `json:"type"`
	Arch      string `json:"arch"`
	Release   string `json:"release"`
	EndOfLife string `json:"end_of_life"`
}

type OperatingSystems []OperatingSystem

type InterfaceType string

const (
	Public  InterfaceType = "public"
	Private InterfaceType = "private"
)

type Interface struct {
	Address         string        `json:"address"`
	MAC             string        `json:"mac"`
	Reverse         string        `json:"reverse"`
	SwitchPortState string        `json:"switch_port_state"`
	Type            InterfaceType `json:"type"`
}

type Location struct {
	Block      string `json:"block"`
	Datacenter string `json:"datacenter"`
	Position   int    `json:"position"`
	Rack       string `json:"rack"`
	Room       string `json:"room"`
}

func (s *Server) InterfaceByType(t InterfaceType) *Interface {
	for _, i := range s.IP {
		if i.Type == t {
			return i
		}
	}

	return nil
}

type SSHKey struct {
	UUID        string `json:"uuid_ref"`
	Description string `json:"description"`
	Fingerprint string `json:"fingerprint"`
}

type SSHKeys []SSHKey
