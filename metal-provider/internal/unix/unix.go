/*
Copyright (c) 2023 T-Systems International GmbH, SAP SE or an SAP affiliate company. All right reserved

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

package unix

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"

	"golang.org/x/sys/unix"
)

func Listen(ctx context.Context, addr string) (net.Listener, error) {
	lc := &net.ListenConfig{}
	ln, err := lc.Listen(ctx, "unix", addr)
	if errors.Is(err, unix.EADDRINUSE) {
		var si os.FileInfo
		si, err = os.Stat(addr)
		if err != nil {
			return nil, fmt.Errorf("could not stat() socket: %w", err)
		}
		if si.Mode().Type()&os.ModeSocket == 0 {
			return nil, fmt.Errorf("file exists but is not a socket")
		}
		err = os.Remove(addr)
		if err != nil {
			return nil, fmt.Errorf("could not unlink socket: %w", err)
		}
		return lc.Listen(ctx, "unix", addr)
	}
	return ln, nil
}
