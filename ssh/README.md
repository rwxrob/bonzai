# Simplified pure Go secure shell (SSH) multi-host client

[![GoDoc](https://godoc.org/github.com/rwxrob/ssh?status.svg)](https://godoc.org/github.com/rwxrob/ssh)
[![License](https://img.shields.io/badge/license-Apache2-brightgreen.svg)](LICENSE)
example

## Testable `runasany` example

This directory also contains a testable example implementation of `ssh.RunOnAny` as the command [`runonany`](cmd/runonany/main.go) that can be used to validate this specific functionality interactively. This example requires the `podman` container engine and command to be installed on the local system.

Here are the steps to conduct an interactive test. In this way, one can monitor the status of a given server and bring it up or down on demand and note how this affects the `ssh.Controller` (inside of `runonany`).

### Build the SSH server/client image

```sh
build image
```

### Start up the three SSH server containers

Then start the three servers listening on ports `2221-2223`.

```sh
build start-servers
```

### Set the `RUNONANY_TARGET` environment variable

The containers all share the same underlying host (and IP address) but they don't know about it. We use the `RUNONANY_TARGET` environment variable to communicate this to the running containers.

```sh
export RUNONANY_TARGET=192.168.1.6
```

### Check the servers by running `ssh` from each client

Confirm server SSH connection is working with `ssh`. First you will need to note the host IP of the `podman` container engine and export it. (This can be obtained any number of ways.)

Now check the servers (or just do individual):

```sh
build check-servers
```

You should see the `ssh hostname` output of each command.

```
Warning: Permanently added '[192.168.1.6]:2221' (ED25519) to the list of known hosts.
ssh-server1
Warning: Permanently added '[192.168.1.6]:2222' (ED25519) to the list of known hosts.
ssh-server2
Warning: Permanently added '[192.168.1.6]:2223' (ED25519) to the list of known hosts.
ssh-server3
```

Note that the primary distinction is the port number. These servers all share the same `user` and credentials. They even share the same `authorized_hosts` key (which we ignore here deliberately for testing).

### Watch `runonany` output

The [`runonany`](cmd/runonany/main.go) binary is a simple program that encapsulates an `ssh.Controller` configured in the [`runonany.yaml`](testdata/runonany.yaml) YAML file.

```sh
build watch
```

This will update every two seconds.

### Interactively stop and start SSH server containers

The containers running as servers (see [`entrypoint`](entrypoint)) can be stopped and started while monitoring the live status using commands similar to the following:

```sh
build stop-server 2
build start-server 2
build stop-servers
build start-servers
```

It is useful to do these commands from one TMUX pane while running `build watch-client-runonany` from another to see the change in `ssh.Controller.Clients` status.

Here are some things to validate:

* Random servers selected are between 1-3.
* Stop one server and note random selected no longer include.
* Stop two servers and not only single server selected.
* Stop all servers and not only the `ERR` section returned.
* Start one server after stopping and note recovery.
* Start two servers and note recovery.
* Start all servers and note recovery.

### Do other things from client container

If you prefer to play around on the client container itself just override the `--entrypoint`.

```sh
podman run -it --rm  --entrypoint bash runonany:latest
```

This will give you a root login.
