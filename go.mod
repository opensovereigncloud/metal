module github.com/onmetal/k8s-inventory

go 1.15

require (
	github.com/d4l3k/messagediff v1.2.1
	github.com/go-logr/logr v0.3.0
	github.com/onmetal/k8s-size v0.0.0-20210421061640-857f330421d7
	github.com/onsi/ginkgo v1.14.1
	github.com/onsi/gomega v1.10.2
	golang.org/x/mod v0.4.2
	k8s.io/apiextensions-apiserver v0.19.2
	k8s.io/apimachinery v0.19.2
	k8s.io/client-go v0.19.2
	sigs.k8s.io/controller-runtime v0.7.0
	sigs.k8s.io/kustomize v2.0.3+incompatible
)
