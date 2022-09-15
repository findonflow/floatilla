package main

import (
	"embed"
	"encoding/csv"
	"fmt"
	"log"
	"os"

	"github.com/bjartek/overflow"
	"github.com/pkg/errors"
	"github.com/psiemens/sconfig"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"go.mitsakis.org/workerpool"
)

//go:embed transactions/award_manually_many.cdc
//go:embed transactions/adminAddKeys.cdc
//go:embed flow.json
var path embed.FS

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

var cmd = &cobra.Command{
	Use: "floatilla",
	Run: func(cmd *cobra.Command, args []string) {

		addresses := readAddresses(conf.File)

		//These needs to be env because of flow json
		if os.Getenv("FLOATILLA_PRIVATE_KEY") == "" {
			os.Setenv("FLOATILLA_PRIVATE_KEY", conf.PrivateKey)
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
			o.Tx("adminAddKeys",
				overflow.WithSigner("admin"),
				overflow.WithArg("number", 100),
				overflow.WithArg("key", publicKey),
			)
		}

		batches := lo.Chunk(addresses, conf.BatchSize)
		p, err := workerpool.NewPoolSimple(conf.Workers, func(job workerpool.Job[[]string], workerID int) error {
			//		o := Overflow(WithNetwork("mainnet"), WithPrintResults())
			addresses := job.Payload

			o.Tx("award_manually_many",
				overflow.WithSigner(fmt.Sprintf("floatilla%d", workerID+1)),
				overflow.WithArg("forHost", conf.Host),
				overflow.WithArg("eventId", conf.EventID),
				overflow.WithAddresses("recipients", addresses...),
			)

			return nil
		})
		if err != nil {
			panic(err)
		}
		for _, batch := range batches {
			p.Submit(batch)
		}
		p.StopAndWait()

		log.Println("Hello world!")
	},
}

func init() {

	log.SetFlags(0)

	err := sconfig.New(&conf).
		FromEnvironment("FLOATILLA").
		BindFlags(cmd.PersistentFlags()).
		Parse()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("init")

}

func readAddresses(file string) []string {
	f, err := os.Open(conf.File)
	if err != nil {
		log.Fatal(errors.Wrapf(err, "could not open file with name %s", conf.File))
	}

	// remember to close the file at the end of the program
	defer f.Close()

	// read csv values using csv.Reader
	csvReader := csv.NewReader(f)
	data, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	addresses := []string{}
	for _, line := range data {
		addresses = append(addresses, line[0])
	}
	return addresses
}

func main() {
	cmd.Execute()
}
