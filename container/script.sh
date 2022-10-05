#!/bin/bash -x

grep Cap /proc/$$/status
capsh --print
netcap
pscap
