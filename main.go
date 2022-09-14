package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"

	. "github.com/bjartek/overflow"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"go.mitsakis.org/workerpool"
)

func main() {

	arg := os.Args[1:]
	batchSize := 100
	workers := 10

	if len(arg) < 3 {
		log.Fatal("send in arguments <file> <eventId> (host)")
	}

	file := arg[0]

	floatEventId, err := strconv.Atoi(arg[1])
	if err != nil {
		log.Fatal(err)
	}

	host := ""
	if len(arg) > 2 {
		host = arg[2]
	}

	f, err := os.Open(file)
	if err != nil {
		log.Fatal(errors.Wrapf(err, "could not open file with name %s", file))
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
	o := Overflow(WithNetwork("mainnet"), WithPrintResults())
	if host == "" {
		host = o.Address("admin")
	}

	batches := lo.Chunk(addresses, batchSize)
	p, err := workerpool.NewPoolSimple(workers, func(job workerpool.Job[[]string], workerID int) error {
		//		o := Overflow(WithNetwork("mainnet"), WithPrintResults())
		p := job.Payload

		o.Tx("award_manually_many",
			WithSigner(fmt.Sprintf("floatilla%d", workerID+1)),
			WithArg("forHost", host),
			WithArg("eventId", floatEventId),
			WithAddresses("recipients", p...),
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

}
