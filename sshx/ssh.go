package sshx

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/crypto/ssh"
)

var (
	RunShellErr = errors.New("run shell error")
)

// SSHSess is session ssh
type SSHSess struct {
	*ssh.Session
	stdin  io.WriteCloser
	stdout io.Reader
	stderr io.Reader
}

// Login a new ssh session
func Login(username, password, hostname, port string, secure bool) (*SSHSess, error) {
	config := &ssh.ClientConfig{}
	if secure {
		// SSH client config
		config = &ssh.ClientConfig{
			User: username,
			Auth: []ssh.AuthMethod{
				ssh.Password(password),
			},
			// Non-production only
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
			Timeout:         time.Second * 2,
		}
	} else {
		// SSH client config
		config = &ssh.ClientConfig{
			User: username,
			Auth: []ssh.AuthMethod{
				ssh.Password(password),
			},
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
			Timeout:         time.Second * 2,
		}
	}

	// Connect to host
	client, err := ssh.Dial("tcp", hostname+":"+port, config)
	if err != nil {
		return nil, errors.Wrapf(err, "dail to host:port [%s:%s]", hostname, port)
	}

	sess, err := client.NewSession()
	if err != nil {
		return nil, errors.Wrap(err, "new session")
	}

	stdBuf, err := sess.StdoutPipe()
	if err != nil {
		return nil, errors.Wrap(err, "get std out")
	}

	stdin, err := sess.StdinPipe()
	if err != nil {
		return nil, errors.Wrap(err, "get std in")
	}

	stderr, err := sess.StderrPipe()
	if err != nil {
		return nil, errors.Wrap(err, "get std err")
	}

	err = sess.Shell()
	if err != nil {
		return nil, errors.Wrap(err, "session shell")
	}

	s := &SSHSess{Session: sess, stdin: stdin, stdout: stdBuf, stderr: stderr}
	return s, nil
}

// SendCommand send a command to an active ssh session and returns the output
func (s *SSHSess) SendCommand(command string) (string, error) {
	_, err := fmt.Fprintf(s.stdin, "%s\n", command)
	if err != nil {
		return "", errors.Wrap(err, "send command")
	}
	r := string(readConnection(s.stdout))
	return r, nil
}

func readConnection(stdBuf io.Reader) []byte {
	buf := make([]byte, 0, 4096) // big buffer
	tmp := make([]byte, 1024)    // using small tmo buffer for demonstrating
	for {
		n, err := stdBuf.Read(tmp)
		if err != nil {
			if err != io.EOF {
				// log.Fatal("read error:", err)
			}
			break
		}
		buf = append(buf, tmp[:n]...)
		if 1024 > n {
			break
		}
	}
	return buf
}

// ConnectAndCommand connect to host and send command
func ConnectAndCommand(username, password, hostname, port, command string, secure bool, timeout int, envs ...map[string]string) (string, error) {
	config := &ssh.ClientConfig{}

	if secure {
		// SSH client config
		config = &ssh.ClientConfig{
			User: username,
			Auth: []ssh.AuthMethod{
				ssh.Password(password),
			},
			// Non-production only
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
			Timeout:         time.Duration(timeout) * time.Second,
		}
	} else {
		// SSH client config
		config = &ssh.ClientConfig{
			User: username,
			Auth: []ssh.AuthMethod{
				ssh.Password(password),
			},
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
			Timeout:         time.Duration(timeout) * time.Second,
		}
	}

	// Connect to host
	client, err := ssh.Dial("tcp", hostname+":"+port, config)
	if err != nil {
		return "", errors.Wrapf(err, "dail to host:port [%s:%s]", hostname, port)
	}

	defer client.Close()

	// Create session
	sess, err := client.NewSession()
	if err != nil {
		return "", errors.Wrap(err, "new session")
	}
	defer sess.Close()
	// set envs
	var envKv []string
	if len(envs) > 0 {
		for _, env := range envs {
			for k, v := range env {
				envKv = append(envKv, fmt.Sprintf("%s=%s", k, v))
			}
		}
	}
	// exe command
	commandEnv := fmt.Sprintf("%s %s", strings.Join(envKv, " "), command)
	res, err := sess.CombinedOutput(commandEnv)
	if err != nil {
		return string(res), errors.Wrap(RunShellErr, err.Error())
	}
	return string(res), nil
}
