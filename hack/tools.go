//go:build tools

// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package hack

import (
	// Use gen-crd-api-reference-docs for doc generation.
	_ "github.com/ahmetb/gen-crd-api-reference-docs"
	// Use addlicense for adding license headers.
	_ "github.com/google/addlicense"
)
