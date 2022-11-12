# Darboard

Darboard is intended to be a remote pinning server for IPFS.

## Build Requirements

 - Go 1.19+
 - Docker

API was generated initially using the following:

oapi-codegen -package api ./ipfs-pinning-service.yaml > api.go
