local ReplicatedStorage = game:GetService("ReplicatedStorage")

local IMPORT_ROOT = ReplicatedStorage:FindFirstChild("go_include")

local go = {}

function go.import(str: string)
    local parts = string.split(str, "/")
    local part = IMPORT_ROOT
    for _, p in parts do
        part = part:FindFirstChild(str)
    end
    return require(part)
end

return go