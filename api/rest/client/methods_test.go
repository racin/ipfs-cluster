package client

import (
	"context"
	"log"
	"testing"
	"time"

	cid "github.com/ipfs/go-cid"
	ma "github.com/multiformats/go-multiaddr"

	"github.com/ipfs/ipfs-cluster/api"
	"github.com/ipfs/ipfs-cluster/test"
)

func TestVersion(t *testing.T) {
	c, api := testClient(t)
	defer api.Shutdown()
	v, err := c.Version()
	if err != nil || v.Version == "" {
		t.Logf("%+v", v)
		t.Log(err)
		t.Error("expected something in version")
	}
}

func TestID(t *testing.T) {
	c, api := testClient(t)
	defer api.Shutdown()
	id, err := c.ID()
	if err != nil {
		t.Fatal(err)
	}
	if id.ID == "" {
		t.Error("bad id")
	}
}

func TestPeers(t *testing.T) {
	c, api := testClient(t)
	defer api.Shutdown()
	ids, err := c.Peers()
	if err != nil {
		t.Fatal(err)
	}
	if len(ids) == 0 {
		t.Error("expected some peers")
	}
}

func TestPeersWithError(t *testing.T) {
	c, api := testClient(t)
	defer api.Shutdown()
	addr, _ := ma.NewMultiaddr("/ip4/127.0.0.1/tcp/44444")
	c, _ = NewClient(&Config{APIAddr: addr, DisableKeepAlives: true})
	ids, err := c.Peers()
	if err == nil {
		t.Fatal("expected error")
	}
	if ids == nil || len(ids) != 0 {
		t.Fatal("expected no ids")
	}
}

func TestPeerAdd(t *testing.T) {
	c, api := testClient(t)
	defer api.Shutdown()

	addr, _ := ma.NewMultiaddr("/ip4/1.2.3.4/tcp/1234/ipfs/" + test.TestPeerID1.Pretty())
	id, err := c.PeerAdd(addr)
	if err != nil {
		t.Fatal(err)
	}
	if id.ID != test.TestPeerID1 {
		t.Error("bad peer")
	}
}

func TestPeerRm(t *testing.T) {
	c, api := testClient(t)
	defer api.Shutdown()

	err := c.PeerRm(test.TestPeerID1)
	if err != nil {
		t.Fatal(err)
	}
}

func TestPin(t *testing.T) {
	c, api := testClient(t)
	defer api.Shutdown()

	ci, _ := cid.Decode(test.TestCid1)
	err := c.Pin(ci, 6, 7, "hello")
	if err != nil {
		t.Fatal(err)
	}
}

func TestUnpin(t *testing.T) {
	c, api := testClient(t)
	defer api.Shutdown()

	ci, _ := cid.Decode(test.TestCid1)
	err := c.Unpin(ci)
	if err != nil {
		t.Fatal(err)
	}
}

func TestAllocations(t *testing.T) {
	c, api := testClient(t)
	defer api.Shutdown()

	pins, err := c.Allocations()
	if err != nil {
		t.Fatal(err)
	}
	if len(pins) == 0 {
		t.Error("should be some pins")
	}
}

func TestAllocation(t *testing.T) {
	c, api := testClient(t)
	defer api.Shutdown()

	ci, _ := cid.Decode(test.TestCid1)
	pin, err := c.Allocation(ci)
	if err != nil {
		t.Fatal(err)
	}
	if pin.Cid.String() != test.TestCid1 {
		t.Error("should be same pin")
	}
}

func TestStatus(t *testing.T) {
	c, api := testClient(t)
	defer api.Shutdown()

	ci, _ := cid.Decode(test.TestCid1)
	pin, err := c.Status(ci, false)
	if err != nil {
		t.Fatal(err)
	}
	if pin.Cid.String() != test.TestCid1 {
		t.Error("should be same pin")
	}
}

func TestStatusAll(t *testing.T) {
	c, api := testClient(t)
	defer api.Shutdown()

	pins, err := c.StatusAll(false)
	if err != nil {
		t.Fatal(err)
	}

	if len(pins) == 0 {
		t.Error("there should be some pins")
	}
}

func TestSync(t *testing.T) {
	c, api := testClient(t)
	defer api.Shutdown()

	ci, _ := cid.Decode(test.TestCid1)
	pin, err := c.Sync(ci, false)
	if err != nil {
		t.Fatal(err)
	}
	if pin.Cid.String() != test.TestCid1 {
		t.Error("should be same pin")
	}
}

func TestSyncAll(t *testing.T) {
	c, api := testClient(t)
	defer api.Shutdown()

	pins, err := c.SyncAll(false)
	if err != nil {
		t.Fatal(err)
	}

	if len(pins) == 0 {
		t.Error("there should be some pins")
	}
}

func TestRecover(t *testing.T) {
	c, api := testClient(t)
	defer api.Shutdown()

	ci, _ := cid.Decode(test.TestCid1)
	pin, err := c.Recover(ci, false)
	if err != nil {
		t.Fatal(err)
	}
	if pin.Cid.String() != test.TestCid1 {
		t.Error("should be same pin")
	}
}

func TestRecoverAll(t *testing.T) {
	c, api := testClient(t)
	defer api.Shutdown()

	_, err := c.RecoverAll(true)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetConnectGraph(t *testing.T) {
	c, api := testClient(t)
	defer api.Shutdown()

	cg, err := c.GetConnectGraph()
	if err != nil {
		t.Fatal(err)
	}
	if len(cg.IPFSLinks) != 3 || len(cg.ClusterLinks) != 3 ||
		len(cg.ClustertoIPFS) != 3 {
		t.Fatal("Bad graph")
	}
}

func TestClient_WaitFor(t *testing.T) {
	ci, _ := cid.Decode(test.TestCid1)
	type args struct {
		ci        *cid.Cid
		local     bool
		target    api.TrackerStatus
		checkFreq time.Duration
	}
	tests := []struct {
		name string
		args args
		want api.TrackerStatus
	}{
		{
			"unpin",
			args{
				ci,
				false,
				api.TrackerStatusUnpinned,
				2 * time.Second,
			},
			api.TrackerStatusUnpinned,
		},
		{
			"pin",
			args{
				ci,
				false,
				api.TrackerStatusPinned,
				2 * time.Second,
			},
			api.TrackerStatusPinned,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, api := testClient(t)
			defer api.Shutdown()

			var err error
			switch tt.name {
			case "unpin":
				err = c.Pin(ci, 1, 1, "test-pin")
			case "pin":
				err = c.Unpin(ci)
			}

			t.Errorf("failed to perform setup %s: %v", tt.name, err)

			ctx := context.Background()
			statusCh, errCh := c.WaitFor(ctx, tt.args.ci, tt.args.local, tt.args.target, tt.args.checkFreq)
			for {
				select {
				case stat := <-statusCh:
					log.Printf("%#v\n", stat)
					return
				case err := <-errCh:
					log.Printf("%#v\n", err)
					return
				}
			}
		})
	}
}
