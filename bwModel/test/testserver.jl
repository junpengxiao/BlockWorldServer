port = 8081
server = listen(port)
println("listening at port: ", port)
while true
    sock = accept(server)
    @async begin 
        try
            while true
                str = readline(sock)
                println(str)
                println(sock, "1 2 5")
                println("Predict Finished")
            end
        catch err 
        end
    end
end