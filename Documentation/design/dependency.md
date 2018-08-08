# Dependency graph

During the installation, the installer will generate a bunch of files, the dependency graph of those files can be found at [resource_dep.svg](./resource_dep.svg).
It is generated from [resource_dep.dot](./resource_dep.dot) by running
```sh
dot -Tsvg resource_dep.dot -o ./resource_dep.svg
```
