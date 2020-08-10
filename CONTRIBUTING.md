# Contributing Guide

This is part of the [Porter][porter] project. If you are a new contributor,
check out our [New Contributor Guide][new-contrib]. The Porter [Contributing
Guide][contrib] also has lots of information about how to interact with the
project.

[porter]: https://github.com/deislabs/porter
[new-contrib]: https://porter.sh/contribute
[contrib]: https://github.com/deislabs/porter/blob/main/CONTRIBUTING.md

---

* [Initial setup](#initial-setup)
* [Makefile explained](#makefile-explained)

---

# Initial setup

You need to have [porter installed](https://porter.sh/install) first. Then run
`make build install`. This will build and install the mixin into your porter
home directory.

## Makefile explained

Here are the most common Makefile tasks

* `build` builds both the runtime and client.
* `install` installs the mixin into **~/.porter/mixins**.
* `test-unit` runs the unit tests.
* `clean-packr` removes extra packr files that were a side-effect of the build.
  Normally this is run automatically but if you run into issues with packr,
  run this command.

