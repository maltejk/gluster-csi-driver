package heketi

import (
	"os"

	"github.com/gluster/gluster-csi-driver/pkg/command"
	"github.com/gluster/gluster-csi-driver/pkg/controller"
	"github.com/gluster/gluster-csi-driver/pkg/identity"
	"github.com/gluster/gluster-csi-driver/pkg/node"

	"github.com/golang/glog"
	csi "github.com/kubernetes-csi/drivers/pkg/csi-common"
	"k8s.io/kubernetes/pkg/util/mount"
)

// Driver implements command.Driver
type Driver struct {
	*command.Config
	Mounter mount.Interface
}

// New returns a new Driver
func New(config *command.Config, mounter mount.Interface) *Driver {
	hd := &Driver{}

	if config != nil {
		hd.Config = config
	} else {
		glog.Error("failed to initialize Heketi driver: config is nil")
		return nil
	}

	if mounter == nil {
		mounter = mount.New("")
	}
	hd.Mounter = mounter

	glog.V(1).Infof("%s initialized", hd.Desc)

	return hd
}

// Run runs the driver
func (d *Driver) Run() {
	client, err := NewClient(d.Config)
	if err != nil {
		glog.Errorf("failed to get heketiClient: %v", err)
		os.Exit(1)
	}

	is := identity.NewServer(d.Config)
	cs := controller.NewServer(d.Config, client)
	ns := node.NewServer(d.Config, d.Mounter)

	srv := csi.NewNonBlockingGRPCServer()
	srv.Start(d.Config.Endpoint, is, cs, ns)
	srv.Wait()
}
