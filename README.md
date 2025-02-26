# Proxmox Backup Server NBDKIT Plugin

Export your virtual machine disk backups from a Proxmox Backup Server via NBD.

To start an export, pass the plugin to nbdkit with required arguments:

```
 nbdkit --filter cow -f ./nbdkit-pbs-plugin.so \
    image=drive-scsi0.img \                     # disk image to mount
    vmid=103 \                                  # vm id
    timestamp="2025-02-25T10:32:54Z" \          # snapshot timestamp
    repo="root@pam@192.168.161.241:test" \      # repository
    fingerprint="EE:2F:53:29:F2:69:61:0C:D5:96:55:50:6B:2B:84:8E:EF:3E:A6:B0:CA:18:5C:F7:92:BA:54:71:15:56:83:5B" \
    password=test \
    namespace=                                  # namespace, optional

    pbsnbd Plugin loaded
    Connecting PBS: [root@pam@192.168.161.241:test] Namespace: []
    Connected to PBS version: [1.4.1 (UNKNOWN)]
    Attempt to open image [vm/103/2025-02-25T15:48:38Z/drive-scsi0.img]
    Successfully opened image [vm/103/2025-02-25T15:48:38Z/drive-scsi0.img]
```

The NBD device is by default reachable via localhost. The COW Filter (--filter
cow) allows read/write operations to the disks without altering the original
data. This way, you can boot off the NBD device directly using Qemu:

```
 qemu-system-x86_64 -m 2000 -hda nbd://localhost -cpu host -enable-kvm
```

You can also start nbdkit without the Filter option and use an overlay
image instead:

```
 qemu-img create -b nbd://localhost -f qcow2 -F raw image.qcow2
 qemu-system-x86_64 -m 2000 -hda image.qcow2 -cpu host -enable-kvm
```

Mount volumes using guestmount:

```
 nbdkit -f ./nbdkit-pbs-plugin.so [..]
 guestmount --ro --format=raw -ia nbd://localhost /empty/
```

To map the nbd backend into an device:

```
 qemu-nbd  -c /dev/nbd0 nbd://localhost
 fdisk -l /dev/nbd0
 Disk /dev/nbd5: 128 GiB, 137438953472 bytes, 268435456 sectors
 [..]
 Device      Boot   Start       End   Sectors   Size Id Type
 /dev/nbd5p1 *       2048    999423    997376   487M 83 Linux
 /dev/nbd5p2       999424   4999167   3999744   1.9G 82 Linux swap / Solaris
 /dev/nbd5p3      4999168 268433407 263434240 125.6G 83 Linux
```

Access via regular nbd capable tools such as qemu-img:

```
# show nbd backend info:

 nbdinfo nbd://localhost
 protocol: newstyle-fixed without TLS, using structured packets
 export="":
        export-size: 137438953472 (128G)
        content: DOS/MBR boot sector
        uri: nbd://localhost:10809/
 [..]

# convert backup to qcow2 file (load plugin without cow backend for enhanced iops)

 qemu-img convert -p nbd://localhost -f raw -O qcow2 image.qcow2

# convert backup to vmdk file:

 qemu-img convert -p nbd://localhost -f raw -O vmdk image.vmdk
```

# building

```
 sudo apt-get install nbdkit-plugin-dev
 make
```

# notes

Use nbdkit options to

 * change listen address / port
 * enhance performance by using multiple threads (-t X)
