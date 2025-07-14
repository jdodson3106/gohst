package ssh

import (
	"log"
	"os"
)

const (
	// The default file names for private key algos
	ID_RSA     = "id_rsa"
	ID_ECDSA   = "id_ecdsa"
	ID_ED25519 = "id_ed25519"

	// the keytypes
	ED25519 = iota
	RSA
	ECDSA

	// Handling Host Auth
	TOFU = iota
	KnownOnly
	PinnedHost
)

type KeyType int

// Connectionconfig is configured on first use and the stored locally
// for re-use.
type ConnectionConfig struct {
	// The IP address of the host machine connecting to
	Host string

	// The port can be provided. If no port provided the
	// default 22 port will be assumed
	Port string

	// this defines how the public key of the host is verified
	// we do not support blind handshakes, however trust on first use (TOFU)
	// is the default and initial connection is added to your known_hosts file
	HostAuthProtocol int

	// if this is provided the host authed is converted to PinnedHost auth and
	// we perform the handshake using this public key for the host.
	// If this is nil and the HostAuthProtocol is set to PinnedHost the handshake
	// will fail and return an error
	HostKey []byte
}

func DefaultConfiguration(host string) *ConnectionConfig {
	return &ConnectionConfig{
		Host:             host,
		Port:             "22",
		HostAuthProtocol: TOFU,
	}
}

type KeyHost struct {
	// The KeyType represents the algo type used to sign the key
	// supported types are ED25519, RSA, and ECDSA
	KeyType int

	// the keyname will be where we grab the private key from in the .ssh
	// directory to perform the handshake. The defaults are assumed unless
	// otherwise specified
	KeyName string

	// The name of your known hosts file. If not provided we assume the default
	// of known_hosts
	KnownHostsName string
}

// DefaultKeyHost returns a defalt config of the ssh key files
// assumes id_ed25519 for the private key and known_hosts for the
// known hosts file
func DefaultKeyHost() *KeyHost {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("unable to locate user home directory: %s", err)
		return &KeyHost{}
	}

	sshDir := home + "/.ssh/"
	return &KeyHost{
		KeyType:        ED25519,
		KeyName:        sshDir + ID_ED25519,
		KnownHostsName: sshDir + "known_hosts",
	}
}

type SSHClient struct {
	// KeyHost is the settings for where the private key and
	// known hosts files live in the user's system
	KeyHost *KeyHost

	// config is the the connection details for how to
	// connect to the host
	config *ConnectionConfig
}

// NewClient creates a new SSHClient with a default configuration
func NewClient(hostIP string) *SSHClient {
	return &SSHClient{
		config:  DefaultConfiguration(hostIP),
		KeyHost: DefaultKeyHost(),
	}
}

// NewClientWithKeyHost returns a new SSHClient with the provided KeyHost
// used when the machine doesn't hold the default file locations
func NewClientWithKeyHost(hostIP string, kh *KeyHost) *SSHClient {
	return &SSHClient{
		config:  DefaultConfiguration(hostIP),
		KeyHost: kh,
	}
}
