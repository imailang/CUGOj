package testcaller

type BundleConfig struct {
	OciVersion string   `json:"ociVersion"`
	Process    Process  `json:"process"`
	Root       Root     `json:"root"`
	Hostname   string   `json:"hostname"`
	Mounts     []Mounts `json:"mounts"`
	Linux      Linux    `json:"linux"`
}
type User struct {
	UID int `json:"uid"`
	Gid int `json:"gid"`
}
type Capabilities struct {
	Bounding    []string `json:"bounding"`
	Effective   []string `json:"effective"`
	Inheritable []string `json:"inheritable"`
	Permitted   []string `json:"permitted"`
	Ambient     []string `json:"ambient"`
}
type Rlimits struct {
	Type string `json:"type"`
	Hard int    `json:"hard"`
	Soft int    `json:"soft"`
}
type Process struct {
	Terminal        bool         `json:"terminal"`
	User            User         `json:"user"`
	Args            []string     `json:"args"`
	Env             []string     `json:"env"`
	Cwd             string       `json:"cwd"`
	Capabilities    Capabilities `json:"capabilities"`
	Rlimits         []Rlimits    `json:"rlimits"`
	NoNewPrivileges bool         `json:"noNewPrivileges"`
}
type Root struct {
	Path     string `json:"path"`
	Readonly bool   `json:"readonly"`
}
type Mounts struct {
	Destination string   `json:"destination"`
	Type        string   `json:"type"`
	Source      string   `json:"source"`
	Options     []string `json:"options,omitempty"`
}
type Devices struct {
	Allow  bool   `json:"allow"`
	Access string `json:"access"`
}
type Resources struct {
	Devices []Devices `json:"devices"`
}
type Namespaces struct {
	Type string `json:"type"`
}
type Linux struct {
	Resources     Resources    `json:"resources"`
	Namespaces    []Namespaces `json:"namespaces"`
	MaskedPaths   []string     `json:"maskedPaths"`
	ReadonlyPaths []string     `json:"readonlyPaths"`
}

func (config *BundleConfig) SetTimeLimit(time_limit int) {

	for i := 0; i < len(config.Process.Rlimits); i++ {
		if config.Process.Rlimits[i].Type == "RLIMIT_CPU" {
			config.Process.Rlimits[i].Hard = time_limit
			config.Process.Rlimits[i].Soft = time_limit
			return
		}
	}
	config.Process.Rlimits = append(config.Process.Rlimits, Rlimits{
		Type: "RLIMIT_CPU",
		Hard: time_limit,
		Soft: time_limit,
	})
}

func (config *BundleConfig) SetMemoryLimit(memory_limit int) {
	for i := 0; i < len(config.Process.Rlimits); i++ {
		if config.Process.Rlimits[i].Type == "RLIMIT_RSS" {
			config.Process.Rlimits[i].Hard = memory_limit
			config.Process.Rlimits[i].Soft = memory_limit
			return
		}
	}
	config.Process.Rlimits = append(config.Process.Rlimits, Rlimits{
		Type: "RLIMIT_RSS",
		Hard: memory_limit,
		Soft: memory_limit,
	})
}

func (config *BundleConfig) SetStackLimit(stack_limit int) {
	for i := 0; i < len(config.Process.Rlimits); i++ {
		if config.Process.Rlimits[i].Type == "STACK" {
			config.Process.Rlimits[i].Hard = stack_limit
			config.Process.Rlimits[i].Soft = stack_limit
			return
		}
	}
	config.Process.Rlimits = append(config.Process.Rlimits, Rlimits{
		Type: "RLIMIT_STACK",
		Hard: stack_limit,
		Soft: stack_limit,
	})
}

func (config *BundleConfig) SetMount(src, des string, readonly bool) {
	mount := Mounts{
		Source:      src,
		Destination: des,
		Type:        "bind",
	}
	mount.Options = append(mount.Options, "rbind")
	if readonly {
		mount.Options = append(mount.Options, "ro")
	} else {
		mount.Options = append(mount.Options, "rw")
	}
	config.Mounts = append(config.Mounts, mount)
}
