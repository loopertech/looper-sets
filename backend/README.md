# Backend

## Local development
### EdgeDB Setup
1. Install
2. CD to backend
3. `edgedb project init`
4. `edgedb migrate`

### Using Air
1. Install air locally using `go install github.com/cosmtrek/air@latest`. Make sure your GO paths are setup correctly.
2. Change directory to `./backend`.
3. Run `air -c ./air/backend.air.toml` to start the backend with Air.