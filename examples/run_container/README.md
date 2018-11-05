# run_container

Shows usage of [containerator](../../README.md) `RunContainer` function.

```bash
./run_container \
    --image my-image:1 --name my-container \
    --volume /tmp:/usr/app \
    --port 50001:3000 \
    --restart always
```
