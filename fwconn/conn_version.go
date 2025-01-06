package fwconn

import (
	"fmt"

	"github.com/fxamacker/cbor/v2"
)

// sendVersionInfo sends the version information over the connection.
func (c *FWConn) sendVersionInfo(v *VersionInfo) error {
	data, err := cbor.Marshal(v)
	if err != nil {
		return err
	}

	// Use the existing Write method to send the data.
	if err := c.Write(data); err != nil {
		return err
	}

	if v.AgreedVersion != "" {
		// Additional logic if needed.
	} else {
		// Additional logic if needed.
	}
	return nil
}

// receiveClientVersionInfo receives the supported versions from the client.
func (c *FWConn) receiveClientVersionInfo(v *VersionInfo) error {
	data, err := c.Read()
	if err != nil {
		return err
	}

	if err := cbor.Unmarshal(data, v); err != nil {
		return err
	}

	return nil
}

// receiveServerResponse receives the agreed version from the server.
func (c *FWConn) receiveServerResponse(v *VersionInfo) error {
	responseData, err := c.Read()
	if err != nil {
		return err
	}

	var response VersionInfo
	if err := cbor.Unmarshal(responseData, &response); err != nil {
		return err
	}

	if response.AgreedVersion == "" {
		return fmt.Errorf("no agreed version received from server")
	}

	v.AgreedVersion = response.AgreedVersion
	return nil
}
