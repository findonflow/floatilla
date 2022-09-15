package main

import (
	"embed"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/bjartek/overflow"
	"github.com/pkg/errors"
	"github.com/psiemens/sconfig"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"go.mitsakis.org/workerpool"
)

//This file is boilerplate to set up a CLI application in go using cobra/sconfig and embed to embed files into the binary

// we embedd the flow.json  and the transactions we need
//
//go:embed transactions/award_manually_many.cdc
//go:embed transactions/adminAddKeys.cdc
//go:embed flow.json
var path embed.FS

// For config we use psiemens brilliant sconfig repo, IMHO a lot easier to use then plain viper/pflags
type Config struct {
	File        string `default:"recipients.csv" flag:"file,f" info:"Path to file of recipients, one address per line"`
	BatchSize   int    `default:"100" flag:"batchSize,b" info:"How many floats to award in a single batch"`
	Workers     int    `default:"100" flag:"workers,w" info:"Workers to paralell mint as"`
	Host        string `flag:"host" info:"Host to mint for if not main"`
	Private_Key string `flag:"private_key,k" info:"REQUIRED: privateKey to sign as, recommend setting as env var FLOATILLA_PRIVATE_KEY"`
	Address     string `flag:"address,a" info:"REQUIRED: Address to mint as, recomend setting as env var FLOATILLA_ADDRESS"`
}

var conf Config

// we set up our one and only command
var cmd = &cobra.Command{
	Use:   "floatilla <eventId>",
	Short: "Send a floatilla of floats with the given `eventId` to the recipient in `file`",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("You need to send in floatEventID as a valid number")
		}

		eventID, err := strconv.ParseUint(args[0], 10, 64)
		if err != nil {
			return fmt.Errorf("You need to send in floatEventID as a valid number")
		}

		if conf.Private_Key == "" || conf.Address == "" {
			return fmt.Errorf("You need to set FLOATILLA_ADDRESS|FLOATILLA_PRIVATE_KEY or call this binary with -k -a flags")
		}
		fmt.Println(eventID)

		addresses, err := readAddresses(conf.File)
		if err != nil {
			return err
		}

		//These needs to be env because of flow json
		if os.Getenv("FLOATILLA_PRIVATE_KEY") == "" {
			os.Setenv("FLOATILLA_PRIVATE_KEY", conf.Private_Key)
		}
		if os.Getenv("FLOATILLA_ADDRESS") == "" {
			os.Setenv("FLOATILLA_ADDRESS", conf.Address)
		}

		o := overflow.Overflow(
			overflow.WithNetwork("mainnet"),
			overflow.WithBasePath(""),
			overflow.WithPrintResults(),
			overflow.WithEmbedFS(path),
		)

		account, err := o.GetAccount("admin")
		publicKey := account.Keys[0].PublicKey.String()
		if conf.Host == "" {
			host := o.Address("admin")
			conf.Host = host
		}

		//should we configure this with a key somehow?
		if len(account.Keys) == 1 {
			result := o.Tx("adminAddKeys",
				overflow.WithSigner("admin"),
				overflow.WithArg("number", 100),
				overflow.WithArg("key", publicKey),
			)
			if result.Err != nil {
				return result.Err
			}
		}

		batches := lo.Chunk(addresses, conf.BatchSize)
		p, err := workerpool.NewPoolSimple(conf.Workers, func(job workerpool.Job[[]string], workerID int) error {
			addresses := job.Payload

			o.Tx("award_manually_many",
				overflow.WithSigner(fmt.Sprintf("floatilla%d", workerID+1)),
				overflow.WithArg("forHost", conf.Host),
				overflow.WithArg("eventId", eventID),
				overflow.WithAddresses("recipients", addresses...),
			)

			return nil
		})
		if err != nil {
			return err
		}
		for _, batch := range batches {
			p.Submit(batch)
		}
		p.StopAndWait()

		return nil
	},
}

func init() {
	//We do not care about timestamps in logs so we just disable that
	log.SetFlags(0)

	//set up sconfig to read using the FLOATILLA prefix
	err := sconfig.New(&conf).
		FromEnvironment("FLOATILLA").
		BindFlags(cmd.PersistentFlags()).
		Parse()

	help := cmd.Flags().Changed("help")
	if err != nil && !help {
		fmt.Println("Required fields are not set")
		cmd.Help()
		log.Fatal(err)
	}
}

// The main method simply executes the command and exits either successfully or not, note that we print errors to stderr not stdout
func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func readAddresses(file string) ([]string, error) {
	f, err := os.Open(conf.File)
	if err != nil {
		return nil, errors.Wrapf(err, "could not open file with name %s", conf.File)
	}

	// remember to close the file at the end of the program
	defer f.Close()

	// read csv values using csv.Reader
	csvReader := csv.NewReader(f)
	data, err := csvReader.ReadAll()
	if err != nil {
		return nil, err
	}

	addresses := []string{}
	for _, line := range data {
		addresses = append(addresses, line[0])
	}
	return addresses, nil
}
