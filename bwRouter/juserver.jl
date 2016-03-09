@async begin
	server = listen(2001)
	while true
		sock = accept(server)
		@async while isopen(sock)
			data = readlm(sock,Int)
			println(sock, length(data))
			for i in data 
				print(sock, i, ', ')
			end
		end
	end
end