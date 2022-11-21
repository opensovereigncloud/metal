/*
Copyright (c) 2021 T-Systems International GmbH, SAP SE or an SAP affiliate company. All right reserved

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

package scheduler

import (
	"bytes"
	"github.com/Masterminds/sprig"
	inventoryv1alpaha1 "github.com/onmetal/metal-api/apis/inventory/v1alpha1"
	"math"
	"net"
	yaml "sigs.k8s.io/yaml"
	"strconv"
	"text/template"
)

func getLabelSelectorForAvailableMachine(instanceType string) map[string]string {
	sizeLabel := inventoryv1alpaha1.CLabelPrefix + instanceType
	return map[string]string{sizeLabel: "true"}
}

func parseTemplate(temp map[string]string, wrapper IgnitionWrapper) (map[string]string, error) {
	tempStr, err := yaml.Marshal(temp)
	if err != nil {
		return nil, err
	}

	t, err := template.New("temporaryTemplate").Funcs(sprig.HermeticTxtFuncMap()).Parse(string(tempStr))
	if err != nil {
		return nil, err
	}

	var b bytes.Buffer
	err = t.Execute(&b, wrapper)
	if err != nil {
		return nil, err
	}

	var tempMap = make(map[string]string)
	if err = yaml.Unmarshal(b.Bytes(), &tempMap); err != nil {
		return nil, err
	}

	return tempMap, nil
}

// TODO(flpeter) duplicate code from switch configurer
func calculateAsn(addr net.IP) string {
	var asn int
	asn += int(addr[13]) * int(math.Pow(2, 16))
	asn += (int(addr[14])) * int(math.Pow(2, 8))
	asn += int(addr[15])
	return strconv.Itoa(IgnitionASNBase + asn)
}
