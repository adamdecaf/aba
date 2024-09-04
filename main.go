package main

import (
	"cmp"
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/moov-io/ach"
	"github.com/moovfinancial/moov-go/pkg/moov"
)

var (
	flagAch  = flag.Bool("ach", false, "Find the routing number in the Fed ACH directory")
	flagWire = flag.Bool("wire", false, "Find the routing number in the Fed Wire directory")
)

func main() {
	flag.Parse()

	routingNumber := normalizeRoutingNumber(flag.Arg(0))

	resp, err := listRoutingNumbers(routingNumber)
	if err != nil {
		fmt.Printf("ERROR: routing number lookup failed: %v\n", err) //nolint:forbidigo
		os.Exit(1)
	}

	switch {
	case *flagAch:
		printAchParticipants(os.Stdout, routingNumber, resp.AchParticipants)
	case *flagWire:
		printWireParticipants(os.Stdout, resp.WireParticipants)
	default:
		fmt.Println(routingNumber) //nolint:forbidigo
	}
}

func normalizeRoutingNumber(input string) string {
	if len(input) == 8 {
		checkDigit := ach.CalculateCheckDigit(input)
		input += fmt.Sprintf("%d", checkDigit)
	}
	return input
}

func listRoutingNumbers(routingNumber string) (*moov.FinancialInstitutions, error) {
	mc, err := moov.NewClient(moov.WithCredentials(moov.CredentialsFromEnv()))
	if err != nil {
		return nil, err
	}

	ctx := context.Background()

	var rail moov.Rail
	switch {
	case *flagAch:
		rail = moov.RailAch
	case *flagWire:
		rail = moov.RailWire
	default:
		return nil, errors.New("no rail specified")
	}

	resp, err := mc.ListInstitutions(ctx, rail,
		moov.WithInstitutionRoutingNumber(routingNumber),
		moov.WithInstitutionLimit(1),
	)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func printAchParticipants(buf io.Writer, routingNumber string, participants []moov.AchParticipant) {
	w := tabwriter.NewWriter(buf, 0, 0, 2, ' ', 0)
	defer w.Flush()

	fmt.Fprintln(w, "Routing Number\tCustomer Name\tPhone Number\tAddress")

	// If nothing was found at least show the routing number
	if len(participants) == 0 {
		if len(routingNumber) == 8 {
			routingNumber = fmt.Sprintf("%s%d", routingNumber, ach.CalculateCheckDigit(routingNumber))
		}
		participants = append(participants, moov.AchParticipant{
			RoutingNumber: routingNumber,
			CustomerName:  "Unknown",
		})
	}

	for _, p := range participants {
		location := fmt.Sprintf("%s %s %s %s",
			p.AchLocation.Address, p.AchLocation.City, p.AchLocation.State, p.AchLocation.PostalCode)

		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", p.RoutingNumber, p.CustomerName, p.PhoneNumber, location)
	}
}

func printWireParticipants(buf io.Writer, participants []moov.WireParticipant) {
	w := tabwriter.NewWriter(buf, 0, 0, 2, ' ', 0)
	defer w.Flush()

	fmt.Fprintln(w, "Routing Number\tTelegraphic Name\tCustomer Name\tFund Transfers\tSettlement Only\tBook Entry Transfers\tAddress")
	for _, p := range participants {
		location := strings.TrimSpace(fmt.Sprintf("%s %s", p.Location.City, p.Location.State))

		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\t%s\n",
			p.RoutingNumber, p.TelegraphicName, p.CustomerName,
			cmp.Or(strings.TrimSpace(p.FundsTransferStatus), "N"),
			cmp.Or(strings.TrimSpace(p.FundsSettlementOnlyStatus), "N"),
			cmp.Or(strings.TrimSpace(p.BookEntrySecuritiesTransferStatus), "N"),
			location)
	}
}
