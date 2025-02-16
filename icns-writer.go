/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package main

import (
	"encoding/binary"
	"errors"
	"os"
)

var icns_magic = []byte{0x69, 0x63, 0x6e, 0x73}

type IcnsHeader struct {
	Magic      []byte // Magic
	FileLength uint32 // File size
}

type IconHeader struct {
	Type        []byte // Icon Type
	BlockLength uint32 // Block length (Icon header + Icon Data)
}

type Icon struct {
	Type string // Icon type
	Data []byte // Data
}

func uint32ToBytesBigEndian(v uint32) []byte {
	return binary.BigEndian.AppendUint32([]byte{}, v)
}

func IcnsWrite(icons []Icon, path string) error {
	var icns []byte

	for _, icon := range icons {
		if err := PngValidate(icon.Data); err != nil {
			return errors.New("png reader: " + err.Error())
		}
	}

	file_length := uint32( // file length without data length
		8 /* Icns header size */ + 8*len(icons), /* Icon headers size */
	)
	for _, icon := range icons {
		file_length += uint32(len(icon.Data))
	}

	header := IcnsHeader{
		Magic:      icns_magic,
		FileLength: file_length,
	}

	var icon_headers []IconHeader
	for _, icon := range icons {
		icon_headers = append(icon_headers,
			IconHeader{
				Type:        []byte(icon.Type),          // Icon type
				BlockLength: 8 + uint32(len(icon.Data)), // Block length (Icon header + Icon data)
			},
		)
	}

	// Write icns

	// header
	icns = append(icns,
		header.Magic..., // Magic
	)
	icns = append(icns,
		uint32ToBytesBigEndian(header.FileLength)..., // File length
	)
	// icons
	for i, icon_header := range icon_headers {
		icon := icons[i]

		icns = append(icns,
			icon_header.Type..., // Icon type
		)
		icns = append(icns,
			uint32ToBytesBigEndian(icon_header.BlockLength)...,
		)
		icns = append(icns,
			icon.Data...,
		)
	}

	if err := os.WriteFile(path, icns, 0644); err != nil {
		return err
	}

	return nil
}
