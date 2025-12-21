<div align="center">
    <pre style="background: none;">
   █████████                      █████   ████ █████ ███████████
  ███░░░░░███                    ░░███   ███░ ░░███ ░█░░░███░░░█
 ███     ░░░   ██████             ░███  ███    ░███ ░   ░███  ░ 
░███          ███░░███ ██████████ ░███████     ░███     ░███    
░███    █████░███ ░███░░░░░░░░░░  ░███░░███    ░███     ░███    
░░███  ░░███ ░███ ░███            ░███ ░░███   ░███     ░███    
 ░░█████████ ░░██████             █████ ░░████ █████    █████   
  ░░░░░░░░░   ░░░░░░             ░░░░░   ░░░░ ░░░░░    ░░░░░  
    </pre>
</div>

<div align="center">
    <h6>Currently under active development and breaking changes are possible</h6>
    <h3>The Go-Kit library offers a collection of common data structures, algorithms, and utilities designed to streamline and simplify the development process.</h3>
</div>

<div>
    <h1 align="center">Features</h1>
    <ul>
        <li>Configuration — generic utils to parse and map .env or .json files</li>
        <li>Logger — logging with different log levels, zero-alloc object logging and async logging </li>
        <li>Collections — high-level interface abstraction over arrays, set, and double-linked list with common functions to operate</li>
        <li>Graph — graph data structure that is implemented using adjacency matrix</li>
        <li>Big Numbers — wrapper over big.Int and big.Float for comfortable usage and mutability handling</li>
        <li>Caching — implementation of LRU and LFU caches</li>
        <li>CGO Memory — a set of functions that allow to work with raw memory</li>
    </ul>
</div>


<h1 align="center">Installation</h1>

```bash
#latest released
go get github.com/0x626f/go-kit
```

```bash
#dev staged
go get github.com/0x626f/go-kit@develop
```

<h1 align="center">Testing</h1>

```bash
go test ./...
```

```bash
go test -bench=. -benchmem ./...
```


<div align="center">
    <h1 align="center">Contributing</h1>
    <h3>Contributions are welcome! Please ensure tests pass and code is properly documented.</h3>
</div>
