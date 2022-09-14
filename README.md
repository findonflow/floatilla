# Floatilla - a bunch of floats

A tool to batch send FLOATS to users


## Setup
1. create an account on flow mainnet that you have the key for, or reuse an existing one
2. export the following env vars:
 - FLOATILLA_ADDRESS
 - FLOATILLA_PUBLIC_KEY
 - FLOATILLA_PRIVATE_KEY
3. call `addKeys.sh` to add 100 more copies of your key to your account (or ensure that the key index  1-100 is there)
4. ensure that your address either has created this float or is a shared minter
5. create a file with 1 address per line that wants to float
6. call `go run main.go -- <fileName> <eventId>`

