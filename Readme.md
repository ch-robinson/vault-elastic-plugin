# vault-elastic-plugin

# Currently a work in progress

### Setup

***NOTE: For Mac or linux, make should be installed by default so you can skip steps 2 and 3 ***

1. Install Golang: [Go Install](https://golang.org/doc/install)
2. Install Make: [GNU Make Install](http://gnuwin32.sourceforge.net/packages/make.htm)
3. Add make to path: ```SET PATH=%PATH%;C:\Program Files (x86)\GnuWin32\bin```
4. Install glide [glide](https://github.com/Masterminds/glide)
    - Windows - install to "C:\glide" and set path ```SET PATH=%PATH%;C:\glide```
    - Mac - ```brew install glide```
5. ```go get github.com/CH-Robinson/vault-elastic-plugin```

### Testing locally

1. Download [Vault](https://www.vaultproject.io/downloads.html) and extract the compressed file to the location of your choosing
2. Create the config.hcl: ```echo 'plugin_dirctory = "<path-to-plugin-binary>"' > config.hcl```
3. Run Vault locally (assuming the executable is in path): ```vault server -dev -config config.hcl```
4. Create hash: ```cd bin && openssl sha256 bin/<system>/vault-elastic-plugin-x86-64```
5. Add the plugin: ```VAULT_ADDR=http://127.0.0.1:8200 <path-to-vault>/vault.exe write sys/plugins/catalog/vault-elastic-plugin \ sha_256=<from-create-hash> command="vault-elastic-plugin.exe"```
6. Run the plugin: ```VAULT_ADDR=http://127.0.0.1:8200 <path-to-vault>/vault write database/config/elastic_test \ connection_url=<elastic-base-uri> username=vault_admin password=<password> plugin_name=vault-elastic-plugin allowed_roles="*"```

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

1. install dependencies
```make depends```
2. build
```make build```
3. test
```make test```

