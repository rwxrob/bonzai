package ssh_test

import (
	"fmt"

	"github.com/rwxrob/ssh"
	"gopkg.in/yaml.v3"
)

func ExampleUser_from_YAMLJSON() {

	var err error

	yml := []byte(`
name: user1
key: |-
  -----BEGIN OPENSSH PRIVATE KEY-----
  b3BlbnNzaC1rZXktdjEAAAAABG5vbmUAAAAEbm9uZQAAAAAAAAABAAAAMwAAAAtzc2gtZW
  QyNTUxOQAAACB0jdh2hglPJchsrVgnJjTb9bVHIjugS5wlJipnIJiO8gAAAJAZeyGhGXsh
  oQAAAAtzc2gtZWQyNTUxOQAAACB0jdh2hglPJchsrVgnJjTb9bVHIjugS5wlJipnIJiO8g
  AAAEDdV9IJ3LNTiK7D0MFz7IR1Cz/VdqqH6SgOtiDz8/5073SN2HaGCU8lyGytWCcmNNv1
  tUciO6BLnCUmKmcgmI7yAAAACXJ3eHJvYkB0dgECAwQ=
  -----END OPENSSH PRIVATE KEY-----
`)

	// YAML
	user1 := new(ssh.User)
	err = yaml.Unmarshal(yml, user1)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(user1.Name)
	fmt.Println(user1.Key)

	jsn := []byte(`{"name": "user2","key":"-----BEGIN OPENSSH PRIVATE KEY-----\nb3BlbnNzaC1rZXktdjEAAAAABG5vbmUAAAAEbm9uZQAAAAAAAAABAAAAMwAAAAtzc2gtZW\nQyNTUxOQAAACB0jdh2hglPJchsrVgnJjTb9bVHIjugS5wlJipnIJiO8gAAAJAZeyGhGXsh\noQAAAAtzc2gtZWQyNTUxOQAAACB0jdh2hglPJchsrVgnJjTb9bVHIjugS5wlJipnIJiO8g\nAAAEDdV9IJ3LNTiK7D0MFz7IR1Cz/VdqqH6SgOtiDz8/5073SN2HaGCU8lyGytWCcmNNv1\ntUciO6BLnCUmKmcgmI7yAAAACXJ3eHJvYkB0dgECAwQ=\n-----END OPENSSH PRIVATE KEY-----"}`)

	// JSON (which *is* YAML)
	user2 := new(ssh.User)
	err = yaml.Unmarshal(jsn, user2)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(user2.Name)
	fmt.Println(user1.Key)

	// Output:
	// user1
	// -----BEGIN OPENSSH PRIVATE KEY-----
	// b3BlbnNzaC1rZXktdjEAAAAABG5vbmUAAAAEbm9uZQAAAAAAAAABAAAAMwAAAAtzc2gtZW
	// QyNTUxOQAAACB0jdh2hglPJchsrVgnJjTb9bVHIjugS5wlJipnIJiO8gAAAJAZeyGhGXsh
	// oQAAAAtzc2gtZWQyNTUxOQAAACB0jdh2hglPJchsrVgnJjTb9bVHIjugS5wlJipnIJiO8g
	// AAAEDdV9IJ3LNTiK7D0MFz7IR1Cz/VdqqH6SgOtiDz8/5073SN2HaGCU8lyGytWCcmNNv1
	// tUciO6BLnCUmKmcgmI7yAAAACXJ3eHJvYkB0dgECAwQ=
	// -----END OPENSSH PRIVATE KEY-----
	// user2
	// -----BEGIN OPENSSH PRIVATE KEY-----
	// b3BlbnNzaC1rZXktdjEAAAAABG5vbmUAAAAEbm9uZQAAAAAAAAABAAAAMwAAAAtzc2gtZW
	// QyNTUxOQAAACB0jdh2hglPJchsrVgnJjTb9bVHIjugS5wlJipnIJiO8gAAAJAZeyGhGXsh
	// oQAAAAtzc2gtZWQyNTUxOQAAACB0jdh2hglPJchsrVgnJjTb9bVHIjugS5wlJipnIJiO8g
	// AAAEDdV9IJ3LNTiK7D0MFz7IR1Cz/VdqqH6SgOtiDz8/5073SN2HaGCU8lyGytWCcmNNv1
	// tUciO6BLnCUmKmcgmI7yAAAACXJ3eHJvYkB0dgECAwQ=
	// -----END OPENSSH PRIVATE KEY-----

}

func ExampleHost_from_YAMLJSON() {

	var err error

	yml := []byte(`
addr: host1
auth: "randomoption ecdsa-sha2-nistp256 AAAAE2VjZHNhLXNoYTItbmlzdHAyNTYAAAAIbmlzdHAyNTYAAABBBI/WBBaaNFajVHCL0+rQqWP3zhpyXo357iPUvl0GGHWrY6t42WTNJ+bk8shRq7eq8KwefZeL4YvsnekcZb8Uq+8="
`)

	// YAML
	host1 := new(ssh.Host)
	err = yaml.Unmarshal(yml, host1)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(host1.Addr)
	fmt.Println(host1.Auth)

	jsn := []byte(`{"addr":"host2","auth":"randomoption ecdsa-sha2-nistp256 AAAAE2VjZHNhLXNoYTItbmlzdHAyNTYAAAAIbmlzdHAyNTYAAABBBI/WBBaaNFajVHCL0+rQqWP3zhpyXo357iPUvl0GGHWrY6t42WTNJ+bk8shRq7eq8KwefZeL4YvsnekcZb8Uq+8="}`)

	// JSON (which *is* YAML)
	host2 := new(ssh.Host)
	err = yaml.Unmarshal(jsn, host2)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(host2.Addr)
	fmt.Println(host2.Auth)

	// Output:
	// host1
	// randomoption ecdsa-sha2-nistp256 AAAAE2VjZHNhLXNoYTItbmlzdHAyNTYAAAAIbmlzdHAyNTYAAABBBI/WBBaaNFajVHCL0+rQqWP3zhpyXo357iPUvl0GGHWrY6t42WTNJ+bk8shRq7eq8KwefZeL4YvsnekcZb8Uq+8=
	// host2
	// randomoption ecdsa-sha2-nistp256 AAAAE2VjZHNhLXNoYTItbmlzdHAyNTYAAAAIbmlzdHAyNTYAAABBBI/WBBaaNFajVHCL0+rQqWP3zhpyXo357iPUvl0GGHWrY6t42WTNJ+bk8shRq7eq8KwefZeL4YvsnekcZb8Uq+8=

}

func ExampleClient_from_YAMLJSON() {
	var err error

	// YAML references are supported.

	yml := []byte(`
immauser: &auser
  name: someuser
  key: somethingsecret

host:
  addr: host1

user: *auser

# timeout and port set to defaults when first used if zero

`)

	client := new(ssh.Client)
	err = yaml.Unmarshal(yml, client)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(client.User.Name)
	fmt.Println(client.Host.Addr)
	fmt.Println(client.Port)   // 0
	fmt.Println(client.Addr()) // port added ...
	fmt.Println(client.Port)   // but not set

	// Output:
	// someuser
	// host1
	// 0
	// host1:22
	// 0

}

func ExampleController_from_YAMLJSON() {

	// YAML references are supported.

	yml := []byte(`
users:
  user: &user
    name: user
    key: |
      -----BEGIN OPENSSH PRIVATE KEY-----
      b3BlbnNzaC1rZXktdjEAAAAABG5vbmUAAAAEbm9uZQAAAAAAAAABAAAAMwAAAAtzc2gtZW
      QyNTUxOQAAACB0jdh2hglPJchsrVgnJjTb9bVHIjugS5wlJipnIJiO8gAAAJAZeyGhGXsh
      oQAAAAtzc2gtZWQyNTUxOQAAACB0jdh2hglPJchsrVgnJjTb9bVHIjugS5wlJipnIJiO8g
      AAAEDdV9IJ3LNTiK7D0MFz7IR1Cz/VdqqH6SgOtiDz8/5073SN2HaGCU8lyGytWCcmNNv1
      tUciO6BLnCUmKmcgmI7yAAAACXJ3eHJvYkB0dgECAwQ=
      -----END OPENSSH PRIVATE KEY-----

hosts:
  localhost: &localhost
    addr: localhost
    auth: |
      randomoption ecdsa-sha2-nistp256 AAAAE2VjZHNhLXNoYTItbmlzdHAyNTYAAAAIbmlzdHAyNTYAAABBBI/WBBaaNFajVHCL0+rQqWP3zhpyXo357iPUvl0GGHWrY6t42WTNJ+bk8shRq7eq8KwefZeL4YvsnekcZb8Uq+8=

clients:
  - host: *localhost
    user: *user
    port: 2221
    timeout: 5m
  - host: *localhost
    user: *user
    port: 2222
    timeout: 5m
  - host: *localhost
    user: *user
    port: 2223
    timeout: 5m
`)

	ctl := new(ssh.Controller)
	err := yaml.Unmarshal(yml, ctl)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(ctl.Clients[0].User.Name)
	fmt.Println(ctl.Clients[1].User.Name)
	fmt.Println(ctl.Clients[2].User.Name)
	fmt.Println(ctl.Clients[0].Port)
	fmt.Println(ctl.Clients[1].Port)
	fmt.Println(ctl.Clients[2].Port)
	fmt.Println(ctl.Clients[0].Timeout)

	// Output:
	// user
	// user
	// user
	// 2221
	// 2222
	// 2223
	// 5m0s

}

func ExampleClient_Run() {

	// Change YAML data and remove one / from Output
	// to test locally. (But never commit actual keys.)

	yml := []byte(`
user:
  name: user
  key: |
    -----BEGIN OPENSSH PRIVATE KEY-----
    b3BlbnNzaC1rZXktdjEAAAAABG5vbmUAAAAEbm9uZQAAAAAAAAABAAAAMwAAAAtzc2gtZW
    QyNTUxOQAAACB0/Xdc30JNJx+H1bvs9oZ7POIBi/YyZ4UBfQ/oEyffOQAAAJDswLYs7MC2
    LAAAAAtzc2gtZWQyNTUxOQAAACB0/Xdc30JNJx+H1bvs9oZ7POIBi/YyZ4UBfQ/oEyffOQ
    AAAEDWFaCmeeFjBMAzJvtf6z24ai1dHf2FSUmuHrONv/5K6XT9d1zfQk0nH4fVu+z2hns8
    4gGL9jJnhQF9D+gTJ985AAAACXJ3eHJvYkB0dgECAwQ=
    -----END OPENSSH PRIVATE KEY-----
host:
    addr: localhost
    auth: |
      |1|r3meMBTG9TZiPoVHg1n+o1N1xJk=|9I891Skl7BcqG/vaT6wXxt6bZUk= ecdsa-sha2-nistp256 AAAAE2VjZHNhLXNoYTItbmlzdHAyNTYAAAAIbmlzdHAyNTYAAABBBDqzw7+sN4aVOQqgTA5tC9pN/+M0KOcib3lRAGQ+MSKk/4MbdJY2REavrwRetreaIZTkZx4ykTAJ3CeCK45IzsY=
`)

	client := new(ssh.Client)
	err := yaml.Unmarshal(yml, client)
	if err != nil {
		fmt.Println(err)
	}

	// basic command (no stdin)
	stdout, stderr, err := client.Run(`printf hello`, nil)
	fmt.Println(`-----`)
	fmt.Println(stdout)
	fmt.Println(stderr)

	// with stdin
	stdout, stderr, _ = client.Run(`cat`, []byte(`i'm a cat`))
	fmt.Println(`-----`)
	fmt.Println(stdout)
	fmt.Println(stderr)

	// with stderr
	stdout, stderr, _ = client.Run(`ls notafile`, nil)
	fmt.Println(`-----`)
	fmt.Println(stdout)
	fmt.Println(stderr)

	/// Output:
	// -----
	// hello
	//
	// -----
	// i'm a cat
	//
	// -----
	//
	// ls: cannot access 'notafile': No such file or directory

}

func ExampleController_RandomClient_none() {

	// idle, unconnected
	c1 := new(ssh.Client)
	c2 := new(ssh.Client)
	c3 := new(ssh.Client)

	ctl := new(ssh.Controller).Init(c1, c2, c3)
	client := ctl.RandomClient()
	fmt.Println(client)

	// Output:
	// <nil>
}

func ExampleController_RandomClient_single_Good() {

	// YAML references are supported.

	yml := []byte(`
users:
  user: &auser
    name: user
    key: |
      -----BEGIN OPENSSH PRIVATE KEY-----
      b3BlbnNzaC1rZXktdjEAAAAABG5vbmUAAAAEbm9uZQAAAAAAAAABAAAAMwAAAAtzc2gtZW
      QyNTUxOQAAACB0/Xdc30JNJx+H1bvs9oZ7POIBi/YyZ4UBfQ/oEyffOQAAAJDswLYs7MC2
      LAAAAAtzc2gtZWQyNTUxOQAAACB0/Xdc30JNJx+H1bvs9oZ7POIBi/YyZ4UBfQ/oEyffOQ
      AAAEDWFaCmeeFjBMAzJvtf6z24ai1dHf2FSUmuHrONv/5K6XT9d1zfQk0nH4fVu+z2hns8
      4gGL9jJnhQF9D+gTJ985AAAACXJ3eHJvYkB0dgECAwQ=
      -----END OPENSSH PRIVATE KEY-----
hosts:
  localhost: &localhost
    addr: localhost
    auth: |
      |1|r3meMBTG9TZiPoVHg1n+o1N1xJk=|9I891Skl7BcqG/vaT6wXxt6bZUk= ecdsa-sha2-nistp256 AAAAE2VjZHNhLXNoYTItbmlzdHAyNTYAAAAIbmlzdHAyNTYAAABBBDqzw7+sN4aVOQqgTA5tC9pN/+M0KOcib3lRAGQ+MSKk/4MbdJY2REavrwRetreaIZTkZx4ykTAJ3CeCK45IzsY=

clients:
  - host: *localhost
    user: *auser
    comment: number one
  - host: *localhost
    user: *auser
    comment: number two
  - host: *localhost
    user: *auser
    comment: number three
`)

	ctl := new(ssh.Controller)
	err := yaml.Unmarshal(yml, ctl)
	if err != nil {
		fmt.Println(err)
	}

	ctl.Connect()
	client := ctl.RandomClient()
	if client != nil {
		fmt.Println(client.Comment)
	}

	/// Output:
	// number one

}

func ExampleController_RunOnAny() {

	yml := []byte(`
pem: &ukey |
  -----BEGIN OPENSSH PRIVATE KEY-----
  b3BlbnNzaC1rZXktdjEAAAAABG5vbmUAAAAEbm9uZQAAAAAAAAABAAAAMwAAAAtzc2gtZW
  QyNTUxOQAAACB0/Xdc30JNJx+H1bvs9oZ7POIBi/YyZ4UBfQ/oEyffOQAAAJDswLYs7MC2
  LAAAAAtzc2gtZWQyNTUxOQAAACB0/Xdc30JNJx+H1bvs9oZ7POIBi/YyZ4UBfQ/oEyffOQ
  AAAEDWFaCmeeFjBMAzJvtf6z24ai1dHf2FSUmuHrONv/5K6XT9d1zfQk0nH4fVu+z2hns8
  4gGL9jJnhQF9D+gTJ985AAAACXJ3eHJvYkB0dgECAwQ=
  -----END OPENSSH PRIVATE KEY-----

user1: &user1
  name: user1
  key: *ukey
user2: &user2
  name: user2
  key: *ukey
user3: &user3
  name: user3
  key: *ukey

hosts:
  localhost: &localhost
    addr: localhost
    auth: |
      |1|r3meMBTG9TZiPoVHg1n+o1N1xJk=|9I891Skl7BcqG/vaT6wXxt6bZUk= ecdsa-sha2-nistp256 AAAAE2VjZHNhLXNoYTItbmlzdHAyNTYAAAAIbmlzdHAyNTYAAABBBDqzw7+sN4aVOQqgTA5tC9pN/+M0KOcib3lRAGQ+MSKk/4MbdJY2REavrwRetreaIZTkZx4ykTAJ3CeCK45IzsY=

clients:
  - host: *localhost
    user: *user1
  - host: *localhost
    user: *user2
  - host: *localhost
    user: *user3
`)

	ctl := new(ssh.Controller)
	err := yaml.Unmarshal(yml, ctl)
	if err != nil {
		fmt.Println(err)
	}

	ctl.Connect()
	if err != nil {
		fmt.Println(err)
	}

	stdin, _, _ := ctl.RunOnAny(`whoami`, nil)
	fmt.Println(stdin)

	/// Output:
	// user1

}

func ExampleController_RunOnAny_single_Good() {

	yml := []byte(`
pem: &ukey |
  -----BEGIN OPENSSH PRIVATE KEY-----
  b3BlbnNzaC1rZXktdjEAAAAABG5vbmUAAAAEbm9uZQAAAAAAAAABAAAAMwAAAAtzc2gtZW
  QyNTUxOQAAACB0/Xdc30JNJx+H1bvs9oZ7POIBi/YyZ4UBfQ/oEyffOQAAAJDswLYs7MC2
  LAAAAAtzc2gtZWQyNTUxOQAAACB0/Xdc30JNJx+H1bvs9oZ7POIBi/YyZ4UBfQ/oEyffOQ
  AAAEDWFaCmeeFjBMAzJvtf6z24ai1dHf2FSUmuHrONv/5K6XT9d1zfQk0nH4fVu+z2hns8
  4gGL9jJnhQF9D+gTJ985AAAACXJ3eHJvYkB0dgECAwQ=
  -----END OPENSSH PRIVATE KEY-----

user1: &user1
  name: user1
  key: *ukey
user2: &user2
  name: user2
  key: *ukey
user3: &user3
  name: user3
  key: *ukey

hosts:
  localhost: &localhost
    addr: localhost
    auth: |
      |1|r3meMBTG9TZiPoVHg1n+o1N1xJk=|9I891Skl7BcqG/vaT6wXxt6bZUk= ecdsa-sha2-nistp256 AAAAE2VjZHNhLXNoYTItbmlzdHAyNTYAAAAIbmlzdHAyNTYAAABBBDqzw7+sN4aVOQqgTA5tC9pN/+M0KOcib3lRAGQ+MSKk/4MbdJY2REavrwRetreaIZTkZx4ykTAJ3CeCK45IzsY=

clients:
  - host: *localhost
    user: *user1
  - host: *localhost
    user: *user2
  - host: *localhost
    user: *user3
`)

	ctl := new(ssh.Controller)
	err := yaml.Unmarshal(yml, ctl)
	if err != nil {
		fmt.Println(err)
	}

	// only connect one of the clients, others are bogus
	err = ctl.Clients[1].Connect()
	if err != nil {
		fmt.Println(err)
	}

	stdin, _, _ := ctl.RunOnAny(`whoami`, nil)
	fmt.Println(stdin)

	// Output:
	// user2

}

func ExampleController_RunOnAny_no_Clients() {

	ctl := new(ssh.Controller)
	_, _, err := ctl.RunOnAny(`echo hello`, nil)
	fmt.Println(err)

	// Output:
	// all SSH client targets are unavailable

}
