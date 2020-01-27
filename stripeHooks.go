package main

import (
	"html/template"
	"net/http"

	//"github.com/gocql/gocql"
	"github.com/julienschmidt/httprouter"
	//geohasher "github.com/mmcloughlin/geohash"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/paymentmethod"
	"github.com/stripe/stripe-go/setupintent"
	//"github.com/stripe/stripe-go/customer"
)

// First step of adding a new credit card
func hAddCreditCardScreen(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	gibs := unpackReqIntoGibs(w, r, P_PUBLIC)
	patron := PatronFactory(&gibs)
	patron.loadOrCreateNewCustomer()
	sessionId := patron.createCheckoutSession()

	tmpl, err := template.ParseFiles("templates/checkoutStep01.html")
	if err != nil {
		panic(errors.Wrap(err, "failed to parse template"))
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tmpl.Execute(w, sessionId)
}

// TODO: attach to webhook at https://dashboard.stripe.com/test/webhooks
// NOTE: this route is called by stripe, so has no gibs user data
func hAddCreditCardWebhook(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	//gibs := unpackReqIntoGibs(w, r, P_PUBLIC)
	var pl StripeCheckoutSessionCompletedPayload
	decodeJson(r, &pl)

	// get setup_intent object
	setupIntent, err := setupintent.Get(pl.Data.Object.SetupIntentId, nil)
	if err != nil {
		panic(errors.Wrap(err, "failed to get setup intent from stripe"))
	}
	log.Info(setupIntent)

	// attach payment_method id to customer
	stripeId, expireTime := cqlStripeIdFromCheckoutSessionId(pl.Id)

	var attachedPaymentMethod *stripe.PaymentMethod
	var err2 error
	if expireTime < timeNowMilli() {
		params := &stripe.PaymentMethodAttachParams{
			Customer: stripe.String(stripeId),
		}
		attachedPaymentMethod, err2 = paymentmethod.Attach(setupIntent.PaymentMethod.ID, params)
		if err2 != nil {
			panic(errors.Wrap(err2, "failed to attach payment_method to customer"))
		}
		cqlDeleteCheckoutSessionId(pl.Id)
	} else {
		// TODO: send notification to user that the checkout session expired
		cqlDeleteCheckoutSessionId(pl.Id)
		return
	}

	// save payment method info in our DB
	dpm := destitutePMFromStripePM(*attachedPaymentMethod, stripeId)
	dpm.save()

	// TODO: send notification to user app that new card has been added; refresh page of card list
}

type StripeCheckoutSessionCompletedPayload struct {
	Id   string                             `json:"id"`
	Data StripeCheckoutSessionCompletedData `json:"data"`
}

type StripeCheckoutSessionCompletedData struct {
	Object StripeCheckoutSessionCompletedObject `json:"object"`
}

type StripeCheckoutSessionCompletedObject struct {
	CheckoutSessionId string `json:"id"`
	SetupIntentId     string `json:"setup_intent"`
}

type StripeSetupIntent struct {
	Id              string `json:"id"`
	PaymentMethodId string `json:"payment_method"`
}

// example checkout.session.completed object
//{
//"id": "evt_1Ep24XHssDVaQm2PpwS19Yt0",
//"object": "event",
//"api_version": "2019-03-14",
//"created": 1561420781,
//"data": {
//"object": {
//"id": "cs_test_MlZAaTXUMHjWZ7DcXjusJnDU4MxPalbtL5eYrmS2GKxqscDtpJq8QM0k",
//"object": "checkout.session",
//"billing_address_collection": null,
//"cancel_url": "https://example.com/cancel",
//"client_reference_id": null,
//"customer": null,
//"customer_email": null,
//"display_items": [],
//"mode": "setup",
//"setup_intent": "seti_1EzVO3HssDVaQm2PJjXHmLlM",
//"submit_type": null,
//"subscription": null,
//"success_url": "https://example.com/success"
//}
//},
//"livemode": false,
//"pending_webhooks": 1,
//"request": {
//"id": null,
//"idempotency_key": null
//},
//"type": "checkout.session.completed"
//}
