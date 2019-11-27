# CubeBeat

Custom Beat of the Elastic Stack to interact with the Polycube-based eBPF cubes.

---

Ensure that this folder is at the following location:
`${GOPATH}/src/gitlab.com/astrid-repositories/wp2/cubebeat`

## Getting Started

### Requirements

* [Golang](https://golang.org/dl/) 1.7


Setting up the dev. environment.

```console
pip install cookiecutter
go get github.com/elastic/beats
cd $GOPATH/src/github.com/elastic/beats
git checkout 6.8
```

### Init Project
To get running with Cubebeat and also install the dependencies, run the following command:

```console
make setup
```

It will create a clean git history for each major step. Note that you can always rewrite the history if you wish before pushing your changes.

To push CubeBeat in the git repository, run the following commands:

```
git remote set-url origin gitlab.com/astrid-repositories/wp2/cube-beat
git push origin master
```

For further development, check out the [beat developer guide](https://www.elastic.co/guide/en/beats/libbeat/current/new-beat.html).

### Build

To build the binary for CubeBeat run the command below. This will generate a binary
in the same directory with the name cube-beat.

```
make
```


### Run

To run CubeBeat with debugging output enabled, run:

```
./cube-beat -c cube-beat.yml -e -d "*"
```


### Test

To test CubeBeat, run the following command:

```
make testsuite
```

alternatively:
```
make unit-tests
make system-tests
make integration-tests
make coverage-report
```

The test coverage is reported in the folder `./build/coverage/`

### Update

Each beat has a template for the mapping in elasticsearch and a documentation for the fields
which is automatically generated based on `fields.yml` by running the following command.

```
make update
```


### Cleanup

To clean  CubeBeat source code, run the following command:

```
make fmt
```

To clean up the build directory and generated artifacts, run:

```
make clean
```


### Clone

To clone CubeBeat from the git repository, run the following commands:

```
mkdir -p ${GOPATH}/src/github.com/gitlab.com/astrid-repositories/wp2/cube-beat
git clone https://gitlab.com/astrid-repositories/wp2/cube-beat ${GOPATH}/src/github.com/gitlab.com/astrid-repositories/wp2/cube-beat
```

For further development, check out the [beat developer guide](https://www.elastic.co/guide/en/beats/libbeat/current/new-beat.html).


## Packaging

The beat frameworks provides tools to crosscompile and package your beat for different platforms. This requires [docker](https://www.docker.com/) and vendoring as described above. To build packages of your beat, run the following command:

```
make release
```

This will fetch and create all images required for the build process. The whole process to finish can take several minutes.
