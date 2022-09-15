# Floatilla - a bunch of floats

A tool to batch send FLOATS to users

Download a release version from https://github.com/findonflow/floatilla/releases or build locally using `make`

usage `flotilla <eventId` will send floats to every users in `recipients.csv`
 
## Prepare

Follow this guideline for setting up your non-custotial account. If you already have one skip this step
https://developers.flow.com/flow/dapp-development/mainnet-account-setup

set up env vars FLOTILLA_PRIVATE_KEY and FLOATILLA_ADDRESS using the generate private key and address

If you need to create an example float to test you can do so with the followig snippet
```
flow transactions send transactions/adminMintFloat.cdc 0x6d1898524fe4a880 FloatillaTest "This is a test of floatilla" "QmWAEsWxedNcMhMvKpdUfmFMdxUfej9U1wvtL9BKb57ns7" "https://github.com/bjartek/floatilla" -n mainnet --signer mainnet-admin
```

## How to use
Create a file `recipients.csv` and add the addresses you want to  or even easier use their .find names!
run 
```
floatilla <eventId>
```

To se available options run `floatilla -h`


## Building

run `make` to build for your os or `make build_all` to build for all architectures.

## Troubleshoot
 
 - The script will add 100 keys to your account of the same FLOATILLA_ADDRESS if it has a single key, if you have a different setup (like a non-custodial blocto) then modify the source or just ensure you have more keys. There is an adminAddKeys transaction in the transactions folder

## Donate
If you like this project feel free to donate me sone flow/fusd over at https://find.xyz/bjartek and say why you love floatilla.

