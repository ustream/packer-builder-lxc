packer-builder-lxc
==================

Lxc builder for packer, with working provisioning. Currently tested on debian wheezy, although ubuntu should work too.
. On debian you have to edit /etc/lxc/default.conf according to you network settings, if you want to do anything during provisioning that needs network access.

# OS Support

As building platform only linux is supported, at this point Ubuntu and Debian is tested. Please consult your lxc templates so you can pass the right parameters. All OSes should be buildable if they are supported by your lxc templates.

## Debian

Wheezy is supported with lxc 1.0.0. An older version is shipped, with a broken lxc-attach command, so it's not supported. A backport deb can be created easily following [this guide](https://wiki.debian.org/SimpleBackportCreation). Debian does not ship network support by default for lxc, you have to configure it manually. See the [vagrant-lxc wiki](https://github.com/fgrehm/vagrant-lxc/wiki/Usage-on-debian-hosts#network-configuration) for a detailed howto.

You'll need to edit /etc/lxc/default.conf to bring the network up on new containers:

```
lxc.network.type=veth
lxc.network.link=lxcbr0
lxc.network.flags=up
```

If your containers do not get an ip address from dhcp you need to turn off checksum offloading on the bridge:

```/sbin/ethtool -K lxcbr0 tx off```

If you're done with all of these you are ready to build containers on wheezy!

## Ubuntu

Everything above saucy should be supported (saucy is tested). The default configuration is good to go!

If your containers do not get an ip address from dhcp you need to turn off checksum offloading on the bridge:

```/sbin/ethtool -K lxcbr0 tx off```


Building from source
====================

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

Example packer templates
========================

Building wheezy on saucy:

```
{
  "builders": [
    {
      "type": "lxc",
      "config_file": "lxc/config",
      "template_name": "debian",
      "template_environment_vars": [
        "MIRROR=http://http.debian.net/debian/",
        "SUITE=wheezy"
      ],
      "target_runlevel": 3
    }
  ]
}
```

The config file is an lxc config file which will be bundled with the machine. You can create your own or just grab the debian or ubuntu one from [this repo](https://github.com/fgrehm/vagrant-lxc-base-boxes/tree/master/conf).


Building wheezy on wheezy:

```
{
  "builders": [
    {
      "type": "lxc",
      "config_file": "lxc/config",
      "template_name": "debian",
      "template_parameters": ["--arch", "amd64", "--release", "wheezy"],
      "template_environment_vars": [
        "MIRROR=http://http.debian.net/debian/"
      ],
      "target_runlevel": 3
    }
  ],
  "provisioners": [
    {
      "type": "shell",
      "only": ["lxc"],
      "environment_vars": [
        "DISTRIBUTION=debian",
        "RELEASE=wheezy"
      ],
      "scripts": [
        "scripts/lxc/base.sh",
        "scripts/lxc/vagrant-lxc-fixes.sh"
      ]
    }
  ],
  "post-processors": [
    {
      "type": "compress",
      "output": "output-vagrant/wheezy64-lxc.box"
    }
  ],
}
```

Note the differences in template parameters/envvars!

Vagrant publishing
==================

The output artifact can be compressed with the compress publisher to create a working vagrant box (see example).