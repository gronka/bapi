package main

import (
	"net/http"

	//"github.com/gocql/gocql"
	"github.com/julienschmidt/httprouter"
	//geohasher "github.com/mmcloughlin/geohash"
	"github.com/pkg/errors"
	//log "github.com/sirupsen/logrus"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/charge"
	"github.com/stripe/stripe-go/customer"
)

type DoPaymentPayload struct {
	StripeToken string  `json:"stripeToken"`
	Amount      float64 `json:"amount"`
}

// On checkout, we either look up the customer information or create the
// customer then return that information
// NOTE: we might be able to collapse this into hGetPaymentMethods
func hGetStripeCustomer(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	gibs := unpackReqIntoGibs(w, r, P_PUBLIC)
	patron := PatronFactory(&gibs)
	patron.loadOrCreateNewCustomer()

	gibs.ares.Info = ACCEPTED
	gibs.ares.Body["patron"] = patron
	gibs.encodeResponse(w)
}

// Return list of payment methods, stored on our server
func hGetPaymentMethods(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	gibs := unpackReqIntoGibs(w, r, P_PUBLIC)
	patron := PatronFactory(&gibs)
	patron.loadOrCreateNewCustomer()
	dpmList := patron.getPaymentMethods()

	gibs.ares.Info = ACCEPTED
	gibs.ares.Body["paymentMethods"] = dpmList
	gibs.encodeResponse(w)
}

func hStripeDoPayment(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	gibs := unpackReqIntoGibs(w, r, P_PUBLIC)
	var pl DoPaymentPayload
	decodeJson(r, &pl)

	email := "test@qwer.ty"
	customerParams := &stripe.CustomerParams{Email: &email}
	customerParams.SetSource(pl.StripeToken)

	newCustomer, err := customer.New(customerParams)
	if err != nil {
		panic(errors.Wrap(err, "failed to create new stripe customer"))
	}

	var amount int64 = 500
	currency := "usd"
	desc := "sample charge"
	chargeParams := &stripe.ChargeParams{
		Amount:      &amount,
		Currency:    &currency,
		Description: &desc,
		Customer:    &newCustomer.ID,
	}

	if _, err := charge.New(chargeParams); err != nil {
		panic(errors.Wrap(err, "failed to create new charge from chargeParams"))
	}

	gibs.ares.Info = ACCEPTED
	gibs.encodeResponse(w)
}
