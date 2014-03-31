packer-builder-lxc
==========

Lxc builder for packer, with working provisioning. At the moment chroot is used for provisioning, because lxc-attach is completely broken on debian to this date. An lxc-attach provisioner is planned for other distributions. Since the building process uses lxc-create it relies on your system's lxc templates, so results of building the same config my vary based on different host os distributions/releases.

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
      "distribution": "debian",
      "release": "wheezy",
      "mirror_url": "http://ftp.hu.debian.org/debian/"
    }
```
