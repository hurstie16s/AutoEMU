package autoEmuErrors

import "fmt"

type FileRelocationError struct {
	Err error
}

type ReadLinkError struct {
	Err error
}

type ISOGenerationError struct {
	Err error
}

type DeletionError struct {
	Machine string
	Err error
}

type UnzipError struct {
	FileName string
	Err error
}

func (f *FileRelocationError) Error() string {
	return fmt.Sprintf("Failed to move firmware into workspace: %v", f.Err)
}

func (r *ReadLinkError) Error() string {
	return fmt.Sprintf("Readlink failure: %v", r.Err)
}

func (i *ISOGenerationError) Error() string {
	return fmt.Sprintf("Failed to generate ISO image from extracted RFS: %v", i.Err)
}

func (d *DeletionError) Error() string {
	return fmt.Sprintf("Failed to delete machine %s: %v", d.Machine, d.Err)
}

func (u *UnzipError) Error() string {
	return fmt.Sprintf("failed to unzip file %s: %v", u.FileName, u.Err)
}