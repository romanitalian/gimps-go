package primenet

import (
	"encoding/binary"
	"fmt"
	"os"
	"sync"
)

const (
	SpoolFileMagic = 0x2c7330a8
	SpoolFileVersion = 1
)

// Spool manages the spool file for PrimeNet messages
type Spool struct {
	filePath string
	mutex    sync.Mutex
}

// NewSpool creates a new spool manager
func NewSpool(filePath string) *Spool {
	return &Spool{
		filePath: filePath,
	}
}

// AddMessage adds a message to the spool file
func (s *Spool) AddMessage(msgType int, data []byte) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	
	file, err := os.OpenFile(s.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open spool file: %w", err)
	}
	defer file.Close()
	
	// Write message type (int16)
	if err := binary.Write(file, binary.LittleEndian, int16(msgType)); err != nil {
		return err
	}
	
	// Write data length (int16)
	if err := binary.Write(file, binary.LittleEndian, int16(len(data))); err != nil {
		return err
	}
	
	// Write data
	if _, err := file.Write(data); err != nil {
		return err
	}
	
	return nil
}

// ReadMessages reads all messages from the spool file
func (s *Spool) ReadMessages() ([]SpoolMessage, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	
	file, err := os.Open(s.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return []SpoolMessage{}, nil
		}
		return nil, fmt.Errorf("failed to open spool file: %w", err)
	}
	defer file.Close()
	
	var messages []SpoolMessage
	
	for {
		var msgType int16
		if err := binary.Read(file, binary.LittleEndian, &msgType); err != nil {
			break
		}
		
		var dataLen int16
		if err := binary.Read(file, binary.LittleEndian, &dataLen); err != nil {
			break
		}
		
		data := make([]byte, dataLen)
		if _, err := file.Read(data); err != nil {
			break
		}
		
		messages = append(messages, SpoolMessage{
			Type: int(msgType),
			Data: data,
		})
	}
	
	return messages, nil
}

// RemoveMessage removes a message from the spool file
func (s *Spool) RemoveMessage(index int) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	
	messages, err := s.ReadMessages()
	if err != nil {
		return err
	}
	
	if index < 0 || index >= len(messages) {
		return fmt.Errorf("invalid message index")
	}
	
	// Remove the message
	messages = append(messages[:index], messages[index+1:]...)
	
	// Rewrite the file
	file, err := os.Create(s.filePath)
	if err != nil {
		return fmt.Errorf("failed to create spool file: %w", err)
	}
	defer file.Close()
	
	for _, msg := range messages {
		if err := binary.Write(file, binary.LittleEndian, int16(msg.Type)); err != nil {
			return err
		}
		if err := binary.Write(file, binary.LittleEndian, int16(len(msg.Data))); err != nil {
			return err
		}
		if _, err := file.Write(msg.Data); err != nil {
			return err
		}
	}
	
	return nil
}

// SpoolMessage represents a message in the spool file
type SpoolMessage struct {
	Type int
	Data []byte
}

