package vmms

import (
	"fmt"

	"github.com/quarksgroup/sparkd/internal/cmd"
)

func (o *Config) setNetwork() error {

	// delete tap device if it exists
	if res, err := cmd.RunNoneSudo(fmt.Sprintf("ip link del %s 2> /dev/null || true", o.tap)); res != 1 && err != nil {
		return fmt.Errorf("failed during deleting tap device: %v", err)
	}

	// create tap device
	if _, err := cmd.RunSudo(fmt.Sprintf("ip tuntap add dev %s mode tap > /dev/net/tun", o.tap)); err != nil {
		return fmt.Errorf("failed creating ip link for tap: %s", err)
	}

	if _, err := cmd.RunSudo(fmt.Sprintf("sysctl -w net.ipv4.conf.%s.proxy_arp=1", o.tap)); err != nil {
		return fmt.Errorf("failed doing first sysctl: %v", err)
	}

	if _, err := cmd.RunSudo(fmt.Sprintf("sysctl -w net.ipv6.conf.%s.disable_ipv6=1", o.tap)); err != nil {
		return fmt.Errorf("failed doing second sysctl: %v", err)
	}

	// set tap device mac address
	if _, err := cmd.RunSudo(fmt.Sprintf("ip addr add %s/16 dev %s > /dev/net/tun", o.fcIP, o.tap)); err != nil {
		return fmt.Errorf("failed to add ip address on tap device: %v", err)
	}

	// set tap device up by activating it
	if _, err := cmd.RunSudo(fmt.Sprintf("ip link set dev %s up", o.tap)); err != nil {
		return fmt.Errorf("failed to set tap device up: %v", err)
	}

	// bind to the interface associated with the address <host>
	if _, err := cmd.RunSudo(fmt.Sprintf("iperf3 -B %s -s > /dev/null 2>&1 &", o.fcIP)); err != nil {
		return fmt.Errorf("failed to bind to the interface associated with the address: %v", err)
	}

	//enable ip forwarding
	if _, err := cmd.RunSudo(" sh -c 'echo 1 > /proc/sys/net/ipv4/ip_forward'"); err != nil {
		return fmt.Errorf("failed to enable ip forwarding: %v", err)
	}

	// add iptables rule to forward packets from tap to eth0
	if _, err := cmd.RunSudo(fmt.Sprintf("iptables -t nat -A POSTROUTING -o %s -j MASQUERADE", o.backBone)); err != nil {
		return fmt.Errorf("failed to add iptables rule to forward packets from tap to eth0: %v", err)
	}

	// add iptables rule to establish connection between tap and eth0 (forward packets from eth0 to tap)
	if _, err := cmd.RunSudo("iptables -A FORWARD -m conntrack --ctstate RELATED,ESTABLISHED -j ACCEPT"); err != nil {
		return fmt.Errorf("failed to add iptables rule to establish connection between tap and eth0: %v", err)
	}

	// add iptables rule to forward packets from eth0 to tap
	if _, err := cmd.RunSudo(fmt.Sprintf("iptables -A FORWARD -i %s -o %s -j ACCEPT", o.tap, o.backBone)); err != nil {
		return fmt.Errorf("failed to add iptables rule to forward packets from eth0 to tap: %v", err)
	}

	return nil

}
