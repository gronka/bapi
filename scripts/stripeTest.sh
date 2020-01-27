#!/bin/sh
 curl -silent https://api.stripe.com/v1/tokens \
 -u pk_test_8PfL5rQ2K1SaD9cSykWHhu17009ZmNVYj7: \
 -d card[number]=4242424242424242 \
 -d card[exp_month]=12 \
 -d card[exp_year]=2019 \
 -d card[cvc]=123 | grep tok_
