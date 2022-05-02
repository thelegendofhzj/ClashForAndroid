//go:build !linux

package tun

import (
	"errors"
	"fmt"
	"os"
	"runtime"

	"github.com/Dreamacro/clash/listener/tun/device"
	"github.com/Dreamacro/clash/listener/tun/device/iobased"

	"golang.zx2c4.com/wireguard/tun"
)

type TUN struct {
	*iobased.Endpoint

	nt     *tun.NativeTun
	mtu    uint32
	name   string
	offset int

	cache []byte
}

func Open(name string, mtu uint32) (_ device.Device, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("open tun: %v", r)
		}
	}()

	var (
		offset     = 4 /* 4 bytes TUN_PI */
		defaultMTU = 1500
	)
	if runtime.GOOS == "windows" {
		offset = 0
		defaultMTU = 0 /* auto */
	}

	t := &TUN{
		name:   name,
		mtu:    mtu,
		offset: offset,
	}

	forcedMTU := defaultMTU
	if t.mtu > 0 {
		forcedMTU = int(t.mtu)
	}

	nt, err := tun.CreateTUN(t.name, forcedMTU) // forcedMTU do not work on wintun, need to be setting by other way

	// retry if abnormal exit on Windows at last time
	if err != nil && runtime.GOOS == "windows" && errors.Is(err, os.ErrExist) {
		nt, err = tun.CreateTUN(t.name, forcedMTU)
	}

	if err != nil {
		return nil, fmt.Errorf("create tun: %w", err)
	}

	t.nt = nt.(*tun.NativeTun)

	tunMTU, err := nt.MTU()
	if err != nil {
		return nil, fmt.Errorf("get mtu: %w", err)
	}
	t.mtu = uint32(tunMTU)

	if t.offset > 0 {
		t.cache = make([]byte, 65535)
	}

	return t, nil
}

func (t *TUN) Read(packet []byte) (int, error) {
	if t.offset == 0 {
		return t.nt.Read(packet, t.offset)
	}

	n, err := t.nt.Read(t.cache, t.offset)

	copy(packet, t.cache[t.offset:t.offset+n])

	return n, err
}

func (t *TUN) Write(packet []byte) (int, error) {
	if t.offset == 0 {
		return t.nt.Write(packet, t.offset)
	}

	packet = append(t.cache[:t.offset], packet...)

	return t.nt.Write(packet, t.offset)
}

func (t *TUN) Close() error {
	if t.Endpoint != nil {
		t.Endpoint.Close()
	}
	return t.nt.Close()
}

func (t *TUN) Name() string {
	name, _ := t.nt.Name()
	return name
}

func (t *TUN) UseEndpoint() error {
	ep, err := iobased.New(t, t.mtu, t.offset)
	if err != nil {
		return fmt.Errorf("create endpoint: %w", err)
	}
	t.Endpoint = ep
	return nil
}

func (t *TUN) UseIOBased() error {
	return nil
}
