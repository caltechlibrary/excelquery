#!/bin/bash

function MakeWebApp() {
    WEBAPP=$1
    NAV=$2
    TMPL=$3
    HTML=$4
    CWD=$(pwd)
    for P in $WEBAPP; do
        cd $P
        gopherjs build
        mkpage "nav=$NAV" $TMPL > $HTML 
        cd "$CWD"
    done
}

# xlquery webapp
MakeWebApp webapp nav.md page.tmpl index.html

