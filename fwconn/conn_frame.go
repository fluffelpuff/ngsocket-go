package fwconn

import (
	"encoding/binary"
	"fmt"
	"hash/crc32"
	"io"
)

// readFrame liest ein einzelnes Frame von der Verbindung.
func (c *FWConn) readFrame() (*_Frame, error) {
	headerSize := 21
	headerBuf := make([]byte, headerSize)

	// Read the 21 header bytes
	n, err := io.ReadFull(c.conn, headerBuf)
	if err != nil {
		if err == io.EOF {
			return nil, fmt.Errorf("connection closed while reading header: %w", err)
		}
		return nil, fmt.Errorf("failed to read frame header (%d bytes read): %w", n, err)
	}

	// Parse Header
	h := &_Header{}
	h.DataLength = binary.BigEndian.Uint32(headerBuf[0:4])
	h.Checksum = binary.BigEndian.Uint32(headerBuf[4:8])
	h.ProcessId = binary.BigEndian.Uint32(headerBuf[8:12])
	h.FrameNo = binary.BigEndian.Uint64(headerBuf[12:20])
	h.LastFrame = headerBuf[20] == 1

	// Read the body if present (DataLength can be 0)
	var bodyBuf []byte
	if h.DataLength > 0 {
		bodyBuf = make([]byte, h.DataLength)

		n, err := io.ReadFull(c.conn, bodyBuf)
		if err != nil {
			if err == io.EOF {
				return nil, fmt.Errorf("connection closed while reading body: %w", err)
			}
			return nil, fmt.Errorf("failed to read frame body (%d bytes read): %w", n, err)
		}

		// Verify checksum
		csum := crc32.ChecksumIEEE(bodyBuf)
		if csum != h.Checksum {
			return nil, fmt.Errorf("invalid checksum: expected %d, got %d", h.Checksum, csum)
		}
	}

	// Construct the frame
	frame := &_Frame{
		Header: h,
		Body:   bodyBuf,
	}

	return frame, nil
}

// writeFrame sends a single frame over the connection.
func (c *FWConn) writeFrame(frame *_Frame) error {
	bytedFrame := frameToBytes(frame)
	total := len(bytedFrame)
	written := 0

	for written < total {
		n, err := c.conn.Write(bytedFrame[written:])
		if err != nil {
			return fmt.Errorf("write error: %w", err)
		}
		if n == 0 {
			return fmt.Errorf("connection might be broken: wrote 0 bytes")
		}
		written += n
	}
	return nil
}
