# Proxmox Backup Server NBDKIT Plugin

The golang implementation in this repository may deadlock if nbdkit is forked
into the background. Golang shared libraries do not work well together with
applications that use fork().

As this limits the usecases with this plugin, i have re-implemented the plugin
in C:

[cpbsnbd](https://github.com/abbbi/cpbsnbd)
