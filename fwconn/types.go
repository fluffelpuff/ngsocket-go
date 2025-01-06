package fwconn

type _ConnType uint8

type Role int

type VersionInfo struct {
	SupportedVersions []string `cbor:"supported_versions"`
	AgreedVersion     string   `cbor:"agreed_version,omitempty"`
}
