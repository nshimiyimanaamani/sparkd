package vmms

import (
	llg "github.com/sirupsen/logrus"
)

type config struct {
	id             string `long:"id" description:"Jailer VMM id"`
	name           string `long:"name" description:"VM name this is provided by the user"`
	vmIndex        int64  `long:"vm-index" description:"VM index"`
	apiSocket      string `long:"socket-path" short:"s" description:"path to use for firecracker socket"`
	fcBinary       string `long:"firecracker-binary" description:"Path to firecracker binary"`
	fcKernelImage  string `long:"kernel" description:"Path to the kernel image"`
	kernelBootArgs string `long:"kernel-opts" description:"Kernel commandline"`
	rootFsImage    string `long:"root-drive" description:"Path to root disk image"`
	tapMacAddr     string `long:"tap-mac-addr" description:"tap macaddress"`
	tapGateWay     string
	tapMask        string
	tap            string `long:"tap-dev" description:"tap device"`
	fcCPUCount     int64  `long:"ncpus" short:"c" description:"Number of CPUs"`
	fcMemSz        int64  `long:"memory" short:"m" description:"VM memory, in MiB"`
	fcIP           string `long:"fc-ip" description:"IP address of the VM"`

	backBone      string `long:"if-name" description:"if name to match your main ethernet adapter,the one that accesses the Internet - check 'ip addr' or 'ifconfig' if you don't know which one to use"` // eg eth0
	providedImage string `long:"provided-image" description:"provided-image is the image that we want to run in the VM"`
	initdPath     string `long:"initd-path" description:"initd-path is the path to the init binary file"`
	logger        *llg.Logger
	logFile       string `long:"log-file" description:"log-file is the path to the log file"`
	runVmFile     string `long:"run-vm-file" description:"run-vm-file is the path to the run.json file"`
}
