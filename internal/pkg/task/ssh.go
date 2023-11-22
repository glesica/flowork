package task

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"
	"sync"
	"time"

	"golang.org/x/crypto/ssh"

	"github.com/glesica/flowork/internal/pkg/files"
	"github.com/glesica/flowork/internal/pkg/option"
)

// A Machine describes a single machine that can be used to run tasks.
// Unless otherwise noted, if any of the fields are omitted, but necessary,
// then the default values set on the SshRunner itself will be used, if
// possible.
//
// TODO: Support toml and json serialization
type Machine struct {
	// Name is the human-readable name to use for this machine. If it is
	// empty, the Addr field will be used for this purpose.
	Name string

	// Addr is the host:port combination to connect to. This field is
	// required.
	Addr string

	// User is the name of the user that will be used to connect and run
	// commands on the remote machine. It will also be used as the default
	// group name.
	User string

	// Pass is the password to use when connecting.
	Pass string

	// KeyPath is the pass to the public / private key combination to use
	// for connecting to the machine. If a key path is provided, it will
	// be used instead of the password, even if a password is also
	// specified.
	KeyPath files.Path

	// WorkDir is the working directory to use on the machine for mounting
	// volumes to stage inputs and capture outputs.
	WorkDir files.Dir

	// Concurrency is the maximum number of tasks to run on this machine
	// at the same time.
	Concurrency int
}

type SshRunner struct {
	addr        string
	user        string
	pass        string
	keyPath     files.Path
	concurrency int

	client   *ssh.Client
	machines map[string]Machine
}

func NewSshRunner(opts ...option.Func[*SshRunner]) (*SshRunner, error) {
	r := &SshRunner{}
	for _, opt := range opts {
		err := opt(r)
		if err != nil {
			return nil, fmt.Errorf("NewSshRunner: failed to apply option: %w", err)
		}
	}

	return r, nil
}

func WithMachines(machines ...Machine) option.Func[*SshRunner] {
	return func(r *SshRunner) error {
		for _, m := range machines {
			if m.Addr == "" {
				return fmt.Errorf("Machine.Addr is a required field")
			}

			if _, present := r.machines[m.Addr]; present {
				slog.Warn("duplicate ssh machine address found", "machine", m)
			}

			r.machines[m.Addr] = m
		}

		return nil
	}
}

func WithPassword(user, pass string) option.Func[*SshRunner] {
	return func(r *SshRunner) error {
		r.user = user
		r.pass = pass
		return nil
	}
}

func WithKey(user string, keyPath files.Path) option.Func[*SshRunner] {
	return func(r *SshRunner) error {
		return nil
	}
}

func (r *SshRunner) CreateVolume(s files.Size) (Volume, error) {
	// TODO implement me
	panic("implement me")
}

func (r *SshRunner) DeleteVolume(v Volume) error {
	// TODO implement me
	panic("implement me")
}

func (r *SshRunner) AddFile(s files.Path, v Volume, name string) error {
	// TODO: This is all a mess, need to run Flowork remotely to pull, or do this if the path is local

	session, _ := r.session()
	defer func() { _ = session.Close() }()

	file, _ := os.Open("filetocopy")
	defer file.Close()
	stat, _ := file.Stat()

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		hostIn, _ := session.StdinPipe()
		defer hostIn.Close()
		fmt.Fprintf(hostIn, "C0664 %d %s\n", stat.Size(), "filecopyname")
		io.Copy(hostIn, file)
		fmt.Fprint(hostIn, "\x00")
		wg.Done()
	}()

	session.Run("/usr/bin/scp -t /remotedirectory/")
	wg.Wait()

	return nil
}

func (r *SshRunner) ExtractFile(s files.Path, v Volume, d files.Dir) error {
	// TODO implement me
	panic("implement me")
}

func (r *SshRunner) Run(t *Instance, v Volume) error {
	command, err := DockerRun(t, v, r.user)
	if err != nil {
		return fmt.Errorf("SshRunner.Run: failed to build docker command: %w", err)
	}

	return r.execute("Run", command, true)
}

func (r *SshRunner) Close() error {
	return r.client.Close()
}

func (r *SshRunner) session() (*ssh.Session, error) {
	// TODO: Need to handle heterogeneous auth for a cluster of machines
	// Alternately, we could have one runner per machine and then multiplex
	// jobs over many runners, that would also allow heterogeneous environments?
	// Also the multiplexing logic would only live in one place. That would
	// change the API a bit, though, since CreateVolume is designed for
	// multiplexing on a single "host".

	// TODO: Probably use a remote docker invocation to copy files

	if r.client != nil {
		return r.client.NewSession()
	}

	config := &ssh.ClientConfig{
		User: r.user,
		Auth: []ssh.AuthMethod{
			ssh.Password(r.pass),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         10 * time.Second,
	}

	client, err := ssh.Dial("tcp", r.addr, config)
	if err != nil {
		return nil, err
	}

	r.client = client

	return r.session()
}

func (r *SshRunner) execute(method string, command []string, capture bool) error {
	session, err := r.session()
	if err != nil {
		return fmt.Errorf("SshRunner.%s: failed to get ssh session: %w", method, err)
	}
	defer func() { _ = session.Close() }()

	if capture {
		stdout, stderr, err := outputWriters("")
		if err != nil {
			return fmt.Errorf("SshRunner.%s: failed to open output files: %w", method, err)
		}
		defer func() {
			_ = stdout.Close()
			_ = stderr.Close()
		}()

		session.Stdout = stdout
		session.Stderr = stderr
	}

	err = session.Run(strings.Join(command, " "))
	if err != nil {
		// TODO: Check for ExitMissing and ExitError
		return fmt.Errorf("SshRunner.%s: failed to run command: %w", method, err)
	}

	return nil
}
