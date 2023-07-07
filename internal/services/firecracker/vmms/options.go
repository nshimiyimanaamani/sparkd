package vmms

import (
	"fmt"
	"net"
	"os"

	firecracker "github.com/firecracker-microvm/firecracker-go-sdk"
	"github.com/firecracker-microvm/firecracker-go-sdk/client/models"
	"github.com/quarksgroup/sparkd/internal/core"
	"github.com/quarksgroup/sparkd/internal/rand"
	log "github.com/sirupsen/logrus"
)

type Options core.Config

var parent_dir = "/sparkd/"

func (o *Options) GenerateOpt(id byte, image, name string) (*Options, error) {

	fc_ip := net.IPv4(174, 138, 44, 160+id).String()
	// gateway_ip := "174.138.44.163"
	// mask_long := "255.255.0.0"
	bootArgs := "ro console=ttyS0 noapic reboot=k panic=1 earlycon pci=off init=init nomodules random.trust_cpu=on tsc=reliable quiet "
	// bootArgs = bootArgs + fmt.Sprintf("ip=%s::%s:%s::eth0:off", fc_ip, gateway_ip, mask_long)

	out := &Options{
		Id:             rand.UUID(),
		Name:           name,
		VmIndex:        int64(id),
		FcBinary:       "firecracker",
		FcKernelImage:  parent_dir + "vmlinux.bin", // make sure that this file exists in the current directory with valid sum5
		KernelBootArgs: bootArgs,
		ProvidedImage:  image,
		TapMacAddr:     fmt.Sprintf("02:FC:00:00:00:%02x", id),
		Tap:            fmt.Sprintf("fc-tap-%d", id),
		FcIP:           fc_ip,
		BackBone:       "eth0", // eth0 or enp7s0,enp0s25
		// ApiSocket:      fmt.Sprintf("/tmp/firecracker-%d.sock", id),
		FcCPUCount: 1,
		FcMemSz:    256,
		Logger:     log.New(),
	}

	roots, err := out.generateRFs(name)
	if err != nil {
		return nil, fmt.Errorf("failed to generate rootfs image, %s", err)
	}
	out.RootFsImage = roots

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

func (opts *Options) getFcConfig() firecracker.Config {

	return firecracker.Config{
		VMID: opts.Id,
		// SocketPath:      opts.ApiSocket,
		KernelImagePath: opts.FcKernelImage,
		KernelArgs:      opts.KernelBootArgs,
		LogLevel:        "debug",
		InitrdPath:      parent_dir + "initrd.cpio",
		Drives: []models.Drive{
			{
				DriveID:      firecracker.String("1"),
				PathOnHost:   &opts.RootFsImage,
				IsRootDevice: firecracker.Bool(true),
				IsReadOnly:   firecracker.Bool(false),
			},
		},

		//for setting up networking tap config vmmd config
		NetworkInterfaces: []firecracker.NetworkInterface{
			{
				StaticConfiguration: &firecracker.StaticNetworkConfiguration{
					MacAddress:  opts.TapMacAddr,
					HostDevName: opts.Tap,
					IPConfiguration: &firecracker.IPConfiguration{
						IPAddr: net.IPNet{
							IP:   net.ParseIP(opts.FcIP),
							Mask: net.CIDRMask(16, 32),
						},
						Gateway: net.ParseIP("174.138.44.163"),
					},
				},
				AllowMMDS: true,
			},
		},

		// ForwardSignals: []os.Signal{
		// 	os.Interrupt,
		// },

		//for specifying the number of cpus and memory
		MachineCfg: models.MachineConfiguration{
			VcpuCount:  firecracker.Int64(1),
			Smt:        firecracker.Bool(false),
			MemSizeMib: firecracker.Int64(526),
		},

		// Enable seccomp as recommended by firecracker-doc
		Seccomp: firecracker.SeccompConfig{
			Enabled: true,
		},

		// Specify the jailer configuration options
		JailerCfg: &firecracker.JailerConfig{
			ID:             opts.Id,
			UID:            firecracker.Int(int(opts.VmIndex)), // Make that uid and gid are same and unique for each vm in order to provide an extra layer of security for their individually owned
			GID:            firecracker.Int(int(opts.VmIndex) + 1),
			NumaNode:       firecracker.Int(0),
			Daemonize:      true,
			ExecFile:       "/usr/bin/" + opts.FcBinary,
			JailerBinary:   "jailer",
			ChrootBaseDir:  "/tmp",
			CgroupVersion:  "1",
			Stdout:         opts.Logger.WithField("vmm_stream", "stdout").WriterLevel(log.DebugLevel),
			Stderr:         opts.Logger.WithField("vmm_stream", "stderr").WriterLevel(log.DebugLevel),
			Stdin:          os.Stdin,
			ChrootStrategy: firecracker.NewNaiveChrootStrategy(parent_dir + "vmlinux.bin"),
		},
		// LogPath: opts.LogFile,
		//VsockDevices:      vsocks,
		//MetricsFifo:       opts.FcMetricsFifo,
		//FifoLogWriter:     fifo,
	}
}
