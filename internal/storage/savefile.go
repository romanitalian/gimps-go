package storage

import (
	"encoding/binary"
	"fmt"
	"os"
)

const (
	LLSaveFileMagic = 0x2c7330a8
	LLSaveFileVersion = 1
)

// SaveFile manages save files for work units
type SaveFile struct {
	filePath string
}

// NewSaveFile creates a new save file manager
func NewSaveFile(filePath string) *SaveFile {
	return &SaveFile{
		filePath: filePath,
	}
}

// LLSaveData represents Lucas-Lehmer save file data
type LLSaveData struct {
	ErrorCount   uint32
	Counter      uint32
	ShiftCount   uint32
	FFTData      []byte
}

// WriteLLSave writes Lucas-Lehmer save file
func (sf *SaveFile) WriteLLSave(data *LLSaveData) error {
	file, err := os.Create(sf.filePath)
	if err != nil {
		return fmt.Errorf("failed to create save file: %w", err)
	}
	defer file.Close()
	
	// Write header (52 bytes - standard header for all work types)
	header := make([]byte, 52)
	binary.LittleEndian.PutUint32(header[0:4], LLSaveFileMagic)
	binary.LittleEndian.PutUint32(header[4:8], LLSaveFileVersion)
	// Add work unit info to header
	if _, err := file.Write(header); err != nil {
		return err
	}
	
	// Write error count
	if err := binary.Write(file, binary.LittleEndian, data.ErrorCount); err != nil {
		return err
	}
	
	// Write iteration counter
	if err := binary.Write(file, binary.LittleEndian, data.Counter); err != nil {
		return err
	}
	
	// Write shift count
	if err := binary.Write(file, binary.LittleEndian, data.ShiftCount); err != nil {
		return err
	}
	
	// Write FFT data length
	if err := binary.Write(file, binary.LittleEndian, uint32(len(data.FFTData))); err != nil {
		return err
	}
	
	// Write FFT data
	if _, err := file.Write(data.FFTData); err != nil {
		return err
	}
	
	return nil
}

// ReadLLSave reads Lucas-Lehmer save file
func (sf *SaveFile) ReadLLSave() (*LLSaveData, error) {
	file, err := os.Open(sf.filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open save file: %w", err)
	}
	defer file.Close()
	
	// Read header
	header := make([]byte, 52)
	if _, err := file.Read(header); err != nil {
		return nil, err
	}
	
	magic := binary.LittleEndian.Uint32(header[0:4])
	if magic != LLSaveFileMagic {
		return nil, fmt.Errorf("invalid save file magic number")
	}
	
	data := &LLSaveData{}
	
	// Read error count
	if err := binary.Read(file, binary.LittleEndian, &data.ErrorCount); err != nil {
		return nil, err
	}
	
	// Read iteration counter
	if err := binary.Read(file, binary.LittleEndian, &data.Counter); err != nil {
		return nil, err
	}
	
	// Read shift count
	if err := binary.Read(file, binary.LittleEndian, &data.ShiftCount); err != nil {
		return nil, err
	}
	
	// Read FFT data length
	var fftLen uint32
	if err := binary.Read(file, binary.LittleEndian, &fftLen); err != nil {
		return nil, err
	}
	
	// Read FFT data
	data.FFTData = make([]byte, fftLen)
	if _, err := file.Read(data.FFTData); err != nil {
		return nil, err
	}
	
	return data, nil
}

// Exists checks if save file exists
func (sf *SaveFile) Exists() bool {
	_, err := os.Stat(sf.filePath)
	return err == nil
}

