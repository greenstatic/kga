# Installation
!!! Note
    Since kga is still being developed, we do not offer any pre-built releases.
    In order to install kga you need to clone the repository and build kga yourself.
    Don't worry the app is written in Go, building actually works.


## Requirements
* Helm v3
* Kustomize v3

### Build Requirements
* make
* go 1.14
* python3 (to build the docs)
* pip3 (to build the docs)


## How to Build
### 1. Clone the repository
```bash
git clone https://github.com/greenstatic/kga.git
```

### 2. Run make
```bash
cd kga
make
```

### 3. Copy kga into your path
Move `./bin/kga` into your PATH


### Updates
When you wish to update kga (if you are have an old version)
```bash
# Make sure you are in the kga repository
git pull
make
# Move ./bin/kga into your PATH
```