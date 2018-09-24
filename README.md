// README.md

##"Drain Random from Chaos."
##Breeze - a new family of fast CB-PRNG

####Dr. Andreas Briese
#####2018/09/24   
#####eduToolbox@Bri-C GmbH, Sarstedt

Breeze32 
========

[![Build Status](https://travis-ci.org/AndreasBriese/ipLocator.png?branch=master)](http://travis-ci.org/AndreasBriese/ipLocator)

Actual version of the breeze PRNG: (That was the old one with many blabla: https://github.com/AndreasBriese/breeze)

import with `go get github.com/AndreasBriese/breeze32`

2018 Easier structure and even faster and more reliable.

As long as breeze32 is running on one cpu, it appeared to be threadsafe even with intensive use. 

If you need a multiprocessor safe prng use the ByteMP function (and/or adapt it for the other variants). And consider applying the mutex in  roundTrip().  

See the test file for usage.

__Benchmarks__ (`go test -bench=.`)

    goos: darwin
    goarch: amd64
    pkg: github.com/AndreasBriese/breeze32
    Benchmark_Init_and_Seed-8               10000000           186 ns/op
    Benchmark_breeze_Uint64-8               200000000            8.45 ns/op
    Benchmark_breeze_Uint32-8               200000000            7.61 ns/op
    Benchmark_breeze_Uint8-8                200000000            6.49 ns/op
    Benchmark_breeze_ByteMP-8               20000000            58.5 ns/op
    Benchmark_breeze_RandIntN_8byte-8       100000000           14.5 ns/op
    Benchmark_breeze_RandIntN_16byte-8      100000000           14.3 ns/op
    Benchmark_breeze_RandIntN_32byte-8      100000000           14.3 ns/op
    Benchmark_breeze_RandIntN_64byte-8      100000000           21.0 ns/op
    Benchmark_mathRand_Intn_8byte-8         100000000           16.3 ns/op
    Benchmark_mathRand_Intn_16byte-8        100000000           16.0 ns/op
    Benchmark_mathRand_Intn_32byte-8        100000000           13.8 ns/op
    Benchmark_mathRand_Intn_64byte-8        30000000            33.9 ns/op
    PASS
    ok      github.com/AndreasBriese/breeze32   30.911s

go version go1.10.2 darwin/amd64 on MBP2013 i7 8GB Ram 

[random Bytes (gray 4kx2k)](https://github.com/AndreasBriese/breeze32/blob/master/rand.png)

__License__
MIT-Style for software use // see License