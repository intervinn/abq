
function append(slice, ...)
    local tbl = table.clone(slice)
    for _,elem in {...} do
        table.insert(tbl, elem)
    end
    return tbl
end

function cap(tbl)
    return #tbl
end

function close(chan)

end

function clear(tbl)
    table.clear(tbl)
end

function delete(tbl, key)
    if tbl == nil then return end
    if not tbl[key] then return end
    tbl[key] = nil
end

function len(tbl)
    return #tbl
end

function make(t, ...)

end

function max(...)

end

function min(...)

end

function new(t)

end

function panic(v)
    
end

function _print(...)

end

function println(...)

end

function real(c)

end

function recover()

end


return {
    append = append,
    cap = cap,
    close = close,
    clear = clear,
    delete = delete,
    len = len,
    make = make,
    len = len,
    max = max,
    min = min,
    new = new,
    panic = panic,
    print = _print,
    println = println,
    real = real,
    recover = recover
}