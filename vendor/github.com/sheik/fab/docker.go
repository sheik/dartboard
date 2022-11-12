package fab

import (
	"fmt"
)

type ContainerObj struct {
	Name  string
	Flags string
}

func Container(name string) *ContainerObj {
	return &ContainerObj{
		Name: name,
	}
}

func ImageExists(image string) Check {
	return Check{
		Func: ImageExistsInterface,
		Args: image,
	}
}

func ImageExistsInterface(args any) bool {
	return ImageExistsFunc(args.(string))
}

// ImageExistsFunc checks to see if a docker image exists in the local registry
// takes one argument, the image name: docker.ImageExists("NameOfImage")
func ImageExistsFunc(image string) bool {
	command := fmt.Sprintf(`bash -c "if [[ \"$(docker images -q %s)\" == \"\" ]]; then exit 1; else exit 0; fi"`, image)
	return Exec(command) == nil
}

func Pull(image string) bool {
	command := "docker pull " + image
	return Exec(command) == nil
}

func (image *ContainerObj) Interactive() *ContainerObj {
	image.Flags += " -it "
	return image
}

func (image *ContainerObj) Mount(source, dest string) *ContainerObj {
	image.Flags += fmt.Sprintf("-v %s:%s", source, dest)
	return image
}

func (image *ContainerObj) Env(key, val string) *ContainerObj {
	image.Flags += fmt.Sprintf(" -e %s=%s ", key, val)
	return image
}

func (image *ContainerObj) Run(formatString string, args ...interface{}) string {
	dockerCommand := fmt.Sprintf(formatString, args...)
	return fmt.Sprintf("docker run %s --rm -v $PWD:/code %s %s", image.Flags, image.Name, dockerCommand)
}
