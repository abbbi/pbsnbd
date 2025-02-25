/*
PBS nbdkit plugin.
Export PBS VM disk image Backups via NBD
Copyright (C) 2025  Michael Ablassmeier <abi@grinser.de>
*/
package main

import (
	"C"
	"fmt"
	"pbsnbd/nbdkit"
	"time"
	"unsafe"

	bps "github.com/elbandi/go-proxmox-backup-client"
)

var pluginName = "pbsnbd"

type PBSDiskPlugin struct {
	nbdkit.Plugin
}

type PBSDiskConnection struct {
	nbdkit.Connection
}

var size uint64
var image string
var timestamp string
var vmid string
var disk []byte

var repo string
var password string
var fingerprint string
var namespace string

var image_set = false
var timestamp_set = false
var vmid_set = false
var repo_set = false
var namespace_set = false
var fingerprint_set = false
var password_set = false

var client *bps.ProxmoxRestore
var imagefh *bps.RestoreImage
var backupTime string

func (p *PBSDiskPlugin) Config(key string, value string) error {
	if key == "image" {
		image = value
		image_set = true
		return nil
	} else if key == "timestamp" {
		timestamp = value
		timestamp_set = true
		return nil
	} else if key == "vmid" {
		vmid = value
		vmid_set = true
		return nil
	} else if key == "repo" {
		repo = value
		repo_set = true
		return nil
	} else if key == "fingerprint" {
		fingerprint = value
		fingerprint_set = true
		return nil
	} else if key == "namespace" {
		namespace = value
		namespace_set = true
		return nil
	} else if key == "password" {
		password = value
		password_set = true
		return nil
	}

	return nbdkit.PluginError{Errmsg: "unknown parameter: " + key}
}

func (p *PBSDiskPlugin) ConfigComplete() error {
	if !image_set {
		return nbdkit.PluginError{Errmsg: "image parameter is required"}
	}
	if !timestamp_set {
		return nbdkit.PluginError{Errmsg: "timestamp parameter is required"}
	}
	if !vmid_set {
		return nbdkit.PluginError{Errmsg: "vmid parameter is required"}
	}
	if !repo_set {
		return nbdkit.PluginError{Errmsg: "repo parameter is required"}
	}
	if !fingerprint_set {
		return nbdkit.PluginError{Errmsg: "fingerprint parameter is required"}
	}
	if !password_set {
		return nbdkit.PluginError{Errmsg: "password parameter is required"}
	}
	if !namespace_set {
		return nbdkit.PluginError{Errmsg: "namespace parameter is required"}
	}
	return nil
}

func (p *PBSDiskPlugin) GetReady() error {
	f, _ := time.Parse(time.RFC3339, timestamp)
	t := uint64(f.Unix())

	var err error
	client, err = bps.NewRestore(
		repo,
		namespace,
		"vm",
		vmid,
		t,
		password,
		fingerprint,
		"",
		"",
	)
	if err != nil {
		return err
	}
	fmt.Printf("Connected to PBS version: [%s]\n", bps.GetVersion())
	fmt.Printf("Attempt to open image [vm/%s/%s/%s]\n", vmid, timestamp, image)
	imagefh, err = client.OpenImage(image + ".fidx")
	if err != nil {
		return err
	}
	fmt.Printf("Sucessfully opened image [vm/%s/%s/%s]\n", vmid, timestamp, image)
	return nil
}

func (p *PBSDiskPlugin) Open(readonly bool) (nbdkit.ConnectionInterface, error) {
	return &PBSDiskConnection{}, nil
}

func (c *PBSDiskConnection) GetSize() (uint64, error) {
	var err error
	size, err = imagefh.Size()
	if err != nil {
		return 0, err
	}
	return size, nil
}

func (c *PBSDiskConnection) PRead(buf []byte, offset uint64, flags uint32) error {
	imagefh.ReadAt(buf, int64(offset))
	return nil
}

func (c *PBSDiskConnection) CanWrite() (bool, error) {
	return false, nil
}

//export plugin_init
func plugin_init() unsafe.Pointer {
	return nbdkit.PluginInitialize(pluginName, &PBSDiskPlugin{})
}

func main() {}
