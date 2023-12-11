# CNCGO
 is an API designed for sending instructions to Computer Numerical Control (CNC) machines utilizing the GRBL firmware. It facilitates communication with the machine via USB, while the program handles the conversion of user commands into GCode machine language.

![This is an image](https://github.com/gildasgatel/CNCGO/blob/master/_data/cncgo.jpg)

 Endpoint
 * GET("/state")
 * POST("/config")
 ```
    curl --location 'localhost:8080/config' \
--header 'Content-Type: text/plain' \
--data '{
    "machine": "grbl", 
    "connection": "usb", 
    "baudrate": 115200,
    "port": "/dev/ttyACM0"
}'
```
 * POST("/command")
 ```
 curl --location 'localhost:8080/command' \
--header 'Content-Type: text/plain' \
--data '{
    "command": "pause", (move, stop, play)
    "axe":"X",
    "distance":"-10"
}'
```
 * POST("/file")
 ```
 curl --location 'localhost:8080/file' \
--header 'Content-Type: text/plain' \
--data '{
    "path": "/PATH/OF/FILE/Example.nc"
    }'
    ```

