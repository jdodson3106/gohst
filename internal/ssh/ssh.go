package ssh

import (
	"log"
	"os"
	"runtime"
)

const (
	// The default file names for private key algos
	RSA     = "id_rsa"
	ECDSA   = "id_ecdsa"
	ED25519 = "id_ed25519"

	// Handling Host Auth
	TOFU = iota
	KnownOnly
	PinnedHost
)

// Connectionconfig is configured on first use and the stored locally
// for re-use.
type ConnectionConfig struct {
	// The IP address of the host machine connecting to
	Host string

	// The port can be provided. If no port provided the
	// default 22 port will be assumed
	Port string

	// the keyname will be where we grab the private key from in the .ssh
	// directory to perform the handshake. The defaults are assumed unless
	// otherwise specified
	KeyName string

	// The name of your known hosts file. If not provided we assume the default
	// of known_hosts
	KnownHostsName string

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
	osType := runtime.GOOS

	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("unable to locate user home directory: %s", err)
	}

	var sshDir string
	if osType == "windows" {
		// TODO: handle setting up on windows
		return nil
	} else {
		sshDir = homeDir + "/.ssh"
	}

	return &ConnectionConfig{
		Host:             host,
		Port:             "22",
		KeyName:          sshDir + "/" + ED25519,
		KnownHostsName:   sshDir + "/known_hosts",
		HostAuthProtocol: TOFU,
	}
}

type SSHClient struct {
	config *ConnectionConfig
}

// NewClient creates a new SSHClient with a default configuration
func NewClient(hostIP string) *SSHClient {
	return &SSHClient{
		config: DefaultConfiguration(hostIP),
	}
}
