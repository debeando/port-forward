package main

import (
	"fmt"
	"io"
	"net"
	"path/filepath"

	"github.com/debeando/go-common/env"
	"github.com/debeando/go-common/log"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
)

var Debug string
var LocalPort string
var RemoteHost string
var RemotePort string
var SSHHost string
var SSHKey string
var SSHPort string
var SSHUser string

func init() {
	Debug = env.Get("DEBUG", "true")
	LocalPort = env.Get("LOCAL_PORT", "3306")
	RemoteHost = env.Get("REMOTE_HOST", "")
	RemotePort = env.Get("REMOTE_PORT", "3306")
	SSHHost = env.Get("SSH_HOST", "127.0.0.1")
	SSHKey = env.Get("SSH_KEY", "")
	SSHPort = env.Get("SSH_PORT", "22")
	SSHUser = env.Get("SSH_USER", "ec2-user")
}

func main() {
	knowHost := filepath.Join(env.Get("HOME", ""), ".ssh", "known_hosts")

	knownHostsCallback, err := knownhosts.New(knowHost)
	if err != nil {
		log.Error(err.Error())
	}

	signer, err := ssh.ParsePrivateKey([]byte(SSHKey))
	if err != nil {
		log.Error(err.Error())
	}

	sshClient := &ssh.ClientConfig{
		User: SSHUser,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback:   knownHostsCallback,
		HostKeyAlgorithms: []string{ssh.KeyAlgoED25519},
	}

	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%s", SSHHost, SSHPort), sshClient)
	if err != nil {
		log.Error(err.Error())
	}
	defer client.Close()

	listener, err := net.Listen("tcp", fmt.Sprintf("localhost:%s", LocalPort))
	if err != nil {
		log.Error(err.Error())
	}
	defer listener.Close()

	for {
		local, err := listener.Accept()
		if err != nil {
			log.Error(err.Error())
		}

		remote, err := client.Dial("tcp", fmt.Sprintf("%s:%s", RemoteHost, RemotePort))
		if err != nil {
			log.Error(err.Error())
		}

		log.Info("Tunnel established with.")
		runTunnel(local, remote)
	}
}

func runTunnel(local, remote net.Conn) {
	defer local.Close()
	defer remote.Close()
	done := make(chan struct{}, 2)

	go func() {
		io.Copy(local, remote)
		done <- struct{}{}
	}()

	go func() {
		io.Copy(remote, local)
		done <- struct{}{}
	}()

	<-done
}
