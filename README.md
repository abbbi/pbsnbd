== Proxmox Backup Server NBDKIT Plugin ==

Export your virtual machine disks backups from a Proxmox Backup Server via NBD.
Using the COW Filter, it can be used to test-boot virtual machines from outside
a Proxmox Backup Server instance.

Example:

```
 nbdkit --filter cow -f -v ./nbdkit-pbs-plugin.so \
    image=drive-scsi0.img \
    vmid=103 \ 
    timestamp="2025-02-25T10:32:54Z" \
    repo="root@pam@192.168.161.241:test" \
    fingerprint="EE:2F:53:29:F2:69:61:0C:D5:96:55:50:6B:2B:84:8E:EF:3E:A6:B0:CA:18:5C:F7:92:BA:54:71:15:56:83:5B" \
    password=test \
    namespace= \
```

Then, you can boot off the NBD device via:

```
 qemu-system-x86_64 -m 2000 -hda nbd://localhost -cpu host -enable-kvm
```

Or, as alternative, map the nbd device locally and access file systems:

```
 qemu-nbd  -c /dev/nbd5 nbd://localhost
 fdisk -l /dev/nbd5 
 Disk /dev/nbd5: 128 GiB, 137438953472 bytes, 268435456 sectors
 [..]
 Device      Boot   Start       End   Sectors   Size Id Type
 /dev/nbd5p1 *       2048    999423    997376   487M 83 Linux
 /dev/nbd5p2       999424   4999167   3999744   1.9G 82 Linux swap / Solaris
 /dev/nbd5p3      4999168 268433407 263434240 125.6G 83 Linux
```

Or access it via regular nbd tools:

```
nbdinfo nbd://localhost
protocol: newstyle-fixed without TLS, using structured packets
export="":
        export-size: 137438953472 (128G)
        content: DOS/MBR boot sector
        uri: nbd://localhost:10809/
```
