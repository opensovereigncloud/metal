module github.com/onmetal/switch-operator

go 1.15

require (
	github.com/apparentlymart/go-cidr v1.1.0
	github.com/go-logr/logr v0.4.0
	github.com/google/uuid v1.2.0
	github.com/onmetal/ipam v0.0.0-20211029144623-1398cd13a1ae
	github.com/onmetal/k8s-inventory v0.0.0-20210914152907-c3b127f0ebe2
	github.com/onsi/ginkgo v1.16.1
	github.com/onsi/gomega v1.11.0
	github.com/pkg/errors v0.9.1
	go.uber.org/zap v1.16.0 // indirect
	golang.org/x/crypto v0.0.0-20210322153248-0c34fe9e7dc2 // indirect
	golang.org/x/mod v0.4.2
	gopkg.in/yaml.v3 v3.0.0-20200615113413-eeeca48fe776
	k8s.io/api v0.20.5
	k8s.io/apimachinery v0.20.5
	k8s.io/client-go v0.20.5
	sigs.k8s.io/controller-runtime v0.8.3
)
