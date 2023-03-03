package controllers

import "github.com/pkg/errors"

var (
	errKubernetesEndpointIsEmpty            = errors.New("kubernetes endpoint subset is empty")
	errKubernetesEndpointAddressIsEmpty     = errors.New("kubernetes endpoint subset address is empty")
	errKubernetesEndpointAddressPortIsEmpty = errors.New("kubernetes endpoint subset address port is empty")
)
