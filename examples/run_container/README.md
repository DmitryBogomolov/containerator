# run_container

Shows usage of [containerator](../../README.md) function that creates and starts docker containers.

```bash
./run_container \
    --image my-image:1 \
    --name my-container \
    --network my-network \
    --restart always \
    --volume /tmp:/usr/app \
    --port 50001:3000 \
    --env A=1 --env B=2 --env C=3
```
