package autoEmuErrors

import (
	"fmt"
)

type KernelOutputReadError struct {
	Err error
}

type KernelBuildError struct {
	Err error
}

func (k *KernelOutputReadError) Error() string {
	return fmt.Sprintf("Failed to read kernel and RFS images from buildroots output: %v", k.Err)
}

func (k *KernelBuildError) Error() string {
	return fmt.Sprintf("Failed to build kernel and RFS images: %v", k.Err)
}