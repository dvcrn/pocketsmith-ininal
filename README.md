# Ininal importer for Pocketsmith

A simple CLI to import Ininal transactions into Pocketsmith. This is highly experimental and requires you to obtain several authentication tokens yourself.

## Usage

You need to obtain the following tokens from Ininal (through inspecting the app traffic):
- Device ID
- Login Token
- User Token
- Login Bearer Token
- Device Signature
- Login Credential (Phone Number)
- Password (App PIN)

Then run with:

```
go run main.go \
  -device-id=YOUR_DEVICE_ID \
  -login-token=YOUR_LOGIN_TOKEN \
  -user-token=YOUR_USER_TOKEN \
  -login-bearer-token=YOUR_LOGIN_BEARER_TOKEN \
  -device-signature=YOUR_DEVICE_SIGNATURE \
  -login-credential=YOUR_PHONE_NUMBER \
  -password=YOUR_APP_PIN \
  -pocketsmith-token=YOUR_POCKETSMITH_TOKEN
```

Or set environment variables:

```
export ININAL_DEVICE_ID=xxx
export ININAL_LOGIN_TOKEN=xxx
export ININAL_USER_TOKEN=xxx
export ININAL_LOGIN_BEARER_TOKEN=xxx
export ININAL_DEVICE_SIGNATURE=xxx
export ININAL_LOGIN_CREDENTIAL=xxx
export ININAL_PASSWORD=xxx
export POCKETSMITH_TOKEN=xxx

go run main.go
```


### Run with docker (recommended)

```
docker run \
  -e ININAL_DEVICE_ID=xxx \
  -e ININAL_LOGIN_TOKEN=xxx \
  -e ININAL_USER_TOKEN=xxx \
  -e ININAL_LOGIN_BEARER_TOKEN=xxx \
  -e ININAL_DEVICE_SIGNATURE=xxx \
  -e ININAL_LOGIN_CREDENTIAL=xxx \
  -e ININAL_PASSWORD=xxx \
  -e POCKETSMITH_TOKEN=xxx \
  dvcrn/pocketsmith-ininal
```

## Features

- Automatically creates Ininal institution and account in Pocketsmith if they don't exist
- Updates account balance
- Imports transactions with reference numbers
- Prevents duplicate transactions by checking reference numbers
- Handles OTP authentication if required

## License

MIT, commercial use excluded