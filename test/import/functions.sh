#!/bin/bash

function processWithNameIsRunning() {
    processName="$1"
    PID=$(pidof $processName)
    echo "PID of $processName $PID"
    if [ "$PID" == "" ]; then
        return 0
    else 
        return 1
    fi
}

function fileExists() {
    path="$1"
    if [ -f "$path" ]; then 
        echo "$path exists"
        return 1
    else
        echo "$path not found"
        return 0
    fi
}