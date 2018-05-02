# vault-elastic-plugin

### Setup

*** NOTE: For Mac or linux, make should be installed by default so you can skip steps 2 and 3 ***

1. Install Golang: [Go Install](https://golang.org/doc/install)
2. Install Make: [GNU Make Install](http://gnuwin32.sourceforge.net/packages/make.htm)
3. Add make to path: ```SET PATH=%PATH%;C:\Program Files (x86)\GnuWin32\bin```
4. Install glide [glide](https://github.com/Masterminds/glide)
    - Windows - install to "C:\glide" and set path ```SET PATH=%PATH%;C:\glide```
    - Mac - ```brew install glide```
5. ```go get github.com/CH-Robinson/vault-elastic-plugin```

### Testing locally

1. Download [Vault](https://www.vaultproject.io/downloads.html) and extract the compressed file to the location of your choosing
2. Add Vault to path
3. In a new terminal run Vault: ```make run-vault```
4. In a new terminal run: 
    - With build and vault DB configuration: ```make test-plugin ELASTIC_BASE_URI=<uri> ELASTIC_PASSWORD=<password> ELASTIC_USERNAME=<username> INCLUDE_BUILD=true ENABLE_VAULT_DB=true```
    - Without build: ```make test-plugin ELASTIC_BASE_URI=<uri> ELASTIC_PASSWORD=<password> ELASTIC_USERNAME=<username> INCLUDE_BUILD=false ENABLE_VAULT_DB=true```

### Build

- Unix based: ```make build```
- Powershell (if make is not in path): ```C:\Program Files (x86)\GnuWin32\bin\make.exe build```
- bash for Windows (if make is not in path): ```/c/Program\ Files\ \(x86\)/GnuWin32/bin/make.exe build```

The executable binary is located ../bin/run

### Unit Testing

At the root of the project, run 
- Unix based: ```make test```
- Powershell (if make is not in path): ```C:\Program Files (x86)\GnuWin32\bin\make.exe test```
- bash for Windows (if make is not in path): ```/c/Program\ Files\ \(x86\)/GnuWin32/bin/make.exe test```

### Dependencies

- Unix based: ```make depends```
- Powershell (if make is not in path): ```C:\Program Files (x86)\GnuWin32\bin\make.exe depends```
- bash for Windows (if make is not in path): ```/c/Program\ Files\ \(x86\)/GnuWin32/bin/make.exe depends```

