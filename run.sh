#!/bin/sh -e

bin/validate-config && bin/sim | bin/plot $@
