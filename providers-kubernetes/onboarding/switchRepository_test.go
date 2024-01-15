package providers_test

import (
	"testing"

	providers "github.com/ironcore-dev/metal/providers-kubernetes/onboarding"
	"github.com/ironcore-dev/metal/providers-kubernetes/onboarding/fake"
	"github.com/stretchr/testify/assert"
)

func TestSwitchRepositoryByChassisIDSuccess(t *testing.T) {
	t.Parallel()

	a := assert.New(t)
	uuid, namespace := "123", "default"
	chassisID := "123"
	portDescription := "test"
	ipWithPrefix := "192.168.1.1/30"
	sw := fake.FakeSwitchObject(
		uuid,
		namespace,
		chassisID,
		portDescription,
		ipWithPrefix,
	)
	fakeClient, _ := fake.NewFakeWithObjects(sw)

	repo := providers.NewSwitchRepository(fakeClient)

	switchInfo, err := repo.ByChassisID(uuid)
	a.Nil(err)
	a.NotEmpty(switchInfo)
	for _, ip := range switchInfo.InterfacesInfo[portDescription].IP {
		a.Equal(ipWithPrefix, ip.String())
	}
}
