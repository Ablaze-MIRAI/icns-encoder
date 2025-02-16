/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
)

var default_sizes = [][]string{
	{"ic07", "128", "128"},   // ic07: 128x128 PNG/JPEG2000
	{"ic08", "256", "256"},   // ic08: 256x256 PNG/JPEG2000
	{"ic09", "512", "512"},   // ic09: 512x512 PNG/JPEG2000
	{"ic10", "1024", "1024"}, // ic10: 1024x1024 PNG/JPEG2000 (512x512@2x)
	{"ic11", "32", "32"},     // ic11: 32x32 PNG/JPEG2000 (16x16@2x)
	{"ic12", "64", "64"},     // ic12: 64x64 PNG/JPEG2000 (32x32@2x)
	{"ic13", "256", "256"},   // ic13: 256x256 PNG/JPEG2000 (128x128@2x)
	{"ic14", "512", "512"},   // ic14: 512x512 PNG/JPEG2000 (256x256@2x)
}

func main() {
	os.Exit(run())
}

func run() int {
	var input_filepath string
	flag.StringVar(&input_filepath, "i", "", "Specify input file path (.png)")

	var output_filepath string
	flag.StringVar(&output_filepath, "o", "", "Specify output file path (.icns)")

	flag.Parse()

	if input_filepath == "" || output_filepath == "" {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
		return 1
	}

	var icons []Icon

	for _, size := range default_sizes {
		width, err := strconv.Atoi(size[1])
		if err != nil {
			panic(err)
		}
		height, err := strconv.Atoi(size[2])
		if err != nil {
			panic(err)
		}

		data, err := PngResize(input_filepath, width, height)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
			return 1
		}

		icons = append(icons,
			Icon{
				Type: size[0],
				Data: data,
			},
		)
	}

	if err := IcnsWrite(icons, output_filepath); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
		return 1
	}

	return 0
}
