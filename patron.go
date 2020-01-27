package main

import (
	"github.com/gocql/gocql"
	"github.com/pkg/errors"
	//log "github.com/sirupsen/logrus"
	"github.com/stripe/stripe-go"
	//"github.com/stripe/stripe-go/charge"
	"github.com/stripe/stripe-go/checkout/session"
	"github.com/stripe/stripe-go/customer"
	//"github.com/stripe/stripe-go/paymentmethod"
)

//type Payment struct {
//Methods []PaymentMethod
//}

type Patron struct {
	UserUuid gocql.UUID `json:"userUuid"`
	//IsSaved bool `json:"isSaved"`
	StripeId string           `json:"stripeId"`
	Customer *stripe.Customer `json:"customer"`
}

type CheckoutSession struct {
	SessionId string
}

func PatronFactory(gibs *Gibs) (patron Patron) {
	patron.UserUuid = gibs.UserUuid
	return
}

func (patron *Patron) loadOrCreateNewCustomer() {
	if patron.isSaved() {
		patron.loadCustomer()
	} else {
		patron.createNewCustomer()
		patron.loadCustomer()
	}
	return
}

func (patron *Patron) isSaved() bool {
	cql := `SELECT stripe_id FROM user_has_stripe_id WHERE user_uuid = ?`
	iter := s.Query(cql, patron.UserUuid).Iter()
	iter.Scan(&patron.StripeId)
	if err := iter.Close(); err != nil {
		panic(errors.Wrap(err, "failed to close iter"))
	}

	if len(patron.StripeId) > 3 {
		return true
	} else {
		return false
	}
}

func (patron *Patron) createNewCustomer() {
	user := UserFromUserUuid(patron.UserUuid)
	// Create on Stripe
	params := &stripe.CustomerParams{
		Description: stripe.String(`{"user_uuid":"` + patron.UserUuid.String() + `"}`),
		Phone:       stripe.String(user.IsoPhone),
		Email:       stripe.String(user.Email),
	}
	cus, err := customer.New(params)
	if err != nil {
		panic(errors.Wrap(err, "error while creating stripe customer"))
	}

	patron.StripeId = cus.ID
	// Create in Kapi tables
	cql := `INSERT INTO  user_has_stripe_id (user_uuid, stripe_id) VALUES (?, ?)`
	if err := s.Query(cql, patron.UserUuid, patron.StripeId).Exec(); err != nil {
		panic(errors.Wrap(err, "error while inserting"))
	}
}

func (patron *Patron) loadCustomer() {
	var err error
	patron.Customer, err = customer.Get(patron.StripeId, nil)
	if err != nil {
		panic(errors.Wrap(err, "error while creating stripe customer"))
	}
}

func (patron *Patron) createCheckoutSession() string {
	uuid, _ := gocql.RandomUUID()
	sessionId := uuid.String()
	successUrl := Conf.stripePrefix + "/success?session_id=" + sessionId
	cancelUrl := Conf.stripePrefix + "/cancel"

	params := &stripe.CheckoutSessionParams{
		PaymentMethodTypes: stripe.StringSlice([]string{
			"card",
		}),
		Mode:       stripe.String(string(stripe.CheckoutSessionModeSetup)),
		SuccessURL: stripe.String(successUrl),
		CancelURL:  stripe.String(cancelUrl),
	}
	sesh, err := session.New(params)
	if err != nil {
		panic(errors.Wrap(err, "error while creating stripe checkout session"))
	}

	// save this checkout session
	// checkout session expires after 30 minutes = 1.8e6
	expireTime := timeNowMilli() + 1.8e6
	stmt := `INSERT INTO  checkout_sessions_owned_by_stripe_id (checkout_session_id, stripe_id, expire_time) VALUES (?, ?, ?)`
	if err := s.Query(stmt, sessionId, patron.StripeId, expireTime).Exec(); err != nil {
		panic(errors.Wrap(err, "error while inserting"))
	}

	return sesh.ID
}

func cqlStripeIdFromCheckoutSessionId(checkoutSessionId string) (stripeId string, expireTime int64) {
	stmt := `SELECT stripe_id, expire_time FROM checkout_sessions_owned_by_stripe_id WHERE checkout_session_id=?`
	iter := s.Query(stmt, checkoutSessionId).Iter()
	iter.Scan(&stripeId, &expireTime)
	return
}

func cqlDeleteCheckoutSessionId(checkoutSessionId string) {
	stmt := `DELETE FROM checkout_sessions_owned_by_stripe_id WHERE checkout_session_id=?`
	if err := s.Query(stmt, checkoutSessionId).Exec(); err != nil {
		panic(errors.Wrap(err, "error while deleting"))
	}
	return
}

func destitutePMFromStripePM(spm stripe.PaymentMethod, stripeId string) (dpm DestitutePaymentMethod) {
	dpm.StripeId = stripeId
	dpm.PaymentMethodId = spm.ID
	dpm.Brand = spm.Card.Brand
	dpm.ExpMonth = spm.Card.ExpMonth
	dpm.ExpYear = spm.Card.ExpYear
	dpm.Last4 = spm.Card.Last4
	return
}

func (dpm *DestitutePaymentMethod) save() {
	stmt := `INSERT INTO user_has_destitute_payment_methods (stripe_id, payment_method_id, brand, exp_month, exp_year, last4) VALUES (?, ?, ?, ?, ?, ?)`
	if err := s.Query(
		stmt,
		dpm.StripeId,
		dpm.PaymentMethodId,
		dpm.Brand,
		dpm.ExpMonth,
		dpm.ExpYear,
		dpm.Last4,
	).Exec(); err != nil {
		panic(errors.Wrap(err, "error while inserting"))
	}
	return
}

func (patron *Patron) getPaymentMethods() []DestitutePaymentMethod {
	stmt := `SELECT * FROM user_has_destitute_payment_methods WHERE stripe_id=?`
	iter := s.Query(stmt, patron.StripeId).Iter()
	var dpmRow DestitutePaymentMethod
	var dpmList []DestitutePaymentMethod

	for {
		row := map[string]interface{}{
			"stripe_id":         &dpmRow.StripeId,
			"payment_method_id": &dpmRow.PaymentMethodId,
			"brand":             &dpmRow.Brand,
			"exp_month":         &dpmRow.ExpMonth,
			"exp_year":          &dpmRow.ExpYear,
			"last4":             &dpmRow.Last4,
		}
		if !iter.MapScan(row) {
			break
		}
		dpmList = append(dpmList, dpmRow)
	}

	if err := iter.Close(); err != nil {
		panic(errors.Wrap(err, "Error closing iter"))
	}
	return dpmList
}

// NOTE: stripe.PaymentMethodCardBrand is basically a string
type DestitutePaymentMethod struct {
	StripeId        string                        `json:"stripeId"`
	PaymentMethodId string                        `json:"paymentMethodId"`
	Brand           stripe.PaymentMethodCardBrand `json:"brand"`
	ExpMonth        uint64                        `json:"expMonth"`
	ExpYear         uint64                        `json:"expYear"`
	Last4           string                        `json:"last4"`
}

//type StripeCard struct {
//Id                 string `json:"id"`
//AddressCity        string `json:"address_city"`
//AddressCountry     string `json:"address_country"`
//AddressLine1       string `json:"address_line1"`
//AddressLine1Check  string `json:"address_line1_check"`
//AddressLine2       string `json:"address_line2"`
//AddressState       string `json:"address_state"`
//AddressZip         string `json:"address_zip"`
//AddressZipCheck    string `json:"address_zip_check"`
//Brand              string `json:"id"`
//Country            string `json:"id"`
//CvcCheck           string `json:"id"`
//DynamicLast4       string `json:"id"`
//ExpMonth           int    `json:"id"`
//ExpYear            int    `json:"id"`
//Funding            string `json:"id"`
//Last4              string `json:"id"`
//Name               string `json:"id"`
//TokenizationMethod string `json:"id"`
////Object             string `json:"id"`  // always "card"
////Customer string `json:"id"`  // customer object
////Metadata {} `json:"id"`  // for custom info
////Fingerprint        string `json:"id"`  // hash of card number
//}

//func (patron *Patron) getCards() (sources []StripeSource) {
//cards
//}

//func StripeListAllCards

//curl "https://api.stripe.com/v1/customers/cus_FurjAMjdoGPXhQ/sources?object=card&limit=3" \
//-u sk_test_acCVWo94APFamzMeHxTgq5Yc00rbNtFLyS: \
//-G
