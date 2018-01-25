package core

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"

	serial "github.com/ipfs/go-ipfs/repo/fsrepo/serialize"
)

/*
 * LocalNode.Dir needs to modified.
 */

type LocalNode struct {
	Dir    string
	PeerID string
}

func (n *LocalNode) GetPeerID() string {
	return n.PeerID
}

// get ipfs path
func (n *LocalNode) envForDaemon() ([]string, error) {
	envs := os.Environ()
	npath := "IPFS_PATH=" + n.Dir
	for i, e := range envs {
		p := strings.Split(e, "=")
		if p[0] == "IPFS_PATH" {
			envs[i] = npath
			return envs, nil
		}
	}

	return append(envs, npath), nil
}

func (n *LocalNode) getPID() (int, error) {
	b, err := ioutil.ReadFile(filepath.Join(n.Dir, "daemon.pid"))
	if err != nil {
		return -1, err
	}
	return strconv.Atoi(string(b))
}

func (n *LocalNode) isAlive() (bool, error) {
	pid, err := n.getPID()
	if os.IsNotExist(err) {
		return false, nil
	} else if err != nil {
		return false, err
	}

	proc, err := os.FindProcess(pid)
	if err != nil {
		return false, nil
	}

	// check the alive of the proc
	err = proc.Signal(syscall.Signal(0))
	if err != nil {
		return false, nil
	}
	return true, nil
}

func setupOpt(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
}

func tryAPICheck(n *LocalNode) error {
	resp, err := http.Get("http://127.0.0.1:5001/api/v0/id")
	if err != nil {
		return err
	}

	out := make(map[string]interface{})
	err = json.NewDecoder(resp.Body).Decode(&out)
	if err != nil {
		return fmt.Errorf("liveness check failed: %s", err)
	}

	id, ok := out["ID"]
	if !ok {
		return fmt.Errorf("liveness check failed: ID field not present in output")
	}

	idstr := id.(string)
	if idstr != n.GetPeerID() {
		return fmt.Errorf("liveness check failed: unexpected peer at endpoint")
	}

	return nil
}

func waitProcess(p *os.Process, ms int) error {
	for i := 0; i < (ms / 10); i++ {
		err := p.Signal(syscall.Signal(0))
		if err != nil {
			return nil
		}
		time.Sleep(time.Millisecond * 10)
	}
	return errors.New("time out")
}

func waitOnAPI(n *LocalNode) error {
	for i := 0; i < 50; i++ {
		err := tryAPICheck(n)
		if err == nil {
			return nil
		}
		time.Sleep(time.Millisecond * 200)
	}
	return fmt.Errorf("node %s failed to come online in given time period", n.GetPeerID())
}

// Init the ipfs path anf dir
func (n *LocalNode) Init() error {
	n.Dir = os.Getenv("HOME") + "/.ipfs/daemon/"
	std_print("daemon path: %s\n", n.Dir)

	err := os.MkdirAll(n.Dir, 0777)
	if err != nil {
		return err
	}

	// Number of bits to use in the generated RSA private key
	cmd := exec.Command("ipfs", "init", "-b=1024")

	// Get the ipfs path
	cmd.Env, err = n.envForDaemon()
	if err != nil {
		return err
	}

	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s: %s", err, string(out))
	}

	std_print("%s", string(out))
	return nil
}

func (n *LocalNode) Start(args []string) error {
	alive, err := n.isAlive()
	if err != nil {
		return err
	}

	if alive {
		return fmt.Errorf("node is already running")
	}

	dir := n.Dir
	dargs := append([]string{"daemon"}, args...)
	cmd := exec.Command("ipfs", dargs...)
	cmd.Dir = dir

	cmd.Env, err = n.envForDaemon()
	if err != nil {
		return err
	}

	setupOpt(cmd)

	stdout, err := os.Create(filepath.Join(dir, "daemon.stdout"))
	if err != nil {
		return err
	}

	stderr, err := os.Create(filepath.Join(dir, "daemon.stderr"))
	if err != nil {
		return err
	}

	cmd.Stdout = stdout
	cmd.Stderr = stderr

	err = cmd.Start()
	if err != nil {
		return err
	}
	pid := cmd.Process.Pid

	std_print("Start daemon %s, pid = %d\n", dir, pid)
	err = ioutil.WriteFile(filepath.Join(dir, "daemon.pid"), []byte(fmt.Sprint(pid)), 0666)
	if err != nil {
		return err
	}

	// Make sure node 0 is up before starting the rest so bootstrapping works properly
	//TODO
	cfg, err := serial.Load(filepath.Join(dir, "config"))
	if err != nil {
		return err
	}

	n.PeerID = cfg.Identity.PeerID

	err = waitOnAPI(n)
	if err != nil {
		return err
	}

	return nil
}

func (n *LocalNode) Kill() error {
	pid, err := n.getPID()
	if err != nil {
		return fmt.Errorf("error killing daemon %s: %s\n", n.Dir, err)
	}

	p, err := os.FindProcess(pid)
	if err != nil {
		return fmt.Errorf("error killing daemon %s: %s\n", n.Dir, err)
	}

	defer func() {
		err := os.Remove(filepath.Join(n.Dir, "daemon.pid"))
		if err != nil && !os.IsNotExist(err) {
			panic(fmt.Errorf("error removing pid file for daemon at %s: %s\n", n.Dir, err))
		}
	}()

	// kill
	err = p.Signal(syscall.SIGTERM)
	if err != nil {
		return fmt.Errorf("error killing daemon %s: %s\n", n.Dir, err)
	}

	err = waitProcess(p, 1000)
	if err == nil {
		return nil
	}

	err = p.Signal(syscall.SIGTERM)
	if err != nil {
		return fmt.Errorf("error killing daemon %s: %s\n", n.Dir, err)
	}

	err = waitProcess(p, 1000)
	if err == nil {
		return nil
	}

	err = p.Signal(syscall.SIGQUIT)
	if err != nil {
		return fmt.Errorf("error killing daemon %s: %s\n", n.Dir, err)
	}

	err = waitProcess(p, 5000)
	if err == nil {
		return nil
	}

	err = p.Signal(syscall.SIGKILL)
	if err != nil {
		return fmt.Errorf("error killing daemon %s: %s\n", n.Dir, err)
	}

	for {
		err := p.Signal(syscall.Signal(0))
		if err != nil {
			break
		}
		time.Sleep(time.Millisecond * 10)
	}

	return nil
}
