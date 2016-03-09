
function convert(str) 
	tmp = split(str)
	ret = Int[]
	for i in tmp
		push!(ret, parse(Int,i))
	end
	return ret
end

function predict(data) 
	ret = 0
	for i in data
		ret+=i
	end
	return ret
end

server = listen(8081)
while true
	sock = accept(server)
	@async begin 
		try
			while true
				str = readline(sock)
				println(str)
				data = convert(str)
				result = predict(data)
				println(sock, result)
			end
		catch err 
		end
	end
end