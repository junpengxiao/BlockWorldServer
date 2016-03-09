using JLD
using CUDArt
device(0)
using Knet: stack_isempty
using Knet

#set augments
modelnames = ["ModelA1.jld", "ModelA2.jld", "ModelA3.jld"]
port = 8081

# These assume one hot columns:
setrow!(x::SparseMatrixCSC,i,j)=(i>0 ? (x.rowval[j] = i; x.nzval[j] = 1) : (x.rowval[j]=1; x.nzval[j]=0); x)
setrow!(x::Array,i,j)=(x[:,j]=0; i>0 && (x[i,j]=1); x)

#convert str into Int format
function convert(str) 
    tmp = split(str)
    ret = Int[]
    for i in tmp
        push!(ret, parse(Int,i))
    end
    return ret
end

#load tradning model
function loadmodels(model1, model2, model3) 
    global net1 = load(model1, "net")
    global net2 = load(model2, "net")
    global net3 = load(model3, "net")
end

function extract(data, xrange; padding=nothing)
    inst = Any[]
    for i=1:size(data,1)
        s = vec(data[i,xrange])
        if padding != nothing
            while s[end]==padding; pop!(s); end
        end
        push!(inst,s)
    end
    return inst
end

function predict(f, data; xrange=1:79, padding=1, xvocab=326, ftype=Float32, xsparse=false)
    reset!(f)
    sentences = extract(data, xrange; padding=1)	# sentences[i][j] = j'th word of i'th sentence
    ypred = Any[]
    eos = xvocab + 1
    x = (xsparse ? sponehot : zeros)(ftype, eos, 1)
    for s in sentences
        for i = 1:length(s)
            setrow!(x, s[i], 1)
            forw(f, x, predict=false)
        end
        setrow!(x, eos, 1)
        y = forw(f, x, predict=true)
        push!(ypred, indmax(to_host(y)))
        reset!(f)
    end
    println(ypred)
end

loadmodels(modelnames[1],modelnames[2], modelnames[3])
server = listen(port)
println("listening at port: ", port)
while true
    sock = accept(server)
    @async begin 
        try
            while true
                str = readline(sock)
                println(str)
                data = convert(str)
                result = Any[]
                push!(result,predict(net1, data))
                push!(result,predict(net2, data))
                push!(result,predict(net3, data))
                println(sock, result)
            end
        catch err 
        end
    end
end