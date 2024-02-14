### The repository module is available as a local path

Used the golang module system to control the path of the module:

- go mod edit -require=github.com/Serares/undertown_v3/repositories/repository@v0.0.0
- go mod edit -replace=github.com/Serares/undertown_v3/repositories/repository=../../../repositories/repository
