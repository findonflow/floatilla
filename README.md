# Floatilla - a bunch of floats

A tool to batch send FLOATS to users

run `make build` to build a target for your env, or use one of the predefined binaries

## Setup
1. Make sure you have an non-custodial account that you know the private key for
2. export the following env vars:
 - FLOATILLA_ADDRESS
 - FLOATILLA_PUBLIC_KEY
 - FLOATILLA_PRIVATE_KEY
3. ensure that your address either has created this float or is a shared minter
4. create a file with 1 address per line that wants to float
5. call `go run main.go -- <fileName> <eventId>`

## Troubleshoot
 
 - The script will add 100 keys to your account of the same FLOATILLA_PUBLIC_KEY if it has a single key, if you have a different setup (like a non-custodial blocto) then modify the source or just ensure you have more keys. There is an adminAddKeys transaction in the transactions folder

 - The file has to be called from the directory where the flow.json file is in



