package controllers

import (
	"context"

	"github.com/onmetal/metal-api/apis/machine/v1alpha2"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ = Describe("IgnitionReconciler", func() {
	var (
		namespace = "default"
	)

	ctx := context.Background()

	It("should create ipxe configmap", func() {
		By("adding template config")
		templateConfigMap := &corev1.ConfigMap{
			TypeMeta: metav1.TypeMeta{
				Kind:       "ConfigMap",
				APIVersion: "v1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      configMapTemplateName,
				Namespace: namespace,
			},
			Data: map[string]string{
				"name":               "|\n    machine-uuid",
				"machine":            "|\n    {{ if .Machine }}isPresent{{ end }}",
				"machine-assignment": "|\n    {{ if .MachineAssignment }}isPresent{{ end }}",
			},
		}
		Expect(k8sClient.Patch(ctx, templateConfigMap, client.Apply, ignitionFieldOwner, client.ForceOwnership)).To(Succeed())

		By("creating a machine")
		machine := &v1alpha2.Machine{
			ObjectMeta: metav1.ObjectMeta{
				GenerateName: "a",
				Namespace:    namespace,
				Labels: map[string]string{
					"machine.onmetal.de/size-m5.metal": "true",
				},
			},
			Spec: v1alpha2.MachineSpec{},
		}
		Expect(k8sClient.Create(ctx, machine)).To(Succeed())

		By("updating machine status")
		machine.Status = v1alpha2.MachineStatus{
			Health:    v1alpha2.MachineStateHealthy,
			OOB:       v1alpha2.ObjectReference{Exist: true},
			Inventory: v1alpha2.ObjectReference{Exist: true},
			Interfaces: []v1alpha2.Interface{
				{Name: "test"},
				{Name: "test2"},
			},
		}
		Expect(k8sClient.Status().Update(ctx, machine)).To(Succeed())

		By("creating machine assignment")
		machineAssignment := &v1alpha2.MachineAssignment{
			ObjectMeta: metav1.ObjectMeta{
				GenerateName: "sr-",
				Namespace:    namespace,
			},
			Spec: v1alpha2.MachineAssignmentSpec{
				MachineClass: corev1.LocalObjectReference{
					Name: "m5.metal",
				},
				Image: "myimage_repo_location",
			},
		}
		Expect(k8sClient.Create(ctx, machineAssignment)).To(Succeed())

		By("waiting for ipxe configmap to be created")
		ipxeCM := &corev1.ConfigMap{}
		Eventually(func(g Gomega) {
			err := k8sClient.Get(ctx, types.NamespacedName{Namespace: namespace, Name: "ipxe-" + machineAssignment.Name}, ipxeCM)
			g.Expect(err).NotTo(HaveOccurred())
		})
	})

	It("should create ipxe secret", func() {
		By("adding secret template config")
		templateConfigMap := &corev1.Secret{
			TypeMeta: metav1.TypeMeta{
				Kind:       "Secret",
				APIVersion: "v1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      secretTemplateName,
				Namespace: namespace,
			},
			StringData: map[string]string{
				"name":               "|\n    machine-uuid",
				"machine":            "|\n    {{ if .Machine }}isPresent{{ end }}",
				"machine-assignment": "|\n    {{ if .MachineAssignment }}isPresent{{ end }}",
			},
		}
		Expect(k8sClient.Patch(ctx, templateConfigMap, client.Apply, ignitionFieldOwner, client.ForceOwnership)).To(Succeed())

		By("creating a machine")
		machine := &v1alpha2.Machine{
			ObjectMeta: metav1.ObjectMeta{
				GenerateName: "a",
				Namespace:    namespace,
				Labels: map[string]string{
					"machine.onmetal.de/size-m5.metal": "true",
				},
			},
			Spec: v1alpha2.MachineSpec{},
		}
		Expect(k8sClient.Create(ctx, machine)).To(Succeed())

		By("updating machine status")
		machine.Status = v1alpha2.MachineStatus{
			Health:    v1alpha2.MachineStateHealthy,
			OOB:       v1alpha2.ObjectReference{Exist: true},
			Inventory: v1alpha2.ObjectReference{Exist: true},
			Interfaces: []v1alpha2.Interface{
				{Name: "test"},
				{Name: "test2"},
			},
		}
		Expect(k8sClient.Status().Update(ctx, machine)).To(Succeed())

		By("creating machine assignment")
		machineAssignment := &v1alpha2.MachineAssignment{
			ObjectMeta: metav1.ObjectMeta{
				GenerateName: "sr-",
				Namespace:    namespace,
			},
			Spec: v1alpha2.MachineAssignmentSpec{
				MachineClass: corev1.LocalObjectReference{
					Name: "m5.metal",
				},
				Image: "myimage_repo_location",
			},
		}
		Expect(k8sClient.Create(ctx, machineAssignment)).To(Succeed())

		By("waiting for ipxe secret to be created")
		ipxeSecret := &corev1.Secret{}
		Eventually(func(g Gomega) {
			err := k8sClient.Get(ctx, types.NamespacedName{Namespace: namespace, Name: "ipxe-" + machineAssignment.Name}, ipxeSecret)
			g.Expect(err).NotTo(HaveOccurred())
		})
	})

})
