# BlockWorld Server
## Introduction
The server contains 2 parts: a go server which handles in-come and out-come JSON, a julia server which handles prediction. To add new prediction model, do following things:

1. Create a go file, such as **ModelX.go** which should provide a function based on signature: `type modelProcess func(bwStruct.BWData) bwStruct.BWData`. The function accepts an input and returns an output.
2. Create a very simple julia server. **ModelA.jl** is an example. Basically it loads models into memory, and for each input sent by **ModelX.go**, it returns an output.
3. Data format transmitted between go and julia server isn't defined. You can do it based on Julia Model input format requirement.

##Launching
For now, it has to be launched manually. Use ModelA as an example

1. go to bwModel folder, execute `nohup julia ModelA.jl &`.
2. back to BlockWorldServer folder, execute `nohup go run main.go &`.
3. Now you can test servers.

Julia server use 8081 by default while go use 8080. You can use this command to send JSON to server in terminal as a debug method:

```shell
curl -H "Content-Type: application/json" -d '{"world":[{"id":1,"loc":[0.8062,0.1,-0.5769]},{"id":2,"loc":[0.4595,0.1,-0.4644]},{"id":3,"loc":[-0.5604,0.1,0.4731]},{"id":4,"loc":[-0.6564,0.1,0.764]},{"id":5,"loc":[0,0.1,0]},{"id":6,"loc":[0,0.1,-0.1667]},{"id":7,"loc":[0.6544,0.1,-0.919]},{"id":8,"loc":[-0.4754,0.1,-0.2636]},{"id":9,"loc":[-0.3238,0.1,0.0131]},{"id":10,"loc":[0.1667,0.1,0.3333]},{"id":11,"loc":[0.1667,0.1,0.1667]},{"id":12,"loc":[0.1667,0.1,0]},{"id":13,"loc":[0.1667,0.1,-0.1667]},{"id":14,"loc":[0.3333,0.1,0.6667]},{"id":15,"loc":[0.3333,0.1,0.5]},{"id":16,"loc":[0.3333,0.1,0.3333]},{"id":17,"loc":[0.5,0.1,0.6667]}],"version":1,"input":"move the adidas block directly diagonally left and below the heineken block .","error":"None"} ' http://localhost:8080/query
```
