package main

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

func ApiRouter() *httprouter.Router {
	r := httprouter.New()
	r.GET("/", hRoot)
	r.POST("/v1/user/signup", hUserSignUp)
	r.POST("/v1/user/signin", hUserSignIn)
	r.POST("/v1/user/get", hUserGet)
	r.POST("/v1/user/remove", hUserRemove)
	r.POST("/v1/user/field.update", hUserFieldUpdate)
	r.POST("/v1/user/location.update", hUserUpdateLocation)

	r.POST("/v1/tac/add", hAddTac)
	r.POST("/v1/tac/get.byUserUuid", hGetTacsByUserUuid)

	r.POST("/v1/event/create", hCreateEvent)
	//r.POST("/v1/event/get.attending.byUserUuid", hGetEventsAttendingByUserUuid)
	//r.POST("/v1/event/get.created.byUserUuid", hGetEventsAttendingByUserUuid)
	r.POST("/v1/event/get", hGetEventByEventUuid)
	r.POST("/v1/event/get.userAdminsEvents", hGetUserAdminsEvents)
	r.POST("/v1/event/field.update", hEventFieldUpdate)
	r.POST("/v1/event/location.change", hEventLocationChange)

	r.POST("/v1/pinGroup/search.byRegion", hPinGroupSearchByRegion)

	r.POST("/v1/mapsac/predictions", hMapsacPredictions)
	r.POST("/v1/mapsac/lookup.byPlaceId", hMapsacLookupByPlaceId)

	r.GET("/v1/users", hUsers)
	r.POST("/v1/users/search.name", hSearchUsersByName)
	r.POST("/v1/users/search.nearest", hSearchUsersByNearest)

	r.POST("/v1/verify/phone.createCode", hVerifyPhoneCreateCode)
	r.POST("/v1/verify/phone.regenCode", hVerifyPhoneRegenCode)
	r.POST("/v1/verify/phone.checkCode", hVerifyPhoneCheckCode)

	//r.POST("/v1/search/courses/name", hSearchUsersByName)
	r.POST("/v1/courses/search.nearest", hSearchCoursesByNearest)

	r.POST("/v1/shout/create", hCreateShout)
	r.POST("/v1/shout/search.nearest", hSearchShoutsByNearest)

	r.POST("/v1/convo/create", hCreateConvo)
	r.POST("/v1/convo/get.recent", hGetRecentConvos)
	r.POST("/v1/convo/get.msgs", hGetConvoMsgs)
	r.POST("/v1/convo/send", hConvoSendMsg)
	//r.POST("/v1/messages/search.nearest", hSearchShoutsByNearest)

	r.POST("/v1/planner/by.ownerUuid", hGetPlannerByOwnerUuid)
	r.POST("/v1/planner/get.updates", hGetPlannerUpdates)
	r.POST("/v1/planner/update.availableTimes", hUpdateAvailableTimes)

	r.POST("/v1/planner/accept.filledTime", hAcceptFilledTime)
	r.POST("/v1/planner/cancel.filledTime", hCancelFilledTime)
	r.POST("/v1/planner/get.filledTimes", hGetFilledTimes)
	r.POST("/v1/planner/reject.filledTime", hRejectFilledTime)
	r.POST("/v1/planner/request.filledTime", hRequestFilledTime)

	r.POST("/v1/patron/get.stripeCustomer", hGetStripeCustomer)
	r.POST("/v1/patron/get.paymentMethods", hGetPaymentMethods)
	r.POST("/v1/patron/stripe.doPayment", hStripeDoPayment)

	r.GET("/v1/patron/addCreditCardScreen", hAddCreditCardScreen)
	r.POST("/v1/patron/addCreditCardWebhook", hAddCreditCardWebhook)

	r.POST("/v1/admin/truncate.all", hTruncateAllTables)

	if Conf.Environment == "dev" {
		r.PanicHandler = catchPanics
	}

	return r
}

func hRoot(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Write([]byte("HI"))
}

func catchPanics(w http.ResponseWriter, r *http.Request, err interface{}) {
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.WriteHeader(http.StatusInternalServerError)
	var ares ApiResponse
	ares.Info = "error"

	// Use %+v for detailed stack trace
	txt := fmt.Sprintf("%+v", err)
	//txt := fmt.Sprintf("%s", err)
	ares.appendFlash(txt, "error")
	log.Error("API panic caught at router: ", txt)

	// TODO: catch when cql session dies and reconnect. example error:
	// ERRO[0005] API panic caught: runtime error: invalid memory address or nil pointer dereference

}

// debug:
func hPanic(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	err := errors.New(fmt.Sprintf("APIURL/panic GET request"))
	panic(errors.Wrap(err, "(debug) API router caught error"))
}

// debug:
func hPanicWrapped(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	err := fooError()
	panic(errors.Wrap(err, "ERROR WRAPPED"))
}

// debug:
func fooError() error {
	return errors.New("FOO ERROR")
}
