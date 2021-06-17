/*
Copyright 2021 T-Systems International GmbH, SAP SE or an SAP affiliate company. All right reserved

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"strings"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gopkg.in/yaml.v3"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ = Describe("Switch Webhook", func() {
	const (
		SwitchNamespace = "onmetal"
		timeout         = time.Second * 30
		interval        = time.Millisecond * 250
	)

	AfterEach(func() {
		Expect(k8sClient.DeleteAllOf(ctx, &Switch{}, client.InNamespace(SwitchNamespace))).To(Succeed())
		Eventually(func() bool {
			list := &SwitchList{}
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

	Context("On Switch creation", func() {
		It("Should set label", func() {
			By("Create Switch resource")
			ctx := context.Background()
			sample := filepath.Join("..", "..", "config", "samples", "spine-0-1.onmetal.de_v1alpha1_switch.yaml")
			rawInfo := make(map[string]interface{})
			sw := &Switch{}
			sampleBytes, err := ioutil.ReadFile(sample)
			Expect(err).NotTo(HaveOccurred())
			err = yaml.Unmarshal(sampleBytes, rawInfo)
			data, err := json.Marshal(rawInfo)
			Expect(err).NotTo(HaveOccurred())
			err = json.Unmarshal(data, sw)
			Expect(err).NotTo(HaveOccurred())
			Expect(k8sClient.Create(ctx, sw)).To(Succeed())
			Expect(sw.Labels).Should(Equal(map[string]string{LabelChassisId: strings.ReplaceAll(sw.Spec.SwitchChassis.ChassisID, ":", "-")}))
		})
	})
})
