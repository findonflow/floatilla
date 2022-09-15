package main

import (
	"embed"
	"fmt"
	"log"
	"os"

	"github.com/psiemens/sconfig"
	"github.com/spf13/cobra"
)

//This file is boilerplate to set up a CLI application in go using cobra/sconfig and embed to embed files into the binary

//we embedd the flow.json  and the transactions we need
//go:embed transactions/award_manually_many.cdc
//go:embed transactions/adminAddKeys.cdc
//go:embed flow.json
var path embed.FS

//For config we use psiemens brilliant sconfig repo, IMHO a lot easier to use then plain viper/pflags
type Config struct {
	File       string `default:"recipients.csv" flag:"file,f" info:"Path to file of recipients, one address per line"`
	BatchSize  int    `default:"100" flag:"batchSize,b" info:"How many floats to award in a single batch"`
	Workers    int    `default:"100" flag:"workers,w" info:"Workers to paralell mint as"`
	EventID    uint64 `required:"need to set eventId" flag:"event,e" info:"Id of event to award"`
	Host       string `flag:"host" info:"Host to mint for if not main"`
	PrivateKey string `required:"need to set privateKey" flag:"privateKey,k" info:"privateKey to mint as"`
	Address    string `required:"need to set address" flag:"address,a" info:"Address to mint as"`
}

var conf Config

//we set up our one and only command
var cmd = &cobra.Command{
	Use:   "floatilla",
	Short: "Send a floatilla of floats with the given `eventId` to the recipient in `file`",
	Long: ` 
	- set FLOATILLA_PRIVATE_KEY and FLOATILLA_ADDRESS env variables 
	- add recipients to recipients.csv

		floatilla -e 123456

	`,
	RunE: func(cmd *cobra.Command, args []string) error {
		//we do not run logic here as it is hard to test, we send the config and the global embed.FS fs to the function so we could test it seperately
		return BatchMintFloat(conf, path)
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

	if err != nil {
		log.Fatal(err)
	}
}

//The main method simply executes the command and exits either successfully or not, note that we print errors to stderr not stdout
func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(0)
	}
	os.Exit(1)
}
