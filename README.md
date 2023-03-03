# GO NVidia FanControl

### Prepare 

```
go mod download
```

### Run

```
go run .
```

### Build

```
go build -o ~/.local/bin/gfancontrol .
```

### Install 

add the following line to your ~/.xprofile

```
$HOME/.local/bin/gfancontrol 2> ~/.gfancontrol.logs &
```
