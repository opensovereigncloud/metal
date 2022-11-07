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
