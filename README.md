# CubeBeat

Custom Beat of the Elastic Stack to interact with the Polycube-based eBPF cubes.

---

Ensure that this folder is at the following location:
`${GOPATH}/src/gitlab.com/astrid-repositories/wp2/cubebeat`

## Getting Started

### Requirements

* [Golang](https://golang.org/dl/) 1.7
* Follows the instructions at [Setting Up Your Dev Environment](https://www.elastic.co/guide/en/beats/devguide/current/beats-contributing.html#setting-up-dev-environment).

### Download the code

```console
mkdir -p ${GOPATH}/src/gitlab.com/astrid-repositories/wp2/
cd ${GOPATH}/src/gitlab.com/astrid-repositories/wp2/
git clone https://gitlab.com/astrid-repositories/wp2/cubebeat.git
```

### Build

To build the binary for ```Cubebeat``` run the command below. This will generate a binary in the same directory with the name cubebeat.

```console
mage build
```

### Run

To run ```CubeBeat``` with debugging output enabled, run:

```console
./cubebeat -c cubebeat.yml -e -d "*"
```

To run ```CubeBeat``` without debugging output enabled, run:

```console
./cubebeat -c cubebeat.yml -e
```

### Configuration

```Cubebeat``` reads the configuration file (default: ```cubebeat.yml```) that is passed as argument.

This file accepts the common beat configurations as described at [Config file format](https://www.elastic.co/guide/en/beats/libbeat/current/config-file-format.html).

In addition, it accept specific configurations as shown in the next example:

```yaml
cubebeat:
  config.inputs:
    path: config/*.yml
    reload:
      enabled: true
      period: 10s
```

#### Load external configuration files

```Cubebeat``` can load external configuration files for inputs and modules, allowing you to separate your configuration into multiple smaller configuration files.

> On systems with ```POSIX``` file permissions, all configuration files are subject to ownership and file permission checks.<br/> For more information, see [Config File Ownership and Permissions](https://www.elastic.co/guide/en/beats/libbeat/7.4/config-file-permissions.html) in the _Beats Platform Reference_.

You specify the ```path``` option in the ```cubebeat.config.inputs``` section of the ```cubebeat.yml```. For example:

```yaml
cubebeat:
  config.inputs:
    path: config.d/*.yml
```

Each file found by the ```path``` Glob must contain a list of one or more input definitions.

> The first line of each external configuration file must be an input definition that starts with ```- name```.

For example:

```yaml
- name: synflood
  enabled: true
  period: 10s
  polycube.api-url: "http://localhost:9000/polycube/v1/synflood/sf/stats/"

- name: packetcapture
  enabled: true
  period: 5s
  polycube.api-url: "http://localhost:9000/polycube/v1/packetcapture/pc"
```

> It is critical that two running inputs DO NOT have same ```name```. If more than one input the same ```name```, only the first one is accepted; while the other ones are discarded.

When the option ```enabled``` is ```true```, the specific cube input periodically interact with the specific Polycube Cube
each time interval defined in ```period``` making an HTTP request to the URL defined in ```polycube.api-url```.

> If the cube is not reachable or there are some error when retrieves the data, ```cubebeat``` will continue to work, trying a new connection after a period of time defined in ```period```.

Each ```period``` of time, the specific cube input send a new ```Elastic``` event to the output as defined in the config file ```cubebeat.yml```

#### Live reloading

You can configure ```cubebeat``` to dynamically reload external configuration files when there are changes.
This feature is available for input configurations that are loaded as external configuration files.
You cannot use this feature to reload the main ```cubebeat.yml``` configuration file.

To configure this feature, you specify a ```path``` (Glob) to watch for configuration changes.
When the files found by the Glob change, new inputs are started and stopped according to changes in the configuration files.

This feature is especially useful in container environments where one container is used to tail logs for services running in other containers on the same host.

To enable dynamic config reloading, you specify the ```path``` and ```reload``` options under ```cubebeat.config.inputs``` section. For example:

```yaml
cubebeat:
  config.inputs:
    path: config/*.yml
    reload:
      enabled: true
      period: 10s
```

Option               | Description
-------------------: | :----------
```path```           | A Glob that defines the files to check for changes.
```reload.enabled``` | When set to true, enables dynamic config reload.
```reload.period```  | Specifies how often the files are checked for changes.<br/>Do not set the ```period``` to less than ```1s``` because the modification time of files is often stored in seconds.<br/>Setting the ```period``` to less than ```1s``` will result in **unnecessary overhead**.

> On systems with ```POSIX``` file permissions, all configuration files are subject to ownership and file permission checks.<br/> For more information, see [Config File Ownership and Permissions](https://www.elastic.co/guide/en/beats/libbeat/7.4/config-file-permissions.html) in the _Beats Platform Reference_.