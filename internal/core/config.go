package core

import (
	"sync"

	llg "github.com/sirupsen/logrus"
)

type Config struct {
	Id             string `long:"id" description:"Jailer VMM id"`
	VmIndex        int64  `long:"vm-index" description:"VM index"`
	ApiSocket      string `long:"socket-path" short:"s" description:"path to use for firecracker socket"`
	IpId           byte   `byte:"id" description:"an ip we use to generate an ip address"`
	FcBinary       string `long:"firecracker-binary" description:"Path to firecracker binary"`
	FcKernelImage  string `long:"kernel" description:"Path to the kernel image"`
	KernelBootArgs string `long:"kernel-opts" description:"Kernel commandline"`
	RootFsImage    string `long:"root-drive" description:"Path to root disk image"`
	TapMacAddr     string `long:"tap-mac-addr" description:"tap macaddress"`
	Tap            string `long:"tap-dev" description:"tap device"`
	FcCPUCount     int64  `long:"ncpus" short:"c" description:"Number of CPUs"`
	FcMemSz        int64  `long:"memory" short:"m" description:"VM memory, in MiB"`
	FcIP           string `long:"fc-ip" description:"IP address of the VM"`

	BackBone      string `long:"if-name" description:"if name to match your main ethernet adapter,the one that accesses the Internet - check 'ip addr' or 'ifconfig' if you don't know which one to use"` // eg eth0
	InitBaseTar   string `long:"init-base-tar" description:"init-base-tar is our init base image file"`                                                                                                   // make sure that this file is currently exists in the current directory by running task extract-init-base-tar
	ProvidedImage string `long:"provided-image" description:"provided-image is the image that we want to run in the VM"`
	InitdPath     string `long:"initd-path" description:"initd-path is the path to the init binary file"`
	Logger        *llg.Logger
}

// JailerConfig represents Jailerspecific configuration options.
type JailerConfig struct {
	sync.Mutex

	BinaryFirecracker string `json:"BinaryFirecracker" mapstructure:"BinaryFirecracker"`
	BinaryJailer      string `json:"BinaryJailer" mapstructure:"BinaryJailer"`
	ChrootBase        string `json:"ChrootBase" mapstructure:"ChrootBase"`

	JailerGID      int `json:"JailerGid" mapstructure:"JailerGid"`
	JailerNumeNode int `json:"JailerNumaNode" mapstructure:"JailerNumaNode"`
	JailerUID      int `json:"JailerUid" mapstructure:"JailerUid"`

	NetNS string `json:"NetNS" mapstructure:"NetNS"`

	VmmID string
}
