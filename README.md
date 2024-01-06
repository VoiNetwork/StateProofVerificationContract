# StateProofVerificationContract

### Table of contents

* [1. Development](#1-development)
  * [1.1. Requirements](#11-requirements)
  * [1.2. Setup](#12-setup)
  * [1.3. Run Tests](#13-run-tests)

## 1. Development

### 1.1. Requirements

* [Docker][docker]
* [Docker Compose v2.5.0+][docker-compose]
* [Go 1.21+][golang]
* [Python v3.10+][python]
* [pip][pip] (optional, but highly recommended)
* [PyTeal][pyteal]

<sup>[Back to top ^][table-of-contents]</sup>

### 1.2. Setup

1. THe easiest way to install PyTeal is using `pip`
```shell
pip3 pyteal
```

2. Install Go dependencies:
```shell
go mod tidy
```

3. Install Algorand private network (optional):
```shell
./bin/install_algorand.sh
```

<sup>[Back to top ^][table-of-contents]</sup>

### 1.3. Run Tests

1. Run the tests:
```shell
./bin/test.sh
```

> ⚠️ **NOTE:** running the tests will also install (if you haven't installed one previously) and start a private Algorand private network.

<sup>[Back to top ^][table-of-contents]</sup>

[docker]: https://docs.docker.com/get-docker/
[docker-compose]: https://docs.docker.com/compose/install/
[golang]: https://golang.org/dl/
[license]: ./LICENSE
[pip]: https://pip.pypa.io/en/stable/installation/
[pyteal]: https://pyteal.readthedocs.io/en/latest/installation.html
[python]: https://www.python.org/downloads/
[table-of-contents]: #table-of-contents
