## examples

Examples provide a basic demonstration of the features of the library, you can
quickly get started with the examples by running:

> ❗️ at the first you should install `protoc-gen-fieldmask`, if you have no idea
> 
> about how to install it, please refer to the [installation guide](../README.md#installation).

```sh
make gen-normal && go test . -run='^Test_FieldMask_.*$' -count=1 -v 
```

### Samples

- [Masking gRPC response fields](./examples/grpc-masked-response/README.md)
- [Incremental updating](./examples/grpc-increase-update/README.md) 
