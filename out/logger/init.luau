local abq = GO.import("github.com/intervinn/abq/transform")
local logger = require('logger')
local Logger = {}
function NewLogger()
	return logger.new()

end
function Logger.Log(l,msg)
_ = logger.msg(l, msg)

end
return {
	Logger = Logger,
	NewLogger = NewLogger
}
