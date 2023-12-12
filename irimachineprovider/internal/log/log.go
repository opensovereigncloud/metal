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

package log

import (
	"context"
	"io"

	"github.com/go-logr/logr"
	"github.com/go-logr/zerologr"
	"github.com/rs/zerolog"
)

func Setup(ctx context.Context, dev, silent bool, writer io.Writer) context.Context {
	var zeroLog zerolog.Logger

	if silent {
		return logr.NewContext(ctx, logr.Discard())
	}

	if dev {
		cw := zerolog.ConsoleWriter{
			Out:           writer,
			TimeFormat:    "2006-01-02 15:04:05 MST",
			FieldsExclude: []string{"v"},
		}
		zeroLog = zerolog.New(cw).Level(zerolog.DebugLevel).With().Timestamp().Logger()
	} else {
		zeroLog = zerolog.New(writer).Level(zerolog.DebugLevel).With().Timestamp().Logger()
	}

	return logr.NewContext(ctx, zerologr.New(&zeroLog))
}

func Info(ctx context.Context, msg string, keysAndValues ...interface{}) {
	logr.FromContextOrDiscard(ctx).Info(msg, keysAndValues...)
}

func Debug(ctx context.Context, msg string, keysAndValues ...interface{}) {
	logr.FromContextOrDiscard(ctx).V(1).Info(msg, keysAndValues...)
}

func Error(ctx context.Context, err error, keysAndValues ...interface{}) {
	logr.FromContextOrDiscard(ctx).Error(err, "", keysAndValues...)
}

func WithValues(ctx context.Context, keysAndValues ...interface{}) context.Context {
	return logr.NewContext(ctx, logr.FromContextOrDiscard(ctx).WithValues(keysAndValues...))
}
