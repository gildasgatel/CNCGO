# CNCGO
 est une API permettant d'envoyer des instructions aux machines à commande numérique (CNC) utilisant le firmware GRBL. La communiquation avec la machine ce fait en USB et le programme est en charge de convertir les commandes utilisateur en langage machine GCode.

![This is an image](https://github.com/gildasgatel/CNCGO/blob/master/_data/cncgo.jpg

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

