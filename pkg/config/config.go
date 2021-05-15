package config

import (
	"fmt"
	"io"
	"os"

	"github.com/urfave/cli/v2"
)

var (
	NullClient  = Client{}
	NullOptions = Options{}
)

// Options takes the build information and provides the CLI with correct information
type Options struct {
	Version  string
	Revision string
	Created  string
	// DisableExitOnHelp configures if --help should exit or not, used with helpPrinter()
	DisableExitOnHelp bool
}

// Client struct
type Client struct {
	Address           string
	Port              int
	disableExitOnHelp bool
	cliReader         io.Reader
	cliWriter         io.Writer
	cliErrWriter      io.Writer
	version           string
	revision          string
	created           string
}

// NewClient returns the Client or error
func NewClient(opts Options) (Client, error) {
	client := newClient(opts)
	generatedCfg, err := client.generateConfig(os.Args)
	if err != nil {
		return NullClient, err
	}

	return generatedCfg, nil
}

// GenerateMarkdown creates a markdown file with documentation for the application
func GenerateMarkdown(filePath string) error {
	client := newClient(NullOptions)
	err := client.generateMarkdown(filePath)
	return err
}

func newClient(opts Options) *Client {
	return &Client{
		disableExitOnHelp: opts.DisableExitOnHelp,
		cliReader:         os.Stdin,
		cliWriter:         os.Stdout,
		cliErrWriter:      os.Stderr,
		version:           opts.Version,
		revision:          opts.Revision,
		created:           opts.Created,
	}
}

func (client *Client) generateConfig(args []string) (Client, error) {
	app := client.newCLIApp()

	err := app.Run(args)
	if err != nil {
		return NullClient, err
	}

	return *client, nil
}

func (client *Client) generateMarkdown(filePath string) error {
	app := client.newCLIApp()

	md, err := app.ToMarkdown()
	if err != nil {
		return err
	}

	err = os.WriteFile(filePath, []byte(md), 0666) // #nosec
	return err
}

func (client *Client) setConfig(cfg Client) {
	client.Address = cfg.Address
	client.Port = cfg.Port
}

func (client *Client) setIO(reader io.Reader, writer io.Writer, errWriter io.Writer) {
	client.cliReader = reader
	client.cliWriter = writer
	client.cliErrWriter = errWriter
}

func (client *Client) newCLIApp() *cli.App {
	cli.VersionPrinter = client.versionHandler
	cli.HelpPrinter = client.helpPrinter

	app := &cli.App{
		Name:    "opa-bundle-api",
		Usage:   "Open Policy Agent Bundle API",
		Version: client.version,
		Flags:   client.newCLIFlags(),
		Action:  client.newCLIAction,
	}

	app.Writer = client.cliWriter
	app.ErrWriter = client.cliErrWriter
	app.Reader = client.cliReader

	return app
}

func (client *Client) newCLIAction(c *cli.Context) error {
	err := client.setConfigFromCLI(c)
	return err
}

func (client *Client) newCLIFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:     "address",
			Usage:    "The listener address",
			Required: false,
			EnvVars:  []string{"ADDRESS"},
			Value:    "0.0.0.0",
		},
		&cli.IntFlag{
			Name:     "port",
			Usage:    "The listener port",
			Required: false,
			EnvVars:  []string{"PORT"},
			Value:    8080,
		},
	}
}

func (client *Client) setConfigFromCLI(cli *cli.Context) error {
	newCfg := Client{
		Address: cli.String("address"),
		Port:    cli.Int("port"),
	}

	client.setConfig(newCfg)

	return nil
}

func (client *Client) versionHandler(c *cli.Context) {
	fmt.Printf("version=%s revision=%s created=%s\n", client.version, client.revision, client.created)
	os.Exit(0)
}

// helpPrinter uses the default HelpPrinterCustom() but adds an os.Exit(0)
func (client *Client) helpPrinter(out io.Writer, templ string, data interface{}) {
	cli.HelpPrinterCustom(out, templ, data, nil)
	if !client.disableExitOnHelp {
		os.Exit(0)
	}
}
