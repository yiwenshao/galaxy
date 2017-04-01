package main

import (
	"flag"
	"net"

	"git.code.oa.com/gaiastack/galaxy/pkg/network/vlan"
	"git.code.oa.com/gaiastack/galaxy/pkg/utils"
	"github.com/containernetworking/cni/pkg/skel"
	"github.com/containernetworking/cni/pkg/types"
	"github.com/golang/glog"
)

var (
	flagDevice  = flag.String("device", "", "The device which has the ip address, eg. eth1 or eth1.12 (A vlan device)")
	flagNetns   = flag.String("netns", "", "The netns path for the container")
	flagIP      = flag.String("ip", "", "The ip in cidr format for the container")
	flagVlan    = flag.Uint("vlan", 0, "The vlan id of the ip")
	flagGateway = flag.String("gateway", "", "The gateway for the ip")
)

/*
 ./setupvlan -logtostderr -device bond1
 ip netns add ctn2; ./setupvlan -logtostderr -device bond1 -netns=/var/run/netns/ctn2 -ip=10.2.1.111/24 -gateway=10.2.1.1
*/
func main() {
	flag.Parse()
	d := &vlan.VlanDriver{}
	if *flagDevice == "" {
		glog.Fatal("device unset")
	}
	d.NetConf = &vlan.NetConf{Device: *flagDevice}
	if err := d.SetupBridge(); err != nil {
		glog.Fatalf("Error setting up bridge %v", err)
	}
	glog.Infof("setuped bridge docker")
	if *flagNetns == "" {
		return
	}
	ip, ipNet, err := net.ParseCIDR(*flagIP)
	if err != nil {
		glog.Fatalf("invalid cidr %s", *flagIP)
	}
	ipNet.IP = ip
	gateway := net.ParseIP(*flagGateway)
	if gateway == nil {
		glog.Fatalf("invalid gateway %s", *flagGateway)
	}
	if *flagVlan != 0 {
		if err := d.CreateVlanDevice(uint16(*flagVlan)); err != nil {
			glog.Fatalf("Error creating vlan device %v", err)
		}
	}
	if err := utils.ConnectsHostWithContainer(&types.Result{
		IP4: &types.IPConfig{
			IP:      *ipNet,
			Gateway: gateway,
			Routes: []types.Route{{
				Dst: net.IPNet{
					IP:   net.IPv4(0, 0, 0, 0),
					Mask: net.IPv4Mask(0, 0, 0, 0),
				},
			}},
		},
	}, &skel.CmdArgs{Netns: *flagNetns, IfName: "eth0"}, d.BridgeNameForVlan(uint16(*flagVlan))); err != nil {
		glog.Fatalf("Error creating veth %v", err)
	}
	glog.Infof("privisioned container %s", *flagNetns)
}