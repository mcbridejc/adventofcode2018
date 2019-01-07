
package main

import (
    "fmt"
)

// The day21 assembly program, transcribed to go for better readability
// This program is essentially generating a sequence of numbers in r4, until 
// a number equal to r0 is found. Is this some pseudo random number generator
// algorithm? 
func SimulateProgram(r0 int) {
    r2 := 0
    r3 := 0
    r4 := 0
    r5 := 0
    seqCount := 0
    r4 = 0
    r3 = 0x10000
    r4 = 3730679
    for {
        r5 = r3 & 0xff
        r4 += r5
        r4 &= 0xffffff
        r4 *= 65899
        r4 &= 0xffffff
        if r3 < 256 {
            fmt.Printf("%d: r4 = 0x%08x (%d) \n", seqCount, r4, r4)
            seqCount++
            if r0 == r4 {
                break
            } else {
                r3 = r4 | 0x10000
                r4 = 3730679
                continue
            }
        }
        r5 = 0
        for {
            r2 = (r5 + 1) * 256
            if r2 > r3 {
                break
            }
            r5 += 1
        }
        r3 = r5
    }
}

// Repackage the program as an iterator to generate sequences
func NewSequenceGenerator() (func ()int) {
    r2 := 0
    r3 := 0
    r4 := 0
    r5 := 0

    r4 = 0
    r3 = 0x10000
    r4 = 3730679

    return func () int {
        for {
            r5 = r3 & 0xff
            r4 += r5
            r4 &= 0xffffff
            r4 *= 65899
            r4 &= 0xffffff
            if r3 < 256 {
                ret := r4
                r3 = r4 | 0x10000
                r4 = 3730679
                return ret
            }
            r5 = 0
            for {
                r2 = (r5 + 1) * 256
                if r2 > r3 {
                    break
                }
                r5 += 1
            }
            r3 = r5
        }
    }
}