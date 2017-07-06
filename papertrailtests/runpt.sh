#!/bin/sh

# if [ "$#" -ne 3 ]; then
#     echo "Arguments: <search string> <2017-06-11 00:00:00 EDT> <2017-06-12 00:00:00 EDT>"
#     exit 0
# fi

MM="07"

btnhrlcsv="btn-hrl-pt-1-3.csv"
restartcsv="restart-pt-1-3.csv"

extract(){
    SEARCH=$1
    STIME=$2
    ETIME=$3

    /usr/local/bin/papertrail -g Deployed "$SEARCH" --min-time "$STIME" --max-time "$ETIME" | cut -f 4 -d ' ' | sort | uniq -c > $4
}

foraday(){
    let d=$1
    let d1=d+1

    ptfile="btn-hrl-2017-$MM-$d.pt"

    echo "extract BUTTON_PRESS_RISKY_LIFTS_GOAL 2017-$MM-$d 00:00:00 EDT 2017-$MM-$d1 00:00:00 EDT $ptfile"
    extract "BUTTON_PRESS_RISKY_LIFTS_GOAL" "2017-$MM-$d 00:00:00 EDT" "2017-$MM-$d1 00:00:00 EDT" $ptfile

    echo "go run *.go $ptfile $d >> $btnhrlcsv"
    go run *.go $ptfile $d >> $btnhrlcsv

    ptfile="restart-2017-$MM-$d.pt"
    echo "extract Started /home/kinetic/device/c/uart 2017-$MM-$d 00:00:00 EDT 2017-$MM-$d1 00:00:00 EDT $ptfile"
    extract "Started /home/kinetic/device/c/uart" "2017-$MM-$d 00:00:00 EDT" "2017-$MM-$d1 00:00:00 EDT" $ptfile

    echo "go run *.go $ptfile $d >> $restartcsv"
    go run *.go $ptfile $d >> $restartcsv
}

for i in 01 02 03
do
    foraday $i
done


# /usr/local/bin/papertrail -g Deployed "BUTTON_PRESS_RISKY_LIFTS_GOAL" --min-time "2017-06-26 00:00:00 EDT" --max-time "2017-06-27 00:00:00 EDT" | cut -f 4 -d ' ' | sort | uniq -c > tmp.b.pt
