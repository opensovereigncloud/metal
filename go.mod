module github.com/onmetal/switch-operator

go 1.15

require (
	github.com/go-logr/logr v0.4.0
	github.com/google/uuid v1.2.0 // indirect
	github.com/onmetal/k8s-inventory v0.0.0-20210519063844-3509e56a2416
	github.com/onmetal/k8s-network-global v0.0.0-20210528142724-3da4d0e4351e
	github.com/onmetal/k8s-subnet v0.0.0-20210609114747-2f4101a2caa5
	github.com/onsi/ginkgo v1.15.1
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
