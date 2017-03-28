#!/bin/bash

function MakeWebApp() {
	WEBAPP="$1"
	NAV="$2"
	TMPL="$3"
	HTML="$4"
	CWD=$(pwd)
	cd "$WEBAPP"
	echo "Generating $HTML"
	mkpage "nav=$NAV" "$TMPL" >"$HTML"
	git add "$HTML"
	echo "Generating $WEBAPP.js and $WEBAPP.js.map"
	gopherjs build
	if [ -f "$WEBAPP.js" ]; then
		git add "$WEBAPP.js"
	fi
	if [ -f "$WEBAPP.js.map" ]; then
		git add "$WEBAPP.js.map"
	fi
	cd "$CWD"
}

# xlquery webapp
MakeWebApp webapp nav.md index.tmpl index.html
