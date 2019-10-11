# Contributing Guide

This is part of the [Porter][porter] project. If you are a new contributor,
check out our [New Contributor Guide][new-contrib]. The Porter [Contributing
Guide][contrib] also has lots of information about how to interact with the
project.

[porter]: https://github.com/deislabs/porter
[new-contrib]: https://porter.sh/contribute
[contrib]: https://github.com/deislabs/porter/blob/master/CONTRIBUTING.md

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
  Normally this is run automatically but if you run into issues with packr and
  dep, run this command.
* `dep-ensure` runs dep ensure for you while taking care of packr properly. Use
  this if your PRs are often failing on `verify-vendor` because of packr. This
  can be avoided entirely if you use `make build`.
* `verify-vendor` cleans up packr generated files and verifies that dep's Gopkg.lock 
   and vendor/ are up-to-date. Use this makefile target instead of running 
   dep check manually.
