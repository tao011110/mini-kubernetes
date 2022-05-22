package image_factory

import (
	"fmt"
	"mini-kubernetes/tools/def"
	"mini-kubernetes/tools/util"
)

func MakeFunctionalImage(function *def.Function) {
	pyString := util.ReadFile(function.Function)
	requirementsString := util.ReadFile(function.Requirements)
	cmdWritePy := EchoFactory(pyString, def.PyHandlerPath)
	cmdWriteRequirements := EchoFactory(requirementsString, def.RequirementsPath)
	imageName := fmt.Sprintf("image_%s_%d", function.Name, function.Version)
	function.Image = imageName
	ImageFactory(def.PyFunctionTemplateImage, imageName, []string{cmdWritePy, cmdWriteRequirements, def.PyFunctionPrepareCmd})
}
