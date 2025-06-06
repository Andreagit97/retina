// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package outputlocation

import (
	"context"
	"io"
	"os"
	"path/filepath"

	"go.uber.org/zap"

	captureConstants "github.com/microsoft/retina/pkg/capture/constants"
	"github.com/microsoft/retina/pkg/log"
)

type HostPath struct {
	l *log.ZapLogger
}

var _ Location = &HostPath{}

func NewHostPath(logger *log.ZapLogger) Location {
	return &HostPath{l: logger}
}

func (hp *HostPath) Name() string {
	return "HostPath"
}

func (hp *HostPath) Enabled() bool {
	hostPath := os.Getenv(string(captureConstants.CaptureOutputLocationEnvKeyHostPath))
	if len(hostPath) == 0 {
		hp.l.Debug("Output location is not enabled", zap.String("location", hp.Name()))
		return false
	}
	return true
}

func (hp *HostPath) Output(_ context.Context, srcFilePath string) error {
	hostPath := os.Getenv(string(captureConstants.CaptureOutputLocationEnvKeyHostPath))
	hp.l.Info("Copy file",
		zap.String("location", hp.Name()),
		zap.String("source file path", srcFilePath),
		zap.String("destination file path", hostPath),
	)

	srcFile, err := os.Open(srcFilePath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	fileName := filepath.Base(srcFilePath)
	fileHostPath := filepath.Join(hostPath, fileName)
	destFile, err := os.Create(fileHostPath)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, srcFile)
	return err
}
