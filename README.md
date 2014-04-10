packer-builder-lxc
==================

Lxc builder for packer, with working provisioning. Currently tested on debian wheezy, although ubuntu should work too.
Lxc 1.0.0 is needed on debian, since lxc-attach is broken in previous versions. A backport deb can be created following [this guide](https://wiki.debian.org/SimpleBackportCreation). On debian you have to edit /etc/lxc/default.conf according to you network settings, if you want to do anything during provisioning that needs network access.

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

Builder config
==============

Example:
```
    {
      "type": "lxc",
      "config_file": "lxc/config",
      "template_name": "debian",
      "template_parameters": ["--arch", "amd64", "--release", "wheezy"],
      "template_environment_vars": [
        "MIRROR=http://http.debian.net/debian/"
      ]
    }
```

The config file is an lxc config file which will be bundled with the machine. You can create your own or just grab the debian or ubuntu one from [this repo](https://github.com/fgrehm/vagrant-lxc-base-boxes/tree/master/conf).

Example packer template
=======================

TODO

Vagrant publisher
=================

The basebox format will be finalized in [vagrant-lxc](https://github.com/fgrehm/vagrant-lxc) 1.0.0,
then we'll try to get a patch into packer to support this builder in the vagrant publisher.
