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

		institutions := []moov.ACHInstitution{
			{
				RoutingNumber: "021201383",
				Name:          "VALLEY NATIONAL BANK",
				Address: &moov.Address{
					AddressLine1:    "ACH DEPARTMENT 4TH FLOOR",
					City:            "WAYNE",
					StateOrProvince: "NJ",
					PostalCode:      "07470",
				},
				Contact: &moov.Contact{
					Phone: &moov.Phone{
						Number: "9733058800",
					},
				},
			},
		}
		printAchInstitutions(&buf, institutions)

		expected := strings.TrimSpace(`
Routing Number  Customer Name         Phone Number  Address
021201383       VALLEY NATIONAL BANK  9733058800    ACH DEPARTMENT 4TH FLOOR WAYNE NJ 07470
`)
		require.Equal(t, expected, strings.TrimSpace(buf.String()))

	})

	t.Run("rtp", func(t *testing.T) {
		var buf bytes.Buffer

		institutions := []moov.RTPInstitution{
			{
				Name:          "Veridian Credit Union",
				RoutingNumber: "273976369",
				Services: moov.RTPServices{
					ReceivePayments:          true,
					ReceiveRequestForPayment: true,
				},
			},
		}
		printRtpInstitutions(&buf, institutions)

		expected := strings.TrimSpace(`
Routing Number  Customer Name          Receive Payments  Receive Request for Payment
273976369       Veridian Credit Union  true              true
`)
		require.Equal(t, expected, strings.TrimSpace(buf.String()))
	})

	t.Run("wire", func(t *testing.T) {
		var buf bytes.Buffer

		institutions := []moov.WireInstitution{
			{
				RoutingNumber: "273976369",
				Name:          "VERIDIAN CREDIT UNION",
				Address:       nil,
				Services: moov.WireServices{
					FundsTransferStatus:               true,
					FundsSettlementOnlyStatus:         false,
					BookEntrySecuritiesTransferStatus: false,
				},
			},
		}
		printWireInstitutions(&buf, institutions)

		expected := strings.TrimSpace(`
Routing Number  Customer Name          Fund Transfers  Settlement Only  Book Entry Transfers  Address
273976369       VERIDIAN CREDIT UNION  true            false            false
`)
		require.Equal(t, expected, strings.TrimSpace(buf.String()))
	})
}
