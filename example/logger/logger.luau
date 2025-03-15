
local logger = {}

function logger.new()
    return setmetatable({}, logger)
end

function logger:log(msg)
    print(msg)
end

return logger