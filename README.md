This is a test repo to see what it's like to write a parser with https://pkg.go.dev/github.com/mna/pigeon

Nothing here is complete, and it may never be completed.

## Well this is probably a failure

Partially implemented parser, no transformer, basic compiler in the go version, already slower than Lua version:


```
(master) ~/code/go/mooonscript > ls -lah big.moon 
-rw-r--r-- 1 leafo leafo 2.1M Nov  4 20:11 big.moon
(master) ~/code/go/mooonscript > time ./moonscript-go big.moon  > /dev/null
real	0m5.158s
user	0m10.425s
sys	0m0.089s
(master) ~/code/go/mooonscript > time moonc -p big.moon  > /dev/null

real	0m3.876s
user	0m3.725s
sys	0m0.140s
(master) ~/code/go/mooonscript > time moonc-luajit -p big.moon  > /dev/null

real	0m2.176s
user	0m2.067s
sys	0m0.100s
```

