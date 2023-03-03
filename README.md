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
go build -o /usr/local/bin/gfancontrol .
```

### Install 

add the following line to your ~/.xprofile

```
gfancontrol 2> ~/.gfancontrol.logs &
```
