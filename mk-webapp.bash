#!/bin/bash

function MakeWebApp() {
    WEBAPP=$1
    NAV=$2
    TMPL=$3
    HTML=$4
    CWD=$(pwd)
    cd $WEBAPP
    echo "Generating $HTML"
    mkpage "nav=$NAV" $TMPL > $HTML 
    echo "Generating $WEBAPP.js and $WEBAPP.js.map"
    gopherjs build
    cd "$CWD"
}

# xlquery webapp
MakeWebApp webapp nav.md index.tmpl index.html

