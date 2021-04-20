package arduino

import (
	"apm/project"
	"context"
	"errors"
	"fmt"
	acli "github.com/arduino/arduino-cli/cli"
	aconfig "github.com/arduino/arduino-cli/configuration"
	"github.com/arduino/arduino-cli/i18n"
	"github.com/phayes/freeport"
	"google.golang.org/grpc"
	"io"

	rpc "github.com/arduino/arduino-cli/rpc/cc/arduino/cli/commands/v1"
	//"github.com/arduino/arduino-cli/rpc/cc/arduino/cli/settings/v1"
	//"google.golang.org/grpc"
	"github.com/spf13/cobra"
	"log"
	"os"
	"strings"
)

func RunCmdInteractive(cmd *cobra.Command, command []string) error {
	err := os.Setenv("ARDUINO_LIBRARY_ENABLE_UNSAFE_INSTALL", "true")
	if err != nil {
		return err
	}
	cmd.SetArgs(command)
	err = cmd.Execute()
	if err != nil {
		return err
	}
	return nil
}

type ArduinoCli struct {
	grpcServerPort int
	cmd *cobra.Command
	client rpc.ArduinoCoreServiceClient
	grpcConn *grpc.ClientConn
	grpcInstance *rpc.Instance
}

func (c *ArduinoCli) Init() error {
	c.cmd = c.getArduinoCliCommand()
	grpcPort, err := c.startCliGrpcServer()
	if err != nil {
		return err
	}
	c.grpcServerPort = grpcPort

	// grpc connection
	conn, err := c.getGrpcConnection();
	if err != nil {
		return err
	}
	c.grpcConn = conn

	// client
	c.client = rpc.NewArduinoCoreServiceClient(conn)

	// init instance
	c.grpcInstance = c.initInstance(c.client)

	return nil
}

func (c *ArduinoCli) Destroy() {
	c.grpcConn.Close()
}

func (c *ArduinoCli) startCliGrpcServer() (int, error) {
	port, err := freeport.GetFreePort()
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		RunCmdInteractive(c.cmd, strings.Split(fmt.Sprintf("daemon --daemonize --port %d", port), " "))
	}()
	return port, nil
}

func (c *ArduinoCli) getArduinoCliCommand() *cobra.Command {
	aconfig.Settings = aconfig.Init(aconfig.FindConfigFileInArgsOrWorkingDirectory(os.Args))
	i18n.Init()
	return acli.NewCommand()
}

func (c *ArduinoCli) SearchLibrary(query string) ([]*rpc.SearchedLibrary, error) {
	maxSizeOption := grpc.MaxCallRecvMsgSize(64*10e6)
	response, err := c.client.LibrarySearch(context.Background(), &rpc.LibrarySearchRequest{
		Instance: c.grpcInstance,
		Query:    query,
	}, maxSizeOption)
	if err != nil {
		return nil, err
	}
	return response.Libraries, nil
}

func (c *ArduinoCli) initInstance(client rpc.ArduinoCoreServiceClient) *rpc.Instance {
	initRespStream, err := client.Init(context.Background(), &rpc.InitRequest{})
	if err != nil {
		log.Fatalf("Error creating server instance: %s", err)
	}

	var instance *rpc.Instance
	// Loop and consume the server stream until all the setup procedures are done.
	for {
		initResp, err := initRespStream.Recv()
		// The server is done.
		if err == io.EOF {
			break
		}

		// There was an error.
		if err != nil {
			log.Fatalf("Init error: %s", err)
		}

		// The server sent us a valid instance, let's print its ID.
		if initResp.GetInstance() != nil {
			instance = initResp.GetInstance()
		}
	}

	return instance
}

func (c *ArduinoCli) getGrpcConnection() (*grpc.ClientConn, error) {
	return grpc.Dial(fmt.Sprintf("localhost:%d", c.grpcServerPort), grpc.WithInsecure(), grpc.WithBlock())
}

func (c *ArduinoCli) InstallBoardCore(details *project.ProjectDetails) error {
	log.Println("Installing board...")
	board := details.Board

	// update board core index
	additionalArgs := ""
	if board.BoardManagerUrl != "" {
		additionalArgs = fmt.Sprintf("--additional-urls %s", board.BoardManagerUrl)
	}
	err := RunCmdInteractive(c.cmd, strings.Split(fmt.Sprintf("core update-index %s", additionalArgs), " "))
	if err != nil {
		return err
	}

	if strings.ToLower(board.Version) == "latest" {
		return RunCmdInteractive(c.cmd, strings.Split(fmt.Sprintf("core install %s:%s %s", board.Package, board.Architecture, additionalArgs), " "))
	}

	return RunCmdInteractive(c.cmd, strings.Split(fmt.Sprintf("core install %s:%s@%s %s", board.Package, board.Architecture, board.Version, additionalArgs), " "))
}

func (c *ArduinoCli) CheckDependencyVersionMismatch(dep project.ProjectDependency, details *project.ProjectDetails) error {
	for _, projectLib := range details.Dependencies {
		if projectLib.Library != dep.Library || (projectLib.Library == dep.Library && projectLib.Version != dep.Version) {
			libs, err := c.SearchLibrary(projectLib.Library)
			if err != nil {
				return err
			}
			for _, lib := range libs {
				if lib.Name == projectLib.Library {
					if strings.ToLower(projectLib.Version) == "latest" {
						projectLib.Version = lib.Latest.Version
					}
					for _, release := range lib.Releases {
						if release.Version == projectLib.Version {
							// check for dep mismatch directly
							if lib.Name == dep.Library && release.Version != dep.Version {
								return errors.New(
									fmt.Sprintf(
										"Version mismatch: %s@%s ->|<- %s@%s", dep.Library, dep.Version, lib.Name, release.Version,
									),
								)
							}

							// check for dep mismatch in found lib dependencies
							for _, dependency := range release.Dependencies {
								if dependency.VersionConstraint == "" {
									dependency.VersionConstraint = "latest"
								}
								if dependency.Name == dep.Library && dependency.VersionConstraint != dep.Version {
									return errors.New(
										fmt.Sprintf(
											"Version mismatch: %s@%s ->|<- %s -> %s@%s", dep.Library, dep.Version, projectLib.Library, dependency.Name, dependency.VersionConstraint,
											),
										)
								}
							}
						}
					}
				}
			}
		}
	}
	return nil
}

func (c *ArduinoCli) InstallDependencies(details *project.ProjectDetails) error {
	log.Println("Installing dependencies...")

	// update library index
	err := RunCmdInteractive(c.cmd, strings.Split("lib update-index", " "))
	if err != nil {
		return err
	}

	// check if we have any dependency mismatch with current libs
	for _, dep := range details.Dependencies {
		err = c.CheckDependencyVersionMismatch(dep, details)
		if err != nil {
			return err
		}
	}

	// install all library dependencies
	for _, dep := range details.Dependencies {
		// we have a library specified
		if dep.Library != "" {
			// construct library install cmd
			if dep.Version == "" {
				return errors.New("please specify a version")
			}
			lib := fmt.Sprintf("%s@%s", dep.Library, dep.Version)
			if strings.ToLower(dep.Version) == "latest" {
				lib = fmt.Sprintf("%s", dep.Library)
			}

			// run lib install
			err = RunCmdInteractive(c.cmd, []string{"lib", "install", lib})
			if err != nil {
				return err
			}
		} else {
			if dep.Git != "" && dep.Zip != "" {
				return errors.New("please specify git or zip, but NOT both")
			}

			// we have git specified
			if dep.Git != "" {
				log.Printf("Installing dependency from GIT repository: %s...\n", dep.Git)
				err = RunCmdInteractive(c.cmd, strings.Split(fmt.Sprintf("lib install --git-url %s", dep.Git), " "))
				if err != nil {
					return err
				}
			}

			// we have zip specified
			if dep.Zip != "" {
				log.Printf("Installing dependency from ZIP file: %s...\n", dep.Zip)
				err = RunCmdInteractive(c.cmd, strings.Split(fmt.Sprintf("lib install --zip-path %s", dep.Zip), " "))
				if err != nil {
					return err
				}
			}
		}

	}
	return nil
}

func (c *ArduinoCli) UninstallDependency(dep string) error {
	log.Printf("Removing dependency '%s'...\n", dep)
	return RunCmdInteractive(c.cmd, []string{"lib", "uninstall", dep})
}



