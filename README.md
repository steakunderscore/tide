# Tide

Go client for [pixelflut](https://github.com/defnull/pixelflut)

Based off of code I wrote at 36c3, because why not!

## Running

use `--help` for more flags


Printing an image:
```
go run tide.go screen.go pixel.go image.go --address 192.168.8.195:1337 --image ~/Downloads/rickroll.jpg
```

or writing a solid colour on the screen

```
go run tide.go screen.go pixel.go image.go --address 192.168.8.195:1337 --rgb-colour FF0000
```

Have fun!
