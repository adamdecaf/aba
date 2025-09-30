package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/moov-io/ach"
	"github.com/moovfinancial/moov-go/pkg/moov"
)

var (
	flagAch    = flag.Bool("ach", false, "Find the routing number in the Fed ACH directory")
	flagFedNow = flag.Bool("fednow", false, "Find the routing number in the FedNOW directory")
	flagRtp    = flag.Bool("rtp", false, "Find the routing number in the RTP participant directory")
	flagWire   = flag.Bool("wire", false, "Find the routing number in the Fed Wire directory")

	flagLimit = flag.Int("limit", 1, "How many institutions to return for each rail")
)

func main() {
	flag.Parse()

	routingNumber := normalizeRoutingNumber(flag.Arg(0))

	resp, err := listRoutingNumbers(routingNumber, *flagLimit)
	if err != nil {
		fmt.Printf("ERROR: routing number lookup failed: %v\n", err) //nolint:forbidigo
		os.Exit(1)
	}

	// Print everything if no flags were provided
	all := !*flagAch && !*flagFedNow && !*flagRtp && !*flagWire

	var buf bytes.Buffer
	if *flagAch || all {
		if all {
			fmt.Fprintln(&buf, "ACH:")
		}
		printAchInstitutions(&buf, resp.Ach)
		if all {
			fmt.Fprintln(&buf, "")
		}
	}
	if *flagFedNow || all {
		if all {
			fmt.Fprintln(&buf, "FedNow:")
		}
		printFedNowInstitutions(&buf, resp.FedNow)
		if all {
			fmt.Fprintln(&buf, "")
		}
	}
	if *flagRtp || all {
		if all {
			fmt.Fprintln(&buf, "RTP:")
		}
		printRtpInstitutions(&buf, resp.Rtp)
		if all {
			fmt.Fprintln(&buf, "")
		}
	}
	if *flagWire || all {
		if all {
			fmt.Fprintln(&buf, "Wire:")
		}
		printWireInstitutions(&buf, resp.Wire)
		if all {
			fmt.Fprintln(&buf, "")
		}
	}
	fmt.Println(buf.String()) //nolint:forbidigo
}

func normalizeRoutingNumber(input string) string {
	if len(input) == 8 {
		checkDigit := ach.CalculateCheckDigit(input)
		if checkDigit >= 0 {
			input += fmt.Sprintf("%d", checkDigit)
		}
	}
	return input
}

func listRoutingNumbers(routingNumber string, limit int) (*moov.InstitutionsSearchResponse, error) {
	mc, err := moov.NewClient(moov.WithCredentials(moov.CredentialsFromEnv()))
	if err != nil {
		return nil, err
	}

	ctx := context.Background()

	// Append query params
	args := []moov.ListInstitutionsFilter{
		moov.WithInstitutionLimit(limit),
	}
	// If routingNumber is not numeric treat it as a name
	if _, err := strconv.ParseInt(routingNumber, 10, 32); err != nil {
		args = append(args, moov.WithInstitutionName(routingNumber))
	} else {
		args = append(args, moov.WithInstitutionRoutingNumber(routingNumber))
	}

	resp, err := mc.SearchInstitutions(ctx, args...)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func printAchInstitutions(buf io.Writer, participants []moov.ACHInstitution) {
	w := tabwriter.NewWriter(buf, 0, 0, 2, ' ', 0)
	defer w.Flush()

	fmt.Fprintln(w, "Routing Number\tCustomer Name\tPhone Number\tAddress")

	for _, p := range participants {
		var address string
		if p.Address != nil {
			address = fmt.Sprintf("%s %s %s %s", p.Address.AddressLine1, p.Address.City, p.Address.StateOrProvince, p.Address.PostalCode)
		}

		var contact string
		if p.Contact != nil {
			if p.Contact.Phone != nil {
				contact = p.Contact.Phone.Number
			}
		}

		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", p.RoutingNumber, p.Name, contact, address)
	}
}

func printFedNowInstitutions(buf io.Writer, participants []moov.FedNowInstitution) {
	w := tabwriter.NewWriter(buf, 0, 0, 2, ' ', 0)
	defer w.Flush()

	fmt.Fprintln(w, "Routing Number\tCustomer Name\tReceive Payments\tSend Payments\tRequest for Payment")

	for _, p := range participants {
		fmt.Fprintf(w, "%s\t%s\t%v\t%v\t%v\n",
			p.RoutingNumber, p.Name,
			p.Services.ReceivePayments,
			p.Services.SendPayments,
			p.Services.RequestForPayment,
		)
	}
}

func printRtpInstitutions(buf io.Writer, participants []moov.RTPInstitution) {
	w := tabwriter.NewWriter(buf, 0, 0, 2, ' ', 0)
	defer w.Flush()

	fmt.Fprintln(w, "Routing Number\tCustomer Name\tReceive Payments\tReceive Request for Payment")

	for _, p := range participants {
		fmt.Fprintf(w, "%s\t%s\t%v\t%v\n",
			p.RoutingNumber, p.Name,
			p.Services.ReceivePayments,
			p.Services.ReceiveRequestForPayment,
		)
	}
}

func printWireInstitutions(buf io.Writer, participants []moov.WireInstitution) {
	w := tabwriter.NewWriter(buf, 0, 0, 2, ' ', 0)
	defer w.Flush()

	fmt.Fprintln(w, "Routing Number\tCustomer Name\tFund Transfers\tSettlement Only\tBook Entry Transfers\tAddress")
	for _, p := range participants {
		var address string
		if p.Address != nil {
			address = strings.TrimSpace(fmt.Sprintf("%s %s", p.Address.City, p.Address.StateOrProvince))
		}

		fmt.Fprintf(w, "%s\t%s\t%v\t%v\t%v\t%s\n",
			p.RoutingNumber, p.Name,
			p.Services.FundsTransferStatus,
			p.Services.FundsSettlementOnlyStatus,
			p.Services.BookEntrySecuritiesTransferStatus,
			address,
		)
	}
}
