package command

import (
	"context"
	"io"
	"os"

	docker_api_types "github.com/docker/docker/api/types/container"
	docker_client "github.com/docker/docker/client"

	"github.com/CoachApplication/api"
	"github.com/CoachApplication/command"
	handler_dockercli "github.com/CoachApplication/handler-dockercli"
)

type Factory struct {
	ContainerIdentifiers []*ContainerIdentifier

	context context.Context

	cli docker_client.APIClient

	networkDriver string

	in io.ReadCloser
	out io.WriteCloser
	err io.WriteCloser
}

func NewFactory(
	ctx context.Context,
	cli docker_client.APIClient,
	networkDriver string,
	in io.ReadCloser,
	out io.WriteCloser,
	err io.WriteCloser,
) *Factory {

	if cli == nil {
		defCli, _ := handler_dockercli.DefaultClient()
		cli = docker_client.APIClient( defCli )
	}

	if ctx == nil {
		ctx = context.Background()
	}

	if in == nil {
		in = os.Stdin
	}
	if out == nil {
		out = os.Stdout
	}
	if err == nil {
		err = os.Stderr
	}

	return &Factory{
		context: ctx,
		cli: cli,
		networkDriver: networkDriver,

		in: in,
		out: out,
		err: err,
	}
}

func NewDefaultFactory() *Factory {
	return NewFactory(nil, nil, "", nil, nil, nil)
}

func (f *Factory) AddIdentifier(contId *ContainerIdentifier) error {
	f.ContainerIdentifiers = append(f.ContainerIdentifiers, contId)
	return nil
}

func (f *Factory) NewCommand(id string, ui api.Ui, usage api.Usage, privileged bool, image string, volumes []string, links []string) command.Command {
	config := docker_api_types.Config{
		//Hostname        string              // Hostname
		//Domainname      string              // Domainname
		//User            string              // User that will run the command(s) inside the container, also support user:group
		//AttachStdin     bool                // Attach the standard input, makes possible user interaction
		//AttachStdout    bool                // Attach the standard output
		//AttachStderr    bool                // Attach the standard error
		//ExposedPorts    nat.PortSet         `json:",omitempty"` // List of exposed ports
		//Tty             bool                // Attach standard streams to a tty, including stdin if it is not closed.
		//OpenStdin       bool                // Open stdin
		//StdinOnce       bool                // If true, close stdin after the 1 attached client disconnects.
		//Env             []string            // List of environment variable to set in the container
		//Cmd             strslice.StrSlice   // Command to run when starting the container
		//Healthcheck     *HealthConfig       `json:",omitempty"` // Healthcheck describes how to check the container is healthy
		//ArgsEscaped     bool                `json:",omitempty"` // True if command is already escaped (Windows specific)
		//Image           string              // Name of the image as it was passed by the operator (e.g. could be symbolic)
		//Volumes         map[string]struct{} // List of volumes (mounts) used for the container
		//WorkingDir      string              // Current directory (PWD) in the command will be launched
		//Entrypoint      strslice.StrSlice   // Entrypoint to run when starting the container
		//NetworkDisabled bool                `json:",omitempty"` // Is network disabled
		//MacAddress      string              `json:",omitempty"` // Mac Address of the container
		//OnBuild         []string            // ONBUILD metadata that were defined on the image Dockerfile
		//Labels          map[string]string   // List of labels set to this container
		//StopSignal      string              `json:",omitempty"` // Signal to stop a container
		//StopTimeout     *int                `json:",omitempty"` // Timeout (in seconds) to stop a container
		//Shell           strslice.StrSlice   `json:",omitempty"` // Shell for shell-form of RUN, CMD, ENTRYPOINT
	}
	hostConfig := docker_api_types.HostConfig{
		//Binds           []string      // List of volume bindings for this container
		//ContainerIDFile string        // File (path) where the containerId is written
		//LogConfig       LogConfig     // Configuration of the logs for this container
		//NetworkMode     NetworkMode   // Network mode to use for the container
		//PortBindings    nat.PortMap   // Port mapping between the exposed port (container) and the host
		//RestartPolicy   RestartPolicy // Restart policy to be used for the container
		//AutoRemove      bool          // Automatically remove container when it exits
		//VolumeDriver    string        // Name of the volume driver used to mount volumes
		//VolumesFrom     []string      // List of volumes to take from other container

		// Applicable to UNIX platforms
		//CapAdd          strslice.StrSlice // List of kernel capabilities to add to the container
		//CapDrop         strslice.StrSlice // List of kernel capabilities to remove from the container
		//DNS             []string          `json:"Dns"`        // List of DNS server to lookup
		//DNSOptions      []string          `json:"DnsOptions"` // List of DNSOption to look for
		//DNSSearch       []string          `json:"DnsSearch"`  // List of DNSSearch to look for
		//ExtraHosts      []string          // List of extra hosts
		//GroupAdd        []string          // List of additional groups that the container process will run as
		//IpcMode         IpcMode           // IPC namespace to use for the container
		//Cgroup          CgroupSpec        // Cgroup to use for the container
		//Links           []string          // List of links (in the name:alias form)
		//OomScoreAdj     int               // Container preference for OOM-killing
		//PidMode         PidMode           // PID namespace to use for the container
		//Privileged      bool              // Is the container in privileged mode
		//PublishAllPorts bool              // Should docker publish all exposed port for the container
		//ReadonlyRootfs  bool              // Is the container root filesystem in read-only
		//SecurityOpt     []string          // List of string values to customize labels for MLS systems, such as SELinux.
		//StorageOpt      map[string]string `json:",omitempty"` // Storage driver options per container.
		//Tmpfs           map[string]string `json:",omitempty"` // List of tmpfs (mounts) used for the container
		//UTSMode         UTSMode           // UTS namespace to use for the container
		//UsernsMode      UsernsMode        // The user namespace to use for the container
		//ShmSize         int64             // Total shm memory usage
		//Sysctls         map[string]string `json:",omitempty"` // List of Namespaced sysctls used for the container
		//Runtime         string            `json:",omitempty"` // Runtime to use with this container

		// Applicable to Windows
		//ConsoleSize [2]uint   // Initial console size (height,width)
		//Isolation   Isolation // Isolation technology of the container (e.g. default, hyperv)

		//Resources // Contains container's resources (cgroups, ulimits)
		//Mounts []mount.Mount `json:",omitempty"` // Mounts specs used by the container
		//Init *bool `json:",omitempty"`// Run a custom init inside the container, if null, use the daemon's configured settings
		//InitPath string `json:",omitempty"`// Custom init path
	}

	//for _, volume := range volumes {
	//
	//}
	//for _, link := range links {
	//
	//}

	return NewCommand(
		id ,
		ui,
		usage,
		f.context,
		f.cli,
		f.networkDriver,
		config,
		hostConfig,
		f.out,
		f.err,
		f.in,
	).Command()
}

type ContainerIdentifier interface {
	ContainerID(reference string) string
}
