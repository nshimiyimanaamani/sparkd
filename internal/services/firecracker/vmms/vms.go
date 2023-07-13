package vmms

import (
	"fmt"
	"net"
	"os"

	firecracker "github.com/firecracker-microvm/firecracker-go-sdk"
	"github.com/firecracker-microvm/firecracker-go-sdk/client/models"
	"github.com/quarksgroup/sparkd/internal/core"
	log "github.com/sirupsen/logrus"
)

type Config struct {
	*config
	kernel, fcBin, dir, initrd, logLevel string
}

// New initiate options instance implementation
func New(dir, kernel, fcBinary, initrd, level string) *Config {
	return &Config{
		kernel:   dir + kernel,
		dir:      dir,
		fcBin:    fcBinary,
		initrd:   dir + initrd,
		logLevel: level,
	}
}

func (o *Config) generateOpt(index byte, image, id, name string) (*Config, error) {

	fc_ip := net.IPv4(174, 138, 44, 160+index).String()
	// gateway_ip := "174.138.44.163"
	// mask_long := "255.255.0.0"
	bootArgs := "ro console=ttyS0 noapic reboot=k panic=1 earlycon pci=off init=init nomodules random.trust_cpu=on tsc=reliable quiet rw "
	// bootArgs = bootArgs + fmt.Sprintf("ip=%s::%s:%s::eth0:off", fc_ip, gateway_ip, mask_long)

	out := &Config{
		config: &config{
			id:             id,
			name:           name,
			vmIndex:        int64(index),
			fcBinary:       o.fcBin,
			fcKernelImage:  o.kernel, // make sure that this file exists in the current directory with valid sum5
			kernelBootArgs: bootArgs,
			providedImage:  image,
			tapMacAddr:     fmt.Sprintf("02:FC:00:00:00:%02x", index),
			tap:            fmt.Sprintf("fc-tap-%d", index),
			fcIP:           fc_ip,
			initdPath:      o.initrd,
			backBone:       "eth0", // eth0 or enp7s0,enp0s25
			// ApiSocket:      fmt.Sprintf("/tmp/firecracker-%d.sock", id),
			fcCPUCount: 1,
			fcMemSz:    256,
			logger:     log.New(),
		},
	}

	roots, err := out.generateRFs(name)
	if err != nil {
		return nil, fmt.Errorf("failed to generate rootfs image, %s", err)
	}
	out.rootFsImage = roots

	// //create log file
	// _, err = cmd.RunNoneSudo(fmt.Sprintf("touch %d-%s.log", out.VmIndex, name))
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to create log file, %s", err)
	// }

	// out.LogFile = fmt.Sprintf("%s%d-%s.log", parent_dir, out.VmIndex, name)

	// if err := syscall.Mkfifo(out.LogFile, 0700); err != nil {
	// 	return nil, fmt.Errorf("failed to create fifo file, %s", err)
	// }

	return out, nil
}

func getFcConfig(opts *Config) firecracker.Config {

	return firecracker.Config{
		VMID: opts.id,
		// SocketPath:      opts.ApiSocket,
		KernelImagePath: opts.fcKernelImage,
		KernelArgs:      opts.kernelBootArgs,
		LogLevel:        opts.logLevel,
		InitrdPath:      opts.initdPath,
		Drives: []models.Drive{
			{
				DriveID:      firecracker.String("1"),
				PathOnHost:   &opts.rootFsImage,
				IsRootDevice: firecracker.Bool(true),
				IsReadOnly:   firecracker.Bool(false),
			},
		},

		//for setting up networking tap config vmmd config
		NetworkInterfaces: []firecracker.NetworkInterface{
			{
				StaticConfiguration: &firecracker.StaticNetworkConfiguration{
					MacAddress:  opts.tapMacAddr,
					HostDevName: opts.tap,
					IPConfiguration: &firecracker.IPConfiguration{
						IPAddr: net.IPNet{
							IP:   net.ParseIP(opts.fcIP),
							Mask: net.CIDRMask(16, 32),
						},
						Gateway: net.ParseIP("174.138.44.163"),
					},
				},
				AllowMMDS: true,
			},
		},

		ForwardSignals: make([]os.Signal, 0),

		//for specifying the number of cpus and memory
		MachineCfg: models.MachineConfiguration{
			VcpuCount:  firecracker.Int64(1),
			Smt:        firecracker.Bool(false),
			MemSizeMib: firecracker.Int64(1024),
		},

		// Enable seccomp as recommended by firecracker-doc
		Seccomp: firecracker.SeccompConfig{
			Enabled: true,
		},

		// Specify the jailer configuration options
		JailerCfg: &firecracker.JailerConfig{
			ID:             opts.id,
			UID:            firecracker.Int(int(opts.vmIndex)), // Make that uid and gid are same and unique for each vm in order to provide an extra layer of security for their individually owned
			GID:            firecracker.Int(int(opts.vmIndex) + 1),
			NumaNode:       firecracker.Int(0),
			Daemonize:      true,
			ExecFile:       "/usr/bin/" + opts.fcBinary,
			JailerBinary:   "jailer",
			ChrootBaseDir:  "/tmp",
			CgroupVersion:  "1",
			Stdout:         opts.logger.WithField("vmm_stream", "stdout").WriterLevel(log.DebugLevel),
			Stderr:         opts.logger.WithField("vmm_stream", "stderr").WriterLevel(log.DebugLevel),
			Stdin:          os.Stdin,
			ChrootStrategy: firecracker.NewNaiveChrootStrategy(opts.fcKernelImage),
		},
		// LogPath: opts.LogFile,
		//VsockDevices:      vsocks,
		//MetricsFifo:       opts.FcMetricsFifo,
		//FifoLogWriter:     fifo,
	}
}

var _ core.MachineService = (*Config)(nil)
