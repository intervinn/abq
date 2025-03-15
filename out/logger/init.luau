function NewLogger()
	return logger.new()
end
function Logger.Log(l,msg)
_ = logger.msg(l, msg)
end
local logger = require(\"logger\")local Logger = {}
