packer-builder-lxc
==========

Lxc builder for packer, with working provisioning.

Building
========

Install [gox](https://github.com/mitchellh/gox)

```
gox -os=linux -arch=amd64 -output=pkg/{{.OS}}_{{.Arch}}/packer-builder-lxc
```

Installation
============

Add the executable to your packer config:
```
{
  "builders": {
    "lxc": "/vagrant/packer/packer-builder-lxc"
  }
}
```

Example builder config
======================
```
    {
      "type": "lxc",
      "template_name": "debian",
      "template_parameters": ["--arch", "amd64", "--release", "wheezy"],
      "template_environment_vars": [
        "MIRROR=http://http.debian.net/debian/"
      ]
    }
```

Example packer template
=======================

TODO

Vagrant publisher
=================

The basebox format will be finalized in [vagrant-lxc](https://github.com/fgrehm/vagrant-lxc) 1.0.0,
then we'll try to get a patch into packer to support this builder in the vagrant publisher.