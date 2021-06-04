package v1alpha1

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/types"
	controllerRuntime "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ = Describe("SwitchAssignment Webhook", func() {
	const (
		SWANamespace        = "onmetal"
		SWALeafRole         = "Leaf"
		SWASpineRole        = "Spine"
		SWAInvalidChassisID = "0Z:0X:0Y:0A:0B:0C"
		SWAValidChassisID   = "02:ff:0f:50:60:70"
		timeout             = time.Second * 30
		interval            = time.Millisecond * 250
	)

	AfterEach(func() {
		Expect(k8sClient.DeleteAllOf(ctx, &SwitchAssignment{}, client.InNamespace(SWANamespace))).To(Succeed())
		Eventually(func() bool {
			list := &SwitchAssignmentList{}
			err := k8sClient.List(ctx, list)
			if err != nil {
				return false
			}
			if len(list.Items) > 0 {
				return false
			}
			return true
		}, timeout, interval).Should(BeTrue())
	})

	Context("On SwitchAssignment creation", func() {
		It("Should not allow to pass invalid fields values", func() {
			By("Create SwitchAssignment resource")
			ctx := context.Background()

			cr := SwitchAssignment{
				ObjectMeta: controllerRuntime.ObjectMeta{
					Name:      "test-switch-assignment",
					Namespace: SWANamespace,
				},
				Spec: SwitchAssignmentSpec{
					Role:      SWALeafRole,
					Serial:    "999999",
					ChassisID: SWAInvalidChassisID,
				},
			}

			Expect(k8sClient.Create(ctx, &cr)).ShouldNot(Succeed())

			cr.Spec.Role = SWASpineRole
			Expect(k8sClient.Create(ctx, &cr)).ShouldNot(Succeed())

			cr.Spec.ChassisID = SWAValidChassisID
			Expect(k8sClient.Create(ctx, &cr)).Should(Succeed())
			Eventually(func() bool {
				swa := SwitchAssignment{}
				namespacedName := types.NamespacedName{
					Namespace: cr.Namespace,
					Name:      cr.Name,
				}
				err := k8sClient.Get(ctx, namespacedName, &swa)
				if err != nil {
					return false
				}
				return true
			})

			By("Update SwitchAssignment resource")
			cr.Spec.Serial = "000001"
			Expect(k8sClient.Update(ctx, &cr)).ShouldNot(Succeed())
		})
	})
})
