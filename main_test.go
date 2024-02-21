package main

import (
	"bytes"
	"strings"
	"testing"

	"github.com/moovfinancial/moov-go/pkg/moov"

	"github.com/stretchr/testify/require"
)

func TestNormalizeRoutingNumber(t *testing.T) {
	input := "12345678"
	require.Equal(t, "123456780", normalizeRoutingNumber(input))

	input = "123456780"
	require.Equal(t, "123456780", normalizeRoutingNumber(input))
}

func TestPrint(t *testing.T) {
	t.Run("ach", func(t *testing.T) {
		var buf bytes.Buffer

		participants := []moov.AchParticipant{
			{
				RoutingNumber:      "021201383",
				OfficeCode:         "O",
				ServicingFRBNumber: "021001208",
				RecordTypeCode:     "1",
				Revised:            "020222",
				NewRoutingNumber:   "000000000",
				CustomerName:       "VALLEY NATIONAL BANK",
				PhoneNumber:        "9733058800",
				StatusCode:         "1",
				ViewCode:           "1",
				AchLocation: moov.AchLocation{
					Address:             "ACH DEPARTMENT 4TH FLOOR",
					City:                "WAYNE",
					State:               "NJ",
					PostalCode:          "07470",
					PostalCodeExtension: "0000",
				},
			},
		}
		printAchParticipants(&buf, participants)

		expected := strings.TrimSpace(`
Routing Number  Customer Name         Phone Number  Address
021201383       VALLEY NATIONAL BANK  9733058800    ACH DEPARTMENT 4TH FLOOR WAYNE NJ 07470
`)
		require.Equal(t, expected, strings.TrimSpace(buf.String()))

	})

	t.Run("wire", func(t *testing.T) {
		var buf bytes.Buffer

		participants := []moov.WireParticipant{
			{
				RoutingNumber:   "273976369",
				TelegraphicName: "VERIDIAN",
				CustomerName:    "VERIDIAN CREDIT UNION",
				Location: moov.WireLocation{
					City:  "",
					State: "",
				},
				FundsTransferStatus:               "Y",
				FundsSettlementOnlyStatus:         " ",
				BookEntrySecuritiesTransferStatus: "N",
				Date:                              "20141107",
			},
		}
		printWireParticipants(&buf, participants)

		expected := strings.TrimSpace(`
Routing Number  Telegraphic Name  Customer Name          Fund Transfers  Settlement Only  Book Entry Transfers  Address
273976369       VERIDIAN          VERIDIAN CREDIT UNION  Y               N                N
`)
		require.Equal(t, expected, strings.TrimSpace(buf.String()))
	})
}
