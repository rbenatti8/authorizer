# authorizer

## Architecture

In this project, I tried to apply hexagonal architecture paradigms, such as dividing adapters into primary (driver) and secondary (driven).

In a nutshell, I believe the authorization process could be broken into two parts, which is what I tried to make clear. The first one would be the one that tells me if I have the balance for the operation and the second one that says if I can spend, called respectively as ledger and spending control.
## Libs

- [gomock](https://github.com/golang/mock) - mocking framework
- [testify](https://github.com/stretchr/testify) - tools for testifying

## Build

```sh
make build
```

## Run

```sh
make run < 'YOUR_FILE'
```