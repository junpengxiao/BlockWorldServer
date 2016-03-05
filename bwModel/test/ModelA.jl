using ArgParse
using JLD
using CUDArt
device(0)
using Knet: stack_isempty
using Knet

function main(args)
    s = ArgParseSettings()
    s.exc_handler=ArgParse.debug_handler
    @add_arg_table s begin
        ("--datafile"; default="test.data")
        ("--loadfile"; help="initialize model from file"; default="ModelA.jld")
        ("--target"; arg_type=Int; default=1; help="which target to predict: 1:source,2:target,3:direction")
        ("--nx"; arg_type=Int; default=79; help="number of input columns in data")
        ("--ny"; arg_type=Int; default=3; help="number of target columns in data")
        ("--batchsize"; arg_type=Int; default=10; help="minibatch size")
        ("--xvocab"; arg_type=Int; default=326; help="vocab size for input columns (all columns assumed equal)")
        ("--yvocab"; arg_type=Int; nargs='+'; default=[20,20,8]; help="vocab sizes for target columns (all columns assumed independent)")
        ("--xsparse"; action = :store_true; help="use sparse inputs, dense arrays used by default")
        ("--ftype"; default = "Float32"; help="floating point type to use: Float32 or Float64")
    end
    isa(args, AbstractString) && (args=split(args))
    o = parse_args(args, s; as_symbols=true); println(o)
    o[:ftype] = eval(parse(o[:ftype]))
    Knet.gpu(false)

    # Read data files: 6003x82, 855x82
    #global rawdata = map(f->readdlm(f,Int), o[:datafiles])
    #Read predict data file
    global rawdata = readdlm("test.data",Int)
    println(rawdata)


    # Minibatch data: data[1]:train, data[2]:dev
    xrange = 1:o[:nx]
    yrange = (o[:nx] + o[:target]):(o[:nx] + o[:target])
    yvocab = o[:yvocab][o[:target]]
    global data = minibatch(rawdata, xrange, yrange, o[:batchsize]; xvocab=o[:xvocab], yvocab=yvocab, ftype=o[:ftype], xsparse=o[:xsparse])
    #global data = map(rawdata) do d
    #    minibatch(d, xrange, yrange, o[:batchsize]; xvocab=o[:xvocab], yvocab=yvocab, ftype=o[:ftype], xsparse=o[:xsparse])
    #end

    # Load or create the model:
    global net = load("ModelA.jld", "net")
    @date devpred = predict(net, rawdata; xrange=xrange, xvocab=o[:xvocab], ftype=o[:ftype], xsparse=o[:xsparse])
    println(devpred)
end

### Minibatched data format:
# data is an array of (x,y,mask) triples
# x[xvocab+1,batchsize] contains one-hot word columns for the n'th word of batchsize sentences
# xvocab+1=eos is used for end-of-sentence
# sentences in a batch are padded at the beginning and get an eos at the end
# mask[batchsize] indicates whether i'th column of x is padding or not
# y is nothing until the very last token of a sentence batch
# y[yvocab,batchsize] contains one-hot target columns with the last token (eos) of a sentence batch

function predict(f, data; xrange=1:79, padding=1, xvocab=326, ftype=Float32, xsparse=false)
    reset!(f)
    sentences = extract(data, xrange; padding=1)    # sentences[i][j] = j'th word of i'th sentence
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

function minibatch(data, xrange, yrange, batchsize; o...)
    x = extract(data, xrange; padding=1)    # x[i][j] = j'th word of i'th sentence
    y = extract(data, yrange)                   # y[i][j] = j'th class of i'th sentence
    s = sortperm(x, by=length)
    batches = Any[]
    for i=1:batchsize:length(x)
        j=min(i+batchsize-1,length(x))
        xx,yy = x[s[i:j]],y[s[i:j]]
        batchsentences(xx, yy, batches; o...)
    end
    return batches
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

function batchsentences(x, y, batches; xvocab=326, yvocab=20, ftype=Float32, xsparse=false)
    @assert maximum(map(maximum,x)) <= xvocab
    @assert maximum(map(maximum,y)) <= yvocab
    eos = xvocab + 1
    batchsize = length(x)                       # number of sentences in batch
    maxlen = maximum(map(length,x))
    for t=1:maxlen+1                            # pad sentences in the beginning and add eos at the end
        xbatch = (xsparse ? sponehot : zeros)(ftype, eos, batchsize)
        mask = zeros(Cuchar, batchsize)         # mask[i]=0 if xbatch[:,i] is padding
        for s=1:batchsize                       # set xbatch[word][s]=1 if x[s][t]=word
            sentence = x[s]
            position = t - maxlen + length(sentence)
            if position < 1
                mask[s] = 0
            elseif position <= length(sentence)
                word = sentence[position]
                setrow!(xbatch, word, s)
                mask[s] = 1
            elseif position == 1+length(sentence)
                word = eos
                setrow!(xbatch, word, s)
                mask[s] = 1
            else
                error("This should not be happening")
            end
        end
        if t <= maxlen
            ybatch = nothing
        else
            ybatch = zeros(ftype, yvocab, batchsize)
            for s=1:batchsize
                answer = y[s][1]
                setrow!(ybatch, answer, s)
            end
        end
        push!(batches, (xbatch, ybatch, mask))
    end
end

# These assume one hot columns:
setrow!(x::SparseMatrixCSC,i,j)=(i>0 ? (x.rowval[j] = i; x.nzval[j] = 1) : (x.rowval[j]=1; x.nzval[j]=0); x)
setrow!(x::Array,i,j)=(x[:,j]=0; i>0 && (x[i,j]=1); x)

!isinteractive() && main(ARGS)
