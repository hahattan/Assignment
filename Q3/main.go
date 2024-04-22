package main

import (
	"fmt"
	"strings"
)

const (
	numDisks   = 4
	numStripes = 3
	stripeSize = 5
)

type RAID struct {
	Disks      [][]byte
	NumStripes int
	StripeSize int
}

func NewRAID(numDisks, numStripes, stripeSize int) *RAID {
	disks := make([][]byte, numDisks)
	for i := range disks {
		disks[i] = make([]byte, numStripes*stripeSize)
	}
	return &RAID{
		Disks:      disks,
		NumStripes: numStripes,
		StripeSize: stripeSize,
	}
}

// Write data to RAID using specified RAID level
func (r *RAID) Write(data []byte, RAIDLevel string) {
	switch RAIDLevel {
	case "RAID0":
		r.RAID0Write(data, r.Disks)
	case "RAID1":
		r.RAID1Write(data, r.Disks)
	case "RAID10":
		r.RAID10Write(data, r.Disks)
	case "RAID5":
		r.RAID5Write(data, r.Disks)
	case "RAID6":
		r.RAID6Write(data)
	default:
		fmt.Println("RAID level not supported")
	}
}

func (r *RAID) Read(RAIDLevel string, length int) []byte {
	switch RAIDLevel {
	case "RAID0":
		return r.RAID0Read(length)
	case "RAID1":
		return r.RAID1Read(length)
	case "RAID10":
		return r.RAID10Read(length)
	case "RAID5":
		return r.RAID5Read(length)
	case "RAID6":
		r.RAID6Read(length)
	default:
	}

	return []byte("not implemented")
}

func (r *RAID) RAID0Write(data []byte, disks [][]byte) {
	if disks == nil {
		disks = r.Disks
	}

	var diskIdx, stripeOffset, stripeIdx int
	for i := 0; i < len(data); i++ {
		disks[diskIdx][stripeOffset+stripeIdx] = data[i]
		stripeIdx++
		if stripeIdx == stripeSize {
			diskIdx++
			if diskIdx == numDisks {
				diskIdx = 0
				stripeOffset += stripeSize
				continue
			}
			stripeIdx = 0
		}
	}
}

func (r *RAID) RAID0Read(length int) []byte {
	res := make([]byte, length)

	var diskIdx, stripeIdx, stripeOffset int
	for i := 0; i < length; i++ {
		res[i] = r.Disks[diskIdx][stripeOffset+stripeIdx]
		stripeIdx++
		if stripeIdx == stripeSize {
			diskIdx++
			if diskIdx == numDisks {
				diskIdx = 0
				stripeOffset += stripeSize
				continue
			}
			stripeIdx = 0
		}
	}

	return res
}

func (r *RAID) RAID1Write(data []byte, disks [][]byte) {
	if disks == nil {
		disks = r.Disks
	}
	for i := 0; i < len(r.Disks); i++ {
		copy(r.Disks[i], data)
	}
}

func (r *RAID) RAID1Read(length int) []byte {
	for i := 0; i < numDisks; i++ {
		if !allZeros(r.Disks[i]) {
			return r.Disks[i][:length]
		}
	}

	return nil
}

func (r *RAID) RAID10Write(data []byte, disks [][]byte) {
	if disks == nil {
		disks = r.Disks
	}

	partition := numDisks / 2
	var diskIdx, stripeIdx, stripeOffset int
	for i := 0; i < len(data); i++ {
		disks[diskIdx][stripeOffset+stripeIdx] = data[i]
		stripeIdx++
		if stripeIdx == stripeSize {
			diskIdx += partition
			if diskIdx >= numDisks {
				diskIdx = 0
				stripeOffset += stripeSize
				continue
			}
			stripeIdx = 0
		}
	}

	for i := 0; i < numDisks; i += partition {
		for j := 1; j < partition; j++ {
			copy(disks[i+j], disks[i])
		}
	}
}

func (r *RAID) RAID10Read(length int) []byte {
	res := make([]byte, length)

	partition := numDisks / 2
	var diskIdx, stripeIdx, stripeOffset int
	for i := 0; i < length; i++ {
		if allZeros(r.Disks[diskIdx]) {
			diskIdx++
		}
		res[i] = r.Disks[diskIdx][stripeOffset+stripeIdx]
		stripeIdx++
		if stripeIdx == stripeSize {
			diskIdx += partition
			if diskIdx >= numDisks {
				diskIdx = diskIdx % numDisks
				stripeOffset += stripeSize
				continue
			}
			stripeIdx = 0
		}
	}

	return res
}

func (r *RAID) RAID5Write(data []byte, disks [][]byte) {
	if disks == nil {
		disks = r.Disks
	}

	for i := 0; i < len(data); i++ {
		stripeIndex := i / (stripeSize * (numDisks - 1))
		offset := i % (stripeSize * (numDisks - 1))

		dataDiskIndex := (stripeIndex*(numDisks-1) + offset/stripeSize) % numDisks
		parityDiskIndex := (dataDiskIndex + 1) % numDisks

		r.Disks[dataDiskIndex][offset%stripeSize] = data[i]

		// Calculate and update parity disk
		for j := range r.Disks[parityDiskIndex] {
			r.Disks[parityDiskIndex][j] ^= data[i]
		}
	}
}

func (r *RAID) RAID5Read(length int) []byte {
	data := make([]byte, length)

	for i := 0; i < length; i++ {
		stripeIndex := i / (stripeSize * (numDisks - 1))
		offset := i % (stripeSize * (numDisks - 1))

		dataDiskIndex := (stripeIndex*(numDisks-1) + offset/stripeSize) % numDisks

		data[i] = r.Disks[dataDiskIndex][offset%stripeSize]
	}

	return data
}

func (r *RAID) RAID6Write(data []byte) {
	// Calculate dual parity and distribute data across disks
	for i := 0; i < len(r.Disks); i++ {
		for j := 0; j < r.StripeSize; j++ {
			if j%(len(r.Disks)-2) != i && j%(len(r.Disks)-1) != i {
				r.Disks[i][j] = data[j]
			}
		}
	}
}

func (r *RAID) RAID6Read(length int) []byte {
	data := make([]byte, length)

	for i := 0; i < length; i++ {
		stripeIndex := i / (stripeSize * (numDisks - 1))
		offset := i % (stripeSize * (numDisks - 1))

		dataDiskIndex := (stripeIndex*(numDisks-1) + offset/stripeSize) % numDisks

		data[i] = r.Disks[dataDiskIndex][offset%stripeSize]
	}

	return data
}

func (r *RAID) ClearDisk(index int) {
	for i := range r.Disks[index] {
		r.Disks[index][i] = 0
	}
}

func allZeros(arr []byte) bool {
	for _, b := range arr {
		if b != 0 {
			return false
		}
	}

	return true
}

// Function to calculate the parity of a byte
func calculateParity(b byte) int {
	count := 0
	for b != 0 {
		count ^= int(b & 1)
		b >>= 1
	}
	return count
}

// Function to calculate the parity of a []byte and return the parity bits as a []byte
func calculateByteArrayParity(data []byte) []byte {
	var parityBits []byte
	for _, b := range data {
		parity := calculateParity(b)
		// Append the parity bit to the slice
		parityBits = append(parityBits, byte(parity))
	}
	return parityBits
}

// Function to reconstruct a []byte from a given parity and size
func reconstructByteArray(parity []byte, size int) []byte {
	byteArray := make([]byte, size)
	for i := range byteArray {
		byteArray[i] = parity[i] // Retrieve the parity bit
	}
	return byteArray
}

func main() {
	data := []byte("Hello, RAID!")
	n := len(data)

	for _, level := range []string{"RAID0", "RAID1", "RAID10", "RAID5", "RAID6"} {
		raid := NewRAID(numDisks, numStripes, stripeSize)

		// Write data to RAID with different RAID levels
		raid.Write(data, level)

		fmt.Println(level, ":")
		var res []string
		for _, disk := range raid.Disks {
			res = append(res, string(disk))
		}
		fmt.Printf("Before: %s\n", strings.Join(res, " | "))
		// Clear one of the disks
		raid.ClearDisk(2)

		// Read data from RAID and print
		fmt.Printf("After: %s\n\n", raid.Read(level, n))
	}
}
