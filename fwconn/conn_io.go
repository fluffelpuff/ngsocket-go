package fwconn

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"time"
)

// Write sends the provided byte slice over the connection by splitting it into frames.
// Enhanced error handling: avoids panic, retries on transient errors, and provides detailed error information.
func (c *FWConn) Write(b []byte) error {
	return c._Write(b)
}
func (c *FWConn) _Write(b []byte) error {
	// Generate Process ID
	procid, err := randomUint32()
	if err != nil {
		// Log the error and return it
		return fmt.Errorf("process ID generation failed: %w", err)
	}

	// Split the payload into frames
	frames := splitIntoFrames(b, fragmentSize, procid)
	if len(frames) > int(c.maxFrameSize) {
		return fmt.Errorf("to many number of frames (%d)", c.maxFrameSize)
	}

	// Transmit individual frames
	for _, frame := range frames {
		var writeErr error
		for attempt := 1; attempt <= maxWriteRetries; attempt++ {
			writeErr = c.writeFrame(frame)
			if writeErr == nil {
				break // Successfully written, proceed to next frame
			}

			// Check if the error is temporary (e.g., network issues)
			if isTemporaryError(writeErr) && attempt < maxWriteRetries {
				waitDuration := time.Duration(attempt) * retryBaseDelay
				time.Sleep(waitDuration)
				continue
			}

			break
		}

		if writeErr != nil {
			// Return the error with context
			return writeErr
		}
	}

	return nil
}

// Read liest Daten von der Verbindung in Fragmenten.
func (c *FWConn) Read() ([]byte, error) {
	return c._Read()
}
func (c *FWConn) _Read() ([]byte, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	var (
		framesRead uint64
		dataBuffer = bytes.Buffer{}
	)

	for {
		if framesRead >= c.maxFrameSize {
			return nil, fmt.Errorf("maximum number of frames (%d) exceeded", c.maxFrameSize)
		}

		frame, err := c.readFrame()
		if err != nil {
			// Prüfe, ob der Fehler ein EOF-Fehler ist
			if errors.Is(err, io.EOF) {
				if framesRead == 0 {
					return nil, io.EOF
				}
				return dataBuffer.Bytes(), nil // Rückgabe der bisher gelesenen Daten
			}
			return nil, err
		}

		// Überprüfe, ob die Frame-Nummer wie erwartet ist
		if frame.Header.FrameNo != framesRead {
			return nil, fmt.Errorf("invalid frame number: expected %d, got %d", framesRead, frame.Header.FrameNo)
		}

		framesRead++
		if _, err := dataBuffer.Write(frame.GetBodyBytes()); err != nil {
			return nil, err
		}

		// Prüfe, ob dies das letzte Frame ist
		if frame.Header.LastFrame {
			break
		}
	}

	return dataBuffer.Bytes(), nil
}

// Close closes the connection, releasing any resources.
func (c *FWConn) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Close the underlying network connection
	return c.conn.Close()
}
