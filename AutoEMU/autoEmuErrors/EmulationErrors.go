package autoEmuErrors

import "fmt"

type FirmwareExtractionError struct {
	Err error
}

type RFSConversionError struct {
	Err error
}

type ArchitectureUnsupportedError struct {
	Err error
}

type MachineExportError struct {
	Machine string
	Err error
}

func (f *FirmwareExtractionError) Error() string {
	return fmt.Sprintf("Extraction script failed to extract firmware from binary file: %v", f.Err)
}

func (r *RFSConversionError) Error() string {
	return fmt.Sprintf("Failed to convert root file system image to correct format: %v", r.Err)
}

func (a *ArchitectureUnsupportedError) Error() string {
	return fmt.Sprintf("Architecture %v is not supported by AutoEMU at this time", a.Err)
}

func (m *MachineExportError) Error() string {
	return fmt.Sprintf("Failed to export machine %s: %v", m.Machine, m.Err)
}