/*
Package ssh is simplified encapsulation of the standard x/crypto/ssh
package with reasonable defaults and multi-client controller that can
coordinate operations on multiple targets.
*/
package ssh

import (
	"bytes"
	"fmt"
	"log"
	"math/rand"
	"net"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"
)

// DefaultPort for SSH connetions.
var DefaultPort = 22

// DefaultTCPTimeout is the default number of seconds to wait to complete a TCP
// connection.
var DefaultTCPTimeout = 300 * time.Second

// ------------------------------ Errors ------------------------------

type AllUnavailable struct{}

func (AllUnavailable) Error() string { return `all SSH client targets are unavailable` }

type CommandLineMissing struct{}

func (CommandLineMissing) Error() string { return `missing argument containing command line` }

// ------------------------------- User -------------------------------

// User represents a single SSH user on the target host authenticated by
// a private key. A User may be safely marshaled/unmarshaled from
// JSON/YAML.
type User struct {

	// Name of user on target system hosting SSH server.
	Name string

	// Private key in PEM format.
	Key string
}

// Signer trims and parses then content of Key and returns a new
// [crypto/ssh.Signer].
func (u *User) Signer() (ssh.Signer, error) {
	trimmed := []byte(strings.TrimSpace(u.Key))
	return ssh.ParsePrivateKey(trimmed)
}

// ------------------------------- Host -------------------------------

// Host represents a single host on the network that is hosting
// a secure shell server. A Host may be safely marshaled/unmarshaled
// to/from JSON/YAML.
type Host struct {

	// Network host name or IP address (required).
	Addr string

	// Complete line taken in the authorized_hosts format (optional). When
	// included triggers returning ssh.FixedHostKey(pubkey) for
	// KeyCallback.
	Auth string
}

// KeyCallback returns ssh.FixedHostKey(pubkey) where pubkey is derived
// from the Auth string if Auth is not nil. Otherwise, returns
// ssh.InsecureIgnoreHostKey().
func (h *Host) KeyCallback() (ssh.HostKeyCallback, error) {
	auth := strings.TrimSpace(h.Auth)
	if len(auth) == 0 {
		return ssh.InsecureIgnoreHostKey(), nil
	}
	// Netkey is an RFC 4234 (section 6.6) which is an ssh.PublicKey
	// but in RFC format (which fails for ssh.FixedHostKey).
	netkey, _, _, _, err := ssh.ParseAuthorizedKey([]byte(auth))
	if err != nil {
		return nil, err
	}
	pubkey, err := ssh.ParsePublicKey(netkey.Marshal())
	if err != nil {
		return nil, err
	}
	return ssh.FixedHostKey(pubkey), nil
}

// ------------------------------ Client ------------------------------

// Client encapsulates an internal ssh.Client and associates a single
// user, host, and port number to target for specific ssh server connection
// adding a Connect method which implicitly dials up a connection
// setting Connected() and internally caching the client as SSHClient().
// Client may be safely marshaled to/from YAML/JSON directly.
type Client struct {

	// Host contains the host name or IP address along with any
	// authorization credentials and other data that would normally be
	// contained in authorized_hosts.
	Host *Host

	// Port is the port the ssh server is listening on on target server.
	// If unset DefaultPort is used.
	Port int

	// User contains the name of the user on target server and contains
	// the PEM private key authentication data as well.
	User *User

	// Timeout is the default number of seconds to wait to complete a TCP
	// connection. If unset DefaultTCPTimeout is used.
	Timeout time.Duration

	// Comment allows information comments about a specific client
	// connection to be persisted with the configuration data.
	Comment string

	sshclient *ssh.Client
	connected bool
	lasterror error
}

// SSHClient returns a pointer to the internal ssh.Client used for all
// connections and sessions. Only set after first call to Connect.
func (c *Client) SSHClient() *ssh.Client { return c.sshclient }

// Connected returns the last connection state of the internal SSH
// client. This is set to true on Connect. This does not guarantee that
// the current connection is still valid, just the last attempt.
func (c *Client) Connected() bool { return c.connected }

// LastError returns the last error (if any) from an attempt to Connect.
// When set Connected is guaranteed to return false.
func (c *Client) LastError() error { return c.lasterror }

// Addr returns network address suitable for use in TCP/IP connection
// strings. If the Port and Host are zero values returns empty host,
// colon, and DefaultPort (without setting it).
func (c *Client) Addr() string {
	port := DefaultPort
	if c.Port > 0 {
		port = c.Port
	}
	return fmt.Sprintf("%v:%v", c.Host.Addr, port)
}

// Dest returns the Addr with the [User.Name] prepended with an at (@)
// sign (if assigned).
func (c *Client) Dest() string {
	it := c.Addr()
	if c.User != nil && len(c.User.Name) > 0 {
		it = c.User.Name + `@` + it
	}
	return it
}

// Connect creates a new ssh.Client using the [Client.Addr] and caches
// it internally (see [SSHClient]). If [Client.Timeout] is zero uses
// [ssh.DefaultTCPTimeout]. If [Client.Port] is zero uses
// [ssh.DefaultPort]. Always reinitializes a new connection even if
// [Connected] is true.  Also see [User.Signer] and [Host.KeyCallback].
// If an attempted connection fails sets [Client.Connected] to false and
// assigns and returns [LastError].
func (c *Client) Connect() error {
	signer, err := c.User.Signer()
	if err != nil {
		return err
	}
	callback, err := c.Host.KeyCallback()
	if err != nil {
		return err
	}
	timeout := DefaultTCPTimeout
	if c.Timeout != 0 {
		timeout = c.Timeout
	}
	c.sshclient, err = ssh.Dial(`tcp`, c.Addr(), &ssh.ClientConfig{
		User:            c.User.Name,
		Auth:            []ssh.AuthMethod{ssh.PublicKeys(signer)},
		HostKeyCallback: callback,
		Timeout:         timeout,
	})
	if err == nil {
		c.connected = true
	} else {
		c.connected = false
		c.lasterror = err
	}
	return err
}

// Run sends the command with optional standard input to the currently
// open client SSH target as a new ssh.Session. If the SSH connection
// has not yet been established (c.Connected is false) [Connect] is
// called to establish a new client connection. Run returns an error if
// one is generated by the [ssh/Session.Run] call or if a new session
// could not be created (including attempting a new session on
// a timed-out connection or one that has been closed for any other
// reason.) It is the responsibility of the called to respond to such
// errors according to controller policy and associated method calls.
func (c *Client) Run(cmd string, stdin []byte) (stdout, stderr string, err error) {
	if c.sshclient == nil {
		err = c.Connect()
		if err != nil {
			return
		}
	}
	var sess *ssh.Session
	sess, err = c.sshclient.NewSession()
	if err != nil {
		return
	}
	if len(stdin) > 0 {
		sess.Stdin = bytes.NewReader(stdin)
	}
	_out := new(strings.Builder)
	_err := new(strings.Builder)
	sess.Stdout = _out
	sess.Stderr = _err
	sess.Run(cmd)
	stdout = _out.String()
	stderr = _err.String()
	sess.Close()
	return
}

// ---------------------------- Controller ----------------------------

// Controller is responsible for coordinating work requests destined for
// the target ssh servers as contained in its list of Clients.
// A Controller zero value (one with nil clients list) is safe to use
// for all methods but [Init] can be called to add clients as
// a convenience.
//
//	ctl := new(ssh.Controller).Init(cl1,cl2)
type Controller struct {
	Clients []*Client
}

// Init returns a pointer to a Controller with the Clients list
// initialized. Init can be called later to set the internal Client list
// to a newly created one. When called on a Controller with an existing
// Clients list replaces it with a new list. If no clients are passed,
// simply initializes the internal Clients list to an empty list.
func (c *Controller) Init(clients ...*Client) *Controller {
	if len(clients) > 0 {
		c.Clients = clients
		return c
	}
	c.Clients = make([]*Client, 0)
	return c
}

func (c *Controller) LogStatus() {
	for _, c := range c.Clients {
		log.Printf("%v %v %v\n", c.Dest(), c.connected, c.lasterror)
	}
}

// Connect synchronously calls Connect on all Clients in order ensuring that all
// have successfully connected before returning. No
// attempt at error checking for successful connections is attempted but
// the [Client.Connected] and [Client.LastError] can
// be checked when needed. A reference to self is returned as
// convenience.
func (c *Controller) Connect() *Controller {
	if c.Clients == nil {
		return nil
	}
	for _, client := range c.Clients {
		client.Connect()
	}
	return c
}

// RandomClient returns a random active client from the Clients list
// skipping any that are not connected. Returns nil if no connected
// clients are available.
func (c *Controller) RandomClient() *Client {
	var tried int
	count := len(c.Clients)
	if count == 0 {
		return nil
	}
	n := rand.Intn(count)
	for {
		client := c.Clients[n]
		if client.connected {
			return client
		}
		tried += 1
		if tried > count {
			return nil
		}
		n += 1
		if n >= count {
			n = 0
		}
	}
	return nil
}

// RunOnAny calls [Client.Run] on a random client from the [Clients]
// list. If error returned is of type [net.OpError] the
// [Client.Connected] is set to false and the next client in the
// [Clients] order is attempted. Then client producing the error has
// [Client.Connect] called in a separate goroutine (which, if successful,
// restores its [Client.Connected] status to true). If none of the
// clients are connected then an [AllUnavailable] error is returned.
// The cmd is always required but stdin may be nil.
func (c *Controller) RunOnAny(cmd string, stdin []byte) (stdout, stderr string, err error) {

	client := c.RandomClient()
	if client == nil {
		err = AllUnavailable{}
		return
	}

	stdout, stderr, err = client.Run(cmd, stdin)
	if _, is := err.(*net.OpError); is {
		client.connected = false
		go client.Connect()
		c.RunOnAny(cmd, stdin)
	}

	return
}
