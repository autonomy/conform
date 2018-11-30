<p align="center">
  <h1 align="center">Conform</h1>
  <p align="center">DRY, hygienic, fast builds.</p>
  <p align="center">
    <a href="https://gitter.im/autonomy/conform"><img alt="Gitter" src="https://img.shields.io/gitter/room/autonomy/conform.svg?style=flat-square"></a>
    <a href="https://godoc.org/github.com/autonomy/conform"><img alt="GoDoc" src="http://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square"></a>
    <a href="https://travis-ci.org/autonomy/conform"><img alt="Travis" src="https://img.shields.io/travis/autonomy/conform.svg?style=flat-square"></a>
    <a href="https://codecov.io/gh/autonomy/conform"><img alt="Codecov" src="https://img.shields.io/codecov/c/github/autonomy/conform.svg?style=flat-square"></a>
    <a href="https://goreportcard.com/report/github.com/autonomy/conform"><img alt="Go Report Card" src="https://goreportcard.com/badge/github.com/autonomy/conform?style=flat-square"></a>
    <a href="https://github.com/autonomy/conform/releases/latest"><img alt="Release" src="https://img.shields.io/github/release/autonomy/conform.svg?style=flat-square"></a>
    <a href="https://github.com/autonomy/conform/releases/latest"><img alt="GitHub (pre-)release" src="https://img.shields.io/github/release/autonomy/conform/all.svg?style=flat-square"></a>
  </p>
</p>

---

**Conform** is a tool for building enforcing policies on your build pipelines.

Some of the policies included are:

- **Convetion Commits**: Enforce [conventional commits](https://www.conventionalcommits.org) for all commit messages.

## Getting Started

Create a file named `.conform.yaml` with the following contents:

```yaml:
  - type: conventionalCommit
    spec:
      types:
        - "type"
      scopes:
        - "scope"
```

In the same directory, run:

```bash
conform enforce
```

### License
[![license](https://img.shields.io/github/license/autonomy/conform.svg?style=flat-square)](https://github.com/autonomy/conform/blob/master/LICENSE)
