driver-projector
================

driver-projector is a Ninja Sphere driver that can control various projectors. It currently uses [go-dell][1], so only a handful of projectors are supported.

Compatibility
=============

Please check [go-dell][1] for an up to date compatibility list, but so far, the following projectors are support and tested:

- Dell s300wi
- Dell s500wi

More projectors may be supported, especially if they have a Crestron UI. If the device supports the Dynamic Device Discovery Protocol, then they should be automatically found by this driver.

[1]: http://github.com/Grayda/go-dell
