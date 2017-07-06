#!/bin/sh

seq 1 1 | xargs -I {} date -u -v- 2017-05-16{}d +%Y-%m-%d | \
    xargs -I {} curl --progress-bar -f --no-include -o {}.tsv.gz \
    -L -H "X-Papertrail-Token: wAOHfYwtG7viLeec9BQ" https://papertrailapp.com/api/v1/archives/{}/download
