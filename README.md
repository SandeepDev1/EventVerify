# This repo is related to this [Project](https://github.com/SandeepDev1/Ekshetra)

### You can approve the event payments from telegram with just one command

## Dependencies
```bash
go 1.16+
```
Download Go 1.16 from [here](https://go.dev/dl/). Can be found in Archived versions

## Configuration
`token` - The Telegram token to run the bot

`dbUrl` - The mongodb URL

`dbName` - The database name

## Build
```bash
# to build a native executable
go build .
```

## Run
```bash
## windows
# 1. Double click on that exe
# 2. Or open cmd in that folder and run `verify.exe`

## MacOS and Linux
# cd into that folder where the executable is present

./verify
```

## Commands

```bash
/approve <event-id> <txn-id>
# Event ID = the event id is passed in the txn note to the user as well as sent in telegram message
# TXN-id = the transaction id of the payment made to save in DB
```