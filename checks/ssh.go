package checks

import (
	"strconv"
	"time"

	"golang.org/x/crypto/ssh"
)

type Ssh struct {
	checkBase
	Port        int
	PubKey      string
	BadAttempts int
	Commands    []string
	Outputs     []string
}

func (c Ssh) Run(teamPrefix string, res chan Result) {
	// if  pubkey
	// var hostKey ssh.PublicKey
	// pubkey
	// else
	username, password := getCreds(c.CredLists)

	// Create client config
	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		// Hardcoded timeout of 8 seconds
		Timeout: time.Duration(8 * time.Second),
	}

	// Connect to ssh server
	conn, err := ssh.Dial("tcp", teamPrefix+c.Suffix+":"+strconv.Itoa(c.Port), config)
	if err != nil {
		res <- Result{
			Status: false,
			Error:  "error connecting or logging in to ssh server",
			Debug:  "error: " + err.Error() + " for creds " + username + ":" + password,
		}
		return
	}
	defer conn.Close()

	// Create a session
	session, err := conn.NewSession()
	if err != nil {
		res <- Result{
			Status: false,
			Error:  "unable to create ssh session",
			Debug:  "error: " + err.Error(),
		}
		return
	}
	defer session.Close()

	// Set up terminal modes
	modes := ssh.TerminalModes{
		ssh.ECHO:          0,     // disable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}

	// Request pseudo terminal
	if err := session.RequestPty("xterm", 40, 80, modes); err != nil {
		res <- Result{
			Status: false,
			Error:  "couldn't allocate pts",
			Debug:  "error: " + err.Error(),
		}
		return
	}

	// Start remote shell
	if err := session.Shell(); err != nil {
		res <- Result{
			Status: false,
			Error:  "failed to start shell",
			Debug:  "error: " + err.Error(),
		}
		return
	}

	// execute commands
	res <- Result{
		Status: true,
		Debug:  "creds used were " + username + ":" + password,
	}
}